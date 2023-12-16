package neo4j

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"neo4j": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"NEO4J_HOST", "NEO4J_USERNAME"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("NEO4J_HOST, NEO4J_USERNAME and optionally NEO4J_PASSWORD must be set for acceptance tests")
		}
	}

}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
	provider "neo4j" {
		host      = "%s"
		username = "%s"
		password = "%s"
	}
	`,
		os.Getenv("NEO4J_HOST"),
		os.Getenv("NEO4J_USERNAME"),
		os.Getenv("NEO4J_PASSWORD"),
	)
}
