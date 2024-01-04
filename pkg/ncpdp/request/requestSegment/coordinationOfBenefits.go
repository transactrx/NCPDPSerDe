package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type CoordinationOfBenefits struct {
	SegmentId ncpdp.SegmentId

	Count      *int `field:"code=4C,order=2"`
	OtherPayer []OtherPayer

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherPayer struct {
	CoverageType  *string    `field:"code=5C,order=3"`
	IdQualifier   *string    `field:"code=6C,order=4"`
	Id            *string    `field:"code=7C,order=5"`
	Date          *time.Time `field:"code=E8,format=YYYYMMdd,order=6"`
	ControlNumber *string    `field:"code=A7,order=7"`

	OtherPayerAmountPaidCount *int `field:"code=HB,order=8"`
	OtherPayerAmountsPaid     []OtherPayerAmountPaid

	OtherPayerRejectCount *int `field:"code=5E,order=11"`
	OtherPayerRejects     []OtherPayerReject

	OtherPayerPatientResponsibilityAmountCount *int `field:"code=NR,order=13"`
	OtherPayerPatientResponsibilityAmounts     []OtherPayerPatientResponsibility

	BenefitStageCount   *int `field:"code=MU,order=16"`
	BenefitStageAmounts []BenefitStage
}

type OtherPayerAmountPaid struct {
	Qualifier  *string  `field:"code=HC,order=9"`
	AmountPaid *float64 `field:"code=DV,decimalPlaces=2,overpunch=true,order=10"`
}

type OtherPayerReject struct {
	RejectCode *string `field:"code=6E,order=12"`
}

type OtherPayerPatientResponsibility struct {
	Qualifier *string  `field:"code=NP,order=14"`
	Amount    *float64 `field:"code=NQ,decimalPlaces=2,overpunch=true,order=15"`
}

type BenefitStage struct {
	Qualifier *string  `field:"code=MV,order=17"`
	Amount    *float64 `field:"code=MW,decimalPlaces=2,overpunch=true,order=18"`
}
