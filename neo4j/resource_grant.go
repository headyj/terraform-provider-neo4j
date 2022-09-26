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
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"graph": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"segment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceGrantImport,
		},
	}
}

func buildQuery(rawAction string, rawResource string, rawGraph string, rawSegment string, role string, revoke bool) (string, error) {

	actions := map[string]string{
		"traverse":               "TRAVERSE",
		"read":                   "READ",
		"match":                  "MATCH",
		"set_property":           "SET PROPERTY",
		"merge":                  "MERGE",
		"create_element":         "CREATE",
		"delete_element":         "DELETE",
		"set_label":              "SET LABEL",
		"remove_label":           "REMOVE LABEL",
		"write":                  "WRITE",
		"graph_actions":          "ALL GRAPH PRIVILEGES",
		"access":                 "ACCESS",
		"start_database":         "START",
		"stop_database":          "STOP",
		"create_index":           "CREATE INDEX",
		"drop_index":             "DROP INDEX",
		"show_index":             "SHOW INDEX",
		"create_constraint":      "CREATE CONSTRAINT",
		"drop_constraint":        "DROP CONSTRAINT",
		"show_constraint":        "SHOW CONSTRAINT",
		"create_propertykey":     "CREATE NEW NAME",
		"show_transaction":       "SHOW TRANSACTION",
		"terminate_transaction":  "TERMINATE TRANSACTION",
		"index":                  "INDEX MANAGEMENT",
		"constraint":             "CONSTRAINT MANAGEMENT",
		"create_label":           "CREATE NEW NODE LABEL",
		"create_reltype":         "CREATE NEW RELATIONSHIP TYPE",
		"name_management":        "NAME MANAGEMENT",
		"database_actions":       "ALL DATABASE PRIVILEGES",
		"transaction_management": "TRANSACTION MANAGEMENT",
	}

	var resourcePrefix = ""
	var resourceSuffix = ""
	var grantType string
	switch rawAction {
	case "traverse", "read", "match", "set_property", "merge":
		grantType = "GRAPH"
		resourcePrefix = "{"
		resourceSuffix = "}"
	case "create_element", "delete_element", "set_label", "remove_label", "write", "graph_actions":
		grantType = "GRAPH"
	case "access", "start_database", "stop_database", "create_index", "drop_index",
		"show_index", "create_constraint", "drop_constraint", "show_constraint",
		"create_propertykey", "show_transaction", "terminate_transaction", "index",
		"constraint", "create_label", "create_reltype",
		"name_management", "database_actions":
		grantType = "DATABASE"
	case "transaction_management":
		grantType = "DATABASE"
		resourcePrefix = "("
		resourceSuffix = ")"
	default:
		return "", errors.New(fmt.Sprintf("Unexpected action: %s", rawAction))
	}

	action := actions[rawAction]
	resource := ""
	segment := ""
	if strings.Contains(rawResource, "property(") || strings.Contains(rawResource, "label(") {
		resource = fmt.Sprintf("%s%s%s ", resourcePrefix, between(rawResource, "(", ")"), resourceSuffix)
	} else if rawResource == "all_properties" || rawResource == "all_labels" {
		resource = fmt.Sprintf("%s%s%s ", resourcePrefix, "*", resourceSuffix)
	}
	if rawSegment != "database" && rawSegment != "" {
		if strings.Contains(rawSegment, "NODE") {
			segment = fmt.Sprintf("NODE %s ", strings.TrimLeft(strings.TrimRight(rawSegment, ")"), "NODE("))
		} else if strings.Contains(rawSegment, "RELATIONSHIP") {
			segment = fmt.Sprintf("RELATIONSHIP %s ", strings.TrimLeft(strings.TrimRight(rawSegment, ")"), "RELATIONSHIP("))
		} else {
			return "", errors.New(fmt.Sprintf("Unexpected segment: %s", rawSegment))
		}
	}
	graph := "*"
	if rawGraph != "*" {
		graph = fmt.Sprintf("`%s`", rawGraph)
	}

	toFrom := "TO"
	if revoke {
		toFrom = "FROM"
	}

	return fmt.Sprintf("GRANT %s %sON %s %s %s%s `%s`", action, resource, grantType, graph, segment, toFrom, role), nil

}

func resourceGrantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	query, err := buildQuery(d.Get("action").(string), d.Get("resource").(string), d.Get("graph").(string), d.Get("segment").(string), d.Get("role").(string), false)
	if err != nil {
		return diag.FromErr(err)
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "system"})
	defer session.Close()

	fmt.Println(query)
	reqCommand := fmt.Sprintf("SHOW ROLE %s PRIVILEGES WHERE action = '%s' AND graph = '%s'", d.Get("role"), d.Get("action"), d.Get("graph"))
	if d.Get("resource") != "" {
		reqCommand += fmt.Sprintf(" AND resource = '%s'", d.Get("resource"))
	}
	if d.Get("segment") != "" {
		reqCommand += fmt.Sprintf(" AND segment = '%s'", d.Get("segment"))
	}
	result, err := neo4j.Collect(session.Run(reqCommand, map[string]interface{}{}))
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) != 0 {
		return diag.Errorf("Privilege already exists")
	}
	_, err = session.Run(query, map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%s:%s:%s", d.Get("action"), d.Get("graph"), d.Get("role"))

	if d.Get("segment") != "" {
		id += fmt.Sprintf(":%s:%s", d.Get("resource"), d.Get("segment"))
	} else {
		if d.Get("resource") != "" {
			id += fmt.Sprintf(":%s", d.Get("resource"))
		}
	}
	d.SetId(id)

	return diags
}

func resourceGrantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: "system"})
	defer session.Close()
	reqCommand := fmt.Sprintf("SHOW ROLE %s PRIVILEGES WHERE action = '%s' AND graph = '%s'", d.Get("role"), d.Get("action"), d.Get("graph"))
	if d.Get("resource") != "" {
		reqCommand += fmt.Sprintf(" AND resource = '%s'", d.Get("resource"))
	}
	if d.Get("segment") != "" {
		reqCommand += fmt.Sprintf(" AND segment = '%s'", d.Get("segment"))
	}

	result, err := neo4j.Collect(session.Run(reqCommand, map[string]interface{}{}))
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) == 0 {
		fmt.Printf("GRANT not found, removing from state")
		d.SetId("")
	}

	return diags
}

func resourceGrantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	query, err := buildQuery(d.Get("action").(string), d.Get("resource").(string), d.Get("graph").(string), d.Get("segment").(string), d.Get("role").(string), true)
	if err != nil {
		return diag.FromErr(err)
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "system"})
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
	var action string
	var role string
	var graph string
	resource := ""
	segment := ""
	idMembers := strings.Split(d.Id(), ":")
	if len(idMembers) < 3 || len(idMembers) > 5 {
		return nil, fmt.Errorf("Wrong ID format")
	}
	action = idMembers[0]
	graph = idMembers[1]
	role = idMembers[2]
	if len(idMembers) == 4 {
		resource = idMembers[3]
	}
	if len(idMembers) == 5 {
		if idMembers[3] != "" {
			resource = idMembers[3]
		}
		segment = idMembers[4]
	}

	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	command := fmt.Sprintf("SHOW ROLE %s PRIVILEGES WHERE action = '%s' AND graph = '%s'", role, action, graph)
	if resource != "" {
		command += fmt.Sprintf(" AND resource = '%s'", resource)
	}
	if segment != "" {
		command += fmt.Sprintf(" AND segment = '%s'", segment)
	}

	result, err := neo4j.Collect(session.Run(command, map[string]interface{}{}))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("Privilege not found")
	}
	d.Set("action", action)
	d.Set("resource", resource)
	d.Set("graph", graph)
	d.Set("role", role)
	if segment != "" {
		d.Set("segment", segment)
	}
	return []*schema.ResourceData{d}, nil

}
