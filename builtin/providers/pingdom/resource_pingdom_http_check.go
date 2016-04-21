package pingdom

import "github.com/hashicorp/terraform/helper/schema"

// resourcePingdomHTTPCheck defines the resource for the pingdom_http_check
// Terraform resource.
func resourcePingdomHTTPCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomHTTPCheckCreate,
		Read:   resourcePingdomHTTPCheckRead,
		Update: resourcePingdomHTTPCheckUpdate,
		Delete: resourcePingdomHTTPCheckDelete,

		Schema: httpCheckSchemaFull(),
	}
}

// resourcePingdomHTTPCheckCreate runs the create portion of the pingdom_http_check
// resource.
func resourcePingdomHTTPCheckCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourcePingdomHTTPCheckRead runs the read portion of the pingdom_http_check
// resource.
func resourcePingdomHTTPCheckRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourcePingdomHTTPCheckUpdate runs the update portion of the pingdom_http_check
// resource.
func resourcePingdomHTTPCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourcePingdomHTTPCheckDelete runs the delete portion of the pingdom_http_check
// resource.
func resourcePingdomHTTPCheckDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
