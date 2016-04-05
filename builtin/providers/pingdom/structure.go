package pingdom

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/pingdom-go-sdk/resource/checks"
)

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
			Set:      intSetHash,
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
		"auto_tags": &schema.Schema{
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
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

// httpCustomCheckSchema returns a map[string]*schema.Schema with all the elements
// that are specific to a Custom HTTP check type.
func httpCustomCheckSchema() map[string]*schema.Schema {
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
	for k, v := range httpCheckSchema() {
		m[k] = v
	}
	return m
}

// httpCustomCheckSchemaFull returns a merged baseCheckSchema + customHTTPCheckSchema.
func httpCustomCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range httpCustomCheckSchema() {
		m[k] = v
	}
	return m
}

// tcpCheckSchemaFull returns a merged baseCheckSchema + tcpCheckSchema.
func tcpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range tcpCheckSchema() {
		m[k] = v
	}
	return m
}

// PingCheckSchemaFull simply returns baseCheckSchema() as there are currently
// no type-specific parameters for Ping checks.
func pingCheckSchemaFull() map[string]*schema.Schema {
	return baseCheckSchema()
}

// dnsCheckSchemaFull returns a merged baseCheckSchema + dnsCheckSchema.
func dnsCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range dnsCheckSchema() {
		m[k] = v
	}
	return m
}

// udpCheckSchemaFull returns a merged baseCheckSchema + udpCheckSchema.
func udpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range udpCheckSchema() {
		m[k] = v
	}
	return m
}

// smtpCheckSchemaFull returns a merged baseCheckSchema + smtpCheckSchema.
func smtpCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range smtpCheckSchema() {
		m[k] = v
	}
	return m
}

// pop3CheckSchemaFull returns a merged baseCheckSchema + pop3CheckSchema.
func pop3CheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range pop3CheckSchema() {
		m[k] = v
	}
	return m
}

// imapCheckSchemaFull returns a merged baseCheckSchema + imapCheckSchema.
func imapCheckSchemaFull() map[string]*schema.Schema {
	m := baseCheckSchema()
	for k, v := range imapCheckSchema() {
		m[k] = v
	}
	return m
}

// expandBaseCheck expands all of the base check fields and returns a
// CheckConfiguration struct with all the appropriate fields set.
func expandBaseCheck(d *schema.ResourceData) checks.CheckConfiguration {
	return checks.CheckConfiguration{
		Name:                     d.Get("name").(string),
		Host:                     d.Get("host").(string),
		Paused:                   d.Get("paused").(bool),
		Resolution:               d.Get("resolution").(int),
		ContactIDs:               intSlice(d.Get("contact_ids").(*schema.Set).List()),
		SendToEmail:              d.Get("send_to_email").(bool),
		SendToSMS:                d.Get("send_to_sms").(bool),
		SendToTwitter:            d.Get("send_to_twitter").(bool),
		SendToIphone:             d.Get("send_to_iphone").(bool),
		SendToAndroid:            d.Get("send_to_android").(bool),
		SendNotificationWhenDown: d.Get("send_notification_when_down").(int),
		NotifyAgainEvery:         d.Get("notify_again_every").(int),
		NotifyWhenBackUp:         d.Get("notify_when_back_up").(bool),
		Tags:                     stringSlice(d.Get("tags").(*schema.Set).List()),
		IPv6:                     d.Get("ipv6").(bool),
	}
}

// flattenBaseCheck takes a DetailedCheckEntry struct and sets all the
// appropriate base check fields for the resource.
func flattenBaseCheck(c checks.DetailedCheckEntry, d *schema.ResourceData) {
	d.Set("name", c.Name)
	d.Set("host", c.Hostname)
	// TODO: think about removing this from resources altogether
	// d.Set("paused", c.Paused)
	d.Set("resolution", c.Resolution)
	d.Set("contact_ids", schema.NewSet(intSetHash, interfaceSlice(c.ContactIDs)))
	d.Set("send_to_email", c.SendToEmail)
	d.Set("send_to_sms", c.SendToSMS)
	d.Set("send_to_twitter", c.SendToTwitter)
	d.Set("send_to_iphone", c.SendToIphone)
	d.Set("send_to_android", c.SendToAndroid)
	d.Set("send_notification_when_down", c.SendNotificationWhenDown)
	d.Set("notify_again_every", c.NotifyAgainEvery)
	d.Set("notify_when_back_up", c.NotifyWhenBackUp)
	d.Set("ipv6", c.IPv6)
}

