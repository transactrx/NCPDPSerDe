package request

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	requestsegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request/requestSegment"
)

type PriorAuthorizationRequestOnly struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.RequestHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Insurance requestsegment.Insurance `segment:"code=AM04"`
	Patient   requestsegment.Patient   `segment:"code=AM01"`

	//Claim Groups
	Authorizations []PriorAuthorizationRequestOnlyRecord `group:"max=4"`
}

type PriorAuthorizationRequestOnlyRecord struct {
	Raw                 string                                             `field:"code=rawgroup"`
	Claim               requestsegment.Claim                               `segment:"code=AM07"`
	Authorization       requestsegment.PriorAuthorizationRequestAndBilling `segment:"code=AM12"`
	Prescriber          requestsegment.Prescriber                          `segment:"code=AM03"`
	WorkersCompensation requestsegment.WorkersCompensation                 `segment:"code=AM06"`
	Dur                 requestsegment.Dur                                 `segment:"code=AM08"`
	Compound            requestsegment.Compound                            `segment:"code=AM10"`
	Clinical            requestsegment.Clinical                            `segment:"code=AM13"`
	DataCollection      requestsegment.DataCollection                      `segment:"code=AMXX"`
}
