package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type OtherAmountClaimSubmitted struct {
	Qualifier       *string  `field:"code=H8"`
	AmountSubmitted *float64 `field:"code=H9,decimalPlaces=2,overpunch=true"`
}

type Pricing struct {
	SegmentId                         ncpdp.SegmentId
	IngredientCostSubmitted           *float64 `field:"code=D9,decimalPlaces=2,overpunch=true"`
	DispensingFeeSubmitted            *float64 `field:"code=DC,decimalPlaces=2,overpunch=true"`
	ProfessionalServiceFeeSubmitted   *float64 `field:"code=BE,decimalPlaces=2,overpunch=true"`
	PatientPaidAmountSubmitted        *float64 `field:"code=DX,decimalPlaces=2,overpunch=true"`
	IncentiveAmountSubmitted          *float64 `field:"code=E3,decimalPlaces=2,overpunch=true"`
	OtherAmountClaimSubmittedCount    *int     `field:"code=H7"`
	OtherAmountClaimSubmitted         []OtherAmountClaimSubmitted
	FlatSalesTaxAmountSubmitted       *float64 `field:"code=HA,decimalPlaces=2,overpunch=true"`
	PercentageSalesTaxAmountSubmitted *float64 `field:"code=GE,decimalPlaces=2,overpunch=true"`
	PercentageSalesTaxRateSubmitted   *float64 `field:"code=HE,decimalPlaces=4,overpunch=true"`
	PercentageSalesTaxBasisSubmitted  *string  `field:"code=JE"`
	UsualAndCustmaryCharge            *float64 `field:"code=DQ,decimalPlaces=2,overpunch=true"`
	GrossAmountDue                    *float64 `field:"code=DU,decimalPlaces=2,overpunch=true"`
	BasisOfCostDetermination          *string  `field:"code=DN"`
	MedicaidPaidAmount                *float64 `field:"code=N3,decimalPlaces=2,overpunch=true"`
}
