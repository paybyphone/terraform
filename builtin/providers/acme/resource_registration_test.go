package acme

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccACMERegistration_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckReg(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckACMERegistrationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccACMERegistrationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACMERegistrationValid("acme_registration.reg"),
				),
			},
		},
	})
}

func testAccCheckACMERegistrationValid(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find ACME registration: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AMI data source ID not set")
		}

		actual := rs.Primary.Attributes["registration.uri"]
		expected := rs.Primary.ID

		if actual != expected {
			return fmt.Errorf("Expected ID to be %s, got %s", expected, actual)
		}
		return nil
	}
}

func testAccCheckACMERegistrationDestroy(s *terraform.State) error {
	// TODO: Fill this in - need to figure out how to query reg in lego
	return nil
}

func testAccPreCheckReg(t *testing.T) {
	if v := os.Getenv("ACME_EMAIL_ADDRESS"); v == "" {
		t.Fatal("ACME_EMAIL_ADDRESS must be set for the registration acceptance test")
	}
}

func testAccACMERegistrationConfig() string {
	return fmt.Sprintf(`
resource "tls_private_key" "private_key" {
    algorithm = "RSA"
}

resource "acme_registration" "reg" {
	server_url = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  email_address = "%s"

}
`, os.Getenv("ACME_EMAIL_ADDRESS"))
}
