package acme

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/providers/dns/cloudflare"
	"github.com/xenolf/lego/providers/dns/digitalocean"
	"github.com/xenolf/lego/providers/dns/dnsimple"
	"github.com/xenolf/lego/providers/dns/dyn"
	"github.com/xenolf/lego/providers/dns/gandi"
	"github.com/xenolf/lego/providers/dns/googlecloud"
	"github.com/xenolf/lego/providers/dns/namecheap"
	"github.com/xenolf/lego/providers/dns/rfc2136"
	"github.com/xenolf/lego/providers/dns/route53"
	"github.com/xenolf/lego/providers/dns/vultr"
)

// baseCheckSchema returns a map[string]*schema.Schema with all the elements
// necessary to build the base elements of an ACME resource schema. Use this,
// along with a schema helper of a specific check type, to return the full
// schema.
func baseCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"account_key_pem": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

// registrationSchema returns a map[string]*schema.Schema with all the elements
// that are specific to an ACME registration resource.
func registrationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Even though the ACME spec does allow for multiple contact types, lego
		// only works with a single email address.
		"email_address": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"registration_body": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"registration_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"registration_new_authz_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"registration_tos_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// certificateSchema returns a map[string]*schema.Schema with all the elements
// that are specific to an ACME certificate resource.
//
// The initial version of this only supports DNS challenges.
func certificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"common_name": &schema.Schema{
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"cert_request_pem"},
		},
		"subject_alternative_names": &schema.Schema{
			Type:          schema.TypeSet,
			Optional:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			Set:           schema.HashString,
			ForceNew:      true,
			ConflictsWith: []string{"cert_request_pem"},
		},
		"key_type": &schema.Schema{
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Default:       "2048",
			ConflictsWith: []string{"cert_request_pem"},
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(string)
				found := false
				for _, w := range []string{"P256", "P384", "2048", "4096", "8192"} {
					if value == w {
						found = true
					}
				}
				if found == false {
					errors = append(errors, fmt.Errorf(
						"Certificate key type must be one of P256, P384, 2048, 4096, or 8192"))
				}
				return
			},
		},
		"cert_request_pem": &schema.Schema{
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"common_name", "subject_alternative_names", "key_type"},
		},
		"min_days_remaining": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  7,
			ForceNew: true,
		},
		"dns_challenge": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Set:      dnsChallengeSetHash,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"provider": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"config": &schema.Schema{
						Type:     schema.TypeMap,
						Optional: true,
						ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
							value := v.(map[string]interface{})
							bad := false
							for _, w := range value {
								switch w.(type) {
								case string:
									continue
								default:
									bad = true
								}
							}
							if bad == true {
								errors = append(errors, fmt.Errorf(
									"DNS challenge config map values must be strings only"))
							}
							return
						},
					},
				},
			},
			ForceNew: true,
		},
		"http_challenge_port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  80,
			ForceNew: true,
		},
		"tls_challenge_port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  443,
			ForceNew: true,
		},
		"registration_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cert_domain": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"cert_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_ref": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_key_pem": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"certificate_pem": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// registrationSchemaFull returns a merged baseCheckSchema + registrationSchema.
func registrationSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range registrationSchema() {
		m[k] = v
	}
	return m
}

// certificateSchemaFull returns a merged baseCheckSchema + certificateSchema.
func certificateSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range certificateSchema() {
		m[k] = v
	}
	return m
}

// acmeUser implements acme.User.
type acmeUser struct {

	// The email address for the account.
	Email string

	// The registration resource object.
	Registration *acme.RegistrationResource

	// The private key for the account.
	key crypto.PrivateKey
}

