package neo4j

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
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

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run("CREATE ROLE $name", map[string]interface{}{"name": d.Get("name")})
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", d.Get("name"))
	d.SetId(d.Get("name").(string))

	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	result, err := neo4j.Single(session.Run("SHOW ROLES WHERE role = $role", map[string]interface{}{"role": d.Id()}))
	name, _ := result.Get("role")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run("DROP ROLE $role", map[string]interface{}{"role": d.Get("name")})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
