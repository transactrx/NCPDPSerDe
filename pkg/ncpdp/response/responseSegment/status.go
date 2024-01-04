package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Status struct {
	SegmentId                  ncpdp.SegmentId
	ResponseStatusCode         *string `field:"code=AN,order=2"`
	AuthorizationNumber        *string `field:"code=F3,order=3"`
	RejectCodeCount            *int    `field:"code=FA,order=4"`
	RejectCodes                []RejectCode
	ApprovalMessageCodeCount   *int `field:"code=5F,order=7"`
	ApprovalMessageCodes       []ApprovalMessageCode
	AdditionalMessageCount     *int `field:"code=UF,order=9"`
	AdditionalMessages         []AdditionalMessage
	HelpDeskPhoneNumber        HelpDeskPhoneNumber
	TransactionReferenceNumber *string                 `field:"code=K5,order=15"`
	InternalControlNumber      *string                 `field:"code=A7,order=16"`
	Url                        *string                 `field:"code=MA,order=17"`
	DynamicFields              []dynamic.DynamicStruct `field:"code=dynamic"`
}

type RejectCode struct {
	Code                *string `field:"code=FB,order=5"`
	OccurrenceIndicator *int    `field:"code=4F,order=6"`
}

type ApprovalMessageCode struct {
	Code *string `field:"code=6F,order=8"`
}

type AdditionalMessage struct {
	Qualifier    *string `field:"code=UH,order=10"`
	Message      *string `field:"code=FQ,order=11"`
	Continuation *string `field:"code=UG,order=12"`
}

type HelpDeskPhoneNumber struct {
	Qualifier *string `field:"code=7F,order=13"`
	Number    *string `field:"code=8F,order=14"`
}
