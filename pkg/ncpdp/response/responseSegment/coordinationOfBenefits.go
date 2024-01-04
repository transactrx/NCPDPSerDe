package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type CoordinationOfBenefits struct {
	SegmentId       ncpdp.SegmentId
	OtherPayerCount *int `field:"code=NT,order=2"`
	OtherPayers     []OtherPayer
	DynamicFields   []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherPayer struct {
	CoverageType            *string    `field:"code=5C,order=3"`
	IdQualifier             *string    `field:"code=6C,order=4"`
	Id                      *string    `field:"code=7C,order=5"`
	Pcn                     *string    `field:"code=MH,order=6"`
	CardholderId            *string    `field:"code=NU,order=7"`
	GroupId                 *string    `field:"code=MJ,order=8"`
	PersonCode              *string    `field:"code=UV,order=9"`
	HelpDeskPhone           *string    `field:"code=UB,order=10"`
	PatientRelationshipCode *string    `field:"code=UW,order=11"`
	BenefitEffectiveDate    *time.Time `field:"code=UX,format=YYYYMMdd,order=12"`
	BenefitTerminationDate  *time.Time `field:"code=UY,format=YYYYMMdd,order=13"`
}
