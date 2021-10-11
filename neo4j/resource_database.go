package neo4j

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run("CREATE DATABASE $database", map[string]interface{}{"database": d.Get("name")})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("name").(string))

	return diags
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	result, err := neo4j.Single(session.Run("SHOW DATABASE $database", map[string]interface{}{"database": d.Id()}))
	if err != nil {
		return diag.FromErr(err)
	}
	name, _ := result.Get("name")

	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", name)

	return diags
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDatabaseRead(ctx, d, m)
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run("DROP DATABASE $database", map[string]interface{}{"database": d.Get("name")})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
