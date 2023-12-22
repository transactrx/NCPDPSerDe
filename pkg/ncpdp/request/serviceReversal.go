package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type ServiceReversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance requestsegment.Insurance `segment:"code=AM04"`

	//Claim Groups
	Claims []ServiceReversalRecord `group:"max=4"`
}

type ServiceReversalRecord struct {
	Raw                    string                                `field:"code=rawgroup"`
	Claim                  requestsegment.Claim                  `segment:"code=AM07"`
	CoordinationOfBenefits requestsegment.CoordinationOfBenefits `segment:"code=AM05"`
	DataCollection         requestsegment.DataCollection         `segment:"code=AMXX"`
}
