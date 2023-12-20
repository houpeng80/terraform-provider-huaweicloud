package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// 命令行参数
	basePath           string
	filterFilePath     string
	outputDir          string
	version            string
	providerSchemaPath string
	provider           string

	// 保存 openstack/instances.{func} : uri
	urlSupportsInUriFile = make(map[string]string)
)

func init() {
	flag.StringVar(&basePath, "basePath", "./", "base Path")
	flag.StringVar(&outputDir, "outputDir", "./api/", "api yaml file output Dir")
	flag.StringVar(&version, "version", "", "provider version")
	flag.StringVar(&filterFilePath, "filterFilePath", "", "Specifies the terraform resource been scan")
	flag.StringVar(&providerSchemaPath, "providerSchemaPath", "schema.json",
		"CMD: terraform providers schema -json >./schema.json")
	flag.StringVar(&provider, "provider", "huaweicloud", "过滤指定provider输出")

}
func main() {
	flag.Parse()
	log.Printf("start dynamic scan")
	resourceMap := huaweicloud.Provider().ResourcesMap
	dataSourceMap := huaweicloud.Provider().DataSourcesMap
	////resource := resourceMap["huaweicloud_rds_sql_audit"]
	dealFiles(resourceMap, dataSourceMap)

	// 重写 vendor/github.com/chnsz/golangsdk/service_client.go 中 Request 方法，不发送请求，而是记录当前请求信息
	clientPath := basePath + "/vendor/github.com/chnsz/golangsdk/service_client.go"
	rewriteClientRequest(clientPath)

	// 重写 huaweicloud/provider.go 中校验登录信息逻辑，直接创建一个client
	providerPath := basePath + "/huaweicloud/provider.go"
	fixProvider(providerPath)

	// 执行所有资源create、update、read、delete方法，找到对应的接口
	generateApiFile()
}

func generateApiFile() {

	p := huaweicloud.Provider()

	for dataSourceName, dataSource := range p.DataSourcesMap {
		p = makeProvider(dataSourceName)
		diffState := &terraform.InstanceDiff{
			Attributes: make(map[string]*terraform.ResourceAttrDiff),
			Meta:       make(map[string]interface{}),
			Destroy:    true,
		}
		// 查询
		dataSource.ReadDataApply(context.Background(), diffState, p.Meta())
	}

	//for resourceName, resource := range p.ResourcesMap {
	//	p.Schema["cloud"] = nil
	//	fmt.Println(resourceName)
	//	fmt.Println(resource)
	//	initState := &terraform.InstanceState{}
	//	diffState := &terraform.InstanceDiff{
	//		Attributes: make(map[string]*terraform.ResourceAttrDiff),
	//		Meta:       make(map[string]interface{}),
	//		Destroy:    true,
	//	}
	//	// 创建
	//	createState, _ := resource.Apply(context.Background(), initState, diffState, p.Meta())
	//	// 更新
	//	updateState, _ := resource.Apply(context.Background(), initState, diffState, p.Meta())
	//
	//}
}

func makeProvider(dataSourceName string) *schema.Provider {
	p := huaweicloud.Provider()
	p.Schema["region"].Default = "generate_api_doc_region"
	p.Schema["project_id"].Default = "generate_api_doc_project_id"
	raw := make(map[string]interface{})
	p.Schema["cloud"].Default = fmt.Sprintf("data_source_%s", dataSourceName)
	p.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	return p
}

//func buildYaml() {
//
//	err := os.WriteFile(outputFile, []byte(yarmStr), 0664)
//	if err == nil {
//		log.Println("写入成功", outputFile)
//	}
//}

func fixProvider(filePath string) {
	resourceFileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileStr := string(resourceFileBytes)
	// 增加 import
	importReg := regexp.MustCompile(fmt.Sprintf("(import \\()"))
	allIfMatch := importReg.FindAllStringSubmatch(fileStr, -1)
	fileStr = strings.ReplaceAll(fileStr, allIfMatch[0][0], fmt.Sprintf("%s\n\t\"huaweisdk \"github.com/chnsz/golangsdk/openstack\"", allIfMatch[0][0]))

	// 替换掉生成client逻辑
	source := `if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}`
	target := `client, _ := huaweisdk.NewClient("")
	config.HwClient = client
	config.DomainClient = client
	config.RegionProjectIDMap[region] = tenantID`
	fileStr = strings.ReplaceAll(fileStr, source, target)
	fmt.Println(fileStr)
	//if err = ioutil.WriteFile(filePath, []byte(fileStr), 0664); err != nil {
	//	log.Fatal(err)
	//}
}

