package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "name", "opentelekomcloud-cce"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "status", "Available"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "cluster_version", "v1.9.2-r2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "authentication_mode", "x509"),
				),
			},
			{
				Config: testAccCCEClusterV3_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_timeout(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists("opentelekomcloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_cce_cluster_v3.cluster_1", "authentication_mode", "rbac"),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating opentelekomcloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

var testAccCCEClusterV3_basic = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "opentelekomcloud-cce"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.9.2-r2"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_update = fmt.Sprintf(`
resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "opentelekomcloud-cce"
  cluster_type="VirtualMachine"
  flavor_id="cce.s1.small"
  cluster_version = "v1.9.2-r2"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  description="new description"
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccCCEClusterV3_timeout = fmt.Sprintf(`
resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name = "opentelekomcloud-cce"
  cluster_type="VirtualMachine"
  flavor_id="cce.s2.small"
  cluster_version = "v1.9.2-r2"
  vpc_id="%s"
  subnet_id="%s"
  eip="${opentelekomcloud_networking_floatingip_v2.fip_1.address}"
  container_network_type="overlay_l2"
  authentication_mode = "rbac"
    timeouts {
    create = "20m"
    delete = "10m"
  }
  multi_az = true
}
`, OS_VPC_ID, OS_NETWORK_ID)
