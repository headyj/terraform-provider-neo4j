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
  password = "password1"
}

module "psl" {
  source = "./neo4j"

}

/*output "database" {
  value = module.psl.database
}

output "databases" {
  value = module.psl.databases
}*/