---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: resource_opentelekomcloud_cts_tracker_v1"
sidebar_current: "docs-opentelekomcloud-resource-cts-tracker-v1"
description: |-
   CTS tracker allows you to collect, store, and query cloud resource operation records and use these records for security analysis, compliance auditing, resource tracking, and fault locating.
---

# opentelekomcloud_cts_tracker_v1

Allows you to collect, store, and query cloud resource operation records.

## Example Usage

 ```hcl
 variable "bucket_name" { }
 variable "topic_id" { }
 
 resource "opentelekomcloud_cts_tracker_v1" "tracker_v1" {
   bucket_name      = "${var.bucket_name}"
   file_prefix_name      = "yO8Q"
   is_support_smn = true
   topic_id = "${var.topic_id}"
   is_send_all_key_operation = false
   operations = ["login"]
   need_notify_user_list = ["user1"]
 }

 ```
## Argument Reference
The following arguments are supported:

* `bucket_name` - (Required) The OBS bucket name for a tracker.

* `file_prefix_name` - (Optional) The prefix of a log that needs to be stored in an OBS bucket. 

* `is_support_smn` - (Required) Specifies whether SMN is supported. When the value is false, topic_id and operations can be left empty.

* `topic_id` - (Required)The theme of the SMN service, Is obtained from SMN and in the format of **urn:smn:([a-z]|[A-Z]|[0-9]|\-){1,32}:([a-z]|[A-Z]|[0-9]){32}:([a-z]|[A-Z]|[0-9]|\-|\_){1,256}**.

* `operations` - (Required) Trigger conditions for sending a notification.

* `is_send_all_key_operation` - (Required) When the value is **false**, operations cannot be left empty.

* `need_notify_user_list` - (Optional) The users using the login function. When these users log in, notifications will be sent.

* `status` - (Optional) The status of a tracker. The value should be **enabled** when creating a tracker, and when updating the value can be enabled or disabled.



## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `tracker_name` - The tracker name. Currently, only tracker **system** is available.


## Import

CTS tracker can be imported using  `tracker_name`, e.g.

```
$ terraform import opentelekomcloud_cts_tracker_v1.tracker system
```




