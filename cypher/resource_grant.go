package cypher

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

func buildQuery2(privilege string, resource string, name string, entity_type string, entity string, role string, revoke bool) (string, string, error) {
	var resourcePrefix = ""
	var resourceSuffix = " "
	var grantType string
	var entityQuery string
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

func buildQuery(d *schema.ResourceData, revoke bool) (string, string, error) {
	resource := d.Get("resource").(string)
	name := d.Get("name").(string)
	entity := d.Get("entity").(string)
	role := d.Get("role").(string)

	var resourcePrefix = ""
	var resourceSuffix = ""
	var privilege string
	//var resourceAllowed bool
	var grantType string
	var entityQuery string
	switch d.Get("privilege") {
	case "TRAVERSE", "READ", "MATCH", "SET PROPERTY", "MERGE":
		grantType = "GRAPH"
		resourcePrefix = "{"
		resourceSuffix = "}"
		privilege = d.Get("privilege").(string)
		if d.Get("entity_type") == "NODE" || d.Get("entity_type") == "RELATIONSHIP" {
			entityQuery = fmt.Sprintf("%s %s", d.Get("entity_type"), entity)
		} else {
			return "", "", errors.New(fmt.Sprintf("Unexpected entity type: %s", d.Get("entity_type").(string)))
		}
	case "CREATE", "DELETE", "SET LABEL", "REMOVE LABEL", "WRITE", "ALL GRAPH PRIVILEGES":
		grantType = "GRAPH"
		privilege = d.Get("privilege").(string)
		if d.Get("entity_type") == "NODE" || d.Get("entity_type") == "RELATIONSHIP" {
			entityQuery = fmt.Sprintf("%s %s", d.Get("entity_type"), entity)
		} else {
			return "", "", errors.New(fmt.Sprintf("Unexpected entity type: %s", d.Get("entity_type").(string)))
		}
	case "ACCESS", "START", "STOP", "CREATE INDEX", "DROP INDEX", "SHOW INDEX", "CREATE CONSTRAINT", "DROP CONSTRAINT", "SHOW CONSTRAINT", "CREATE NEW NAME", "SHOW TRANSACTION", "TERMINATE TRANSACTION":
		grantType = "DATABASE"
		privilege = d.Get("privilege").(string)
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
	default:
		return "", "", errors.New(fmt.Sprintf("Unexpected privilege: %s", d.Get("privilege").(string)))
	}

	if name != "*" {
		name = fmt.Sprintf("`%s`", d.Get("name").(string))
	}

	toFrom := "TO"
	if revoke {
		toFrom = "FROM"
	}
	return fmt.Sprintf("GRANT %s %s%s%s ON %s %s %s %s `%s`", privilege, resourcePrefix, resource, resourceSuffix, grantType, name, entityQuery, toFrom, role), privilege, nil
}

func resourceGrantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	query, privilege, err := buildQuery2(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), false)
	//query, privilege, err := buildQuery(d, false)
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

	if d.Get("entity_type") != "" {
		d.SetId(fmt.Sprintf("%s:%s:%s:%s:%s:%s", privilege, d.Get("resource"), d.Get("name"), d.Get("entity_type"), d.Get("entity"), d.Get("role")))
	} else if d.Get("resource") != "" {
		d.SetId(fmt.Sprintf("%s:%s:%s:%s", privilege, d.Get("resource"), d.Get("name"), d.Get("role")))
	} else {
		d.SetId(fmt.Sprintf("%s:%s:%s", privilege, d.Get("name"), d.Get("role")))
	}

	return diags
}

func resourceGrantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	query, _, err := buildQuery2(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), false)
	//query, _, err := buildQuery(d, false)
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
	query, _, err := buildQuery2(d.Get("privilege").(string), d.Get("resource").(string), d.Get("name").(string), d.Get("entity_type").(string), d.Get("entity").(string), d.Get("role").(string), true)

	//query, _, err := buildQuery(d, true)
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

	idMembers := strings.Split(d.Id(), ":")

	var role string
	var name string
	var resource string
	var entity_type string
	var entity string
	privilege := idMembers[0]

	if len(idMembers) == 3 {
		name = idMembers[1]
		role = idMembers[2]
	} else if len(idMembers) == 4 {
		resource = idMembers[1]
		name = idMembers[2]
		role = idMembers[3]
	} else if len(idMembers) == 6 {
		resource = idMembers[1]
		name = idMembers[2]
		entity_type = idMembers[3]
		entity = idMembers[4]
		role = idMembers[5]
	} else {
		return nil, fmt.Errorf("Wrong ID format")
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	command, _, err := buildQuery2(privilege, resource, name, entity_type, entity, role, false)
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
