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
					resource.TestCheckResourceAttr("neo4j_database.test", "name", "mydatabase2"),
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
		name = "mydatabase2"
	}
	`)
}
