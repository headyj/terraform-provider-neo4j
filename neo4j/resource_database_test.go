package neo4j

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
					resource.TestCheckResourceAttr("data.neo4j_databases.all", "name", "mysql"),
					resource.TestCheckResourceAttr("data.neo4j_databases.all", "pattern", "%"),
				),
			},
		},
	})
}

func testResourceDatabaseConfig_basic() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_database" "test" {
		name = "test"
	}
	`)
}
