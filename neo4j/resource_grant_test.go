package neo4j

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestResourceGrantV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceGrantV2Config_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccGrantV2Exists("neo4j_role.test"),
					resource.TestCheckResourceAttr("neo4j_role.test", "name", "testRole"),
				),
			},
		},
	})
}

func TestImportResourceGrantV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testResourceGrantV2Config_import(),
				ResourceName:  "neo4j_grant.reader_match_all",
				ImportState:   true,
				ImportStateId: "access:*:reader:database",
			},
			{
				Config:        testResourceGrantV2Config_import(),
				ResourceName:  "neo4j_grant.reader_access",
				ImportState:   true,
				ImportStateId: "match:*:reader",
			},
			{
				Config:        testResourceGrantV2Config_import(),
				ResourceName:  "neo4j_grant.admin_access_all",
				ImportState:   true,
				ImportStateId: "access:*:admin",
			},
			{
				Config:        testResourceGrantV2Config_import(),
				ResourceName:  "neo4j_grant.admin_transaction_management_all",
				ImportState:   true,
				ImportStateId: "transaction_management:*:admin",
			},
		},
	})
}

func testAccGrantV2Exists(rn string) resource.TestCheckFunc {
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
		session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close()
		_, err = neo4j.Single(session.Run("SHOW ROLES WHERE role = $role", map[string]interface{}{"role": rs.Primary.ID}))
		if err != nil {
			return err
		}

		return nil
	}
}

func testResourceGrantV2Config_basic() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_role" "test" {
		name = "testRole"
	}
	resource "neo4j_user" "test" {
		name = "testUser"
		password = "test"
		roles = [
			neo4j_role.test.name
		]
	}
	resource "neo4j_grant" "remove_label" {
		role = "${neo4j_role.test.name}"
		action = "remove_label"
		resource = "all_labels"
		graph = "*"
	}

	resource "neo4j_grant" "set_label" {
		role = "${neo4j_role.test.name}"
		action = "set_label"
		resource = "all_labels"
		graph = "*"
	}
	resource "neo4j_grant" "set_label_resource" {
		role = "${neo4j_role.test.name}"
		action = "set_label"
		resource = "label(toto)"
		graph = "*"
	}

	resource "neo4j_grant" "match" {
		role = "${neo4j_role.test.name}"
		action = "match"
		resource = "all_properties"
		graph = "neo4j"
		segment = "NODE(toto)"
	}
	resource "neo4j_grant" "match_property" {
		role = "${neo4j_role.test.name}"
		action = "match"
		resource = "property(toto)"
		graph = "neo4j"
		segment = "NODE(toto)"
	}

	resource "neo4j_grant" "merge" {
		role = "${neo4j_role.test.name}"
		action = "merge"
		resource = "all_properties"
		graph = "system"
		segment = "RELATIONSHIP(toto)"
	}

	resource "neo4j_grant" "set_property" {
		role = "${neo4j_role.test.name}"
		action = "set_property"
		resource = "all_properties"
		graph = "system"
		segment = "RELATIONSHIP(*)"
	}

	resource "neo4j_grant" "delete_element" {
		role = "${neo4j_role.test.name}"
		action = "delete_element"
		graph = "system"
	}












	resource "neo4j_grant" "test" {
		role = "${neo4j_role.test.name}"
		action = "create_element"
		resource = "graph"
		graph = "*"
		segment = "NODE(toto)"
	}
	resource "neo4j_grant" "test2" {
		role = "${neo4j_role.test.name}"
		action = "write"
		resource = "graph"
		graph = "*"
	}
	
	`)
}

func testResourceGrantV2Config_import() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host     = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_role" "reader" {
		name = "reader"
	}
	resource "neo4j_role" "admin" {
		name = "admin"
	}
	resource "neo4j_grant" "reader" {
		role = "${neo4j_role.reader.name}"
		action = "access"
		graph = "*"
	}
	resource "neo4j_grant" "reader_match_all" {
		action   = "match"
		graph    = "*"
		role     = neo4j_role.reader.name
		resource = "database"
	
	}

	resource "neo4j_grant" "reader_access" {
		action = "match"
		graph  = "*"
		role   = neo4j_role.reader.name
	}
	resource "neo4j_grant" "admin_access_all" {
		action = "access"
		graph  = "*"
		role   = "${neo4j_role.admin.name}"
	}
	resource "neo4j_grant" "admin_transaction_management_all" {
		action = "transaction_management"
		graph  = "*"
		role   = "${neo4j_role.admin.name}"
	}
	`)
}
