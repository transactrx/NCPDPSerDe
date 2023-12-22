package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type ControlledSubstanceReportingRebill struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Patient requestsegment.Patient `segment:"code=AM01"`

	//Claim Groups
	Claims []ControlledSubstanceReportingRebillRecord `group:"max=4"`
}

type ControlledSubstanceReportingRebillRecord struct {
	Raw            string                          `field:"code=rawgroup"`
	Claim          requestsegment.Claim            `segment:"code=AM07"`
	Pharmacy       requestsegment.PharmacyProvider `segment:"code=AM02"`
	Prescriber     requestsegment.Prescriber       `segment:"code=AM03"`
	DataCollection requestsegment.DataCollection   `segment:"code=AMXX"`
}
