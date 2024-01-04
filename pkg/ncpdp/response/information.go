package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type Information struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message         responsesegment.Message   `segment:"code=AM20,order=1"`
	Insurance       responsesegment.Insurance `segment:"code=AM25,order=2"`
	Patient         responsesegment.Patient   `segment:"code=AM29,order=3"`
	DynamicSegments []dynamic.DynamicStruct   `segment:"code=dynamic,order=100"`

	//Claim Groups
	Claims []InformationRecord `group:"max=4"`
}

type InformationRecord struct {
	Raw             string                  `field:"code=rawgroup"`
	Status          responsesegment.Status  `segment:"code=AM21,order=1"`
	Claim           responsesegment.Claim   `segment:"code=AM22,order=2"`
	Dur             responsesegment.Dur     `segment:"code=AM24,order=3"`
	DynamicSegments []dynamic.DynamicStruct `segment:"code=dynamic,order=100"`

	//TOOD: Include price update info
}
