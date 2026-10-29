package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Intermediary struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	AuthorizationCount *int `field:"code=8R,order=2,countfor=Authorizations,sinceVersion=F6"`
	Authorizations     []IntermediaryAuthorization

	Messages []IntermediaryMessage

	HelpDeskSupportTypeCount *int `field:"code=KC,order=6,countfor=HelpDeskSupportTypes,sinceVersion=F6"`
	HelpDeskSupportTypes     []IntermediaryHelpDeskSupportType

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type IntermediaryAuthorization struct {
	TypeId *string `field:"code=8S,order=3,sinceVersion=F6"`
	Id     *string `field:"code=8T,order=4,sinceVersion=F6"`
}

type IntermediaryMessage struct {
	Message *string `field:"code=8Q,order=5,sinceVersion=F6"`
}

// IntermediaryHelpDeskSupportType is one occurrence of the F6 Intermediary
// Help Desk Support loop. Per the vF6 implementation guide ("The Intermediary
// Help Desk Business Unit Type loop will repeat within an occurrence of the
// Intermediary Help Desk Support Type loop"), each support type carries its
// own business unit count (G9) and contacts.
type IntermediaryHelpDeskSupportType struct {
	Type *string `field:"code=KB,order=7,sinceVersion=F6"`

	BusinessUnitCount *int `field:"code=G9,order=8,countfor=BusinessUnits,sinceVersion=F6"`
	BusinessUnits     []IntermediaryHelpDeskBusinessUnit
}

type IntermediaryHelpDeskBusinessUnit struct {
	Type                        *string `field:"code=G8,order=9,sinceVersion=F6"`
	ContactInformationQualifier *string `field:"code=KA,order=10,sinceVersion=F6"`
	ContactInformation          *string `field:"code=JP,order=11,sinceVersion=F6"`
	ContactInformationExtension *string `field:"code=JR,order=12,sinceVersion=F6"`
}
