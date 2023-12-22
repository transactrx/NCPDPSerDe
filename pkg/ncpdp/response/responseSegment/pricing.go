package responsesegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Pricing struct {
	SegmentId                                                    ncpdp.SegmentId
	PatientPayAmount                                             *float64 `field:"code=F5,decimalPlaces=2,overpunch=true"`
	IngredientCostPaid                                           *float64 `field:"code=F6,decimalPlaces=2,overpunch=true"`
	DispensingFeePaid                                            *float64 `field:"code=F7,decimalPlaces=2,overpunch=true"`
	TaxExemptIndicator                                           *string  `field:"code=AV"`
	FlatSalesTaxAmountPaid                                       *float64 `field:"code=AW,decimalPlaces=2,overpunch=true"`
	PercentageSalesTaxAmountPaid                                 *float64 `field:"code=AX,decimalPlaces=2,overpunch=true"`
	PercentageSalesTaxRatePaid                                   *float64 `field:"code=AY,decimalPlaces=4,overpunch=true"`
	PercentageSalesTaxBasisPaid                                  *string  `field:"code=AZ"`
	IncentiveAmountPaid                                          *float64 `field:"code=FL,decimalPlaces=2,overpunch=true"`
	ProfessionalServiceFeePaid                                   *float64 `field:"code=J1,decimalPlaces=2,overpunch=true"`
	OtherAmountPaidCount                                         *int     `field:"code=J2"`
	OtherAmountsPaid                                             []OtherAmountPaid
	OtherPayerAmountRecognized                                   *float64 `field:"code=J5,decimalPlaces=2,overpunch=true"`
	TotalAmountPaid                                              *float64 `field:"code=F9,decimalPlaces=2,overpunch=true"`
	BasisOfReimbursementDetermination                            *int     `field:"code=FM"`
	AmountAttributedToSalesTax                                   *float64 `field:"code=FN,decimalPlaces=2,overpunch=true"`
	AccumulatedDeductibleAmount                                  *float64 `field:"code=FC,decimalPlaces=2,overpunch=true"`
	RemainingDeductibleAmount                                    *float64 `field:"code=FD,decimalPlaces=2,overpunch=true"`
	RemainingBenefitAmount                                       *float64 `field:"code=FE,decimalPlaces=2,overpunch=true"`
	AmountAppliedToPeriodicDeductible                            *float64 `field:"code=FH,decimalPlaces=2,overpunch=true"`
	CopayAmount                                                  *float64 `field:"code=FI,decimalPlaces=2,overpunch=true"`
	AmountExceedingPeriodicBenefitMax                            *float64 `field:"code=FK,decimalPlaces=2,overpunch=true"`
	BasisOfCalculation                                           BasisOfCalculation
	AmountAttributedToProcessorFee                               *float64 `field:"code=NZ,decimalPlaces=2,overpunch=true"`
	PatientSalesTaxAmount                                        *float64 `field:"code=EQ,decimalPlaces=2,overpunch=true"`
	PlanSalesTaxAmount                                           *float64 `field:"code=2Y,decimalPlaces=2,overpunch=true"`
	CoinsuranceAmount                                            *float64 `field:"code=4U,decimalPlaces=2,overpunch=true"`
	BenefitStageCount                                            *int     `field:"code=MU"`
	BenefitStageAmounts                                          []BenefitStage
	EstimatedGenericSavings                                      *float64 `field:"code=G3,decimalPlaces=2,overpunch=true"`
	SpendingAccountAmountRemaining                               *float64 `field:"code=UC,decimalPlaces=2,overpunch=true"`
	HealthPlanFundedAssistanceAmount                             *float64 `field:"code=UD,decimalPlaces=2,overpunch=true"`
	AmountAttributedToNetworkSelection                           *float64 `field:"code=UJ,decimalPlaces=2,overpunch=true"`
	AmountAttributedToBrandProductSelection                      *float64 `field:"code=UK,decimalPlaces=2,overpunch=true"`
	AmountAttributedToNonPreferredFormularyProductSelection      *float64 `field:"code=UM,decimalPlaces=2,overpunch=true"`
	AmountAttributedToBrandNonPreferredFormularyProductSelection *float64 `field:"code=UN,decimalPlaces=2,overpunch=true"`
	AmountAttributedToCoverageGap                                *float64 `field:"code=UP,decimalPlaces=2,overpunch=true"`
	IngredientCostContractedReimbursableAmount                   *float64 `field:"code=U8,decimalPlaces=2,overpunch=true"`
	DispensingFeeContractedReimbursableAmount                    *float64 `field:"code=U9,decimalPlaces=2,overpunch=true"`
}

type OtherAmountPaid struct {
	Qualifier  *string  `field:"code=J3"`
	AmountPaid *float64 `field:"code=J4,decimalPlaces=2,overpunch=true"`
}

type BasisOfCalculation struct {
	DispensingFee      *string `field:"code=HH"`
	Copay              *string `field:"code=HJ"`
	FlatSalesTax       *string `field:"code=HK"`
	PercentageSalesTax *string `field:"code=HM"`
	Coinsurance        *string `field:"code=4V"`
}

type BenefitStage struct {
	Qualifier *string  `field:"code=MV"`
	Amount    *float64 `field:"code=MW,decimalPlaces=2,overpunch=true"`
}
