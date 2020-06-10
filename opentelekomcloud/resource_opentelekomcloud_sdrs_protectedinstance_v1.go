package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sdrs/v1/protectedinstances"
	"log"
	"time"
)

func resourceSdrsProtectedinstancesV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsProtectedinstanceV1Create,
		Read:   resourceSdrsProtectedinstanceV1Read,
		Update: resourceSdrsProtectedinstanceV1Update,
		Delete: resourceSdrsProtectedinstanceV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"server_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			}, "server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			}, "name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}
func resourceSdrsProtectedinstanceV1Create(d *schema.ResourceData, meta interface{}) error{
	config := meta.(*Config)
	sdrsClient, err := config.SdrsV1Client(GetRegion(d, config))

	fmt.Println(sdrsClient)
	if err!=nil{
		log.Panic(err)
	}
	createOpts := protectedinstances.CreateOpts{
		GroupID: d.Get("server_group_id").(string),
		ServerID:  d.Get("server_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SubnetID:    d.Get("primary_subnet_id").(string),
		IpAddress:    d.Get("primary_ip_address").(string),

	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := protectedinstances.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil{
		return fmt.Errorf("Error creating OpenTelekomcomCloud SDRS Protectedinstances: %s", err)
	}

	if err := protectedinstances.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := protectedinstances.GetJobEntity(sdrsClient, n.JobID, "protected_instance_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceSdrsProtectedinstanceV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSdrsProtectedinstancesV1Create.")

}
func resourceSdrsProtectedinstanceV1Read(d *schema.ResourceData, meta interface{})error{
	config := meta.(*Config)
	sdrsClient, err := config.SdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS client: %s", err)
	}
	n, err := protectedinstances.Get(sdrsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud SDRS Protectedinstance: %s", err)
	}
	d.Set("server_group_id",n.GroupID)
	d.Set("server_id",n.Id)
	d.Set("name", n.Name)
	d.Set("description", n.Description)

	return nil
}
func resourceSdrsProtectedinstanceV1Update(d *schema.ResourceData, meta interface{}) error{
	config := meta.(*Config)
	sdrsClient, err := config.SdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS Client: %s", err)
	}
	var updateOpts protectedinstances.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
	_, err = protectedinstances.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud SDRS Protectedinstance: %s", err)
	}
	return resourceSdrsProtectedinstanceV1Read(d, meta)
}
func resourceSdrsProtectedinstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := config.SdrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud SDRS client: %s", err)
	}
	var deleteOpts protectedinstances.DeleteOpts
	n, err := protectedinstances.Delete(sdrsClient, d.Id(),deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud SDRS Protectedinstance: %s", err)
	}

	if err := protectedinstances.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}