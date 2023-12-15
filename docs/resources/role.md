# neo4j_role

the `neo4j_role` resource creates and manage roles on a Neo4j server.

## Example Usage

```hcl
resource "neo4j_role" "my_role" {
  name ="myRole"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.

## Attribute Reference

The following arguments are exported:

* `name` - The name of the role.

## Import

neo4j_database resource can be importe using the resource name, e.g.

```bash
terraform import neo4j_role.my_role myRole
```
