package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccFWRuleV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_fw_rule_v2.rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2_basic_2,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
