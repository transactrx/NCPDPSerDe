package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type PriorAuthorizationReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance requestsegment.Insurance `segment:"code=AM04"`

	//Claim Groups
	Authorizations []PriorAuthorizationReversalRecord `group:"max=4"`
}

type PriorAuthorizationReversalRecord struct {
	Raw            string                                             `field:"code=rawgroup"`
	Authorization  requestsegment.PriorAuthorizationRequestAndBilling `segment:"code=AM12"`
	DataCollection requestsegment.DataCollection                      `segment:"code=AMXX"`
}
