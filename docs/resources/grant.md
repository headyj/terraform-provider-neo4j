# neo4j_grant

the `neo4j_grant` resource creates and manage grants on a Neo4j server.

## Example Usage

```hcl
resource "neo4j_role" "my_role" {
  name = "myRole"
}
resource "neo4j_grant" "my_grant" {
  action   = "match"
  graph    = "neo4j"
  role     = neo4j_role.my_role.name
  resource = "all_properties"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required) The privilege name of the grant. It can be both related to databases or graphs. See [available actions](#available-actions) below for valid values. Please read [the offical documentation](https://neo4j.com/docs/cypher-manual/current/access-control/manage-roles/) for more information.
* `graph` - (Required) The name of the database or graph associated with the grant. it can be "*" or the specific database or graph name.
* `role` - (Required) The role associated with the grant.
* `resource` - (Optional) The resource associated with the grant. Valid values are (depending on the type of action) `all_labels`, `all_properties`,`graph`,`database`,`label(<value>)`,`property(<value>)`.
* `segment` - (Optional) In the case of graph related grant, you can specify the segment of the grant. Valid values are `NODE(*)`, `RELATIONSHIP(*)`, `NODE(<value>)`, `RELATIONSHIP(<value>)`.

### Available actions
* traverse (TRAVERSE)
* read (READ)
* match (MATCH)
* set_property (SET PROPERTY)
* merge (MERGE)
* create_element (CREATE)
* delete_element (DELETE)
* set_label (SET LABEL)
* remove_label (REMOVE LABEL)
* write (WRITE)
* graph_actions (ALL GRAPH PRIVILEGES)
* access (ACCESS)
* start_database (START)
* stop_database (STOP)
* create_index (CREATE INDEX)
* drop_index (DROP INDEX)
* show_index (SHOW INDEX)
* create_constraint (CREATE CONSTRAINT)
* drop_constraint (DROP CONSTRAINT)
* show_constraint (SHOW CONSTRAINT)
* create_propertykey (CREATE NEW NAME)
* show_transaction (SHOW TRANSACTION)
* terminate_transaction (TERMINATE TRANSACTION)
* index (INDEX MANAGEMENT)
* constraint (CONSTRAINT MANAGEMENT)
* create_label (CREATE NEW NODE LABEL)
* create_reltype (CREATE NEW RELATIONSHIP TYPE)
* name_management (NAME MANAGEMENT)
* database_actions (ALL DATABASE PRIVILEGES)
* transaction_management (TRANSACTION MANAGEMENT)

## Attribute Reference

The following arguments are exported:

* `action` - The privilege name of the grant.
* `graph` - The name of the database or graph associated with the grant.
* `role` - The role associated with the grant.
* `resource` - The resource associated with the grant.
* `segment` - the segment of the grant

## Import

neo4j_grant resource can be importe using the following pattern:
action:graph:role:resource:segment

action, name and role are mandatory. resource and segments are optional

Not that if there is no resource but a segment, you will have to put 2 separators.

Here are some examples how to import default reader and admin roles on a Neo4j cluster

```bash
terraform import neo4j_grant.reader_match_all access:*:reader:database
terraform import neo4j_grant.admin_access_all access:database:admin
terraform import neo4j_grant.admin_transaction_management_all transaction_management:database:admin
```