package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"golang.org/x/crypto/ocsp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xenolf/lego/acme"
)

func TestAccACMECertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCert(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckACMECertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccACMECertificateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACMECertificateValid("acme_certificate.certificate"),
				),
			},
		},
	})
}

func testAccCheckACMECertificateValid(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find ACME certificate: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ACME certificate ID not set")
		}

		cert := rs.Primary.Attributes["certificate_pem"]
		key := rs.Primary.Attributes["private_key_pem"]

		x509Certs, err := parsePEMBundle([]byte(cert))
		if err != nil {
			return err
		}
		x509Cert := x509Certs[0]
		privateKey, err := privateKeyFromPEM([]byte(key))
		if err != nil {
			return err
		}

		var privPub crypto.PublicKey

		switch v := privateKey.(type) {
		case *rsa.PrivateKey:
			privPub = v.Public()
		case *ecdsa.PrivateKey:
			privPub = v.Public()
		}

		if reflect.DeepEqual(x509Cert.PublicKey, privPub) != true {
			return fmt.Errorf("Public key for cert and private key don't match: %#v, %#v", x509Cert.PublicKey, privPub)
		}

		// TODO: domains
		return nil
	}
}

func testAccCheckACMECertificateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "acme_certificate" {
			continue
		}

		cert := rs.Primary.Attributes["certificate_pem"]

		_, resp, err := acme.GetOCSPForCert([]byte(cert))
		if err != nil {
			return fmt.Errorf("Bad: %s", err.Error())
		}

		if resp.Status != ocsp.Revoked {
			return fmt.Errorf("Expected status to be revoked, got OCSP status %d", resp.Status)
		}

		return nil
	}
	return fmt.Errorf("acme_certificate resource not found")
}

func testAccPreCheckCert(t *testing.T) {
	if v := os.Getenv("ACME_EMAIL_ADDRESS"); v == "" {
		t.Fatal("ACME_EMAIL_ADDRESS must be set for the certificate acceptance test")
	}
	if v := os.Getenv("ACME_CERT_DOMAIN"); v == "" {
		t.Fatal("ACME_CERT_DOMAIN must be set for the certificate acceptance test")
	}
	if v := os.Getenv("AWS_PROFILE"); v == "" {
		if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
			t.Fatal("AWS_ACCESS_KEY_ID must be set for the certificate acceptance test")
		}
		if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
			t.Fatal("AWS_SECRET_ACCESS_KEY must be set for the certificate acceptance test")
		}
	}
	if v := os.Getenv("AWS_DEFAULT_REGION"); v == "" {
		log.Println("[INFO] Test: Using us-west-2 as test region")
		os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
	}
}

func testAccACMECertificateConfig() string {
	return fmt.Sprintf(`
resource "tls_private_key" "private_key" {
    algorithm = "RSA"
}

resource "acme_registration" "reg" {
	server_url = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  email_address = "%s"
}

resource "acme_certificate" "certificate" {
	server_url = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  common_name = "%s"
	dns_challenge {
		provider = "route53"
	}
	registration_uri = "${acme_registration.reg.registration_uri}"
}
`, os.Getenv("ACME_EMAIL_ADDRESS"), os.Getenv("ACME_CERT_DOMAIN"))
}
