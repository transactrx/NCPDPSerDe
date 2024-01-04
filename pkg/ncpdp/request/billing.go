package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type Billing struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=01"`
	Patient         requestsegment.Patient   `segment:"code=AM01,order=02"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`

	//Claim Groups
	Claims []BillingRecord `group:"max=4"`
}

type BillingRecord struct {
	Raw                     string                                 `field:"code=rawgroup"`
	Claim                   requestsegment.Claim                   `segment:"code=AM07,order=01"`
	Pricing                 requestsegment.Pricing                 `segment:"code=AM11,order=02"`
	Pharmacy                requestsegment.PharmacyProvider        `segment:"code=AM02,order=03"`
	Prescriber              requestsegment.Prescriber              `segment:"code=AM03,order=04"`
	CoordinationOfBenefits  requestsegment.CoordinationOfBenefits  `segment:"code=AM05,order=05"`
	WorkersCompensation     requestsegment.WorkersCompensation     `segment:"code=AM06,order=06"`
	Dur                     requestsegment.Dur                     `segment:"code=AM08,order=07"`
	Coupon                  requestsegment.Coupon                  `segment:"code=AM09,order=08"`
	Compound                requestsegment.Compound                `segment:"code=AM10,order=09"`
	Clinical                requestsegment.Clinical                `segment:"code=AM13,order=10"`
	AdditionalDocumentation requestsegment.AdditionalDocumentation `segment:"code=AM14,order=11"`
	Facility                requestsegment.Facility                `segment:"code=AM15,order=12"`
	Narrative               requestsegment.Narrative               `segment:"code=AM16,order=13"`
	DataCollection          requestsegment.DataCollection          `segment:"code=AMXX,order=99"`
	DynamicSegments         []dynamic.DynamicStruct                `segment:"code=dynamic,order=100"`
}

func (c Billing) RawValue() string {
	return ncpdp.Empty
}
