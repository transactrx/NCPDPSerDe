package ncpdp_test

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
	"github.com/transactrx/NCPDPSerDe/pkg/serde"
)

// The JSON schema (schemas/ncpdp-schemas.json) marks every F6-only field with
// a description ending in " - F6". The Go structs mark the same fields with a
// sinceVersion=F6 field tag: code-bearing leaves carry it inside their tag,
// and repeating-group slices flagged as F6 composites carry a code-less
// field:"sinceVersion=F6" tag. This test keeps the two in sync so a field
// added to one side without the other fails CI instead of silently drifting.

var schemaDefTypes = map[string]reflect.Type{
	"request.Insurance":                           reflect.TypeOf(requestsegment.Insurance{}),
	"request.Patient":                             reflect.TypeOf(requestsegment.Patient{}),
	"request.Claim":                               reflect.TypeOf(requestsegment.Claim{}),
	"request.Pricing":                             reflect.TypeOf(requestsegment.Pricing{}),
	"request.PharmacyProvider":                    reflect.TypeOf(requestsegment.PharmacyProvider{}),
	"request.Prescriber":                          reflect.TypeOf(requestsegment.Prescriber{}),
	"request.CoordinationOfBenefits":              reflect.TypeOf(requestsegment.CoordinationOfBenefits{}),
	"request.WorkersCompensation":                 reflect.TypeOf(requestsegment.WorkersCompensation{}),
	"request.Dur":                                 reflect.TypeOf(requestsegment.Dur{}),
	"request.Coupon":                              reflect.TypeOf(requestsegment.Coupon{}),
	"request.Compound":                            reflect.TypeOf(requestsegment.Compound{}),
	"request.Clinical":                            reflect.TypeOf(requestsegment.Clinical{}),
	"request.AdditionalDocumentation":             reflect.TypeOf(requestsegment.AdditionalDocumentation{}),
	"request.Facility":                            reflect.TypeOf(requestsegment.Facility{}),
	"request.Narrative":                           reflect.TypeOf(requestsegment.Narrative{}),
	"request.PriorAuthorizationRequestAndBilling": reflect.TypeOf(requestsegment.PriorAuthorizationRequestAndBilling{}),
	"request.Intermediary":                        reflect.TypeOf(requestsegment.Intermediary{}),
	"request.LastKnown4Rx":                        reflect.TypeOf(requestsegment.LastKnown4Rx{}),
	"request.NTransactionPayerIdentification":     reflect.TypeOf(requestsegment.NTransactionPayerIdentification{}),
	"request.DataCollection":                      reflect.TypeOf(requestsegment.DataCollection{}),
	"response.Message":                            reflect.TypeOf(responsesegment.Message{}),
	"response.Insurance":                          reflect.TypeOf(responsesegment.Insurance{}),
	"response.InsuranceAdditionalInformation":     reflect.TypeOf(responsesegment.InsuranceAdditionalInformation{}),
	"response.Patient":                            reflect.TypeOf(responsesegment.Patient{}),
	"response.Status":                             reflect.TypeOf(responsesegment.Status{}),
	"response.Claim":                              reflect.TypeOf(responsesegment.Claim{}),
	"response.Pricing":                            reflect.TypeOf(responsesegment.Pricing{}),
	"response.Dur":                                reflect.TypeOf(responsesegment.Dur{}),
	"response.PriorAuthorization":                 reflect.TypeOf(responsesegment.PriorAuthorization{}),
	"response.CoordinationOfBenefits":             reflect.TypeOf(responsesegment.CoordinationOfBenefits{}),
	"response.Intermediary":                       reflect.TypeOf(responsesegment.Intermediary{}),
	"response.OtherRelatedBenefitDetail":          reflect.TypeOf(responsesegment.OtherRelatedBenefitDetail{}),
	"response.Provider":                           reflect.TypeOf(responsesegment.Provider{}),
}

var fieldCodePattern = regexp.MustCompile(`\(([A-Z0-9]{1,3})\)`)

// F6 expectations extracted from one schema definition.
type schemaF6Facts struct {
	// leafF6 maps a field code to whether the schema flags it as F6-only.
	// Codes flagged inconsistently within one definition are recorded in
	// ambiguous and excluded from comparison.
	leafF6     map[string]bool
	ambiguous  map[string]bool
	composites map[string]bool // code-less object/array property names flagged F6
}

