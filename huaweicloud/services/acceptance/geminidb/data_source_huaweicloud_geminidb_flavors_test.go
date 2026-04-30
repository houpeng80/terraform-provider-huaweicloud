package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBFlavors_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.az_status.%"),

					resource.TestCheckOutput("is_engine_name_filter_useful", "true"),
					resource.TestCheckOutput("is_engine_and_mode_filter_useful", "true"),
					resource.TestCheckOutput("is_engine_and_product_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBFlavors_basic() string {
	return `
data "huaweicloud_geminidb_flavors" "test" {}

data "huaweicloud_geminidb_flavors" "engine_name_filter" {
  engine_name = "redis"
}

output "is_engine_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_flavors.engine_name_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_flavors.engine_name_filter.flavors[*].engine_name :
    v == "redis"]
  )
}

data "huaweicloud_geminidb_flavors" "engine_and_mode_filter" {
  engine_name = "influxdb"
  mode        = "EnhancedCluster"
}

output "is_engine_and_mode_filter_useful" {
  value = length(data.huaweicloud_geminidb_flavors.engine_and_mode_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_flavors.engine_and_mode_filter.flavors[*].engine_name :
    v == "influxdb"]
  )
}

data "huaweicloud_geminidb_flavors" "engine_and_product_type_filter" {
  engine_name  = "redis"
  mode         = "CloudNativeCluster"
  product_type = "Standard"
}

output "is_engine_and_product_type_filter_useful" {
  value = length(data.huaweicloud_geminidb_flavors.engine_and_product_type_filter.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_flavors.engine_and_product_type_filter.flavors[*].engine_name :
    v == "redis"]
  )
}
`
}
