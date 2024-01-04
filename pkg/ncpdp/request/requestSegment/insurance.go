package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Insurance struct {
	SegmentId                   ncpdp.SegmentId
	Cardholder                  Cardholder
	HomePlan                    *string `field:"code=CE,order=5"`
	PlanId                      *string `field:"code=FO,order=6"`
	EligbilityClarificationCode *string `field:"code=C9,order=7"`
	GroupId                     *string `field:"code=C1,order=8"`
	PersonCode                  *string `field:"code=C3,order=9"`
	PatientRelationshipCode     *string `field:"code=C6,order=10"`
	OtherPayer                  OtherInsurancePayer
	MedigapId                   *string `field:"code=2A,order=15"`
	Medicaid                    Medicaid
	ProviderAcceptAssignment    *string                 `field:"code=2D,order=17"`
	CmsPartDQualifiedFacility   *string                 `field:"code=G2,order=18"`
	DynamicFields               []dynamic.DynamicStruct `field:"code=dynamic"`
}

type Cardholder struct {
	Id        *string `field:"code=C2,order=2"`
	FirstName *string `field:"code=CC,order=3"`
	LastName  *string `field:"code=CD,order=4"`
}

type OtherInsurancePayer struct {
	Bin          *string `field:"code=MG,order=11"`
	Pcn          *string `field:"code=MH,order=12"`
	CardholderId *string `field:"code=NU,order=13"`
	GroupId      *string `field:"code=MJ,order=14"`
}

type Medicaid struct {
	Indicator    *string `field:"code=2B,order=16"`
	Id           *string `field:"code=N5,order=19"`
	AgencyNumber *string `field:"code=N6,order=20"`
}
