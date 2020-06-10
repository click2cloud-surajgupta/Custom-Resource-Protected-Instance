package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOTCRtsSoftwareConfigV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsSoftwareConfigV1DataSourceID("data.opentelekomcloud_rts_software_config_v1.configs"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_config_v1.configs", "name", "opentelekomcloud-config"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_rts_software_config_v1.configs", "group", "script"),
				),
			},
		},
	})
}

func testAccCheckRtsSoftwareConfigV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find software config data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS software config data source ID not set ")
		}

		return nil
	}
}

var testAccRtsSoftwareConfigV1DataSource_basic = `
resource "opentelekomcloud_rts_software_config_v1" "config_1" {
  name = "opentelekomcloud-config"
  output_values = [{
    type = "String"
    name = "result"
    error_output = "false"
    description = "value1"
  }]
  input_values = [{
    default = "0"
    type = "String"
    name = "foo"
    description = "value2"
  }]
  group = "script"
}

data "opentelekomcloud_rts_software_config_v1" "configs" {
  id = "${opentelekomcloud_rts_software_config_v1.config_1.id}"
}
`
