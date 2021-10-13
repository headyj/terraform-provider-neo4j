resource "neo4j_role" "reader" {
  name = "reader"
}
resource "neo4j_role" "admin" {
  name = "admin"
}
resource "neo4j_grant" "reader" {
  role   = neo4j_role.reader.name
  action = "access"
  graph  = "*"
}
resource "neo4j_grant" "reader_match_all" {
  action   = "match"
  graph    = "*"
  role     = neo4j_role.reader.name
  resource = "all_properties"

}
resource "neo4j_grant" "reader_access" {
  action = "match"
  graph  = "*"
  role   = neo4j_role.reader.name
}
resource "neo4j_grant" "admin_access_all" {
  action = "access"
  graph  = "*"
  role   = "admin"
}
resource "neo4j_grant" "admin_transaction_management_all" {
  action = "transaction_management"
  graph  = "*"
  role   = "admin"
}