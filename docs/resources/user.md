---
layout: "neo4j"
page_title: "Neo4j: neo4j_user"
description: |-
  Neo4j user
---

# neo4j_user

the `neo4j_user` resource creates and manage users on a Neo4j server.

## Example Usage

```hcl
resource "neo4j_role" "my_role" {
		name ="myRole"
	}
resource "neo4j_user" "my_user" {
  name = "myUser"
  password = "password"
  roles = [
    neo4j_role.my_role.name
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user.
* `password` - (Required) The password of the user.
* `roles` - (Optional) The list of user roles associated.

## Attribute Reference

The following arguments are exported:

* `name` - The name of the user.
* `roles` - The list of user roles associated.

## Import

neo4j_user resource can be importe using the name, e.g.

```bash
terraform import neo4j_user.my_user myUser
```