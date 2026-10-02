package claimserializer

import (
	"strings"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response"
)

func intermediaryWithOneId() requestsegment.Intermediary {
	typeCode := "01"
	id := "INTID"
	return requestsegment.Intermediary{
		Ids: []requestsegment.IntermediaryId{{TypeCode: &typeCode, Id: &id}},
	}
}

// An F6-only segment (Intermediary, AM19) hand-populated on a D0 request must
// vanish entirely: every field is version-gated and buildSegment drops a
// segment that renders as a bare header.
func Test_D0RequestOmitsF6OnlySegment(t *testing.T) {
	i, err := claimdeserializer.Deserialize(REQUEST_B1)
	if err != nil {
		t.Fatal(err)
	}
	item, ok := i.(request.Billing)
	if !ok {
		t.Fatalf("expected request.Billing.  Got: %T", i)
	}

	item.Claims[0].Intermediary = intermediaryWithOneId()

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(serialized, "AM19") || strings.Contains(serialized, string(ncpdp.FIELD)+"8M") {
		t.Errorf("D0 request must not contain the F6-only Intermediary segment (AM19).\nGot: %q", serialized)
	}
}

// The same segment is emitted, counter included, when the transmission is F6.
func Test_F6RequestKeepsF6OnlySegment(t *testing.T) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &obj); err != nil {
		t.Fatal(err)
	}

	obj.Claims[0].Intermediary = intermediaryWithOneId()

	serialized, err := Serialize(&obj)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"AM19", string(ncpdp.FIELD) + "8G1", string(ncpdp.FIELD) + "8MINTID"} {
		if !strings.Contains(serialized, want) {
			t.Errorf("F6 request must contain %q (Intermediary segment).\nGot: %q", want, serialized)
		}
	}
}

// Response-side counterpart: the F6-only Other Related Benefit Detail segment
// (AM39) must be dropped from a D0 response.
func Test_D0ResponseOmitsF6OnlySegment(t *testing.T) {
	i, err := claimdeserializer.DeserializeResponse(RESPONSE_B1)
	if err != nil {
		t.Fatal(err)
	}
	item, ok := i.(response.Billing)
	if !ok {
		t.Fatalf("expected response.Billing.  Got: %T", i)
	}

	planType := "PBM"
	item.Claims[0].OtherRelatedBenefitDetail.PlanType = &planType

	serialized, err := Serialize(&item)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(serialized, "AM39") || strings.Contains(serialized, string(ncpdp.FIELD)+"KS") {
		t.Errorf("D0 response must not contain the F6-only Other Related Benefit Detail segment (AM39).\nGot: %q", serialized)
	}
}
