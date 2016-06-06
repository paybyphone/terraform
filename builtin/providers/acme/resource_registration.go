package acme

import "github.com/hashicorp/terraform/helper/schema"

func resourceACMERegistration() *schema.Resource {
	return &schema.Resource{
		Create: resourceACMERegistrationCreate,
		Delete: resourceACMERegistrationDelete,

		Schema: registrationSchema(),
	}
}

func resourceACMERegistrationCreate(d *schema.ResourceData, meta interface{}) error {
	// register and agree to the TOS
	client, err := expandACMEClient(d)
	if err != nil {
		return err
	}
	reg, err := client.Register()
	if err != nil {
		return err
	}
	err = client.AgreeToTOS()
	if err != nil {
		return err
	}

	// save the reg
	err = saveACMERegistration(d, reg)
	if err != nil {
		return err
	}

	return nil
}
func resourceACMERegistrationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := expandACMEClient(d)
	if err != nil {
		return err
	}
	err = client.DeleteRegistration()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
