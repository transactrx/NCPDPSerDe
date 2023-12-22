package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Cardholder struct {
	Id        *string `field:"code=C2"`
	FirstName *string `field:"code=CC"`
	LastName  *string `field:"code=CD"`
}

type OtherInsurancePayer struct {
	Bin          *string `field:"code=MG"`
	Pcn          *string `field:"code=MH"`
	CardholderId *string `field:"code=NU"`
	GroupId      *string `field:"code=MJ"`
}

type Medicaid struct {
	Indicator    *string `field:"code=2B"`
	Id           *string `field:"code=N5"`
	AgencyNumber *string `field:"code=N6"`
}

type Insurance struct {
	SegmentId                   ncpdp.SegmentId
	Cardholder                  Cardholder
	HomePlan                    *string `field:"code=CE"`
	PlanId                      *string `field:"code=FO"`
	EligbilityClarificationCode *string `field:"code=C9"`
	GroupId                     *string `field:"code=C1"`
	PersonCode                  *string `field:"code=C3"`
	PatientRelationshipCode     *string `field:"code=C6"`
	OtherPayer                  OtherInsurancePayer
	MedigapId                   *string `field:"code=2A"`
	Medicaid                    Medicaid
	ProviderAcceptAssignment    *string `field:"code=2D"`
	CmsPartDQualifiedFacility   *string `field:"code=G2"`
}
