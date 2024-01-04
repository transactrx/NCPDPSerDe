package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type InformationReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=1"`
	Patient         requestsegment.Patient   `segment:"code=AM01,order=2"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`

	//Claim Groups
	Claims []InformationReversalRecord `group:"max=4"`
}

type InformationReversalRecord struct {
	Raw             string                        `field:"code=rawgroup"`
	Claim           requestsegment.Claim          `segment:"code=AM07,order=1"`
	DataCollection  requestsegment.DataCollection `segment:"code=AMXX,order=99"`
	DynamicSegments []dynamic.DynamicStruct       `segment:"code=dynamic,order=100"`
}
