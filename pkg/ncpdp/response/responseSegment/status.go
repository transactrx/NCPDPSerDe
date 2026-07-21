package responsesegment

import (
	"fmt"
	"strings"
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	reflectionutils "github.com/transactrx/NCPDPSerDe/pkg/reflectionUtils"
	stringutils "github.com/transactrx/NCPDPSerDe/pkg/stringUtils"
)

const (
	MaxMessageCount  = 25
	MaxMessageLength = 40
)

type Status struct {
	SegmentId                  ncpdp.SegmentId `json:"-"`
	ResponseStatusCode         *string         `field:"code=AN,order=2"`
	AuthorizationNumber        *string         `field:"code=F3,order=3"`
	RejectCodeCount            *int            `field:"code=FA,order=4,countfor=RejectCodes"`
	RejectCodes                []RejectCode
	ApprovalMessageCodeCount   *int `field:"code=5F,order=7,countfor=ApprovalMessageCodes"`
	ApprovalMessageCodes       []ApprovalMessageCode
	AdditionalMessageCount     *int `field:"code=UF,order=9,countfor=AdditionalMessages"`
	AdditionalMessages         []AdditionalMessage
	HelpDeskPhoneNumber        HelpDeskPhoneNumber
	TransactionReferenceNumber *string `field:"code=K5,order=15"`
	InternalControlNumber      *string `field:"code=A7,order=16"`
	Url                        *string `field:"code=MA,order=17"`

	ReconciliationId *string `field:"code=34,order=18"`

	HelpDeskSupportTypeCount *int `field:"code=BH,order=19,countfor=HelpDeskSupportTypes"`
	HelpDeskSupportTypes     []HelpDeskSupportType

	HelpDeskBusinessUnitTypeCount *int `field:"code=BB,order=21,countfor=HelpDeskBusinessUnits"`
	HelpDeskBusinessUnits         []HelpDeskBusinessUnit

	HelpDeskContactInformationExtension *string    `field:"code=BD,order=25"`
	AdjudicatedProgramType              *string    `field:"code=ZR,order=26"`
	NextAvailableFillDate               *time.Time `field:"code=BT,format=YYYYMMdd,order=27"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type HelpDeskSupportType struct {
	Type *string `field:"code=BG,order=20"`
}

type HelpDeskBusinessUnit struct {
	Type                        *string `field:"code=BA,order=22"`
	ContactInformationQualifier *string `field:"code=BF,order=23"`
	ContactInformation          *string `field:"code=BC,order=24"`
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

// Append message to additional messages array, chunking message into the appropriate
// data size and using continuation characters where required.
// Returns error when unable to append message.
func (segment *Status) AppendMessage(verbiage string) error {
	if segment == nil || len(verbiage) == 0 {
		return nil
	}

	// Uppercase it
	verbiage = strings.ToUpper(verbiage)

	// Check message count
	if len(segment.AdditionalMessages) >= MaxMessageCount {
		return fmt.Errorf("maximum message count exceeded")
	}

	// Break into chunks by max field length
	chunks := stringutils.Chunk(verbiage, MaxMessageLength)

	if len(chunks)+len(segment.AdditionalMessages) >= MaxMessageCount {
		return fmt.Errorf("maximum message count exceeded")
	}

	for i := 0; i < len(chunks); i++ {
		msgText := chunks[i]

		addtlMessage := AdditionalMessage{
			Qualifier: reflectionutils.ToPointer(fmt.Sprintf("%02d", len(segment.AdditionalMessages)+1)),
			Message:   reflectionutils.ToPointer(msgText),
		}

		if i < len(chunks)-1 {
			addtlMessage.Continuation = reflectionutils.ToPointer("+")
		}

		segment.AdditionalMessages = append(segment.AdditionalMessages, addtlMessage)
		segment.AdditionalMessageCount = reflectionutils.ToPointer(len(segment.AdditionalMessages))
	}

	return nil
}