func rewriteClientRequest(filePath string) {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, filePath, nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse package %s: %s\n", filePath, err)
		os.Exit(1)
	}
	resourceFileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileStr := string(resourceFileBytes)
	// 增加 import
	importReg := regexp.MustCompile(fmt.Sprintf("(import \\()"))
	allIfMatch := importReg.FindAllStringSubmatch(fileStr, -1)
	fileStr = strings.ReplaceAll(fileStr, allIfMatch[0][0], fmt.Sprintf("%s\n\t\"fmt\"\n\t\"log\"\n\t\"os\"", allIfMatch[0][0]))

	// 替换request函数
	allResourceFileFunc := findAllFunc(file)
	for _, fn := range allResourceFileFunc {
		if fn.Name.Name == "Request" {
			startIndex := fSet.Position(fn.Pos()).Offset
			endIndex := fSet.Position(fn.End()).Offset
			funcSrc := string(resourceFileBytes[startIndex:endIndex])
			fmt.Println(funcSrc)
			funcTarget := `func (client *ServiceClient) Request1(method, url string, options *RequestOpts) (*http.Response, error) {
	parts := strings.Split(url, ".")
	product := strings.Split(parts[0], "//")[1]
	path := parts[2]
	resourceName := path[:strings.Index(path, "/")]
	url = path[strings.Index(path, "/"):]

	fileName := fmt.Sprintf("./temp/%s.txt", resourceName)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := fmt.Sprintf("%s %s %s\n", url, "GET", product)
	if _, err = file.WriteString(data); err != nil {
		log.Fatal(err)
	}
	return &http.Response{}, nil
}`
			fileStr = strings.ReplaceAll(fileStr, funcSrc, funcTarget)
			fmt.Println(fileStr)
		}
	}
	if err = ioutil.WriteFile(filePath, []byte(fileStr), 0664); err != nil {
		log.Fatal(err)
	}
}

func dealFiles(resourceMap, dataSourceMap map[string]*schema.Resource) {
	subPackagePath := basePath + provider + "/"
	err := filepath.Walk(subPackagePath, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			log.Printf("scan path %s failed: %s\n", path, err)
			return err
		}

		if fInfo.IsDir() && !isSkipDirectory(path) {

			dealFile(path, resourceMap, dataSourceMap)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("ERROR: scan path failed: %s\n", err)
	}
}

