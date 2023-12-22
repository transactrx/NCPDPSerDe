package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Dur struct {
	SegmentId ncpdp.SegmentId

	Items []DurItem
}

type DurItem struct {
	Counter                  *int       `field:"code=J6"`
	ReasonForServiceCode     *string    `field:"code=E4"`
	ClinicalSignificanceCode *string    `field:"code=FS"`
	OtherPharmacyIndicator   *string    `field:"code=FT"`
	PreviousFillDate         *time.Time `field:"code=FU,format=YYYYMMdd"`
	PreviousFillDateQuantity *float64   `field:"code=FV,decimalPlaces=3"`
	DatabaseIndicator        *string    `field:"code=FW"`
	OtherPrescriberIndicator *string    `field:"code=FX"`
	FreeText                 *string    `field:"code=FY"`
	AdditionalText           *string    `field:"code=NS"`
}
