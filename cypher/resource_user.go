package cypher

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"change_password_required": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		//Importer: &schema.ResourceImporter{
		//	State: resourceUserImport,
		//},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	command := "CREATE USER $username SET PASSWORD $password"

	if !d.Get("change_password_required").(bool) {
		command += " CHANGE NOT REQUIRED"
	}

	_, err = session.Run(command, map[string]interface{}{"username": d.Get("name"), "password": d.Get("password")})
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", d.Get("name"))
	d.SetId(d.Get("name").(string))

	roles := d.Get("roles").(*schema.Set).List()
	var roleList []string
	if len(roles) > 0 {
		for _, role := range roles {
			roleList = append(roleList, role.(string))
		}
		_, err = session.Run("GRANT ROLE $roles TO $username", map[string]interface{}{"roles": strings.Join(roleList, ","), "username": d.Get("name")})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	result, err := neo4j.Single(session.Run("SHOW USERS YIELD user WHERE user = $username", map[string]interface{}{"username": d.Id()}))
	if err != nil {
		return diag.FromErr(err)
	}
	name, _ := result.Get("user")

	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	/*if d.HasChange("name") {
		current_username, new_username := d.GetChange("name")
		_, err = session.Run("RENAME USER $current_username TO $new_username", map[string]interface{}{"current_username": current_username, "new_username": new_username})
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("name", new_username)

	}*/

	if d.HasChange("password") {
		_, err = session.Run("ALTER USER $username SET PASSWORD $password", map[string]interface{}{"username": d.Get("name"), "password": d.Get("password")})
		if err != nil {
			if !strings.Contains(fmt.Sprint(err), "Old password and new password cannot be the same") {
				return diag.FromErr(err)
			}
		}
		d.GetChange("password")
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run("DROP USER $username", map[string]interface{}{"username": d.Get("name")})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
