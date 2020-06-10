package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOpenTelekomCloudImagesV2ImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_cirros,
			},
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.opentelekomcloud_images_image_v2.image_1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_images_image_v2.image_1", "name", "CirrOS-tf"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_images_image_v2.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_images_image_v2.image_1", "min_ram_mb", "0"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_images_image_v2.image_1", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_images_image_v2.image_1", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccOpenTelekomCloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_cirros,
			},
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_queryTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.opentelekomcloud_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.opentelekomcloud_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.opentelekomcloud_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccOpenTelekomCloudImagesV2ImageDataSource_cirros,
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

// Standard CirrOS image
const testAccOpenTelekomCloudImagesV2ImageDataSource_cirros = `
resource "opentelekomcloud_images_image_v2" "image_1" {
	name = "CirrOS-tf"
	container_format = "bare"
	disk_format = "qcow2"
	image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
	tags = ["cirros-tf"]
}
`

var testAccOpenTelekomCloudImagesV2ImageDataSource_basic = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
	most_recent = true
	name = "${opentelekomcloud_images_image_v2.image_1.name}"
}
`, testAccOpenTelekomCloudImagesV2ImageDataSource_cirros)

var testAccOpenTelekomCloudImagesV2ImageDataSource_queryTag = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "cirros-tf"
}
`, testAccOpenTelekomCloudImagesV2ImageDataSource_cirros)

var testAccOpenTelekomCloudImagesV2ImageDataSource_querySizeMin = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccOpenTelekomCloudImagesV2ImageDataSource_cirros)

var testAccOpenTelekomCloudImagesV2ImageDataSource_querySizeMax = fmt.Sprintf(`
%s

data "opentelekomcloud_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccOpenTelekomCloudImagesV2ImageDataSource_cirros)
