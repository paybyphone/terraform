package pingdom

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"app_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["app_key"],
			},

			"email_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["email_address"],
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["password"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"pingdom_http_check": resourcePingdomHTTPCheck(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"app_key":       "The application key to connect to Pingdom with.",
		"email_address": "The email address of the Pingdom account.",
		"password":      "The password of the Pingdom account.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AppKey:       d.Get("app_key").(string),
		EmailAddress: d.Get("email_address").(string),
		Password:     d.Get("password").(string),
	}
	return config.Client()
}
