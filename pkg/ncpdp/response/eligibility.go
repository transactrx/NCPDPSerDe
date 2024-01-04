package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type Eligibility struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared segments
	Message                        responsesegment.Message                        `segment:"code=AM20,order=1"`
	Insurance                      responsesegment.Insurance                      `segment:"code=AM25,order=2"`
	InsuranceAdditionalInformation responsesegment.InsuranceAdditionalInformation `segment:"code=AM27,order=3"`
	Patient                        responsesegment.Patient                        `segment:"code=AM29,order=4"`
	DynamicSegments                []dynamic.DynamicStruct                        `segment:"code=dynamic,order=100"`

	Groups []EligibilityRecord `group:"max=1"`
}

type EligibilityRecord struct {
	Raw                    string                                 `field:"code=rawgroup"`
	Status                 responsesegment.Status                 `segment:"code=AM21,order=1"`
	CoordinationOfBenefits responsesegment.CoordinationOfBenefits `segment:"code=AM28,order=2"`
	DynamicSegments        []dynamic.DynamicStruct                `segment:"code=dynamic,order=100"`
}
