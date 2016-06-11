---
layout: "acme"
page_title: "ACME: certificate"
sidebar_current: "docs-acme-resource-certificate"
description: |-
  Creates and manages an ACME certificate.
---

# acme\_certificate

Use this resource to create and manage an ACME TLS certificate.

~> **NOTE:** Note that the example uses the
[Let's Encrypt staging environment][1]. If you are using Let's Encrypt, make
sure you change the URL to the correct endpoint (currently
`https://acme-v01.api.letsencrypt.org`).

~> **NOTE:** Some current ACME CA implementations like [boulder][2] strip
most of the organization information out of a certificate request's subject,
so you may wish to confirm with the CA what behaviour to expect when using the
`cert_request_pem` argument with this resource.

## Example

**Full example with `common_name` and `subject_alternative_names`**

```
# Create the private key for the registration (not the certificate)
resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

# Set up a registration using a private key from tls_private_key
resource "acme_registration" "reg" {
  server_url      = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  email_address   = "nobody@example.com"
}

# Create a certificate
resource "acme_certificate" "certificate" {
  server_url                = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem           = "${tls_private_key.private_key.private_key_pem}"
  common_name               = "www.example.com"
  subject_alternative_names = ["www2.example.com"]

  dns_challenge {
    provider = "route53"
  }

  registration_url = "${acme_registration.reg.id}"
}
```

**Full example with `cert_request_pem`**

```
resource "tls_private_key" "reg_private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  server_url      = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem = "${tls_private_key.private_key.private_key_pem}"
  email_address   = "nobody@example.com"
}

resource "tls_private_key" "cert_private_key" {
  algorithm = "RSA"
}

resource "tls_cert_request" "req" {
  key_algorithm   = "RSA"
  private_key_pem = "${tls_private_key.cert_private_key.private_key_pem}"
  dns_names       = ["www.example.com", "www2.example.com"]

  subject {
    common_name  = "www.example.com"
  }
}

resource "acme_certificate" "certificate" {
  server_url       = "https://acme-staging.api.letsencrypt.org/directory"
  account_key_pem  = "${tls_private_key.reg_private_key.private_key_pem}"
  cert_request_pem = "${tls_cert_request.req.cert_request_pem}"

  dns_challenge {
    provider = "route53"
  }

  registration_url = "${acme_registration.reg.id}"
}
```

## Argument Reference

The resource takes the following arguments:

 * `server_url` (Required) - The URL of the ACME directory endpoint.
 * `account_key_pem` (Required) - The private key used to sign requests. This
    will be the private key that will be registered to the account.

**WIP - Fill me in!!**

## Attribute Reference

The following attributes are exported:

 * `id` - The full URL of the certificate. Same as `cert_url`.

**WIP - Fill me in!!**

[1]: https://letsencrypt.org/docs/staging-environment/
[2]: https://github.com/letsencrypt/boulder
