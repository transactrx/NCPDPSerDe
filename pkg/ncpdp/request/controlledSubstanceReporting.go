package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type ControlledSubstanceReporting struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Patient         requestsegment.Patient  `segment:"code=AM01,order=1"`
	DynamicSegments []dynamic.DynamicStruct `segment:"code=dynamic,order=100"`

	//Claim Groups
	Claims []ControlledSubstanceReportingRecord `group:"max=4"`
}

type ControlledSubstanceReportingRecord struct {
	Raw             string                          `field:"code=rawgroup"`
	Claim           requestsegment.Claim            `segment:"code=AM07,order=1"`
	Pharmacy        requestsegment.PharmacyProvider `segment:"code=AM02,order=2"`
	Prescriber      requestsegment.Prescriber       `segment:"code=AM03,order=3"`
	DataCollection  requestsegment.DataCollection   `segment:"code=AMXX,order=99"`
	DynamicSegments []dynamic.DynamicStruct         `segment:"code=dynamic,order=100"`
}
