---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_database_operation"
description: |-
  Use this resource to perform operations on GeminiDB database within HuaweiCloud.
---

# huaweicloud_geminidb_database_operation

Use this resource to perform operations on GeminiDB database within HuaweiCloud.

## Example Usage

### Flush All Data

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_database_operation" "test" {
  instance_id = var.instance_id
  action      = "flush"
}
```

### Flush Specific Database

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_database_operation" "test" {
  instance_id = var.instance_id
  action      = "flush"
  db_id       = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new rds instance resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB Redis instance.

* `action` - (Required, String, NonUpdatable) Specifies the operation to perform on the database.

* `db_id` - (Optional, Int, NonUpdatable) Specifies the DB ID to be cleared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
