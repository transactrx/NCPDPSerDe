package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
	"github.com/transactrx/NCPDPSerDe/pkg/serde"
)

// The version gate applies only to code-bearing fields: a code-less
// sinceVersion tag (repeating-group slices like Patient.Ids) is metadata for
// reflection consumers and must never suppress serialization — dual-mapped
// slices carry the D0 scalar's data. Version comparison is rank-based
// (ncpdp.VersionAtLeast), so an F6-scoped field is kept for every version at
// or after F6, an unknown transmission version ranks oldest, and a field
// scoped to a version this build does not know is always omitted.
func Test_FieldOmittedForVersion(t *testing.T) {
	leaf := &serde.FieldAttribute{Code: "S8", SinceVersion: "F6"}
	composite := &serde.FieldAttribute{SinceVersion: "F6"}
	unscoped := &serde.FieldAttribute{Code: "CX"}
	futureScoped := &serde.FieldAttribute{Code: "XX", SinceVersion: "ZZ"}

	cases := []struct {
		name    string
		attr    *serde.FieldAttribute
		version string
		want    bool
	}{
		{"F6-only leaf omitted from D0", leaf, ncpdp.D0, true},
		{"F6-only leaf kept for F6", leaf, ncpdp.F6, false},
		{"F6-only leaf omitted for unknown transmission version", leaf, "ZZ", true},
		{"F6-only leaf omitted for blank transmission version", leaf, "", true},
		{"leaf scoped to an unmodeled version always omitted", futureScoped, ncpdp.F6, true},
		{"code-less composite tag never gates D0", composite, ncpdp.D0, false},
		{"code-less composite tag never gates F6", composite, ncpdp.F6, false},
		{"unscoped leaf kept for D0", unscoped, ncpdp.D0, false},
		{"nil attribute kept", nil, ncpdp.D0, false},
	}

	for _, c := range cases {
		if got := fieldOmittedForVersion(c.attr, c.version); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// Request-side counterpart of the response-path leaf-gate tests: an F6-only
// field (Species, S8) set by application code must be omitted from a D0
// request.
func Test_D0RequestOmitsHandPopulatedF6Leaf(t *testing.T) {
	i, err := claimdeserializer.Deserialize(REQUEST_B1)
	if err != nil {
		t.Fatal(err)
	}
	item, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing.  Got: %T", i)
	}

	species := "01"
	item.Patient.Species = &species

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(serialized, string(ncpdp.FIELD)+"S8") {
		t.Errorf("D0 request must not contain the F6-only Species (S8) field.\nGot: %q", serialized)
	}
}

// The same field is emitted when the transmission is F6.
func Test_F6RequestKeepsHandPopulatedF6Leaf(t *testing.T) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &obj); err != nil {
		t.Fatal(err)
	}

	species := "01"
	obj.Patient.Species = &species

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(serialized, string(ncpdp.FIELD)+"S801") {
		t.Errorf("F6 request must contain the Species (S8) field.\nGot: %q", serialized)
	}
}

// Framing trim on the response path: leading whitespace and trailing CR/LF are
// stripped before the transaction code is read from fixed byte offsets.
func Test_ResponseFramingTrimmed(t *testing.T) {
	i, err := claimdeserializer.DeserializeResponse("\r\n  " + RESPONSE_B1 + "\r\n")
	if err != nil {
		t.Fatal(err)
	}
	item, ok := i.(response.Billing)
	if !ok {
		t.Fatalf("expected response.Billing.  Got: %T", i)
	}

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(serialized, RESPONSE_B1[:20]) {
		t.Errorf("response header not parsed from position zero.\nGot: %q", serialized[:30])
	}
	if strings.ContainsAny(serialized, "\r\n") {
		t.Errorf("CR/LF framing leaked into serialized output: %q", serialized)
	}
}

// Framing trim on an F6 transmission: version detection reads byte zero, so a
// framed F6 message must still be recognized as a request.
func Test_F6FramingTrimmed(t *testing.T) {
	i, err := claimdeserializer.Deserialize("\n" + buildF6BillingRequest() + "\r\n")
	if err != nil {
		t.Fatal(err)
	}
	item, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing.  Got: %T", i)
	}

	if item.Header.Value.Version != ncpdp.F6 {
		t.Errorf("Version mismatch. Wanted: %q   Got: %q", ncpdp.F6, item.Header.Value.Version)
	}
}

// Framing trim on the DeserializeType path (deserializeRaw), which is not
// covered by the Deserialize/DeserializeRequest entry points.
func Test_DeserializeTypeFramingTrimmed(t *testing.T) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType("  "+REQUEST_B1+"\r\n", &obj); err != nil {
		t.Fatal(err)
	}

	if obj.Header.Value.Bin != "880151" {
		t.Errorf("BIN mismatch. Wanted: %q   Got: %q", "880151", obj.Header.Value.Bin)
	}
}
