package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type Eligibility struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Segments
	Insurance               requestsegment.Insurance               `segment:"code=AM04"`
	Patient                 requestsegment.Patient                 `segment:"code=AM01"`
	Pharmacy                requestsegment.PharmacyProvider        `segment:"code=AM02"`
	Prescriber              requestsegment.Prescriber              `segment:"code=AM03"`
	AdditionalDocumentation requestsegment.AdditionalDocumentation `segment:"code=AM14"`
	DataCollection          requestsegment.DataCollection          `segment:"code=AMXX"`
}
