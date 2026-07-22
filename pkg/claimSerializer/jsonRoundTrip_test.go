package claimserializer

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

// serializeConcrete serializes a deserialized claim of any supported concrete type.
func serializeConcrete(i any) (string, error) {
	switch item := i.(type) {
	case request.Billing:
		return Serialize(&item)
	case request.Reversal:
		return Serialize(&item)
	case request.Rebill:
		return Serialize(&item)
	case request.Eligibility:
		return Serialize(&item)
	case response.Billing:
		return Serialize(&item)
	case response.Reversal:
		return Serialize(&item)
	case response.Rebill:
		return Serialize(&item)
	case response.Eligibility:
		return Serialize(&item)
	default:
		return "", fmt.Errorf("unknown type %T", i)
	}
}

// unmarshalAndSerialize unmarshals a claim's JSON into a fresh struct of the
// given concrete type and serializes the result.
func unmarshalAndSerialize[V any](jsonBytes []byte) (string, error) {
	fresh := new(V)
	if err := json.Unmarshal(jsonBytes, fresh); err != nil {
		return "", err
	}
	return Serialize(fresh)
}

// jsonRoundTripSerialize marshals a deserialized claim to JSON, unmarshals it
// into a fresh struct of the same concrete type, and serializes the result.
func jsonRoundTripSerialize(i any) (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	switch i.(type) {
	case request.Billing:
		return unmarshalAndSerialize[request.Billing](jsonBytes)
	case request.Reversal:
		return unmarshalAndSerialize[request.Reversal](jsonBytes)
	case request.Rebill:
		return unmarshalAndSerialize[request.Rebill](jsonBytes)
	case request.Eligibility:
		return unmarshalAndSerialize[request.Eligibility](jsonBytes)
	case response.Billing:
		return unmarshalAndSerialize[response.Billing](jsonBytes)
	case response.Reversal:
		return unmarshalAndSerialize[response.Reversal](jsonBytes)
	case response.Rebill:
		return unmarshalAndSerialize[response.Rebill](jsonBytes)
	case response.Eligibility:
		return unmarshalAndSerialize[response.Eligibility](jsonBytes)
	default:
		return "", fmt.Errorf("unknown type %T", i)
	}
}

// A claim that went through a JSON round-trip must serialize identically to
// the original deserialized claim. Dynamic segments and fields lose their
// generated reflect.Type in JSON (DynamicType is json:"-"); DynamicStruct's
// UnmarshalJSON reconstructs it from the encoded field names.
func Test_CanSerializeAfterJsonRoundTrip(t *testing.T) {
	tests := append(append([]serializerTest{}, dynamicTests...), extraSegmentsRequestTests...)

	for _, test := range tests {
		i, err := claimdeserializer.Deserialize(test.rawData)
		if err != nil {
			t.Error(err)
			continue
		}

		direct, err := serializeConcrete(i)
		if err != nil {
			t.Error(err)
			continue
		}

		roundTripped, err := jsonRoundTripSerialize(i)
		if err != nil {
			t.Error(err)
			continue
		}

		if roundTripped != direct {
			t.Errorf("JSON round-trip serialization mismatch.\nWanted: %q\nGot:    %q", direct, roundTripped)
		}
	}
}

// A single other payer with multiple rejects must re-serialize with the
// original payer count (4C) and reject count (5E) instead of being split
// into multiple payer occurrences.
func Test_SerializePreservesNestedRepeatingGroups(t *testing.T) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(REQUEST_B1, &obj); err != nil {
		t.Fatal(err)
	}

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	fs := string(ncpdp.FIELD)

	for _, want := range []string{
		fs + "4C1" + fs,
		fs + "5E2" + fs + "6E70 " + fs + "6EA5 ",
	} {
		if !strings.Contains(serialized, want) {
			t.Errorf("serialized COB segment missing %q", want)
		}
	}

	assertFieldCount(t, serialized, "5C", 1)
	assertFieldCount(t, serialized, "5E", 1)
}
