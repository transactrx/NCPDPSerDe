package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type InformationReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance requestsegment.Insurance `segment:"code=AM04"`
	Patient   requestsegment.Patient   `segment:"code=AM01"`

	//Claim Groups
	Claims []InformationReversalRecord `group:"max=4"`
}

type InformationReversalRecord struct {
	Raw            string                        `field:"code=rawgroup"`
	Claim          requestsegment.Claim          `segment:"code=AM07"`
	DataCollection requestsegment.DataCollection `segment:"code=AMXX"`
}
