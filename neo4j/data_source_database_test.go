package neo4j

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDatabases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDatabasesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.neo4j_databases.all", "name", "mysql"),
					resource.TestCheckResourceAttr("data.neo4j_databases.all", "pattern", "%"),
					testDatabasesCount("data.neo4j_databases.all", "tables.#", func(rn string, table_count int) error {
						if table_count < 1 {
							return fmt.Errorf("%s: tables not found", rn)
						}

						return nil
					}),
				),
			},
		},
	})
}

func testDatabasesCount(rn string, key string, check func(string, int) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]

		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		value, ok := rs.Primary.Attributes[key]

		if !ok {
			return fmt.Errorf("%s: attribute '%s' not found", rn, key)
		}

		table_count, err := strconv.Atoi(value)

		if err != nil {
			return err
		}

		return check(rn, table_count)
		//return fmt.Errorf("TOTO %s %s", rn, strconv.Itoa(table_count))
	}
}

func testDatabasesConfig_basic() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	data "neo4j_databases" "all" {}
	`)
}
