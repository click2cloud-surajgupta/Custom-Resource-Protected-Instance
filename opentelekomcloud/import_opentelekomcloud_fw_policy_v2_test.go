package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccFWPolicyV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_fw_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_addRules,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
