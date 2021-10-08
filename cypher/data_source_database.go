package cypher

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func dataSourceDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabasesRead,
		Schema: map[string]*schema.Schema{
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDatabasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c, err := m.(*Neo4jConfiguration).GetDbConn()
	if err != nil {
		return diag.FromErr(err)
	}
	defer c.Close()

	session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	//toto2, err := session.Run("CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)", map[string]interface{}{"message": "hello, world"})
	toto2, err := session.Run("SHOW DATABASES", nil)
	var databases []string
	//err = json.NewDecoder(toto2.).Decode(&databases)
	//session.Run
	for toto2.Next() {
		databases = append(databases, toto2.Record().Values[0].(string))
		//databases = append(databases, "toto")
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("databases", databases); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
