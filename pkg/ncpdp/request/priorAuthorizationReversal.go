package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type PriorAuthorizationReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=1"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`

	//Claim Groups
	Authorizations []PriorAuthorizationReversalRecord `group:"max=4"`
}

type PriorAuthorizationReversalRecord struct {
	Raw             string                                             `field:"code=rawgroup"`
	Authorization   requestsegment.PriorAuthorizationRequestAndBilling `segment:"code=AM12,order=1"`
	DataCollection  requestsegment.DataCollection                      `segment:"code=AMXX,order=99"`
	DynamicSegments []dynamic.DynamicStruct                            `segment:"code=dynamic,order=100"`
}
