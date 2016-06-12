package acme

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func registrationResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: registrationSchemaFull(),
	}
	d := r.TestResourceData()

	d.SetId("regurl")
	d.Set("server_url", "https://acme-staging.api.letsencrypt.org/directory")
	d.Set("account_key_pem", "acctkey")
	d.Set("email_address", "nobody@example.com")
	d.Set("registration_body", "regbody")
	d.Set("registration_url", "regurl")
	d.Set("registration_new_authz_url", "new-authz")
	d.Set("registration_tos_url", "tosurl")

	return d
}

func TestACME_registrationSchemaFull(t *testing.T) {
	m := registrationSchemaFull()
	fields := []string{"email_address", "registration_body", "registration_url", "registration_new_authz_url", "registration_tos_url"}
	for _, v := range fields {
		if _, ok := m[v]; ok == false {
			t.Fatalf("Expected %s to be present", v)
		}
	}
}

func TestACME_certificateSchema(t *testing.T) {
	m := registrationSchemaFull()
	fields := []string{
		"common_name",
		"subject_alternative_names",
		"key_type",
		"cert_request_pem",
		"min_days_remaining",
		"dns_challenge",
		"http_challenge_port",
		"tls_challenge_port",
		"registration_url",
		"cert_domain",
		"cert_url",
		"account_ref",
		"private_key_pem",
		"certificate_pem",
	}
	for _, v := range fields {
		if _, ok := m[v]; ok == false {
			t.Fatalf("Expected %s to be present", v)
		}
	}
}
