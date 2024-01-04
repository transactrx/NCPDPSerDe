package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type PredeterminationOfBenefits struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance       requestsegment.Insurance `segment:"code=AM04,order=1"`
	Patient         requestsegment.Patient   `segment:"code=AM01,order=2"`
	DynamicSegments []dynamic.DynamicStruct  `segment:"code=dynamic,order=100"`

	//Claim Groups
	Claims []PredeterminationOfBenefitsRecord `group:"max=4"`
}

type PredeterminationOfBenefitsRecord struct {
	Raw             string                          `field:"code=rawgroup"`
	Claim           requestsegment.Claim            `segment:"code=AM07,order=1"`
	Pricing         requestsegment.Pricing          `segment:"code=AM11,order=2"`
	Pharmacy        requestsegment.PharmacyProvider `segment:"code=AM02,order=3"`
	Prescriber      requestsegment.Prescriber       `segment:"code=AM03,order=4"`
	Dur             requestsegment.Dur              `segment:"code=AM08,order=5"`
	Compound        requestsegment.Compound         `segment:"code=AM10,order=6"`
	Clinical        requestsegment.Clinical         `segment:"code=AM13,order=7"`
	Facility        requestsegment.Facility         `segment:"code=AM15,order=8"`
	DataCollection  requestsegment.DataCollection   `segment:"code=AMXX,order=99"`
	DynamicSegments []dynamic.DynamicStruct         `segment:"code=dynamic,order=100"`
}
