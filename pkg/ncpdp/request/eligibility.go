package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type Eligibility struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Segments
	Insurance               requestsegment.Insurance               `segment:"code=AM04,order=1"`
	Patient                 requestsegment.Patient                 `segment:"code=AM01,order=2"`
	Pharmacy                requestsegment.PharmacyProvider        `segment:"code=AM02,order=3"`
	Prescriber              requestsegment.Prescriber              `segment:"code=AM03,order=4"`
	AdditionalDocumentation requestsegment.AdditionalDocumentation `segment:"code=AM14,order=5"`
	DataCollection          requestsegment.DataCollection          `segment:"code=AMXX,order=99"`
	DynamicSegments         []dynamic.DynamicStruct                `segment:"code=dynamic,order=100"`
}
