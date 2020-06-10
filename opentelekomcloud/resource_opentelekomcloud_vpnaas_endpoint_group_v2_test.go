package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func TestAccVpnGroupV2_basic(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"opentelekomcloud_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
				),
			},
			{
				Config: testAccEndpointGroupV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"opentelekomcloud_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
				),
			},
		},
	})
}

func testAccCheckEndpointGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpnaas_group" {
			continue
		}
		_, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EndpointGroup (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckEndpointGroupV2Exists(n string, group *endpointgroups.EndpointGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		var found *endpointgroups.EndpointGroup

		found, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*group = *found

		return nil
	}
}

var testAccEndpointGroupV2_basic = `
	resource "opentelekomcloud_vpnaas_endpoint_group_v2" "group_1" {
		name = "Group 1"
		type = "cidr"
		endpoints = ["10.3.0.0/24",
			"10.2.0.0/24",]
	}
`

var testAccEndpointGroupV2_update = `
	resource "opentelekomcloud_vpnaas_endpoint_group_v2" "group_1" {
		name = "Updated Group 1"
		type = "cidr"
		endpoints = ["10.2.0.0/24",
			"10.3.0.0/24",]
	}
`
