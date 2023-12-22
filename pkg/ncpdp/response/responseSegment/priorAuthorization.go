package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PriorAuthorization struct {
	SegmentId           ncpdp.SegmentId
	ProcessedDate       *time.Time `field:"code=PR,format=YYYYMMdd"`
	EffectiveDate       *time.Time `field:"code=PS,format=YYYYMMdd"`
	ExpirationDate      *time.Time `field:"code=PT,format=YYYYMMdd"`
	Quantity            *float64   `field:"code=RA,decimalPlaces=3"`
	AmountAuthorized    *float64   `field:"code=RB,decimalPlaces=2,overpunch=true"`
	RefillsAuthorized   *int       `field:"code=PW"`
	QuantityAccumulated *float64   `field:"code=PX,decimalPlaces=3"`
	NumberAssigned      *string    `field:"code=PY"`
}
