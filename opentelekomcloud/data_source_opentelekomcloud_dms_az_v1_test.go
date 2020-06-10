package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.opentelekomcloud_dms_az_v1.az1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "name", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "port", "8002"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "code", OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_dms_az_v1" "az1" {
  name = "%s"
  port = "8002"
}
`, OS_AVAILABILITY_ZONE)
