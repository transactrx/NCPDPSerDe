package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Pricing struct {
	SegmentId                         ncpdp.SegmentId
	IngredientCostSubmitted           *float64 `field:"code=D9,decimalPlaces=2,overpunch=true,order=2"`
	DispensingFeeSubmitted            *float64 `field:"code=DC,decimalPlaces=2,overpunch=true,order=3"`
	ProfessionalServiceFeeSubmitted   *float64 `field:"code=BE,decimalPlaces=2,overpunch=true,order=4"`
	PatientPaidAmountSubmitted        *float64 `field:"code=DX,decimalPlaces=2,overpunch=true,order=5"`
	IncentiveAmountSubmitted          *float64 `field:"code=E3,decimalPlaces=2,overpunch=true,order=6"`
	OtherAmountClaimSubmittedCount    *int     `field:"code=H7,order=7"`
	OtherAmountClaimSubmitted         []OtherAmountClaimSubmitted
	FlatSalesTaxAmountSubmitted       *float64                `field:"code=HA,decimalPlaces=2,overpunch=true,order=10"`
	PercentageSalesTaxAmountSubmitted *float64                `field:"code=GE,decimalPlaces=2,overpunch=true,order=11"`
	PercentageSalesTaxRateSubmitted   *float64                `field:"code=HE,decimalPlaces=4,overpunch=true,order=12"`
	PercentageSalesTaxBasisSubmitted  *string                 `field:"code=JE,order=13"`
	UsualAndCustmaryCharge            *float64                `field:"code=DQ,decimalPlaces=2,overpunch=true,order=14"`
	GrossAmountDue                    *float64                `field:"code=DU,decimalPlaces=2,overpunch=true,order=15"`
	BasisOfCostDetermination          *string                 `field:"code=DN,order=16"`
	MedicaidPaidAmount                *float64                `field:"code=N3,decimalPlaces=2,overpunch=true,order=17"`
	DynamicFields                     []dynamic.DynamicStruct `field:"code=dynamic"`
}

type OtherAmountClaimSubmitted struct {
	Qualifier       *string  `field:"code=H8,order=8"`
	AmountSubmitted *float64 `field:"code=H9,decimalPlaces=2,overpunch=true,order=9"`
}
