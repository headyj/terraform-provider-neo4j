# Neo4j Provider

Neo4j is a graph database. This provider gives the ability to create, update, delete and import resources from a Neo4j Entreprise server

> :warning: **Neo4j Enterprise is a licensed product. Please read the [official license documentation](https://neo4j.com/licensing)**

## Example Usage

```hcl
provider "neo4j" {
  host      = "neo4j://localhost:7687"
  username = "neo4j"
  password = "password1"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required) The address of the Neo4j server to use, using the format "neo4j://hostname:port". Can also be sourced from the `NEO4J_HOST` environment variable.
* `username` - (Required) The username to connect to the server. Can also be sourced from the `NEO4J_USERNAME` environment variable.
* `password` - (Required) The password associated with the username. Can also be sourced from the `NEO4J_PASSWORD` environment variable.