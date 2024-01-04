package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Pricing struct {
	SegmentId                                                    ncpdp.SegmentId
	PatientPayAmount                                             *float64 `field:"code=F5,decimalPlaces=2,overpunch=true,order=2"`
	IngredientCostPaid                                           *float64 `field:"code=F6,decimalPlaces=2,overpunch=true,order=3"`
	DispensingFeePaid                                            *float64 `field:"code=F7,decimalPlaces=2,overpunch=true,order=4"`
	TaxExemptIndicator                                           *string  `field:"code=AV,order=5"`
	FlatSalesTaxAmountPaid                                       *float64 `field:"code=AW,decimalPlaces=2,overpunch=true,order=6"`
	PercentageSalesTaxAmountPaid                                 *float64 `field:"code=AX,decimalPlaces=2,overpunch=true,order=7"`
	PercentageSalesTaxRatePaid                                   *float64 `field:"code=AY,decimalPlaces=4,overpunch=true,order=8"`
	PercentageSalesTaxBasisPaid                                  *string  `field:"code=AZ,order=9"`
	IncentiveAmountPaid                                          *float64 `field:"code=FL,decimalPlaces=2,overpunch=true,order=10"`
	ProfessionalServiceFeePaid                                   *float64 `field:"code=J1,decimalPlaces=2,overpunch=true,order=11"`
	OtherAmountPaidCount                                         *int     `field:"code=J2,order=12"`
	OtherAmountsPaid                                             []OtherAmountPaid
	OtherPayerAmountRecognized                                   *float64 `field:"code=J5,decimalPlaces=2,overpunch=true,order=15"`
	TotalAmountPaid                                              *float64 `field:"code=F9,decimalPlaces=2,overpunch=true,order=16"`
	BasisOfReimbursementDetermination                            *int     `field:"code=FM,order=17"`
	AmountAttributedToSalesTax                                   *float64 `field:"code=FN,decimalPlaces=2,overpunch=true,order=18"`
	AccumulatedDeductibleAmount                                  *float64 `field:"code=FC,decimalPlaces=2,overpunch=true,order=19"`
	RemainingDeductibleAmount                                    *float64 `field:"code=FD,decimalPlaces=2,overpunch=true,order=20"`
	RemainingBenefitAmount                                       *float64 `field:"code=FE,decimalPlaces=2,overpunch=true,order=21"`
	AmountAppliedToPeriodicDeductible                            *float64 `field:"code=FH,decimalPlaces=2,overpunch=true,order=22"`
	CopayAmount                                                  *float64 `field:"code=FI,decimalPlaces=2,overpunch=true,order=23"`
	AmountExceedingPeriodicBenefitMax                            *float64 `field:"code=FK,decimalPlaces=2,overpunch=true,order=24"`
	BasisOfCalculation                                           BasisOfCalculation
	AmountAttributedToProcessorFee                               *float64 `field:"code=NZ,decimalPlaces=2,overpunch=true,order=29"`
	PatientSalesTaxAmount                                        *float64 `field:"code=EQ,decimalPlaces=2,overpunch=true,order=30"`
	PlanSalesTaxAmount                                           *float64 `field:"code=2Y,decimalPlaces=2,overpunch=true,order=31"`
	CoinsuranceAmount                                            *float64 `field:"code=4U,decimalPlaces=2,overpunch=true,order=32"`
	BenefitStageCount                                            *int     `field:"code=MU,order=34"`
	BenefitStageAmounts                                          []BenefitStage
	EstimatedGenericSavings                                      *float64                `field:"code=G3,decimalPlaces=2,overpunch=true,order=37"`
	SpendingAccountAmountRemaining                               *float64                `field:"code=UC,decimalPlaces=2,overpunch=true,order=38"`
	HealthPlanFundedAssistanceAmount                             *float64                `field:"code=UD,decimalPlaces=2,overpunch=true,order=39"`
	AmountAttributedToNetworkSelection                           *float64                `field:"code=UJ,decimalPlaces=2,overpunch=true,order=40"`
	AmountAttributedToBrandProductSelection                      *float64                `field:"code=UK,decimalPlaces=2,overpunch=true,order=41"`
	AmountAttributedToNonPreferredFormularyProductSelection      *float64                `field:"code=UM,decimalPlaces=2,overpunch=true,order=42"`
	AmountAttributedToBrandNonPreferredFormularyProductSelection *float64                `field:"code=UN,decimalPlaces=2,overpunch=true,order=43"`
	AmountAttributedToCoverageGap                                *float64                `field:"code=UP,decimalPlaces=2,overpunch=true,order=44"`
	IngredientCostContractedReimbursableAmount                   *float64                `field:"code=U8,decimalPlaces=2,overpunch=true,order=45"`
	DispensingFeeContractedReimbursableAmount                    *float64                `field:"code=U9,decimalPlaces=2,overpunch=true,order=46"`
	DynamicFields                                                []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherAmountPaid struct {
	Qualifier  *string  `field:"code=J3,order=13"`
	AmountPaid *float64 `field:"code=J4,decimalPlaces=2,overpunch=true,order=14"`
}

type BasisOfCalculation struct {
	DispensingFee      *string `field:"code=HH,order=25"`
	Copay              *string `field:"code=HJ,order=26"`
	FlatSalesTax       *string `field:"code=HK,order=27"`
	PercentageSalesTax *string `field:"code=HM,order=28"`
	Coinsurance        *string `field:"code=4V,order=33"`
}

type BenefitStage struct {
	Qualifier *string  `field:"code=MV,order=35"`
	Amount    *float64 `field:"code=MW,decimalPlaces=2,overpunch=true,order=36"`
}
