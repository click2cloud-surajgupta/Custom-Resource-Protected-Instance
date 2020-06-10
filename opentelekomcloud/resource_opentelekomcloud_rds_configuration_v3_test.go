package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rds/v3/configurations"
)

func TestAccRdsConfigurationV3_basic(t *testing.T) {
	var config configurations.Configuration

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsConfigV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfigV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists("opentelekomcloud_rds_parametergroup_v3.pg_1", &config),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_parametergroup_v3.pg_1", "name", "pg_1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_parametergroup_v3.pg_1", "description", "description_1"),
				),
			},
			{
				Config: testAccRdsConfigV3_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsConfigV3Exists("opentelekomcloud_rds_parametergroup_v3.pg_1", &config),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_parametergroup_v3.pg_1", "name", "pg_update"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rds_parametergroup_v3.pg_1", "description", "description_update"),
				),
			},
		},
	})
}

func testAccCheckRdsConfigV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	rdsClient, err := config.rdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud RDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rds_parametergroup_v3" {
			continue
		}

		_, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Rds configuration still exists")
		}
	}

	return nil
}

func testAccCheckRdsConfigV3Exists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		rdsClient, err := config.rdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud RDS client: %s", err)
		}

		found, err := configurations.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Rds configuration not found")
		}

		*configuration = *found

		return nil
	}
}

const testAccRdsConfigV3_basic = `
resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
	name = "pg_1"
	description = "description_1"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "mysql"
		version = "5.6"
	}
}
`

const testAccRdsConfigV3_update = `
resource "opentelekomcloud_rds_parametergroup_v3" "pg_1" {
	name = "pg_update"
	description = "description_update"
	values = {
		max_connections = "10"
		autocommit = "OFF"
	}
	datastore {
		type = "mysql"
		version = "5.6"
	}
}
`
