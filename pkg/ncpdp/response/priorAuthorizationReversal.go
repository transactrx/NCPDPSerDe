package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type PriorAuthorizationReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message responsesegment.Message `segment:"code=AM20"`

	//Authorization Groups
	Authorizations []PriorAuthorizationReversalRecord `group:"max=4"`
}

type PriorAuthorizationReversalRecord struct {
	Raw    string                 `field:"code=rawgroup"`
	Status responsesegment.Status `segment:"code=AM21"`

	//TOOD: Include price update info
}
