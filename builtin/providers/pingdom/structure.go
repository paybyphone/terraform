package pingdom

import "github.com/hashicorp/terraform/helper/schema"

// baseCheckSchema returns a map[string]*schema.Schema with all the elements
// necessary to build the base elements of a check schema. Use this, along
// with a schema helper of a specific check type, to return the full schema.
func baseCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"host": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"paused": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"resolution": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"contact_ids": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Set:      contactIdsHash,
		},
		"send_to_email": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_to_sms": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_to_twitter": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_to_iphone": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_to_android": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"send_notification_when_down": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"notify_again_every": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"notify_when_back_up": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"tags": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
		"ipv6": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

// httpCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a HTTP check type.
func httpCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"encryption": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"auth": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"should_contain": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"should_not_contain": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"post_data": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"request_headers": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
	}
}

// customHTTPCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a Custom HTTP check type.
func customHTTPCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"encryption": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"auth": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"additional_urls": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
		},
	}
}

// tcpCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a TCP check type.
func tcpCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"string_to_send": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"string_to_expect": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// dnsCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a DNS check type.
func dnsCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"expected_ip": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"nameserver": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// udpCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a UDP check type.
func udpCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"string_to_send": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"string_to_expect": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// smtpCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a SMTP check type.
func smtpCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"auth": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"string_to_expect": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"encryption": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

// pop3CheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a POP3 check type.
func pop3CheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string_to_expect": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"encryption": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

// imapCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a IMAP check type.
func imapCheckSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string_to_expect": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"encryption": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

// httpCheckSchemaFull returns a merged baseCheckSchema + httpCheckSchema.
func httpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range httpCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// customHTTPCheckSchemaFull returns a merged baseCheckSchema + customHTTPCheckSchema.
func customHTTPCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range customHTTPCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// tcpCheckSchemaFull returns a merged baseCheckSchema + tcpCheckSchema.
func tcpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range tcpCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// PingCheckSchemaFull simply returns baseCheckSchema() as there are currently
// no type-specific parameters for Ping checks.
func PingCheckSchemaFull() map[string]*schema.Schema {
	return baseCheckSchema()
}

// dnsCheckSchemaFull returns a merged baseCheckSchema + dnsCheckSchema.
func dnsCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range dnsCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// udpCheckSchemaFull returns a merged baseCheckSchema + udpCheckSchema.
func udpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range udpCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// smtpCheckSchemaFull returns a merged baseCheckSchema + smtpCheckSchema.
func smtpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range smtpCheckSchemaFull() {
		m[k] = v
	}
	return m
}

// pop3CheckSchemaFull returns a merged baseCheckSchema + pop3CheckSchema.
func pop3CheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range pop3CheckSchemaFull() {
		m[k] = v
	}
	return m
}

// imapCheckSchemaFull returns a merged baseCheckSchema + imapCheckSchema.
func imapCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range imapCheckSchemaFull() {
		m[k] = v
	}
	return m
}

func contactIdsHash(v interface{}) int {
	return v.(int)
}
