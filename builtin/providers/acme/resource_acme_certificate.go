package acme

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xenolf/lego/acme"
)

func resourceACMECertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceACMECertificateCreate,
		Read:   resourceACMECertificateRead,
		Delete: resourceACMECertificateDelete,

		Schema: certificateSchemaFull(),
	}
}

func resourceACMECertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := expandACMEClient(d, d.Get("registration_url").(string))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("dns_challenge"); ok {
		setDNSChallenge(client, v.(*schema.Set).List()[0].(map[string]interface{}))
	} else {
		client.SetHTTPAddress(":" + strconv.Itoa(d.Get("http_challenge_port").(int)))
		client.SetTLSAddress(":" + strconv.Itoa(d.Get("tls_challenge_port").(int)))
	}

	var cert acme.CertificateResource
	var errs map[string]error

	if v, ok := d.GetOk("cert_request_pem"); ok {
		csr, err := csrFromPEM([]byte(v.(string)))
		if err != nil {
			return err
		}
		cert, errs = client.ObtainCertificateForCSR(*csr, true)
	} else {
		cn := d.Get("common_name").(string)
		domains := []string{cn}
		if s, ok := d.GetOk("subject_alternative_names"); ok {
			for _, v := range stringSlice(s.(*schema.Set).List()) {
				if v == cn {
					return fmt.Errorf("common name %s should not appear in SAN list", v)
				}
				domains = append(domains, v)
			}
		}

		cert, errs = client.ObtainCertificate(domains, true, nil)
	}

	if len(errs) > 0 {
		messages := []string{}
		for k, v := range errs {
			messages = append(messages, fmt.Sprintf("%s: %s", k, v))
		}
		return fmt.Errorf("Errors were encountered creating the certificate:\n    %s", strings.Join(messages, "\n    "))
	}

	// done! save the cert
	saveCertificateResource(d, cert)

	return nil
}

// resourceACMECertificateRead renews the certificate if it is close to expiry.
// This value is controlled by the min_days_remaining attribute - if this value
// less than zero, the certificate is never renewed.
func resourceACMECertificateRead(d *schema.ResourceData, meta interface{}) error {
	mindays := d.Get("min_days_remaining").(int)
	if mindays < 0 {
		log.Printf("[WARN] min_days_remaining is set to less than 0, certificate will never be renewed")
		return nil
	}

	client, _, err := expandACMEClient(d, d.Get("registration_url").(string))
	if err != nil {
		return err
	}

	cert := expandCertificateResource(d)
	remaining, err := certDaysRemaining(cert)
	if err != nil {
		return err
	}

	if int64(mindays) >= remaining {
		if v, ok := d.GetOk("dns_challenge"); ok {
			setDNSChallenge(client, v.(*schema.Set).List()[0].(map[string]interface{}))
		} else {
			client.SetHTTPAddress(":" + strconv.Itoa(d.Get("http_challenge_port").(int)))
			client.SetTLSAddress(":" + strconv.Itoa(d.Get("tls_challenge_port").(int)))
		}
		newCert, err := client.RenewCertificate(cert, true)
		if err != nil {
			return err
		}
		saveCertificateResource(d, newCert)
	}

	return nil
}

// resourceACMECertificateDelete "deletes" the certificate by revoking it.
func resourceACMECertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := expandACMEClient(d, d.Get("registration_url").(string))
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("certificate_pem"); ok {
		err = client.RevokeCertificate([]byte(v.(string)))
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