func (u acmeUser) GetEmail() string {
	return u.Email
}
func (u acmeUser) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}
func (u acmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// expandACMEUser creates a new instance of an ACME user from set
// email_address and private_key_pem fields, and a registration
// if one exists.
func expandACMEUser(d *schema.ResourceData) (*acmeUser, error) {
	key, err := privateKeyFromPEM([]byte(d.Get("account_key_pem").(string)))
	if err != nil {
		return nil, err
	}

	user := &acmeUser{
		key: key,
	}

	// only set these fields if they are in the schema.
	if v, ok := d.GetOk("email_address"); ok {
		user.Email = v.(string)
	}
	//if reg, ok := expandACMERegistration(d); ok {
	//user.Registration = reg
	//}

	return user, nil
}

// expandACMERegistration takes a *schema.ResourceData and builds an
// *acme.RegistrationResource, complete with body, for use in lego calls.
//
// Note that this function will return nil, false if any of the computed
// registration fields are un-readable, to allow non-fatal behaviour when the
// data does not exist.
func expandACMERegistration(d *schema.ResourceData) (*acme.RegistrationResource, bool) {
	reg := acme.RegistrationResource{}
	var v interface{}
	var ok bool

	if v, ok = d.GetOk("registration_body"); ok == false {
		return nil, false
	}
	body := acme.Registration{}
	err := json.Unmarshal([]byte(v.(string)), &body)
	if err != nil {
		log.Printf("[DEBUG] Error reading JSON for registration body: %s", err.Error())
		return nil, false
	}
	reg.Body = body

	if v, ok = d.GetOk("registration_url"); ok == false {
		return nil, false
	}
	reg.URI = v.(string)

	if v, ok = d.GetOk("registration_new_authz_url"); ok == false {
		return nil, false
	}
	reg.NewAuthzURL = v.(string)

	// TOS url can be blank, ensure we don't fail on this
	if v, ok = d.GetOk("registration_tos_url"); ok == false {
		log.Printf("[DEBUG] ACME reg %s: registration_tos_url is blank", reg.URI)
	}
	reg.TosURL = v.(string)

	return &reg, true
}

// saveACMERegistration takes an *acme.RegistrationResource and sets the appropriate fields
// for a registration resource.
func saveACMERegistration(d *schema.ResourceData, reg *acme.RegistrationResource) error {
	d.SetId(reg.URI)

	body, err := json.Marshal(reg)
	if err != nil {
		return fmt.Errorf("error reading registration body: %s", err.Error())
	}
	d.Set("registration_body", string(body))

	d.Set("registration_url", reg.URI)
	d.Set("registration_new_authz_url", reg.NewAuthzURL)
	d.Set("registration_tos_url", reg.TosURL)

	return nil
}

// expandACMEClient creates a connection to an ACME server from resource data,
// and also returns the user.
//
// If regURL is supplied, the registration information is loaded in to the user's registration
// via the registration URL.
func expandACMEClient(d *schema.ResourceData, regURL string) (*acme.Client, *acmeUser, error) {
	user, err := expandACMEUser(d)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting user data: %s", err.Error())
	}

	// Note this function is used by both the registration and certificate
	// resources, but key type is not necessary during registration, so
	// it's okay if it's empty for that.
	var keytype string
	if v, ok := d.GetOk("key_type"); ok {
		keytype = v.(string)
	}

	client, err := acme.NewClient(d.Get("server_url").(string), user, acme.KeyType(keytype))
	if err != nil {
		return nil, nil, err
	}

	if regURL != "" {
		user.Registration = &acme.RegistrationResource{
			URI: regURL,
		}
		reg, err := client.QueryRegistration()
		if err != nil {
			return nil, nil, err
		}
		user.Registration = reg
	}

	return client, user, nil
}

// expandCertificateResource takes saved state in the certificate resource
// and returns an acme.CertificateResource.
func expandCertificateResource(d *schema.ResourceData) acme.CertificateResource {
	cert := acme.CertificateResource{
		Domain:      d.Get("cert_domain").(string),
		CertURL:     d.Get("cert_url").(string),
		AccountRef:  d.Get("account_ref").(string),
		PrivateKey:  []byte(d.Get("private_key_pem").(string)),
		Certificate: []byte(d.Get("certificate_pem").(string)),
		CSR:         []byte(d.Get("cert_request_pem").(string)),
	}
	return cert
}

// saveCertificateResource takes an acme.CertificateResource and sets fields.
func saveCertificateResource(d *schema.ResourceData, cert acme.CertificateResource) {
	d.SetId(cert.CertURL)
	d.Set("cert_domain", cert.Domain)
	d.Set("cert_url", cert.CertURL)
	d.Set("account_ref", cert.AccountRef)
	d.Set("private_key_pem", string(cert.PrivateKey))
	d.Set("certificate_pem", string(cert.Certificate))
	// note that CSR is skipped as it is not computed.
}

// certDaysRemaining takes an acme.CertificateResource, parses the
// certificate, and computes the days that it has remaining.
func certDaysRemaining(cert acme.CertificateResource) (int64, error) {
	x509Certs, err := parsePEMBundle(cert.Certificate)
	if err != nil {
		return 0, err
	}

	c := x509Certs[0]

	expiry := c.NotAfter.Unix()
	now := time.Now().Unix()

	return (expiry - now) / 86400, nil
}

