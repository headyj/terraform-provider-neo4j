# neo4j_database

the `neo4j_database` resource creates and manage databases on a Neo4j server.

## Example Usage

```hcl
resource "neo4j_database" "my_database" {
  name ="myDatabase"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database.

## Attribute Reference

The following arguments are exported:

* `name` - The name of the database.

## Import

neo4j_database resource can be importe using the resource name, e.g.

```bash
terraform import neo4j_user.my_database myDatabase
```