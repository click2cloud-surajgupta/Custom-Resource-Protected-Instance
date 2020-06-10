package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
)

func TestAccVpnIPSecPolicyV2_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "tenant_id", &policy.TenantID),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "transform_protocol", &policy.TransformProtocol),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "encapsulation_mode", &policy.EncapsulationMode),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "encryption_algorithm", &policy.EncryptionAlgorithm),
				),
			},
			{
				Config: testAccIPSecPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVpnIPSecPolicyV2_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_vpnaas_ipsec_policy_v2" {
			continue
		}
		_, err = ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IPSec policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIPSecPolicyV2Exists(n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
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

		found, err := ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIPSecPolicyV2_basic = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
}
`

const testAccIPSecPolicyV2_Update = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIPSecPolicyV2_withLifetime = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "md5"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIPSecPolicyV2_withLifetimeUpdate = `
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "md5"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
