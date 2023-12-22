package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type CoordinationOfBenefits struct {
	SegmentId       ncpdp.SegmentId
	OtherPayerCount *int `field:"code=NT"`
	OtherPayers     []OtherPayer
}

type OtherPayer struct {
	CoverageType            *string    `field:"code=5C"`
	IdQualifier             *string    `field:"code=6C"`
	Id                      *string    `field:"code=7C"`
	Pcn                     *string    `field:"code=MH"`
	CardholderId            *string    `field:"code=NU"`
	GroupId                 *string    `field:"code=MJ"`
	PersonCode              *string    `field:"code=UV"`
	HelpDeskPhone           *string    `field:"code=UB"`
	PatientRelationshipCode *string    `field:"code=UW"`
	BenefitEffectiveDate    *time.Time `field:"code=UX,format=YYYYMMdd"`
	BenefitTerminationDate  *time.Time `field:"code=UY,format=YYYYMMdd"`
}
