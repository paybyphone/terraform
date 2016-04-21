package pingdom

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paybyphone/pingdom-go-sdk/resource/checks"
)

func baseCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: baseCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("name", "My check")
	d.Set("host", "example.com")
	d.Set("resolution", 1)
	d.Set("contact_ids", schema.NewSet(intSetHash, interfaceSlice([]int{1234, 5678})))
	d.Set("send_to_email", true)
	d.Set("send_to_sms", true)
	d.Set("send_to_twitter", true)
	d.Set("send_to_iphone", true)
	d.Set("send_to_android", true)
	d.Set("send_notification_when_down", 2)
	d.Set("notify_again_every", 1)
	d.Set("notify_when_back_up", true)
	d.Set("tags", schema.NewSet(schema.HashString, interfaceSlice([]string{"bar", "foo"})))
	d.Set("ipv6", false)

	return d
}

func checkConfigurationData() checks.CheckConfiguration {
	return checks.CheckConfiguration{
		Name:                     "My check",
		Host:                     "example.com",
		Resolution:               1,
		ContactIDs:               []int{1234, 5678},
		SendToEmail:              true,
		SendToSMS:                true,
		SendToTwitter:            true,
		SendToIphone:             true,
		SendToAndroid:            true,
		SendNotificationWhenDown: 2,
		NotifyAgainEvery:         1,
		NotifyWhenBackUp:         true,
		Tags:                     []string{"bar", "foo"},
		IPv6:                     false,
	}
}

func getDetailedCheckEntryData() checks.DetailedCheckEntry {
	return checks.DetailedCheckEntry{
		ID:                       85975,
		Name:                     "My check",
		Hostname:                 "example.com",
		Resolution:               1,
		ContactIDs:               []int{1234, 5678},
		SendToEmail:              true,
		SendToSMS:                true,
		SendToTwitter:            true,
		SendToIphone:             true,
		SendToAndroid:            true,
		SendNotificationWhenDown: 2,
		NotifyAgainEvery:         1,
		NotifyWhenBackUp:         true,
		IPv6:                     false,
		Status:                   "up",
		LastErrorTime:            1293143467,
		LastTestTime:             1294064823,
		LastResponseTime:         1294064824,
		Created:                  1240394682,
	}
}

func flattenBaseCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"name":                        "My check",
		"host":                        "example.com",
		"resolution":                  1,
		"contact_ids":                 []int{1234, 5678},
		"send_to_email":               true,
		"send_to_sms":                 true,
		"send_to_twitter":             true,
		"send_to_iphone":              true,
		"send_to_android":             true,
		"send_notification_when_down": 2,
		"notify_again_every":          1,
		"notify_when_back_up":         true,
		"ipv6":                        false,
		"last_error_time":             1293143467,
		"last_test_time":              1294064823,
		"last_response_time":          1294064824,
		"created":                     1240394682,
		"status":                      "up",
	}
}

func getDetailedCheckEntryTagsData() []checks.CheckListEntryTags {
	return []checks.CheckListEntryTags{
		checks.CheckListEntryTags{
			Name:  "apache",
			Type:  "a",
			Count: 2,
		},
		checks.CheckListEntryTags{
			Name:  "nginx",
			Type:  "u",
			Count: 1,
		},
	}
}

func flattenBaseCheckTagsExpected() map[string]interface{} {
	return map[string]interface{}{
		"tags":      []string{"nginx"},
		"auto_tags": []string{"apache"},
	}
}

func httpCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: httpCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("url", "/test")
	d.Set("encryption", true)
	d.Set("port", 443)
	d.Set("auth", "foo:bar")
	d.Set("should_contain", "foo")
	d.Set("should_not_contain", "bar")
	d.Set("post_data", "baz")
	d.Set("request_headers", schema.NewSet(schema.HashString, interfaceSlice([]string{"X-Header1:foo", "X-Header2:bar", "X-Header3:baz"})))

	return d
}

func checkConfigurationHTTPData() checks.CheckConfigurationHTTP {
	return checks.CheckConfigurationHTTP{
		URL:              "/test",
		Encryption:       true,
		Port:             443,
		Auth:             "foo:bar",
		ShouldContain:    "foo",
		ShouldNotContain: "bar",
		PostData:         "baz",
		RequestHeaders:   []string{"X-Header1:foo", "X-Header2:bar", "X-Header3:baz"},
	}
}

