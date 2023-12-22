package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type Reversal struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance requestsegment.Insurance `segment:"code=AM04"`

	//Claim Groups
	Claims []ReversalRecord `group:"max=4"`
}

type ReversalRecord struct {
	Raw                    string                                `field:"code=rawgroup"`
	Claim                  requestsegment.Claim                  `segment:"code=AM07"`
	Pricing                requestsegment.Pricing                `segment:"code=AM11"`
	CoordinationOfBenefits requestsegment.CoordinationOfBenefits `segment:"code=AM05"`
	Dur                    requestsegment.Dur                    `segment:"code=AM08"`
	DataCollection         requestsegment.DataCollection         `segment:"code=AMXX"`
}