func collectSchemaFacts(t *testing.T, properties map[string]any, facts *schemaF6Facts) {
	t.Helper()

	for name, raw := range properties {
		prop, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		desc, _ := prop["description"].(string)
		flagged := strings.HasSuffix(strings.TrimSpace(desc), "- F6")

		code := ""
		if m := fieldCodePattern.FindStringSubmatch(desc); m != nil {
			code = m[1]
		}

		nested, _ := prop["properties"].(map[string]any)
		if items, ok := prop["items"].(map[string]any); ok && nested == nil {
			nested, _ = items["properties"].(map[string]any)
		}

		switch {
		case nested != nil:
			if flagged && code == "" {
				facts.composites[name] = true
			}
			collectSchemaFacts(t, nested, facts)
		case code != "":
			if existing, seen := facts.leafF6[code]; seen && existing != flagged {
				facts.ambiguous[code] = true
			}
			facts.leafF6[code] = flagged
		}
	}
}

// Codes with structural meaning that the schema does not describe as fields.
var skipCodes = map[string]bool{"dynamic": true, "AM": true, "rawgroup": true, "rawsegment": true}

func checkStructTags(t *testing.T, defName string, structType reflect.Type, facts *schemaF6Facts, visited map[reflect.Type]bool) {
	t.Helper()

	if visited[structType] || structType.String() == "time.Time" {
		return
	}
	visited[structType] = true

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if tag, found := field.Tag.Lookup(serde.FieldTag); found {
			attr, err := serde.GetFieldAttribute(tag)
			if err != nil {
				t.Errorf("%s.%s: bad field tag: %v", defName, field.Name, err)
				continue
			}

			tagged := attr.SinceVersion != ""
			if tagged && !ncpdp.KnownVersion(attr.SinceVersion) {
				t.Errorf("%s.%s: sinceVersion=%q is not in ncpdp's version rank table — the serializer would omit it from every transmission",
					defName, field.Name, attr.SinceVersion)
			}
			switch {
			case skipCodes[attr.Code]:
				// structural, not schema-described
			case attr.Code != "":
				want, known := facts.leafF6[attr.Code]
				if known && !facts.ambiguous[attr.Code] && want != tagged {
					t.Errorf("%s.%s (code %s): schema F6=%v but sinceVersion tag present=%v",
						defName, field.Name, attr.Code, want, tagged)
				}
			case tagged:
				if !facts.composites[field.Name] {
					t.Errorf("%s.%s: code-less sinceVersion tag but schema does not flag it as an F6 composite",
						defName, field.Name)
				}
			}
		}

		elem := field.Type
		for elem.Kind() == reflect.Pointer || elem.Kind() == reflect.Slice {
			elem = elem.Elem()
		}
		if elem.Kind() == reflect.Struct {
			checkStructTags(t, defName, elem, facts, visited)
		}
	}
}

func TestSinceVersionTagsMatchSchema(t *testing.T) {
	raw, err := os.ReadFile("schemas/ncpdp-schemas.json")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	var doc struct {
		Defs map[string]struct {
			Properties map[string]any `json:"properties"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	for defName, structType := range schemaDefTypes {
		def, found := doc.Defs[defName]
		if !found {
			t.Errorf("schema definition %s not found", defName)
			continue
		}

		facts := &schemaF6Facts{
			leafF6:     map[string]bool{},
			ambiguous:  map[string]bool{},
			composites: map[string]bool{},
		}
		collectSchemaFacts(t, def.Properties, facts)
		checkStructTags(t, defName, structType, facts, map[reflect.Type]bool{})

		// Reverse direction: every schema-flagged composite must exist as a
		// top-level slice/struct field carrying a code-less sinceVersion tag.
		for name := range facts.composites {
			field, found := structType.FieldByName(name)
			if !found {
				t.Errorf("%s: schema flags composite %s as F6 but the struct has no such field", defName, name)
				continue
			}
			attr, err := serde.GetFieldAttribute(field.Tag.Get(serde.FieldTag))
			if err != nil || attr.SinceVersion == "" || attr.Code != "" {
				t.Errorf("%s.%s: schema flags it as an F6 composite but it lacks a code-less sinceVersion tag", defName, name)
			}
		}
	}
}
