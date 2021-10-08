package cypher

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
		"cypher": testAccProvider,
	}
}

//func TestProvider(t *testing.T) {
//	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
//		t.Fatalf("err: %s", err)
//	}
//}
//
//func TestProvider_impl(t *testing.T) {
//	var _ terraform.ResourceProvider = Provider()
//}

func testAccPreCheck(t *testing.T) {
	for _, name := range []string{"NEO4J_URI", "NEO4J_USERNAME"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("NEO4J_URI, NEO4J_USERNAME and optionally NEO4J_PASSWORD must be set for acceptance tests")
		}
	}

}

func testAccProviderConfig() string {
	return fmt.Sprint(`
	provider "cypher" {
		uri      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password1"
	}
	`)
}
