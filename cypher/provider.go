package cypher

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jConfiguration struct {
	Uri      string
	Username string
	Password string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_USERNAME", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cypher_database": resourceDatabase(),
			"cypher_user":     resourceUser(),
			"cypher_role":     resourceRole(),
			"cypher_grant":    resourceGrant(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cypher_databases": dataSourceDatabases(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func (c *Neo4jConfiguration) GetDbConn() (neo4j.Driver, error) {

	client, err := neo4j.NewDriver(c.Uri, neo4j.BasicAuth(c.Username, c.Password, ""))

	return client, err
}

func (c *Neo4jConfiguration) GetDbUri() string {

	return c.Uri
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	uri := d.Get("uri").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//	if (username != "") && (password != "") {
	//
	//		if err != nil {
	//			diags = append(diags, diag.Diagnostic{
	//				Severity: diag.Error,
	//				Summary:  "Unable to create Cyher client",
	//				Detail:   "Unable to auth user for authenticated Cyher client",
	//			})
	//			return nil, diags
	//		}
	//
	//		return conf, diags
	//	}
	//
	//	if err != nil {
	//		diags = append(diags, diag.Diagnostic{
	//			Severity: diag.Error,
	//			Summary:  "Unable to create Cyher client",
	//			Detail:   "Unable to auth user for authenticated HashiCups client",
	//		})
	//		return nil, diags
	//	}

	conf := &Neo4jConfiguration{
		Uri:      uri,
		Username: username,
		Password: password,
	}

	return conf, diags
}
