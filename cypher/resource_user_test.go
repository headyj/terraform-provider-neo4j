package cypher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestResourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceUserConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("cypher_user.test"),
					resource.TestCheckResourceAttr("cypher_user.test", "name", "testUser"),
					resource.TestCheckResourceAttr("cypher_user.test", "password", "test"),
				),
			},
		},
	})
}

func testAccUserExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("user id not set")
		}

		c, err := testAccProvider.Meta().(*Neo4jConfiguration).GetDbConn()
		if err != nil {
			return err
		}
		defer c.Close()
		session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()
		result, err := neo4j.Single(session.Run("SHOW USER $username PRIVILEGES YIELD user RETURN user LIMIT 1", map[string]interface{}{"username": rs.Primary.ID}))

		fmt.Print(result)

		return nil
	}
}

func testResourceUserConfig_basic() string {
	return fmt.Sprint(`
	provider "cypher" {
		uri      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	resource "cypher_user" "test" {
		name = "testUser"
		password = "test"
		roles = [
			cypher_role.test_role.name
		]
	}
	resource "cypher_role" "test_role" {
		name ="testRole"
	}
	`)
}
