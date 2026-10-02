package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type OtherRelatedBenefitDetail struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	PlanType                     *string    `field:"code=KS,order=2,sinceVersion=F6"`
	LisLevel                     *string    `field:"code=KF,order=3,sinceVersion=F6"`
	LisEffectiveDate             *time.Time `field:"code=KD,format=YYYYMMdd,order=4,sinceVersion=F6"`
	LisTerminationDate           *time.Time `field:"code=KG,format=YYYYMMdd,order=5,sinceVersion=F6"`
	DisabilityEffectiveDate      *time.Time `field:"code=AH,format=YYYYMMdd,order=6,sinceVersion=F6"`
	EsrdIndicator                *string    `field:"code=A5,order=7,sinceVersion=F6"`
	EsrdEffectiveDate            *time.Time `field:"code=AJ,format=YYYYMMdd,order=8,sinceVersion=F6"`
	EsrdTerminationDate          *time.Time `field:"code=A6,format=YYYYMMdd,order=9,sinceVersion=F6"`
	HospiceEffectiveDate         *time.Time `field:"code=G4,format=YYYYMMdd,order=10,sinceVersion=F6"`
	HospiceTerminationDate       *time.Time `field:"code=G7,format=YYYYMMdd,order=11,sinceVersion=F6"`
	HospiceProviderNumber        *string    `field:"code=G6,order=12,sinceVersion=F6"`
	HospiceFacilityName          *string    `field:"code=G5,order=13,sinceVersion=F6"`
	HospiceTelephoneNumber       *string    `field:"code=A8,order=14,sinceVersion=F6"`
	InstitutionalIndicator       *string    `field:"code=BJ,order=15,sinceVersion=F6"`
	InstitutionalEffectiveDate   *time.Time `field:"code=BK,format=YYYYMMdd,order=16,sinceVersion=F6"`
	InstitutionalTerminationDate *time.Time `field:"code=GD,format=YYYYMMdd,order=17,sinceVersion=F6"`

	OtherBenefitCount *int `field:"code=M8,order=18,countfor=OtherBenefits,sinceVersion=F6"`
	OtherBenefits     []OtherBenefit

	OtherBenefitDetailInformationCount *int `field:"code=N8,order=26,countfor=OtherBenefitDetailInformation,sinceVersion=F6"`
	OtherBenefitDetailInformation      []OtherBenefitDetailInformation

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherBenefit struct {
	TypeCode                *string    `field:"code=PN,order=19,sinceVersion=F6"`
	EffectiveDate           *time.Time `field:"code=MZ,format=YYYYMMdd,order=20,sinceVersion=F6"`
	TerminationDate         *time.Time `field:"code=NN,format=YYYYMMdd,order=21,sinceVersion=F6"`
	StateProvinceAddress    *string    `field:"code=TA,order=22,sinceVersion=F6"`
	TypeId                  *string    `field:"code=N9,order=23,sinceVersion=F6"`
	FacilityName            *string    `field:"code=N1,order=24,sinceVersion=F6"`
	FacilityTelephoneNumber *string    `field:"code=N7,order=25,sinceVersion=F6"`
}

type OtherBenefitDetailInformation struct {
	Indicator               *string    `field:"code=MS,order=27,sinceVersion=F6"`
	EffectiveDate           *time.Time `field:"code=MM,format=YYYYMMdd,order=28,sinceVersion=F6"`
	TerminationDate         *time.Time `field:"code=MX,format=YYYYMMdd,order=29,sinceVersion=F6"`
	ProviderNumber          *string    `field:"code=MR,order=30,sinceVersion=F6"`
	FacilityName            *string    `field:"code=MN,order=31,sinceVersion=F6"`
	FacilityTelephoneNumber *string    `field:"code=MP,order=32,sinceVersion=F6"`
}