// parsePEMBundle parses a certificate bundle from top to bottom and returns
// a slice of x509 certificates. This function will error if no certificates are found.
//
// Note that this is taken directly from lego - may look at exporting it there.
func parsePEMBundle(bundle []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate
	var certDERBlock *pem.Block

	for {
		certDERBlock, bundle = pem.Decode(bundle)
		if certDERBlock == nil {
			break
		}

		if certDERBlock.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				return nil, err
			}
			certificates = append(certificates, cert)
		}
	}

	if len(certificates) == 0 {
		return nil, errors.New("No certificates were found while parsing the bundle.")
	}

	return certificates, nil
}

// setDNSChallenge takes an *acme.Client and the DNS challenge complex
// structure as a map[string]interface{}, and configues the client to only
// allow a DNS challenge with the configured provider.
func setDNSChallenge(client *acme.Client, m map[string]interface{}) error {
	var provider acme.ChallengeProvider
	var err error
	var providerName string

	if v, ok := m["provider"]; ok && v.(string) != "" {
		providerName = v.(string)
	} else {
		return fmt.Errorf("DNS challenge provider not defined")
	}
	// Config only needs to be set if it's defined, otherwise existing env/SDK
	// defaults are fine.
	if v, ok := m["config"]; ok {
		for k, v := range v.(map[string]interface{}) {
			os.Setenv(k, v.(string))
		}
	}

	// The below list was taken from lego's cli_handlers.go from this writing
	// (June 2016). As such, this might not be a full list of supported
	// providers. If a specific provider is wanted, check lego first, and then
	// write the provider for lego and put in a PR.
	switch providerName {
	case "cloudflare":
		provider, err = cloudflare.NewDNSProvider()
	case "digitalocean":
		provider, err = digitalocean.NewDNSProvider()
	case "dnsimple":
		provider, err = dnsimple.NewDNSProvider()
	case "dyn":
		provider, err = dyn.NewDNSProvider()
	case "gandi":
		provider, err = gandi.NewDNSProvider()
	case "gcloud":
		provider, err = googlecloud.NewDNSProvider()
	case "manual":
		provider, err = acme.NewDNSProviderManual()
	case "namecheap":
		provider, err = namecheap.NewDNSProvider()
	case "route53":
		provider, err = route53.NewDNSProvider()
	case "rfc2136":
		provider, err = rfc2136.NewDNSProvider()
	case "vultr":
		provider, err = vultr.NewDNSProvider()
	default:
		return fmt.Errorf("%s: unsupported DNS challenge provider", providerName)
	}
	if err != nil {
		return err
	}
	client.SetChallengeProvider(acme.DNS01, provider)
	client.ExcludeChallenges([]acme.Challenge{acme.HTTP01, acme.TLSSNI01})

	return nil
}

// dnsChallengeSetHash computes the hash for the DNS challenge.
func dnsChallengeSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["provider"].(string)))
	for k, v := range m["config"].(map[string]interface{}) {
		buf.WriteString(fmt.Sprintf("%s-%s-", k, v.(string)))
	}
	return hashcode.String(buf.String())
}

// stringSlice converts an interface slice to a string slice.
func stringSlice(src []interface{}) []string {
	var dst []string
	for _, v := range src {
		dst = append(dst, v.(string))
	}
	return dst
}

// privateKeyFromPEM converts a PEM block into a crypto.PrivateKey.
func privateKeyFromPEM(pemData []byte) (crypto.PrivateKey, error) {
	var result *pem.Block
	rest := pemData
	for {
		result, rest = pem.Decode([]byte(rest))
		if result == nil {
			return nil, fmt.Errorf("Cannot decode supplied PEM data")
		}
		switch result.Type {
		case "RSA PRIVATE KEY":
			return x509.ParsePKCS1PrivateKey(result.Bytes)
		case "EC PRIVATE KEY":
			return x509.ParseECPrivateKey(result.Bytes)
		}
	}
}

// csrFromPEM converts a PEM block into an *x509.CertificateRequest.
func csrFromPEM(pemData []byte) (*x509.CertificateRequest, error) {
	var result *pem.Block
	rest := pemData
	for {
		result, rest = pem.Decode([]byte(rest))
		if result == nil {
			return nil, fmt.Errorf("Cannot decode supplied PEM data")
		}
		if result.Type == "CERTIFICATE REQUEST" {
			return x509.ParseCertificateRequest(result.Bytes)
		}
	}
}
