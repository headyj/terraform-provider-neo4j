package neo4j

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jConfiguration struct {
	Host     string
	Username string
	Password string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEO4J_USERNAME", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEO4J_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NEO4J_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"neo4j_database": resourceDatabase(),
			"neo4j_user":     resourceUser(),
			"neo4j_role":     resourceRole(),
			"neo4j_grant":    resourceGrant(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"neo4j_databases": dataSourceDatabases(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func (c *Neo4jConfiguration) GetDbConn() (neo4j.Driver, error) {

	client, err := neo4j.NewDriver(c.Host, neo4j.BasicAuth(c.Username, c.Password, ""))

	return client, err
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	conf := &Neo4jConfiguration{
		Host:     host,
		Username: username,
		Password: password,
	}

	return conf, diags
}
