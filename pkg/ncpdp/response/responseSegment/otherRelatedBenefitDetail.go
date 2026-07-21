package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type OtherRelatedBenefitDetail struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	PlanType                     *string    `field:"code=KS,order=2"`
	LisLevel                     *string    `field:"code=KF,order=3"`
	LisEffectiveDate             *time.Time `field:"code=KD,format=YYYYMMdd,order=4"`
	LisTerminationDate           *time.Time `field:"code=KG,format=YYYYMMdd,order=5"`
	DisabilityEffectiveDate      *time.Time `field:"code=AH,format=YYYYMMdd,order=6"`
	EsrdIndicator                *string    `field:"code=A5,order=7"`
	EsrdEffectiveDate            *time.Time `field:"code=AJ,format=YYYYMMdd,order=8"`
	EsrdTerminationDate          *time.Time `field:"code=A6,format=YYYYMMdd,order=9"`
	HospiceEffectiveDate         *time.Time `field:"code=G4,format=YYYYMMdd,order=10"`
	HospiceTerminationDate       *time.Time `field:"code=G7,format=YYYYMMdd,order=11"`
	HospiceProviderNumber        *string    `field:"code=G6,order=12"`
	HospiceFacilityName          *string    `field:"code=G5,order=13"`
	HospiceTelephoneNumber       *string    `field:"code=A8,order=14"`
	InstitutionalIndicator       *string    `field:"code=BJ,order=15"`
	InstitutionalEffectiveDate   *time.Time `field:"code=BK,format=YYYYMMdd,order=16"`
	InstitutionalTerminationDate *time.Time `field:"code=GD,format=YYYYMMdd,order=17"`

	OtherBenefitCount *int `field:"code=M8,order=18,countfor=OtherBenefits"`
	OtherBenefits     []OtherBenefit

	OtherBenefitDetailInformationCount *int `field:"code=N8,order=26,countfor=OtherBenefitDetailInformation"`
	OtherBenefitDetailInformation      []OtherBenefitDetailInformation

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherBenefit struct {
	TypeCode                *string    `field:"code=PN,order=19"`
	EffectiveDate           *time.Time `field:"code=MZ,format=YYYYMMdd,order=20"`
	TerminationDate         *time.Time `field:"code=NN,format=YYYYMMdd,order=21"`
	StateProvinceAddress    *string    `field:"code=TA,order=22"`
	TypeId                  *string    `field:"code=N9,order=23"`
	FacilityName            *string    `field:"code=N1,order=24"`
	FacilityTelephoneNumber *string    `field:"code=N7,order=25"`
}

type OtherBenefitDetailInformation struct {
	Indicator               *string    `field:"code=MS,order=27"`
	EffectiveDate           *time.Time `field:"code=MM,format=YYYYMMdd,order=28"`
	TerminationDate         *time.Time `field:"code=MX,format=YYYYMMdd,order=29"`
	ProviderNumber          *string    `field:"code=MR,order=30"`
	FacilityName            *string    `field:"code=MN,order=31"`
	FacilityTelephoneNumber *string    `field:"code=MP,order=32"`
}
