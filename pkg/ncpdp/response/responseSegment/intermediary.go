package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Intermediary struct {
	SegmentId ncpdp.SegmentId

	AuthorizationCount *int `field:"code=8R,order=2"`
	Authorizations     []IntermediaryAuthorization

	Messages []IntermediaryMessage

	HelpDeskSupportTypeCount *int `field:"code=KC,order=6"`
	HelpDeskSupportTypes     []IntermediaryHelpDeskSupportType

	HelpDeskBusinessUnitTypeCount *int `field:"code=G9,order=8"`
	HelpDeskBusinessUnits         []IntermediaryHelpDeskBusinessUnit

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type IntermediaryAuthorization struct {
	TypeId *string `field:"code=8S,order=3"`
	Id     *string `field:"code=8T,order=4"`
}

type IntermediaryMessage struct {
	Message *string `field:"code=8Q,order=5"`
}

type IntermediaryHelpDeskSupportType struct {
	Type *string `field:"code=KB,order=7"`
}

type IntermediaryHelpDeskBusinessUnit struct {
	Type                        *string `field:"code=G8,order=9"`
	ContactInformationQualifier *string `field:"code=KA,order=10"`
	ContactInformation          *string `field:"code=JP,order=11"`
	ContactInformationExtension *string `field:"code=JR,order=12"`
}
