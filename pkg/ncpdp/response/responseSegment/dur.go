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
}
