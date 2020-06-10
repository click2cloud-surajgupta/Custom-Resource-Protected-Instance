package opentelekomcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/recordsets"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"
)

func resourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordSetV2Create,
		Read:   resourceDNSRecordSetV2Read,
		Update: resourceDNSRecordSetV2Update,
		Delete: resourceDNSRecordSetV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"records": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  300,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceDNSRecordSetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
	}

	recordsraw := d.Get("records").(*schema.Set).List()
	records := make([]string, len(recordsraw))
	for i, recordraw := range recordsraw {
		records[i] = recordraw.(string)
	}

	createOpts := RecordSetCreateOpts{
		recordsets.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Records:     records,
			TTL:         d.Get("ttl").(int),
			Type:        d.Get("type").(string),
		},
		MapValueSpecs(d),
	}

	zoneID := d.Get("zone_id").(string)

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}

	id := fmt.Sprintf("%s/%s", zoneID, n.ID)
	d.SetId(id)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
		if err != nil {
			return fmt.Errorf("Error getting resource type of DNS record set %s: %s", n.ID, err)
		}

		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(dnsClient, resourceType, n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of DNS record set %s: %s", n.ID, tagErr)
		}
	}

	log.Printf("[DEBUG] Created OpenTelekomCloud DNS record set %s: %#v", n.ID, n)
	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "record_set")
	}

	log.Printf("[DEBUG] Retrieved  record set %s: %#v", recordsetID, n)

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("ttl", n.TTL)
	d.Set("type", n.Type)
	if err := d.Set("records", n.Records); err != nil {
		return fmt.Errorf("[DEBUG] Error saving records to state for OpenTelekomCloud DNS record set (%s): %s", d.Id(), err)
	}
	d.Set("region", GetRegion(d, config))
	d.Set("zone_id", zoneID)

	// save tags
	resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("Error getting resource type of DNS record set %s: %s", recordsetID, err)
	}
	resourceTags, err := tags.Get(dnsClient, resourceType, recordsetID).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching OpenTelekomCloud DNS record set tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for OpenTelekomCloud DNS record set %s: %s", recordsetID, err)
	}

	return nil
}

func resourceDNSRecordSetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
	}

	var updateOpts recordsets.UpdateOpts
	if d.HasChange("ttl") {
		updateOpts.TTL = d.Get("ttl").(int)
	}

	// `records` is required attribute for update request
	recordsRaw := d.Get("records").(*schema.Set).List()
	records := make([]string, len(recordsRaw))
	for i, recordRaw := range recordsRaw {
		records[i] = recordRaw.(string)
	}
	updateOpts.Records = records

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating  record set %s with options: %#v", recordsetID, updateOpts)

	_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud DNS  record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to update", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become ACTIVE for updation: %s",
			recordsetID, err)
	}

	// update tags
	resourceType, err := getDNSRecordSetResourceType(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("Error getting resource type of DNS record set %s: %s", d.Id(), err)
	}

	tagErr := UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
	if tagErr != nil {
		return fmt.Errorf("Error updating tags of DNS record set %s: %s", d.Id(), tagErr)
	}

	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to be deleted", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become DELETED for deletion: %s",
			recordsetID, err)
	}

	d.SetId("")
	return nil
}

func waitForDNSRecordSet(dnsClient *golangsdk.ServiceClient, zoneID, recordsetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recordset, err := recordsets.Get(dnsClient, zoneID, recordsetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return recordset, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud DNS record set (%s) current status: %s", recordset.ID, recordset.Status)
		return recordset, parseStatus(recordset.Status), nil
	}
}

func parseDNSV2RecordSetID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("Unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

// get resource type of DNS record set from zone_id
func getDNSRecordSetResourceType(client *golangsdk.ServiceClient, zone_id string) (string, error) {
	zone, err := zones.Get(client, zone_id).Extract()
	if err != nil {
		return "", err
	}

	zoneType := zone.ZoneType
	if zoneType == "public" {
		return "DNS-public_recordset", nil
	} else if zoneType == "private" {
		return "DNS-private_recordset", nil
	}
	return "", fmt.Errorf("invalid zone type: %s", zoneType)
}
