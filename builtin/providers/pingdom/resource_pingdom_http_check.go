package pingdom

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/pingdom-go-sdk/resource/checks"
)

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
	svc := meta.(*ProviderPingdomClient).checksconn
	params := checks.CreateCheckInput{
		CheckConfiguration:     expandBaseCheck(d),
		CheckConfigurationHTTP: expandHTTPCheck(d),
	}
	params.Type = "http"

	out, err := svc.CreateCheck(params)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(out.Check.ID))

	return resourcePingdomHTTPCheckRead(d, meta)
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
