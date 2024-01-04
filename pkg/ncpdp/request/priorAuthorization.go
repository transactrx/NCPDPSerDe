package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type PriorAuthorization struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=1"`
	Patient         requestsegment.Patient   `segment:"code=AM01,order=2"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`

	//Claim Groups
	Authorizations []PriorAuthorizationRecord `group:"max=4"`
}

type PriorAuthorizationRecord struct {
	Raw                     string                                             `field:"code=rawgroup"`
	Claim                   requestsegment.Claim                               `segment:"code=AM07,order=1"`
	Pricing                 requestsegment.Pricing                             `segment:"code=AM11,order=2"`
	Authorization           requestsegment.PriorAuthorizationRequestAndBilling `segment:"code=AM12,order=3"`
	Pharmacy                requestsegment.PharmacyProvider                    `segment:"code=AM02,order=4"`
	Prescriber              requestsegment.Prescriber                          `segment:"code=AM03,order=5"`
	CoordinationOfBenefits  requestsegment.CoordinationOfBenefits              `segment:"code=AM05,order=6"`
	WorkersCompensation     requestsegment.WorkersCompensation                 `segment:"code=AM06,order=7"`
	Dur                     requestsegment.Dur                                 `segment:"code=AM08,order=8"`
	Compound                requestsegment.Compound                            `segment:"code=AM10,order=9"`
	Clinical                requestsegment.Clinical                            `segment:"code=AM13,order=10"`
	AdditionalDocumentation requestsegment.AdditionalDocumentation             `segment:"code=AM14,order=11"`
	Facility                requestsegment.Facility                            `segment:"code=AM15,order=12"`
	Narrative               requestsegment.Narrative                           `segment:"code=AM16,order=13"`
	DataCollection          requestsegment.DataCollection                      `segment:"code=AMXX,order=99"`
	DynamicSegments         []dynamic.DynamicStruct                            `segment:"code=dynamic,order=100"`
}
