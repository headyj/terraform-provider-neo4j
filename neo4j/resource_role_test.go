package neo4j

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestResourceRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceRoleConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccRoleExists("neo4j_role.test"),
					resource.TestCheckResourceAttr("neo4j_role.test", "name", "test"),
				),
			},
		},
	})
}

func testAccRoleExists(rn string) resource.TestCheckFunc {
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

func testResourceRoleConfig_basic() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_role" "test" {
		name = "test"
	}
	`)
}
