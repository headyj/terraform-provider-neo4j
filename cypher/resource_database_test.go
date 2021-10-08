package cypher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceDatabaseConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cypher_databases.all", "name", "mysql"),
					resource.TestCheckResourceAttr("data.cypher_databases.all", "pattern", "%"),
				),
			},
		},
	})
}

func testResourceDatabaseConfig_basic() string {
	return fmt.Sprint(`
	provider "cypher" {
		uri      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	resource "cypher_database" "test" {
		name = "test"
	}
	`)
}
