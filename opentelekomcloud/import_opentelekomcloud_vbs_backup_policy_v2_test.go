package opentelekomcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVBSBackupPolicyV2_importBasic(t *testing.T) {
	resourceName := "opentelekomcloud_vbs_backup_policy_v2.vbs"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
