package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBDatabaseOperation_basic(t *testing.T) {
	resourceName := "huaweicloud_geminidb_database_operation.test"
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBDatabaseOperation_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "flush"),
				),
			},
			{
				Config: testAccGeminiDBDatabaseOperation_withDbId(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "flush"),
					resource.TestCheckResourceAttr(resourceName, "db_id", "1"),
				),
			},
		},
	})
}

func testAccGeminiDBDatabaseOperation_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_database_operation" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  action      = "flush"
}
`, testAccGeminiDbInstance_basic(name))
}

func testAccGeminiDBDatabaseOperation_withDbId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_database_operation" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  action      = "flush"
  db_id       = 1
}
`, testAccGeminiDbInstance_basic(name))
}
