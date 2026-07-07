package config_load

import (
	"strings"
	"testing"
)

func TestOnlyDeviceName_SingleMatch(t *testing.T) {
	devices := map[string]Device{
		"switch1": {Type: "mds"},
		"ucsm1":   {Type: "ucsm"},
	}

	got, err := onlyDeviceName(devices, "mds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "switch1" {
		t.Errorf("got %q, want %q", got, "switch1")
	}
}

func TestOnlyDeviceName_CaseInsensitiveType(t *testing.T) {
	devices := map[string]Device{
		"switch1": {Type: "MDS"},
	}

	got, err := onlyDeviceName(devices, "mds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "switch1" {
		t.Errorf("got %q, want %q", got, "switch1")
	}
}

func TestOnlyDeviceName_NoMatch(t *testing.T) {
	devices := map[string]Device{
		"ucsm1": {Type: "ucsm"},
	}

	_, err := onlyDeviceName(devices, "mds")
	if err == nil {
		t.Fatal("expected an error when no device matches the type, got nil")
	}
}

func TestOnlyDeviceName_MultipleMatches(t *testing.T) {
	devices := map[string]Device{
		"switch-b": {Type: "mds"},
		"switch-a": {Type: "mds"},
	}

	_, err := onlyDeviceName(devices, "mds")
	if err == nil {
		t.Fatal("expected an error when multiple devices match the type, got nil")
	}
	// names should be sorted in the error message so the user gets a stable, readable list
	if !strings.Contains(err.Error(), "switch-a, switch-b") {
		t.Errorf("error %q does not contain sorted device names %q", err.Error(), "switch-a, switch-b")
	}
}

func TestValidateDevice_AllFieldsPresent(t *testing.T) {
	device := Device{
		Type:     "mds",
		Host:     "10.0.0.1",
		Port:     "443",
		Username: "admin",
		Password: "secret",
	}

	if err := validateDevice("switch1", device); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDevice_MissingFields(t *testing.T) {
	device := Device{
		Type: "mds",
		Host: "10.0.0.1",
		// Port, Username, Password left empty
	}

	err := validateDevice("switch1", device)
	if err == nil {
		t.Fatal("expected an error for missing required fields, got nil")
	}
	if !strings.Contains(err.Error(), "port, username, password") {
		t.Errorf("error %q does not list missing fields in expected order", err.Error())
	}
}

func TestValidateDevice_MissingSingleField(t *testing.T) {
	device := Device{
		Type:     "mds",
		Host:     "10.0.0.1",
		Port:     "443",
		Username: "admin",
		// Password left empty
	}

	err := validateDevice("switch1", device)
	if err == nil {
		t.Fatal("expected an error for missing password, got nil")
	}
	if !strings.Contains(err.Error(), "password") {
		t.Errorf("error %q does not mention the missing password field", err.Error())
	}
}
