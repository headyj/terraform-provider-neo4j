package neo4j

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func resourceGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrantCreate,
		ReadContext:   resourceGrantRead,
		DeleteContext: resourceGrantDelete,
		Schema: map[string]*schema.Schema{
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privilege": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"entity_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"entity": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			//StateContext: schema.ImportStatePassthroughContext,
			State: resourceGrantImport,
		},
	}
}

func buildQuery(privilege string, resource string, name string, entity_type string, entity string, role string, revoke bool) (string, string, error) {
	var resourcePrefix = ""
	var resourceSuffix = " "
	var grantType string
	var entityQuery string
	privilege = strings.ReplaceAll(privilege, "-", " ")
	switch privilege {
	case "TRAVERSE", "READ", "MATCH", "SET PROPERTY", "MERGE":
		grantType = "GRAPH"
		resourcePrefix = "{"
		resourceSuffix = "} "
		if entity_type == "NODE" || entity_type == "RELATIONSHIP" {
			entityQuery = fmt.Sprintf("%s %s ", entity_type, entity)
		} else {
			return "", "", errors.New(fmt.Sprintf("Unexpected entity type: %s", entity_type))
		}
	case "CREATE", "DELETE", "SET LABEL", "REMOVE LABEL", "WRITE", "ALL GRAPH PRIVILEGES":
		grantType = "GRAPH"
		if entity_type == "NODE" || entity_type == "RELATIONSHIP" {
			entityQuery = fmt.Sprintf("%s %s ", entity_type, entity)
		} else {
			return "", "", errors.New(fmt.Sprintf("Unexpected entity type: %s", entity_type))
		}
	case "ACCESS", "START", "STOP", "CREATE INDEX", "DROP INDEX", "SHOW INDEX", "CREATE CONSTRAINT", "DROP CONSTRAINT", "SHOW CONSTRAINT", "CREATE NEW NAME", "SHOW TRANSACTION", "TERMINATE TRANSACTION":
		grantType = "DATABASE"
	case "INDEX", "INDEX MANAGEMENT":
		grantType = "DATABASE"
		privilege = "INDEX MANAGEMENT"
	case "CONSTRAINT", "CONSTRAINT MANAGEMENT":
		grantType = "DATABASE"
		privilege = "CONSTRAINT MANAGEMENT"
	case "CREATE NEW LABEL", "CREATE NEW NODE LABEL":
		grantType = "DATABASE"
		privilege = "CREATE NEW NODE LABEL"
	case "CREATE NEW TYPE", "CREATE NEW RELATIONSHIP TYPE":
		grantType = "DATABASE"
		privilege = "CREATE NEW RELATIONSHIP TYPE"
	case "NAME", "NAME MANAGEMENT":
		grantType = "DATABASE"
		privilege = "NAME MANAGEMENT"
	case "ALL", "ALL DATABASE", "ALL DATABASE PRIVILEGES":
		grantType = "DATABASE"
		privilege = "ALL DATABASE PRIVILEGES"
	case "TRANSACTION", "TRANSACTION MANAGEMENT":
		grantType = "DATABASE"
		privilege = "TRANSACTION MANAGEMENT"
		resourcePrefix = "("
		resourceSuffix = ") "
	default:
		return "", "", errors.New(fmt.Sprintf("Unexpected privilege: %s", privilege))
	}

	if name != "*" {
		name = fmt.Sprintf("`%s`", name)
	}

	if resource == "" {
		resourceSuffix = ""
	}

	toFrom := "TO"
	if revoke {
		toFrom = "FROM"
	}
	return fmt.Sprintf("GRANT %s %s%s%sON %s %s %s%s `%s`", privilege, resourcePrefix, resource, resourceSuffix, grantType, name, entityQuery, toFrom, role), privilege, nil

	//"GRANT ACCESS ON DATABASE *  TO `reader`"

}
func resourceGrantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	query, privilege, err := buildQuery(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), false)
	if err != nil {
		return diag.FromErr(err)
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	fmt.Print(query)
	_, err = session.Run(query, map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	/**************** ID REWORK TODO ******************
	  privilege:name:role
	  privilege:name:role|resource
	  privilege:name:role||entity_type:entity
	  privilege:name:role|resource|entity:entity_type
	*/

	id := fmt.Sprintf("%s:%s:%s", strings.ReplaceAll(privilege, " ", "-"), d.Get("name"), d.Get("role"))

	if d.Get("entity_type") != "" && d.Get("entity") != "" {
		id += fmt.Sprintf("_%s_%s:%s", d.Get("resource"), d.Get("entity_type"), d.Get("entity"))
	} else {
		if d.Get("resource") != "" {
			id += fmt.Sprintf("_%s", d.Get("resources"))
		}
	}
	d.SetId(id)

	//	d.SetId(fmt.Sprintf("%s:%s:%s:%s:%s:%s", privilege, d.Get("resource"), d.Get("name"), d.Get("entity_type"), d.Get("entity"), d.Get("role")))
	//} else if d.Get("resource") != "" {
	//	d.SetId(fmt.Sprintf("%s:%s:%s:%s", privilege, d.Get("resource"), d.Get("name"), d.Get("role")))
	//} else {
	//	d.SetId(fmt.Sprintf("%s:%s:%s", privilege, d.Get("name"), d.Get("role")))
	//}

	return diags
}

func resourceGrantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	query, _, err := buildQuery(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), false)
	if err != nil {
		return diag.FromErr(err)
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	reqCommand := fmt.Sprintf("SHOW ROLE %s PRIVILEGES AS COMMAND WHERE command = '%s'", d.Get("role"), query)
	result, err := neo4j.Single(session.Run(reqCommand, map[string]interface{}{}))
	if result == nil {
		fmt.Printf("GRANT not found, removing from state")
		d.SetId("")
	}

	return diags
}

func resourceGrantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	return resourceGrantRead(ctx, d, m)
}

func resourceGrantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	query, _, err := buildQuery(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), true)

	if err != nil {
		return diag.FromErr(err)
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	reqCommand := fmt.Sprintf("REVOKE %s", query)
	_, err = session.Run(reqCommand, map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)

	}

	d.SetId("")
	return diags
}

func resourceGrantImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var privilege string
	var role string
	var name string
	resource := ""
	entity_type := ""
	entity := ""
	idGroups := strings.Split(d.Id(), "_")
	if len(idGroups) > 0 {
		mandatoryMembers := strings.Split(idGroups[0], ":")

		if len(mandatoryMembers) != 3 {
			return nil, fmt.Errorf("Wrong ID format")
		}
		privilege = strings.ReplaceAll(mandatoryMembers[0], "-", " ")
		name = mandatoryMembers[1]
		role = mandatoryMembers[2]
		if len(idGroups) > 1 {
			resource = idGroups[1]
		}
		if len(idGroups) > 2 {
			entities := strings.Split(idGroups[2], ":")
			entity_type = entities[0]
			entity = entities[1]
		}
	} else {
		return nil, fmt.Errorf("Wrong ID format")
	}
	idMembers := strings.Split(d.Id(), "_")

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	command, _, err := buildQuery(privilege, resource, name, entity_type, entity, role, false)
	result, err := neo4j.Single(session.Run("SHOW ROLE $rolename PRIVILEGES AS COMMANDS where command = $command", map[string]interface{}{"rolename": role, "command": command}))
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("Privilege not found")
	}
	d.Set("privilege", privilege)
	d.Set("resource", resource)
	d.Set("name", name)
	if len(idMembers) == 4 {
		d.Set("role", role)
	} else {
		d.Set("entity_type", entity_type)
		d.Set("entity", entity)
		d.Set("role", role)
	}
	return []*schema.ResourceData{d}, nil

}
