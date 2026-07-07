package receiver

import "testing"

func TestParseServers_RackUnitsBeforeBlades(t *testing.T) {
	input := []byte(`
		<configResolveClass>
			<outConfigs>
				<computeRackUnit dn="sys/rack-unit-1" model="UCSC-C220-M5" serial="RACK123" operState="ok"/>
				<computeBlade dn="sys/chassis-1/blade-1" model="UCSB-B200-M5" serial="BLADE456" operState="ok"/>
			</outConfigs>
		</configResolveClass>
	`)

	got, err := ParseServers(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Server{
		{DN: "sys/rack-unit-1", Model: "UCSC-C220-M5", Serial: "RACK123", OperState: "ok"},
		{DN: "sys/chassis-1/blade-1", Model: "UCSB-B200-M5", Serial: "BLADE456", OperState: "ok"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d servers, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("server %d = %+v, want %+v (rack units must come before blades)", i, got[i], w)
		}
	}
}

func TestParseServers_OnlyBlades(t *testing.T) {
	input := []byte(`
		<configResolveClass>
			<outConfigs>
				<computeBlade dn="sys/chassis-1/blade-1" model="UCSB-B200-M5" serial="BLADE111" operState="ok"/>
				<computeBlade dn="sys/chassis-1/blade-2" model="UCSB-B200-M5" serial="BLADE222" operState="discovering"/>
			</outConfigs>
		</configResolveClass>
	`)

	got, err := ParseServers(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d servers, want 2", len(got))
	}
	if got[1].Serial != "BLADE222" || got[1].OperState != "discovering" {
		t.Errorf("got %+v, want serial=BLADE222 operState=discovering", got[1])
	}
}

func TestParseServers_Empty(t *testing.T) {
	input := []byte(`
		<configResolveClass>
			<outConfigs></outConfigs>
		</configResolveClass>
	`)

	got, err := ParseServers(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d servers, want 0", len(got))
	}
}

func TestParseServers_InvalidXML(t *testing.T) {
	_, err := ParseServers([]byte(`not xml`))
	if err == nil {
		t.Fatal("expected an error for malformed XML, got nil")
	}
}
