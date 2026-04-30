package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDbAccountsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_geminidb_accounts.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDbAccountsDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.privilege"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.databases.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccGeminiDbAccountsDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_geminidb_accounts" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id

  depends_on = [huaweicloud_geminidb_account.test]
}

data "huaweicloud_geminidb_accounts" "name_filter" {
  instance_id = huaweicloud_geminidb_instance.test.id
  name        = "%[2]s"

  depends_on = [huaweicloud_geminidb_account.test]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_accounts.name_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_accounts.name_filter.users[*].name :
    v == "%[2]s"]
  )
}
`, testAccGeminiDbAccount_basic(rName), rName)
}