// flattenBaseCheckTags takes a CheckListEntryTags slice, sets user-tagged tags,
// and adds auto-tagged tags as computed values.
func flattenBaseCheckTags(s []checks.CheckListEntryTags, d *schema.ResourceData) {
	var tags, autotags []string
	for _, v := range s {
		switch v.Type {
		case "a":
			autotags = append(autotags, v.Name)
		case "u":
			tags = append(tags, v.Name)
		}
	}
	d.Set("tags", schema.NewSet(schema.HashString, interfaceSlice(tags)))
	d.Set("auto_tags", schema.NewSet(schema.HashString, interfaceSlice(autotags)))
}

// expandHTTPCheck expands all of the base check fields and returns a
// CheckConfigurationHTTP struct with all the appropriate fields set.
func expandHTTPCheck(d *schema.ResourceData) checks.CheckConfigurationHTTP {
	return checks.CheckConfigurationHTTP{
		URL:              d.Get("url").(string),
		Encryption:       d.Get("encryption").(bool),
		Port:             d.Get("port").(int),
		Auth:             d.Get("auth").(string),
		ShouldContain:    d.Get("should_contain").(string),
		ShouldNotContain: d.Get("should_not_contain").(string),
		PostData:         d.Get("post_data").(string),
		RequestHeaders:   stringSlice(d.Get("request_headers").(*schema.Set).List()),
	}
}

// flattenHTTPCheck takes a DetailedCheckEntryHTTP struct and sets all the
// appropriate type-specific check fields for a HTTP check.
func flattenHTTPCheck(c checks.DetailedCheckEntryHTTP, d *schema.ResourceData) {
	d.Set("url", c.URL)
	d.Set("encryption", c.Encryption)
	d.Set("port", c.Port)
	d.Set("auth", fmt.Sprintf("%s:%s", c.Username, c.Password))
	d.Set("should_contain", c.ShouldContain)
	d.Set("should_not_contain", c.ShouldNotContain)
	d.Set("post_data", c.PostData)
	d.Set("request_headers", schema.NewSet(schema.HashString, interfaceSlice(c.RequestHeaders)))
}

// expandHTTPCustomCheck expands all of the base check fields and returns a
// CheckConfigurationHTTPCustom struct with all the appropriate fields set.
func expandHTTPCustomCheck(d *schema.ResourceData) checks.CheckConfigurationHTTPCustom {
	return checks.CheckConfigurationHTTPCustom{
		URL:            d.Get("url").(string),
		Encryption:     d.Get("encryption").(bool),
		Port:           d.Get("port").(int),
		Auth:           d.Get("auth").(string),
		AdditionalURLs: stringSlice(d.Get("additional_urls").(*schema.Set).List()),
	}
}

// flattenHTTPCustomCheck takes a DetailedCheckEntryHTTPCustom struct and sets all the
// appropriate type-specific check fields for a Custom HTTP check.
func flattenHTTPCustomCheck(c checks.DetailedCheckEntryHTTPCustom, d *schema.ResourceData) {
	d.Set("url", c.URL)
	d.Set("encryption", c.Encryption)
	d.Set("port", c.Port)
	d.Set("auth", fmt.Sprintf("%s:%s", c.Username, c.Password))
	d.Set("additional_urls", schema.NewSet(schema.HashString, interfaceSlice(c.AdditionalURLs)))
}

// expandTCPCheck expands all of the base check fields and returns a
// CheckConfigurationTCP struct with all the appropriate fields set.
func expandTCPCheck(d *schema.ResourceData) checks.CheckConfigurationTCP {
	return checks.CheckConfigurationTCP{
		Port:           d.Get("port").(int),
		StringToSend:   d.Get("string_to_send").(string),
		StringToExpect: d.Get("string_to_expect").(string),
	}
}

// flattenTCPCheck takes a DetailedCheckEntryTCP struct and sets all the
// appropriate type-specific check fields for a TCP check.
func flattenTCPCheck(c checks.DetailedCheckEntryTCP, d *schema.ResourceData) {
	d.Set("port", c.Port)
	d.Set("string_to_send", c.StringToSend)
	d.Set("string_to_expect", c.StringToExpect)
}

// expandDNSCheck expands all of the base check fields and returns a
// CheckConfigurationDNS struct with all the appropriate fields set.
func expandDNSCheck(d *schema.ResourceData) checks.CheckConfigurationDNS {
	return checks.CheckConfigurationDNS{
		NameServer: d.Get("nameserver").(string),
		ExpectedIP: d.Get("expected_ip").(string),
	}
}

// flattenDNSCheck takes a DetailedCheckEntryDNS struct and sets all the
// appropriate type-specific check fields for a DNS check.
func flattenDNSCheck(c checks.DetailedCheckEntryDNS, d *schema.ResourceData) {
	d.Set("nameserver", c.DNSServer)
	d.Set("expected_ip", c.ExpectedIP)
}

