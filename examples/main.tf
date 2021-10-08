terraform {
  required_providers {
    cypher = {
      version = "0.2"
      source  = "hashicorp.com/edu/cypher"
    }
  }
}

provider "cypher" {
  uri      = "neo4j://localhost:7687"
  username = "neo4j"
  password = "password1"
}

module "psl" {
  source = "./cypher"

}

/*output "database" {
  value = module.psl.database
}

output "databases" {
  value = module.psl.databases
}*/