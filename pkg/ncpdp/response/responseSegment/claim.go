package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Claim struct {
	SegmentId                                ncpdp.SegmentId
	PrescriptionServiceReference             PrescriptionServiceReference
	PreferredProductCount                    *int `field:"code=9F,order=4"`
	PreferredProducts                        []PreferredProduct
	MedicaidSubrogationInternalControlNumber *string `field:"code=N4,order=10"`

	PlanBenefitOverrideIndicator  *string `field:"code=RC,order=20"`
	PlanBenefitOverrideValueCount *int    `field:"code=RD,order=21"`
	PlanBenefitOverrideValues     []PlanBenefitOverrideValue

	MaximumAgeQualifier              *string    `field:"code=F8,order=23"`
	MaximumAge                       *int       `field:"code=GA,order=24"`
	MinimumAgeQualifier              *string    `field:"code=GQ,order=25"`
	MinimumAge                       *int       `field:"code=GR,order=26"`
	MinimumAmountQualifier           *string    `field:"code=M2,order=27"`
	MinimumAmount                    *float64   `field:"code=M1,decimalPlaces=3,order=28"`
	MaximumAmountQualifier           *string    `field:"code=GC,order=29"`
	MaximumAmount                    *float64   `field:"code=GB,decimalPlaces=3,order=30"`
	MaximumAmountTimePeriod          *string    `field:"code=GF,order=31"`
	MaximumAmountTimePeriodEndDate   *time.Time `field:"code=GH,format=YYYYMMdd,order=32"`
	MaximumAmountTimePeriodStartDate *time.Time `field:"code=GG,format=YYYYMMdd,order=33"`
	MaximumAmountTimePeriodUnits     *int       `field:"code=GJ,order=34"`
	RemainingAmountQualifier         *string    `field:"code=M7,order=35"`
	RemainingAmount                  *float64   `field:"code=M6,decimalPlaces=3,order=36"`

	BenefitTypeOpportunityCount *int `field:"code=AF,order=37"`
	BenefitTypeOpportunities    []BenefitTypeOpportunity

	SubrogationRequestorsReconciliationId *string `field:"code=KY,order=39"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PrescriptionServiceReference struct {
	Qualifier *string `field:"code=EM,order=2"`
	Number    *string `field:"code=D2,order=3"`
}

type PreferredProduct struct {
	IdQualifier        *string  `field:"code=AP,order=5"`
	Id                 *string  `field:"code=AR,order=6"`
	Incentive          *float64 `field:"code=AS,decimalPlaces=2,overpunch=true,order=7"`
	CostShareIncentive *float64 `field:"code=AT,decimalPlaces=2,overpunch=true,order=8"`
	Description        *string  `field:"code=AU,order=9"`

	EffectiveDate   *time.Time `field:"code=ZO,format=YYYYMMdd,order=11"`
	PlanBenefitTier *string    `field:"code=PV,order=12"`
	ReasonCode      *string    `field:"code=PZ,order=13"`

	RequiredTherapyIndicatorCount *int `field:"code=P0,order=14"`
	RequiredTherapyIndicators     []RequiredTherapyIndicator

	RequiredTherapyTimePeriodQualifier *string    `field:"code=P2,order=16"`
	RequiredTherapyTimePeriodDuration  *string    `field:"code=P3,order=17"`
	RequiredTherapyTimePeriodStartDate *time.Time `field:"code=P4,format=YYYYMMdd,order=18"`
	RequiredTherapyTimePeriodEndDate   *time.Time `field:"code=P5,format=YYYYMMdd,order=19"`
}

type RequiredTherapyIndicator struct {
	Indicator *string `field:"code=P1,order=15"`
}

type PlanBenefitOverrideValue struct {
	Value *string `field:"code=RF,order=22"`
}

type BenefitTypeOpportunity struct {
	Opportunity *string `field:"code=AE,order=38"`
}
