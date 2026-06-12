package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Dur struct {
	SegmentId ncpdp.SegmentId

	Items []DurItem

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type DurItem struct {
	Counter                  *int       `field:"code=J6,order=2"`
	ReasonForServiceCode     *string    `field:"code=E4,order=3"`
	ClinicalSignificanceCode *string    `field:"code=FS,order=4"`
	OtherPharmacyIndicator   *string    `field:"code=FT,order=5"`
	PreviousFillDate         *time.Time `field:"code=FU,format=YYYYMMdd,order=6"`
	PreviousFillDateQuantity *float64   `field:"code=FV,decimalPlaces=3,order=7"`
	DatabaseIndicator        *string    `field:"code=FW,order=8"`
	OtherPrescriberIndicator *string    `field:"code=FX,order=9"`
	FreeText                 *string    `field:"code=FY,order=10"`
	AdditionalText           *string    `field:"code=NS,order=11"`

	CoAgentIdQualifier *string `field:"code=J9,order=12"`
	CoAgentId          *string `field:"code=H6,order=13"`
	CoAgentDescription *string `field:"code=ZC,order=14"`

	OtherPharmacyIdQualifier *string `field:"code=Z9,order=15"`
	OtherPharmacyId          *string `field:"code=Z8,order=16"`
	OtherPharmacyName        *string `field:"code=Z7,order=17"`
	OtherPharmacyTelephone   *string `field:"code=Z6,order=18"`

	UnitOfMeasureForPreviousDispensedQuantity *string `field:"code=ZA,order=19"`

	OtherPrescriberIdQualifier     *string `field:"code=Z4,order=20"`
	OtherPrescriberId              *string `field:"code=Z3,order=21"`
	OtherPrescriberLastName        *string `field:"code=Z5,order=22"`
	OtherPrescriberTelephoneNumber *string `field:"code=Z2,order=23"`

	CompoundProductIdQualifier *string `field:"code=Z0,order=24"`
	CompoundProductId          *string `field:"code=Z1,order=25"`

	MaximumDailyDoseQuantity      *float64 `field:"code=YO,decimalPlaces=3,order=26"`
	MaximumDailyDoseUnitOfMeasure *string  `field:"code=YL,order=27"`
	MinimumDailyDoseQuantity      *float64 `field:"code=YJ,decimalPlaces=3,order=28"`
	MinimumDailyDoseUnitOfMeasure *string  `field:"code=YI,order=29"`
}