// expandUDPCheck expands all of the base check fields and returns a
// CheckConfigurationUDP struct with all the appropriate fields set.
func expandUDPCheck(d *schema.ResourceData) checks.CheckConfigurationUDP {
	return checks.CheckConfigurationUDP{
		Port:           d.Get("port").(int),
		StringToSend:   d.Get("string_to_send").(string),
		StringToExpect: d.Get("string_to_expect").(string),
	}
}

// flattenUDPCheck takes a DetailedCheckEntryUDP struct and sets all the
// appropriate type-specific check fields for a UDP check.
func flattenUDPCheck(c checks.DetailedCheckEntryUDP, d *schema.ResourceData) {
	d.Set("port", c.Port)
	d.Set("string_to_send", c.StringToSend)
	d.Set("string_to_expect", c.StringToExpect)
}

// expandSMTPCheck expands all of the base check fields and returns a
// CheckConfigurationSMTP struct with all the appropriate fields set.
func expandSMTPCheck(d *schema.ResourceData) checks.CheckConfigurationSMTP {
	return checks.CheckConfigurationSMTP{
		Port:           d.Get("port").(int),
		Auth:           d.Get("auth").(string),
		Encryption:     d.Get("encryption").(bool),
		StringToExpect: d.Get("string_to_expect").(string),
	}
}

// flattenSMTPCheck takes a DetailedCheckEntrySMTP struct and sets all the
// appropriate type-specific check fields for a SMTP check.
func flattenSMTPCheck(c checks.DetailedCheckEntrySMTP, d *schema.ResourceData) {
	d.Set("port", c.Port)
	d.Set("auth", fmt.Sprintf("%s:%s", c.Username, c.Password))
	d.Set("encryption", c.Encryption)
	d.Set("string_to_expect", c.StringToExpect)
}

// expandPOP3Check expands all of the base check fields and returns a
// CheckConfigurationPOP3 struct with all the appropriate fields set.
func expandPOP3Check(d *schema.ResourceData) checks.CheckConfigurationPOP3 {
	return checks.CheckConfigurationPOP3{
		Port:           d.Get("port").(int),
		Encryption:     d.Get("encryption").(bool),
		StringToExpect: d.Get("string_to_expect").(string),
	}
}

// flattenPOP3Check takes a DetailedCheckEntryPOP3 struct and sets all the
// appropriate type-specific check fields for a POP3 check.
func flattenPOP3Check(c checks.DetailedCheckEntryPOP3, d *schema.ResourceData) {
	d.Set("port", c.Port)
	d.Set("encryption", c.Encryption)
	d.Set("string_to_expect", c.StringToExpect)
}

// expandIMAPCheck expands all of the base check fields and returns a
// CheckConfigurationIMAP struct with all the appropriate fields set.
func expandIMAPCheck(d *schema.ResourceData) checks.CheckConfigurationIMAP {
	return checks.CheckConfigurationIMAP{
		Port:           d.Get("port").(int),
		Encryption:     d.Get("encryption").(bool),
		StringToExpect: d.Get("string_to_expect").(string),
	}
}

// flattenIMAPCheck takes a DetailedCheckEntryIMAP struct and sets all the
// appropriate type-specific check fields for a IMAP check.
func flattenIMAPCheck(c checks.DetailedCheckEntryIMAP, d *schema.ResourceData) {
	d.Set("port", c.Port)
	d.Set("encryption", c.Encryption)
	d.Set("string_to_expect", c.StringToExpect)
}

// intSetHash just returns the number in a specific element for a int set.
func intSetHash(v interface{}) int {
	return v.(int)
}

// stringSlice converts an interface slice to a string slice.
func stringSlice(src []interface{}) []string {
	var dst []string
	for _, v := range src {
		dst = append(dst, v.(string))
	}
	return dst
}

// intSlice converts an interface slice to an int slice.
func intSlice(src []interface{}) []int {
	var dst []int
	for _, v := range src {
		dst = append(dst, v.(int))
	}
	return dst
}

// interfaceSlice converts a slice of string or int back to an interface
// slice.
func interfaceSlice(src interface{}) []interface{} {
	var dst []interface{}
	switch w := src.(type) {
	case []string:
		for _, v := range w {
			dst = append(dst, v)
		}
	case []int:
		for _, v := range w {
			dst = append(dst, v)
		}
	case map[string]string:
		for _, v := range w {
			dst = append(dst, v)
		}
	default:
		panic("Unsupported type for interaceSlice() conversion")
	}
	return dst
}
