package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type PriorAuthorizationRequestOnly struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message responsesegment.Message `segment:"code=AM20"`

	//Authorization Groups
	Authorizations []PriorAuthorizationRequestOnlyRecord `group:"max=4"`
}

type PriorAuthorizationRequestOnlyRecord struct {
	Raw                    string                                 `field:"code=rawgroup"`
	Status                 responsesegment.Status                 `segment:"code=AM21"`
	Claim                  responsesegment.Claim                  `segment:"code=AM22"`
	PriorAuthorization     responsesegment.PriorAuthorization     `segment:"code=AM26"`
	CoordinationOfBenefits responsesegment.CoordinationOfBenefits `segment:"code=AM28"`

	//TOOD: Include price update info
}