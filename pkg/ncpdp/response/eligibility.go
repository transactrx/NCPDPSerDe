package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type Eligibility struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared segments
	Message                        responsesegment.Message                        `segment:"code=AM20"`
	Insurance                      responsesegment.Insurance                      `segment:"code=AM25"`
	InsuranceAdditionalInformation responsesegment.InsuranceAdditionalInformation `segment:"code=AM27"`
	Patient                        responsesegment.Patient                        `segment:"code=AM29"`

	Groups []EligibilityRecord `group:"max=1"`
}

type EligibilityRecord struct {
	Raw                    string                                 `field:"code=rawgroup"`
	Status                 responsesegment.Status                 `segment:"code=AM21"`
	CoordinationOfBenefits responsesegment.CoordinationOfBenefits `segment:"code=AM28"`
}