func getDetailedCheckEntryHTTPData() checks.DetailedCheckEntryHTTP {
	return checks.DetailedCheckEntryHTTP{
		URL:              "/test",
		Encryption:       true,
		Port:             443,
		Username:         "foo",
		Password:         "bar",
		ShouldContain:    "foo",
		ShouldNotContain: "bar",
		PostData:         "baz",
		RequestHeaders: map[string]string{
			"X-Header1": "foo",
			"X-Header2": "bar",
			"X-Header3": "baz",
		},
	}
}

func flattenHTTPCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"url":                "/test",
		"encryption":         true,
		"port":               443,
		"auth":               "foo:bar",
		"should_contain":     "foo",
		"should_not_contain": "bar",
		"post_data":          "baz",
		"request_headers":    []string{"X-Header1:foo", "X-Header2:bar", "X-Header3:baz"},
	}
}

func httpCustomCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: httpCustomCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("url", "/test")
	d.Set("encryption", true)
	d.Set("port", 443)
	d.Set("auth", "foo:bar")
	d.Set("additional_urls", schema.NewSet(schema.HashString, interfaceSlice([]string{"www.mysite.com", "www.myothersite.com"})))

	return d
}

func checkConfigurationHTTPCustomData() checks.CheckConfigurationHTTPCustom {
	return checks.CheckConfigurationHTTPCustom{
		URL:            "/test",
		Encryption:     true,
		Port:           443,
		Auth:           "foo:bar",
		AdditionalURLs: []string{"www.mysite.com", "www.myothersite.com"},
	}
}

func getDetailedCheckEntryHTTPCustomData() checks.DetailedCheckEntryHTTPCustom {
	return checks.DetailedCheckEntryHTTPCustom{
		URL:            "/test",
		Encryption:     true,
		Port:           443,
		Username:       "foo",
		Password:       "bar",
		AdditionalURLs: []string{"www.mysite.com", "www.myothersite.com"},
	}
}

func flattenHTTPCustomCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"url":             "/test",
		"encryption":      true,
		"port":            443,
		"auth":            "foo:bar",
		"additional_urls": []string{"www.mysite.com", "www.myothersite.com"},
	}
}

func tcpCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: tcpCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("port", 22)
	d.Set("string_to_send", "foo")
	d.Set("string_to_expect", "bar")

	return d
}

func checkConfigurationTCPData() checks.CheckConfigurationTCP {
	return checks.CheckConfigurationTCP{
		Port:           22,
		StringToSend:   "foo",
		StringToExpect: "bar",
	}
}

func getDetailedCheckEntryTCPData() checks.DetailedCheckEntryTCP {
	return checks.DetailedCheckEntryTCP{
		Port:           22,
		StringToSend:   "foo",
		StringToExpect: "bar",
	}
}

func flattenTCPCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"port":             22,
		"string_to_send":   "foo",
		"string_to_expect": "bar",
	}
}

func dnsCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: dnsCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("nameserver", "ns1.example.com")
	d.Set("expected_ip", "127.0.0.1")

	return d
}

func checkConfigurationDNSData() checks.CheckConfigurationDNS {
	return checks.CheckConfigurationDNS{
		NameServer: "ns1.example.com",
		ExpectedIP: "127.0.0.1",
	}
}

func getDetailedCheckEntryDNSData() checks.DetailedCheckEntryDNS {
	return checks.DetailedCheckEntryDNS{
		DNSServer:  "ns1.example.com",
		ExpectedIP: "127.0.0.1",
	}
}

func flattenDNSCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"nameserver":  "ns1.example.com",
		"expected_ip": "127.0.0.1",
	}
}

func udpCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: udpCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("port", 53)
	d.Set("string_to_send", "foo")
	d.Set("string_to_expect", "bar")

	return d
}

func checkConfigurationUDPData() checks.CheckConfigurationUDP {
	return checks.CheckConfigurationUDP{
		Port:           53,
		StringToSend:   "foo",
		StringToExpect: "bar",
	}
}

func getDetailedCheckEntryUDPData() checks.DetailedCheckEntryUDP {
	return checks.DetailedCheckEntryUDP{
		Port:           53,
		StringToSend:   "foo",
		StringToExpect: "bar",
	}
}

func flattenUDPCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"port":             53,
		"string_to_send":   "foo",
		"string_to_expect": "bar",
	}
}

func smtpCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: smtpCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("port", 587)
	d.Set("auth", "foo:bar")
	d.Set("encryption", true)
	d.Set("string_to_expect", "foobar")

	return d
}

func checkConfigurationSMTPData() checks.CheckConfigurationSMTP {
	return checks.CheckConfigurationSMTP{
		Port:           587,
		Auth:           "foo:bar",
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func getDetailedCheckEntrySMTPData() checks.DetailedCheckEntrySMTP {
	return checks.DetailedCheckEntrySMTP{
		Port:           587,
		Username:       "foo",
		Password:       "bar",
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func flattenSMTPCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"port":             587,
		"auth":             "foo:bar",
		"encryption":       true,
		"string_to_expect": "foobar",
	}
}

func pop3CheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: pop3CheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("port", 993)
	d.Set("encryption", true)
	d.Set("string_to_expect", "foobar")

	return d
}

func checkConfigurationPOP3Data() checks.CheckConfigurationPOP3 {
	return checks.CheckConfigurationPOP3{
		Port:           993,
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func getDetailedCheckEntryPOP3Data() checks.DetailedCheckEntryPOP3 {
	return checks.DetailedCheckEntryPOP3{
		Port:           993,
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func flattenPOP3CheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"port":             993,
		"encryption":       true,
		"string_to_expect": "foobar",
	}
}

func imapCheckResourceData() *schema.ResourceData {
	r := &schema.Resource{
		Schema: imapCheckSchema(),
	}
	d := r.TestResourceData()

	d.Set("port", 995)
	d.Set("encryption", true)
	d.Set("string_to_expect", "foobar")

	return d
}

func checkConfigurationIMAPData() checks.CheckConfigurationIMAP {
	return checks.CheckConfigurationIMAP{
		Port:           995,
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func getDetailedCheckEntryIMAPData() checks.DetailedCheckEntryIMAP {
	return checks.DetailedCheckEntryIMAP{
		Port:           995,
		Encryption:     true,
		StringToExpect: "foobar",
	}
}

func flattenIMAPCheckExpected() map[string]interface{} {
	return map[string]interface{}{
		"port":             995,
		"encryption":       true,
		"string_to_expect": "foobar",
	}
}

func TestInterfaceSlice_string(t *testing.T) {
	in := []string{"one", "two"}
	expected := []interface{}{"one", "two"}
	out := interfaceSlice(in)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestInterfaceSlice_int(t *testing.T) {
	in := []int{1, 2}
	expected := []interface{}{1, 2}
	out := interfaceSlice(in)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestInterfaceSlice_unsupported(t *testing.T) {
	failed := false
	in := []bool{true, false}
	defer func() {
		if r := recover(); r != nil {
			failed = true
		}
	}()
	_ = interfaceSlice(in)

	if failed == false {
		t.Fatal("Expected a panic, but did not get one")
	}
}

func TestIntSlice(t *testing.T) {
	in := []interface{}{1, 2}
	expected := []int{1, 2}
	out := intSlice(in)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestStringSlice(t *testing.T) {
	in := []interface{}{"one", "two"}
	expected := []string{"one", "two"}
	out := stringSlice(in)

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestIntSetHash(t *testing.T) {
	in := 1
	expected := 1
	out := intSetHash(in)

	if expected != out {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestHTTPCheckSchemaFull(t *testing.T) {
	m := httpCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["request_headers"]; ok == false {
		t.Fatal("Expected request_headers to be present")
	}
}

func TestHTTPCustomCheckSchemaFull(t *testing.T) {
	m := httpCustomCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["additional_urls"]; ok == false {
		t.Fatal("Expected additional_urls to be present")
	}
}

func TestTCPCheckSchemaFull(t *testing.T) {
	m := tcpCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["port"]; ok == false {
		t.Fatal("Expected port to be present")
	}
}

func TestDNSCheckSchemaFull(t *testing.T) {
	m := dnsCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["nameserver"]; ok == false {
		t.Fatal("Expected nameserver to be present")
	}
}

func TestUDPCheckSchemaFull(t *testing.T) {
	m := udpCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["port"]; ok == false {
		t.Fatal("Expected port to be present")
	}
}

func TestSMTPCheckSchemaFull(t *testing.T) {
	m := smtpCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["auth"]; ok == false {
		t.Fatal("Expected auth to be present")
	}
}

func TestPOP3CustomCheckSchemaFull(t *testing.T) {
	m := pop3CheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["encryption"]; ok == false {
		t.Fatal("Expected encryption to be present")
	}
}

func TestIMAPCheckSchemaFull(t *testing.T) {
	m := imapCheckSchemaFull()
	if _, ok := m["name"]; ok == false {
		t.Fatal("Expected name to be present")
	}
	if _, ok := m["string_to_expect"]; ok == false {
		t.Fatal("Expected string_to_expect to be present")
	}
}

func TestExapndBaseCheck(t *testing.T) {
	in := baseCheckResourceData()
	out := expandBaseCheck(in)
	expected := checkConfigurationData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenBaseCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: baseCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryData()
	flattenBaseCheck(in, out)

	if out.Id() != strconv.Itoa(85975) {
		t.Fatalf("expected ID to be %s, got %s", strconv.Itoa(85975), out.Id())
	}

	for k, v := range flattenBaseCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestFlattenBaseCheckTags(t *testing.T) {
	r := &schema.Resource{
		Schema: baseCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryTagsData()
	flattenBaseCheckTags(in, out)

	for k, v := range flattenBaseCheckTagsExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndHTTPCheck(t *testing.T) {
	in := httpCheckResourceData()
	out := expandHTTPCheck(in)
	expected := checkConfigurationHTTPData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenHTTPCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: httpCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryHTTPData()
	flattenHTTPCheck(in, out)

	for k, v := range flattenHTTPCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndHTTPCustomCheck(t *testing.T) {
	in := httpCustomCheckResourceData()
	out := expandHTTPCustomCheck(in)
	expected := checkConfigurationHTTPCustomData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenHTTPCustomCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: httpCustomCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryHTTPCustomData()
	flattenHTTPCustomCheck(in, out)

	for k, v := range flattenHTTPCustomCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndTCPCheck(t *testing.T) {
	in := tcpCheckResourceData()
	out := expandTCPCheck(in)
	expected := checkConfigurationTCPData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenTCPCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: tcpCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryTCPData()
	flattenTCPCheck(in, out)

	for k, v := range flattenTCPCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndDNSCheck(t *testing.T) {
	in := dnsCheckResourceData()
	out := expandDNSCheck(in)
	expected := checkConfigurationDNSData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenDNSCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: dnsCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryDNSData()
	flattenDNSCheck(in, out)

	for k, v := range flattenDNSCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndUDPCheck(t *testing.T) {
	in := udpCheckResourceData()
	out := expandUDPCheck(in)
	expected := checkConfigurationUDPData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenUDPCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: udpCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryUDPData()
	flattenUDPCheck(in, out)

	for k, v := range flattenUDPCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndSMTPCheck(t *testing.T) {
	in := smtpCheckResourceData()
	out := expandSMTPCheck(in)
	expected := checkConfigurationSMTPData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenSMTPCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: smtpCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntrySMTPData()
	flattenSMTPCheck(in, out)

	for k, v := range flattenSMTPCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndPOP3Check(t *testing.T) {
	in := pop3CheckResourceData()
	out := expandPOP3Check(in)
	expected := checkConfigurationPOP3Data()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenPOP3Check(t *testing.T) {
	r := &schema.Resource{
		Schema: pop3CheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryPOP3Data()
	flattenPOP3Check(in, out)

	for k, v := range flattenPOP3CheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}

func TestExapndIMAPCheck(t *testing.T) {
	in := imapCheckResourceData()
	out := expandIMAPCheck(in)
	expected := checkConfigurationIMAPData()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %+v, got %+v", expected, out)
	}
}

func TestFlattenIMAPCheck(t *testing.T) {
	r := &schema.Resource{
		Schema: imapCheckSchema(),
	}
	out := r.TestResourceData()

	in := getDetailedCheckEntryIMAPData()
	flattenIMAPCheck(in, out)

	for k, v := range flattenIMAPCheckExpected() {
		switch w := out.Get(k).(type) {
		case *schema.Set:
			if reflect.DeepEqual(w.List(), interfaceSlice(v)) == false {
				t.Fatalf("expected %s to be %+v, got %+v", k, v, w.List())
			}
		default:
			if w != v {
				t.Fatalf("expected %s to be %v, got %v", k, v, w)
			}
		}
	}
}
