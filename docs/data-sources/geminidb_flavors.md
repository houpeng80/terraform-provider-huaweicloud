---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_flavors"
description: |-
  Use this data source to query the list of GeminiDB instance flavors within HuaweiCloud.
---

# huaweicloud_geminidb_flavors

Use this data source to query the list of GeminiDB instance flavors within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_geminidb_flavors" "test" {}
```

### Filter by Database Engine

```hcl
data "huaweicloud_geminidb_flavors" "test" {
  engine_name = "redis"
}
```

### Filter by Engine and Mode

```hcl
data "huaweicloud_geminidb_flavors" "test" {
  engine_name = "redis"
  mode        = "EnhancedCluster"
}
```

### Filter by Engine, Mode and Product Type

```hcl
data "huaweicloud_geminidb_flavors" "test" {
  engine_name  = "redis"
  mode         = "CloudNativeCluster"
  product_type = "Standard"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the flavors.
  If omitted, the provider-level region will be used.

* `engine_name` - (Optional, String) Specifies the database engine name.
  The valid values are as follows:
  + **cassandra**: GeminiDB Cassandra.
  + **mongodb**: GeminiDB Mongo.
  + **influxdb**: GeminiDB Influx.
  + **redis**: GeminiDB Redis.

* `mode` - (Optional, String) Specifies the instance mode.
  The valid values are as follows:
  + **CloudNativeCluster**: Cloud-native deployment mode.
  + **EnhancedCluster**: Enhanced cluster mode for GeminiDB Influx.

* `product_type` - (Optional, String) Specifies the product type.
  This parameter is required when querying GeminiDB Redis cloud-native cluster flavors.
  The valid values are as follows:
  + **Standard**: Standard type.
  + **Capacity**: Capacity type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of GeminiDB instance flavors.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `engine_name` - The database engine name.

* `engine_version` - The database engine version.

* `vcpus` - The number of CPU cores.

* `ram` - The memory size in megabytes.

* `spec_code` - The resource specification code.
  For example: `geminidb.cassandra.8xlarge.4`.

* `az_status` - The status of the flavor in availability zones.
  The valid values are as follows:
  + **normal**: The flavor is available.
  + **unsupported**: The flavor is not supported.
  + **sellout**: The flavor is sold out.
