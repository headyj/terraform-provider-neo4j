terraform {
  required_providers {
    cypher = {
      version = "0.2"
      source  = "hashicorp.com/edu/cypher"
    }
  }
}
//data "cypher_databases" "all" {}

/*output "databases" {
  value = data.cypher_databases.all.databases
}

output "database" {
  value = resource.cypher_database.jonas
}*/

resource "cypher_database" "neo4j" {
  name = "neo4j"
}

resource "cypher_user" "neo4j" {
  name     = "neo4j"
  password = "password1"
}

resource "cypher_user" "jonas" {
  name     = "jonas"
  password = "password1"

  roles = [
    cypher_role.jonas.name
  ]
}

resource "cypher_role" "jonas" {
  name = "jonas"

}

// TO IMPORT
resource "cypher_grant" "reader" {
  role      = "reader"
  privilege = "ACCESS" // ForceNew
  name      = "*"      // can be * or a graph name
}
resource "cypher_grant" "reader_match_all" {
  role        = "reader"
  privilege   = "MATCH" // ForceNew
  resource    = "*"
  name        = "*"
  entity_type = "NODE" // can be * or a graph name
  entity      = "*"
}
resource "cypher_grant" "admin_transaction_management_all" {
  role      = "admin"
  privilege = "TRANSACTION MANAGEMENT"
  resource  = "*"
  name      = "*"
}
/*
resource "cypher_grant" "jonas_grant" {
  role        = cypher_role.jonas.name
  privilege   = "READ"  // ForceNew
  resource    = "toto"  // ForceNew can be * or a resource name
  name        = "neo4j" // can be * or a graph name
  entity      = "*"     // can be NODE, RELATIONSHIP
  entity_type = "NODE"  // can be NODE, RELATIONSHIP
}
*/