func dealFile(path string, resourceMap, dataSourceMap map[string]*schema.Resource) {
	fSet := token.NewFileSet()
	packs, err := parser.ParseDir(fSet, path, nil, 0)
	if err != nil {
		fmt.Printf("Failed to parse package %s: %s\n", path, err)
		os.Exit(1)
	}
	for _, pack := range packs {
		packageName := pack.Name
		log.Printf("package name: %s, file count: %d\n", packageName, len(pack.Files))
		for filePath, file := range pack.Files {
			// 获得文件名并去除版本号
			resourceName := filePath[strings.LastIndex(filePath, "\\")+1 : len(filePath)-3]
			re, _ := regexp.Compile(`_v\d+$`)
			resourceName = re.ReplaceAllString(resourceName, "")

			if !isExportResource(resourceName, resourceMap, dataSourceMap) {
				continue
			}

			resourceFileBytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatal(err)
			}
			fileStr := string(resourceFileBytes)
			//resFileStr := fmt.Sprintf("package %s\n", file.Name.Name)
			//resFileStr = addImports(resFileStr, file)
			allResourceFileFunc := findAllFunc(file)

			for _, fn := range allResourceFileFunc {
				startIndex := fSet.Position(fn.Pos()).Offset
				endIndex := fSet.Position(fn.End()).Offset
				funcSrc := string(resourceFileBytes[startIndex:endIndex])
				funcTarget := dealFunc(funcSrc)
				fileStr = strings.ReplaceAll(fileStr, funcSrc, funcTarget)
			}
			// 重新写回文件
			if err = ioutil.WriteFile(filePath, []byte(fileStr), 0664); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func dealFunc(funcSrc string) string {
	funcTarget := funcSrc
	// return都使用 fmt.println替代, 只保留最后一个
	returnReg := regexp.MustCompile(fmt.Sprintf("(return.*\n\t*})"))
	allReturnMatch := returnReg.FindAllStringSubmatch(funcTarget, -1)
	for i, returnMatch := range allReturnMatch {
		if i == len(allReturnMatch)-1 {
			break
		}
		fmt.Println(funcTarget)
		fmt.Println(returnMatch[0])
		s := fmt.Sprintf("fmt.Println(%s)}", returnMatch[0][len("return"):strings.LastIndex(returnMatch[0], "}")])
		funcTarget = strings.Replace(funcTarget, returnMatch[0], s, 1)
	}

	// 如果修改了return，就需要判断下是否需要添加import fmt
	if len(allReturnMatch) > 1 {

	}

	// 替换 `if expression {` 为 `if expression || true` {
	ifReg := regexp.MustCompile(fmt.Sprintf("(if\\s.*\\s{)"))
	allIfMatch := ifReg.FindAllStringSubmatch(funcTarget, -1)
	for _, ifMatch := range allIfMatch {
		funcTarget = strings.Replace(funcTarget, ifMatch[0], fmt.Sprintf(`%s || true {`, ifMatch[0][:len(ifMatch[0])-1]), 1)
	}

	// 替换 `} else {` 为 `}\n\t if  true {`
	elseReg := regexp.MustCompile(fmt.Sprintf("(}\\s*else\\s*{)"))
	allElseMatch := elseReg.FindAllStringSubmatch(funcTarget, -1)
	for _, elseMatch := range allElseMatch {
		funcTarget = strings.Replace(funcTarget, elseMatch[0], "}\n\tif true {", 1)
	}

	// 替换 `} if else expression {` 为 `}\n\t if expression || true {`
	// 经过前边的替换 if 语句，此时的 `if else expression` 已经被替换为 `if expression || true`, 所以只需要将else替换为换行即可
	ifElseReg := regexp.MustCompile(fmt.Sprintf("(}\\s*else if\\s.*\\s{)"))
	allIfElseMatch := ifElseReg.FindAllStringSubmatch(funcTarget, -1)
	for _, ifElseMatch := range allIfElseMatch {
		funcTarget = strings.Replace(funcTarget, ifElseMatch[0], fmt.Sprintf("}\n\t%s", ifElseMatch[0][strings.Index(ifElseMatch[0], "else")+4:]), 1)
	}
	fmt.Println(funcTarget)

	return funcTarget
}

func addImports(resFileStr string, file *ast.File) string {
	resFileStr = fmt.Sprintf("%simport (\n", resFileStr)
	for _, v := range file.Imports {
		resFileStr = fmt.Sprintf("%s\t%s\n", resFileStr, v.Path.Value)
	}
	resFileStr = fmt.Sprintf("%s)", resFileStr)
	return resFileStr
}

func findAllFunc(f *ast.File) []*ast.FuncDecl {
	funcs := make([]*ast.FuncDecl, 0)

	for _, d := range f.Decls {
		if fn, isFn := d.(*ast.FuncDecl); isFn {
			funcs = append(funcs, fn)
		}
	}

	return funcs
}

func isExportResource(resourceName string, resourceMap, dataSourceMap map[string]*schema.Resource) bool {
	if strings.HasPrefix(resourceName, "resource_") {
		if _, ok := resourceMap[resourceName[len("resource_"):]]; ok {
			return true
		}
		return false
	}

	if strings.HasPrefix(resourceName, "data_source_") {
		if _, ok := dataSourceMap[resourceName[len("data_source_"):]]; ok {
			return true
		}
		return false
	}
	return false
}

func isSkipDirectory(path string) bool {
	var skipKeys = []string{
		"acceptance", "utils", "internal", "helper", "deprecated",
	}

	for _, sub := range skipKeys {
		if strings.Contains(path, sub) {
			return true
		}
	}
	return false
}
