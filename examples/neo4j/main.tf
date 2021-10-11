terraform {
  required_providers {
    neo4j = {
      version = "0.2"
      source  = "headyj/neo4j"
    }
  }
}

provider "neo4j" {
  host     = "neo4j://localhost:7687"
  username = "neo4j"
  password = "password"
}

resource "neo4j_database" "my_database" {
  name = "myDatabase"
}

resource "neo4j_role" "my_role" {
  name = "myRole"
}

resource "neo4j_user" "my_user" {
  name     = "myUser"
  password = "password"

  roles = [
    neo4j_role.my_role.name
  ]
}