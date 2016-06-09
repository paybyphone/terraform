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
	"time"

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
					testAccCheckACMECertificateValid("acme_certificate.certificate", "www", "www2"),
				),
			},
		},
	})
}

func TestAccACMECertificate_CSR(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCert(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckACMECertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccACMECertificateCSRConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckACMECertificateValid("acme_certificate.certificate", "www3", "www4"),
				),
			},
		},
	})
}

func testAccCheckACMECertificateValid(n, cn, san string) resource.TestCheckFunc {
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

		// Skip the private key test if we have an empty key. This is a legit case
		// that comes up when a CSR is supplied instead of creating a cert from
		// scratch.
		if key != "" {
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
		}

		// domains
		domain := "." + os.Getenv("ACME_CERT_DOMAIN")
		expectedCN := cn + domain
		expectedSANs := []string{cn + domain, san + domain}

		actualCN := x509Cert.Subject.CommonName
		actualSANs := x509Cert.DNSNames

		if expectedCN != actualCN {
			return fmt.Errorf("Expected common name to be %s, got %s", expectedCN, actualCN)
		}

		if reflect.DeepEqual(expectedSANs, actualSANs) != true {
			return fmt.Errorf("Expected SANs to be %#v, got %#v", expectedSANs, actualSANs)
		}

		return nil
	}
}

func testAccCheckACMECertificateCSRSubject(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find ACME certificate: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ACME certificate ID not set")
		}

		cert := rs.Primary.Attributes["certificate_pem"]

		x509Certs, err := parsePEMBundle([]byte(cert))
		if err != nil {
			return err
		}
		x509Cert := x509Certs[0]

		expectedOrg := []string{"ACME Examples, Inc"}
		actualOrg := x509Cert.Subject.Organization

		if reflect.DeepEqual(expectedOrg, actualOrg) != true {
			return fmt.Errorf("Expected org to be %#v, got %#v", expectedOrg, actualOrg)
		}

		return nil
	}
}

func testAccCheckACMECertificateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "acme_certificate" {
			continue
		}

		cert := rs.Primary.Attributes["certificate_pem"]

		// Add a state waiter for the OCSP status of the cert, as we have been
		// seeing sparodic failures, so it's possible revociation is not
		// instant.
		state := &resource.StateChangeConf{
			Pending:    []string{"Good"},
			Target:     []string{"Revoked"},
			Refresh:    testAccCheckACMECertificateDestroyRefreshFunc(cert),
			Timeout:    3 * time.Minute,
			MinTimeout: 15 * time.Second,
			Delay:      5 * time.Second,
		}

		_, err := state.WaitForState()
		if err != nil {
			return fmt.Errorf("Cert did not revoke: %s", err.Error())
		}

		return nil
	}
	return fmt.Errorf("acme_certificate resource not found")
}

func testAccCheckACMECertificateDestroyRefreshFunc(cert string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		_, resp, err := acme.GetOCSPForCert([]byte(cert))
		if err != nil {
			return nil, "", fmt.Errorf("Bad: %s", err.Error())
		}
		switch resp.Status {
		case ocsp.Revoked:
			return cert, "Revoked", nil
		case ocsp.Good:
			return cert, "Good", nil
		default:
			return nil, "", fmt.Errorf("Bad status: OCSP status %d", resp.Status)
		}
	}
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
variable "email_address" {
  default = "%s"
}

variable "domain" {
  default = "%s"
}

resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  server_url      = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  email_address   = "${var.email_address}"
}

resource "acme_certificate" "certificate" {
  server_url                = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem           = "${tls_private_key.private_key.private_key_pem}"
  common_name               = "www.${var.domain}"
  subject_alternative_names = ["www2.${var.domain}"]

  dns_challenge {
    provider = "route53"
  }

  registration_uri = "${acme_registration.reg.registration_uri}"
}
`, os.Getenv("ACME_EMAIL_ADDRESS"), os.Getenv("ACME_CERT_DOMAIN"))
}

func testAccACMECertificateCSRConfig() string {
	return fmt.Sprintf(`
variable "email_address" {
  default = "%s"
}

variable "domain" {
  default = "%s"
}

resource "tls_private_key" "reg_private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  server_url      = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.reg_private_key.private_key_pem}"
  email_address   = "${var.email_address}"
}

resource "tls_private_key" "cert_private_key" {
  algorithm = "RSA"
}

resource "tls_cert_request" "req" {
  key_algorithm   = "RSA"
  private_key_pem = "${tls_private_key.cert_private_key.private_key_pem}"
  dns_names       = ["www3.${var.domain}", "www4.${var.domain}"]

  subject {
    common_name  = "www3.${var.domain}"
  }
}

resource "acme_certificate" "certificate" {
  server_url       = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem  = "${tls_private_key.reg_private_key.private_key_pem}"
  cert_request_pem = "${tls_cert_request.req.cert_request_pem}"

  dns_challenge {
    provider = "route53"
  }

  registration_uri = "${acme_registration.reg.registration_uri}"
}
`, os.Getenv("ACME_EMAIL_ADDRESS"), os.Getenv("ACME_CERT_DOMAIN"))
}
