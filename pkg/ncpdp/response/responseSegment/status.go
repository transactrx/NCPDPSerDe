package responsesegment

import (
	"fmt"
	"strconv"
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

	// Additional Message Information Qualifier (132-UH) values "01"-"09" are the
	// only values defined for free text; higher numeric and alpha values are
	// ECL-coded meanings (e.g. "12", "38" = Max Fair Price) and must never be
	// generated or reassigned by this library.
	maxFreeTextQualifier = 9
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

	ReconciliationId *string `field:"code=34,order=18,sinceVersion=F6"`

	HelpDeskSupportTypeCount *int `field:"code=BH,order=19,countfor=HelpDeskSupportTypes,sinceVersion=F6"`
	HelpDeskSupportTypes     []HelpDeskSupportType

	AdjudicatedProgramType *string    `field:"code=ZR,order=26,sinceVersion=F6"`
	NextAvailableFillDate  *time.Time `field:"code=BT,format=YYYYMMdd,order=27,sinceVersion=F6"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

// HelpDeskSupportType is one occurrence of the F6 Help Desk Support loop. Per
// the vF6 implementation guide ("The Help Desk Business Unit Type loop will
// repeat within an occurrence of the Help Desk Support Type loop"), each
// support type carries its own business unit count (BB) and contacts.
type HelpDeskSupportType struct {
	Type *string `field:"code=BG,order=20,sinceVersion=F6"`

	BusinessUnitCount *int `field:"code=BB,order=21,countfor=BusinessUnits,sinceVersion=F6"`
	BusinessUnits     []HelpDeskBusinessUnit
}

type HelpDeskBusinessUnit struct {
	Type                        *string `field:"code=BA,order=22,sinceVersion=F6"`
	ContactInformationQualifier *string `field:"code=BF,order=23,sinceVersion=F6"`
	ContactInformation          *string `field:"code=BC,order=24,sinceVersion=F6"`
	ContactInformationExtension *string `field:"code=BD,order=25,sinceVersion=F6"`
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

	if len(chunks)+len(segment.AdditionalMessages) > MaxMessageCount {
		return fmt.Errorf("maximum message count exceeded")
	}

	// Qualifiers must occur no more than once per transaction, so number the new
	// chunks after the highest free-text qualifier already present rather than
	// by slice position (which would collide with, or renumber past, coded
	// values such as "12" or "38").
	next := segment.nextFreeTextQualifier()
	if next+len(chunks)-1 > maxFreeTextQualifier {
		return fmt.Errorf("maximum free text message count exceeded")
	}

	newMessages := make([]AdditionalMessage, 0, len(chunks))

	for i := 0; i < len(chunks); i++ {
		msgText := chunks[i]

		addtlMessage := AdditionalMessage{
			Qualifier: reflectionutils.ToPointer(fmt.Sprintf("%02d", next+i)),
			Message:   reflectionutils.ToPointer(msgText),
		}

		if i < len(chunks)-1 {
			addtlMessage.Continuation = reflectionutils.ToPointer("+")
		}

		newMessages = append(newMessages, addtlMessage)
	}

	// Qualifier values must remain ordered within the transaction (numeric
	// characters precede alpha), so insert the new block ahead of any
	// higher-valued coded qualifiers instead of appending at the end.
	idx := segment.qualifierInsertIndex(fmt.Sprintf("%02d", next))
	segment.AdditionalMessages = append(segment.AdditionalMessages[:idx], append(newMessages, segment.AdditionalMessages[idx:]...)...)
	segment.AdditionalMessageCount = reflectionutils.ToPointer(len(segment.AdditionalMessages))

	return nil
}

// SetMessage sets free-text messages to the given message, chunking it into the appropriate data
// size and using continuation characters where required. An empty message removes any existing one.
// FQ entries with qualifiers greater than 09 are always retained.
// Returns an error when unable to set the message.
func (segment *Status) SetMessage(message string) error {
	if segment == nil {
		return nil
	}

	// Preserve original values in case we need to roll back
	previous := segment.AdditionalMessages
	previousCount := segment.AdditionalMessageCount

	// Preserve only non-free-text messages
	segment.AdditionalMessages = nil
	for _, msg := range previous {
		if _, ok := freeTextQualifier(msg.Qualifier); !ok {
			segment.AdditionalMessages = append(segment.AdditionalMessages, msg)
		}
	}

	segment.AdditionalMessageCount = nil
	if len(segment.AdditionalMessages) > 0 {
		segment.AdditionalMessageCount = reflectionutils.ToPointer(len(segment.AdditionalMessages))
	}

	if err := segment.AppendMessage(message); err != nil {
		segment.AdditionalMessages = previous
		segment.AdditionalMessageCount = previousCount
		return err
	}

	return nil
}

// Next unused free-text qualifier value (1-9), continuing after the highest
// free-text qualifier already present. Coded (>09 or alpha) and missing
// qualifiers are ignored.
func (segment *Status) nextFreeTextQualifier() int {
	next := 1

	for _, msg := range segment.AdditionalMessages {
		n, ok := freeTextQualifier(msg.Qualifier)
		if !ok {
			continue
		}

		if n+1 > next {
			next = n + 1
		}
	}

	return next
}

// freeTextQualifier returns the numeric value of a free-text qualifier (1-9).
// Returns false for structured (greater than 09), non-numeric, or missing qualifiers
func freeTextQualifier(qualifier *string) (int, bool) {
	if qualifier == nil {
		return 0, false
	}

	n, err := strconv.Atoi(strings.TrimSpace(*qualifier))
	if err != nil || n < 1 || n > maxFreeTextQualifier {
		return 0, false
	}

	return n, true
}

// Index at which a message with the given qualifier belongs: before the first
// existing message whose qualifier sorts higher. Two-character qualifiers sort
// correctly as strings (ASCII digits precede alpha, matching the NCPDP 0-9
// then A-Z ordering rule).
func (segment *Status) qualifierInsertIndex(qualifier string) int {
	for i, msg := range segment.AdditionalMessages {
		if msg.Qualifier == nil {
			continue
		}

		if normalizeQualifier(*msg.Qualifier) > qualifier {
			return i
		}
	}

	return len(segment.AdditionalMessages)
}

// Pad numeric qualifier values to two characters so a payer-sent "1" sorts
// with "01"; non-numeric values are compared as sent.
func normalizeQualifier(qualifier string) string {
	qualifier = strings.TrimSpace(qualifier)

	if n, err := strconv.Atoi(qualifier); err == nil {
		return fmt.Sprintf("%02d", n)
	}

	return qualifier
}
