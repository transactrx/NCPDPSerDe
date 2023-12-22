package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type Rebill struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message   responsesegment.Message   `segment:"code=AM20"`
	Insurance responsesegment.Insurance `segment:"code=AM25"`
	Patient   responsesegment.Patient   `segment:"code=AM29"`

	//Claim Groups
	Claims []RebillRecord `group:"max=4"`
}

type RebillRecord struct {
	Raw                    string                                 `field:"code=rawgroup"`
	Status                 responsesegment.Status                 `segment:"code=AM21"`
	Claim                  responsesegment.Claim                  `segment:"code=AM22"`
	Pricing                responsesegment.Pricing                `segment:"code=AM23"`
	Dur                    responsesegment.Dur                    `segment:"code=AM24"`
	PriorAuthorization     responsesegment.PriorAuthorization     `segment:"code=AM26"`
	CoordinationOfBenefits responsesegment.CoordinationOfBenefits `segment:"code=AM28"`

	//TOOD: Include price update info
}
