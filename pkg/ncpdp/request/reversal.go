package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type Reversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=1"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`	

	//Claim Groups
	Claims []ReversalRecord `group:"max=4"`
}

type ReversalRecord struct {
	Raw                    string                                `field:"code=rawgroup"`
	Claim                  requestsegment.Claim                  `segment:"code=AM07,order=1"`
	Pricing                requestsegment.Pricing                `segment:"code=AM11,order=3"`
	CoordinationOfBenefits requestsegment.CoordinationOfBenefits `segment:"code=AM05,order=4"`
	Dur                    requestsegment.Dur                    `segment:"code=AM08,order=2"`
	DataCollection         requestsegment.DataCollection         `segment:"code=AMXX,order=99"`
	DynamicSegments        []dynamic.DynamicStruct               `segment:"code=dynamic,order=100"`
}
