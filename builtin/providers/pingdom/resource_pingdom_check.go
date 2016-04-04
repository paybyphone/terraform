package pingdom

import "github.com/hashicorp/terraform/helper/schema"

func resourcePingdomHTTPCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomHTTPCheckCreate,
		Read:   resourcePingdomHTTPCheckRead,
		Update: resourcePingdomHTTPCheckUpdate,
		Delete: resourcePingdomHTTPCheckDelete,

		Schema: map[string]*schema.Schema{},
	}
}

func resourcePingdomHTTPCheckCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePingdomHTTPCheckRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePingdomHTTPCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePingdomHTTPCheckDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
