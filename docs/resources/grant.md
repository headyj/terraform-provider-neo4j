# neo4j_database

the `neo4j_grant` resource creates and manage grants on a Neo4j server.

## Example Usage

```hcl
resource "neo4j_database" "my_database" {
  name ="myDatabase"
}
```

## Argument Reference

The following arguments are supported:

* `privilege` - (Required) The privilege name of the grant. It can be both related to databases or graphs. Please read [the offical documentation](https://neo4j.com/docs/cypher-manual/current/access-control/manage-roles/) for more information
* `name` - (Required) The name of the database or graph associated with the grant. it can be "*" or the specific database or graph name
* `role` - (Required) The role associated with the grant.
* `resource` - (Optional) The role associated with the grant.
* `entity_type` - (Optional) In the case of graph related grant, you can specify the entity type of the grant. It can be "NODE" or "RELATIONSHIP"
* `entity` - (Optional) In the case of graph related grant, you can specify the entity of the grant. It can be "*" or a specific entity name

## Attribute Reference

The following arguments are exported:

* `privilege` - The privilege name of the grant.
* `name` - The name of the database or graph associated with the grant.
* `role` - The role associated with the grant.
* `resource` - The role associated with the grant.
* `entity_type` - the entity_type of the grant
* `entity` - tne entity of the grant

## Import

neo4j_grant resource can be importe using the following pattern:
firstGroup_secondGroup_thirdGroup

* The first group is mandatory and represents the privilege, name and role associated with the grant:
privilege:name:role

If the privilege is composed of multiple words (e.g. "TRANSACTION MANAGEMENT"), just put a dash between them

* The second group is optional and represents the associated resource
privilege:name:role_resource

* The third group is optional and represents the entity type and entity name. If there is no resource, 2 separators have to be specified
privilege:name:role__entity_type:entity
privilege:name:role_resource_entity_type:entity

Here are some examples how to import default reader and admin roles on a Neo4j cluster

```bash
terraform import neo4j_grant.reader_match_all ACCESS:*:reader
terraform import neo4j_grant.reader_match_all_node MATCH:*:reader_*_NODE:*
terraform import neo4j_grant.reader_match_all_relationship MATCH:*:reader_*_RELATIONSHIP:*
terraform import neo4j_grant.admin_access_all ACCESS:*:admin
terraform import neo4j_grant.admin_transaction_management_all TRANSACTION-MANAGEMENT:*:admin_*
```