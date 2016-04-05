package pingdom

import (
	"reflect"
	"testing"
)

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

func TestInterfaceSlice_map(t *testing.T) {
	in := map[string]string{
		"a": "one",
		"b": "two",
	}
	expected := []interface{}{"one", "two"}
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
