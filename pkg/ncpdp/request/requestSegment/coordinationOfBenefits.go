package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type OtherPayer struct {
	CoverageType  *string    `field:"code=5C"`
	IdQualifier   *string    `field:"code=6C"`
	Id            *string    `field:"code=7C"`
	Date          *time.Time `field:"code=E8,format=YYYYMMdd"`
	ControlNumber *string    `field:"code=A7"`
}

type OtherPayerAmountPaid struct {
	Qualifier  *string  `field:"code=HC"`
	AmountPaid *float64 `field:"code=DV,decimalPlaces=2,overpunch=true"`
}

type OtherPayerReject struct {
	RejectCode *string `field:"code=6E"`
}

type OtherPayerPatientResponsibility struct {
	Qualifier *string  `field:"code=NP"`
	Amount    *float64 `field:"code=NQ,decimalPlaces=2,overpunch=true"`
}

type BenefitStage struct {
	Qualifier *string  `field:"code=MV"`
	Amount    *float64 `field:"code=MW,decimalPlaces=2,overpunch=true"`
}

type CoordinationOfBenefits struct {
	SegmentId ncpdp.SegmentId

	Count      *int `field:"code=4C"`
	OtherPayer []OtherPayer

	OtherPayerAmountPaidCount *int `field:"code=HB"`
	OtherPayerAmountsPaid     []OtherPayerAmountPaid

	OtherPayerRejectCount *int `field:"code=5E"`
	OtherPayerRejects     []OtherPayerReject

	OtherPayerPatientResponsibilityAmountCount *int `field:"code=NR"`
	OtherPayerPatientResponsibilityAmounts     []OtherPayerPatientResponsibility

	BenefitStageCount   *int `field:"code=MU"`
	BenefitStageAmounts []BenefitStage
}
