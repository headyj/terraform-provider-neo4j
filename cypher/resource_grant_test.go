package cypher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestResourceGrant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceGrantConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccGrantExists("cypher_role.test"),
					resource.TestCheckResourceAttr("cypher_role.test", "name", "testRole"),
				),
			},
		},
	})
}

func TestImportResourceGrant(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "cypher_grant.reader",
				ImportState:   true,
				ImportStateId: "ACCESS:*:reader",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "cypher_grant.reader_match_all_node",
				ImportState:   true,
				ImportStateId: "MATCH:*:*:NODE:*:reader",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "cypher_grant.reader_match_all_relationship",
				ImportState:   true,
				ImportStateId: "MATCH:*:*:RELATIONSHIP:*:reader",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "cypher_grant.admin_access_all",
				ImportState:   true,
				ImportStateId: "ACCESS:*:admin",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "cypher_grant.admin_transaction_management_all",
				ImportState:   true,
				ImportStateId: "TRANSACTION MANAGEMENT:*:*:admin",
			},
		},
	})
}

func testAccGrantExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role id not set")
		}

		c, err := testAccProvider.Meta().(*Neo4jConfiguration).GetDbConn()
		if err != nil {
			return err
		}
		defer c.Close()
		session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()
		result, err := neo4j.Single(session.Run("SHOW ROLE $name PRIVILEGES YIELD role RETURN role LIMIT 1", map[string]interface{}{"username": rs.Primary.ID}))

		fmt.Print(result)

		return nil
	}
}

func testResourceGrantConfig_basic() string {
	return fmt.Sprint(`
	provider "cypher" {
		uri      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	resource "cypher_role" "test" {
		name = "testRole"
	}
	resource "cypher_user" "test" {
		name = "testUser"
		password = "test"
		roles = [
			cypher_role.test.name
		]
	}
	resource "cypher_grant" "test" {
		role = "${cypher_role.test.name}"
		privilege = "READ"
		resource = "*"
		name = "*"
		entity_type = "NODE"
		entity = "*"
	}
	`)
}

func testResourceGrantConfig_import() string {
	return fmt.Sprint(`
	provider "cypher" {
		uri      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	resource "cypher_role" "reader" {
		name = "reader"
	}
	resource "cypher_role" "admin" {
		name = "admin"
	}
	resource "cypher_grant" "reader" {
		role = "${cypher_role.reader.name}"
		privilege = "ACCESS"
		name = "*"
	}
	resource "cypher_grant" "reader_match_all_node" {
		role        = "${cypher_role.reader.name}"
		privilege   = "MATCH"
		resource    = "*"
		name        = "*"
		entity_type = "NODE"
	}
	resource "cypher_grant" "reader_match_all_relationship" {
		role        = "${cypher_role.reader.name}"
		privilege   = "MATCH"
		resource    = "*"
		name        = "*"
		entity_type = "RELATONSHIPO"
	}
	resource "cypher_grant" "admin_access_all" {
		role        = "admin"
		privilege   = "ACCESS"
		name        = "*"
	}
	resource "cypher_grant" "admin_transaction_management_all" {
		role        = "admin"
		privilege   = "ACCESS"
		resource = "*"
		name        = "*"
	}
	`)
}
