---
layout: "opentelekomcloud"
page_title: "OpentelekomCloud: opentelekomcloud_nat_dnat_rule_v2"
sidebar_current: "docs-opentelekomcloud-resource-nat-dnat-rule-v2"
description: |-
  Manages a V2 dnat rule resource within OpentelekomCloud Nat.
---

# opentelekomcloud\_nat\_dnat\_rule_v2


## Example Usage

### Dnat

```hcl
resource "opentelekomcloud_nat_dnat_rule_v2" "dnat_1" {
  floating_ip_id = "2bd659ab-bbf7-43d7-928b-9ee6a10de3ef"
  nat_gateway_id = "bf99c679-9f41-4dac-8513-9c9228e713e1"
  private_ip = "10.0.0.12"
  internal_service_port = 993
  protocol = "tcp"
  external_service_port = 242
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_id` - (Required) Specifies the ID of the floating IP address.
  Changing this creates a new resource.

* `internal_service_port` - (Required) Specifies port used by ECSs or BMSs
  to provide services for external systems. Changing this creates a new resource.

* `nat_gateway_id` - (Required) ID of the nat gateway this dnat rule belongs to.
   Changing this creates a new dnat rule.

* `port_id` - (Optional) Specifies the port ID of an ECS or a BMS.
  This parameter and private_ip are alternative. Changing this creates a
  new dnat rule.

* `private_ip` - (Optional) Specifies the private IP address of a
  user, for example, the IP address of a VPC for dedicated connection.
  This parameter and port_id are alternative.
  Changing this creates a new dnat rule.

* `protocol` - (Required) Specifies the protocol type. Currently,
  TCP, UDP, and ANY are supported.
  Changing this creates a new dnat rule.

* `external_service_port` - (Required) Specifies port used by ECSs or
  BMSs to provide services for external systems.
  Changing this creates a new dnat rule.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` - Dnat rule creation time.

* `status` - Dnat rule status.

* `floating_ip_address` - The actual floating IP address.

## Import

Dnat can be imported using the following format:

```
$ terraform import opentelekomcloud_nat_dnat_rule_v2.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```
