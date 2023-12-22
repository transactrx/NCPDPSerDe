package responsesegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Status struct {
	SegmentId                  ncpdp.SegmentId
	ResponseStatusCode         *string `field:"code=AN"`
	AuthorizationNumber        *string `field:"code=F3"`
	RejectCodeCount            *int    `field:"code=FA"`
	RejectCodes                []RejectCode
	ApprovalMessageCodeCount   *int `field:"code=5F"`
	ApprovalMessageCodes       []ApprovalMessageCode
	AdditionalMessageCount     *int `field:"code=UF"`
	AdditionalMessages         []AdditionalMessage
	HelpDeskPhoneNumber        HelpDeskPhoneNumber
	TransactionReferenceNumber *string `field:"code=K5"`
	InternalControlNumber      *string `field:"code=A7"`
	Url                        *string `field:"code=MA"`
}

type RejectCode struct {
	Code                *string `field:"code=FB"`
	OccurrenceIndicator *int    `field:"code=4F"`
}

type ApprovalMessageCode struct {
	Code *string `field:"code=6F"`
}

type AdditionalMessage struct {
	Qualifier    *string `field:"code=UH"`
	Message      *string `field:"code=FQ"`
	Continuation *string `field:"code=UG"`
}

type HelpDeskPhoneNumber struct {
	Qualifier *string `field:"code=7F"`
	Number    *string `field:"code=8F"`
}
