package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PriorAuthorization struct {
	SegmentId           ncpdp.SegmentId
	ProcessedDate       *time.Time              `field:"code=PR,format=YYYYMMdd,order=2"`
	EffectiveDate       *time.Time              `field:"code=PS,format=YYYYMMdd,order=3"`
	ExpirationDate      *time.Time              `field:"code=PT,format=YYYYMMdd,order=4"`
	Quantity            *float64                `field:"code=RA,decimalPlaces=3,order=5"`
	AmountAuthorized    *float64                `field:"code=RB,decimalPlaces=2,overpunch=true,order=6"`
	RefillsAuthorized   *int                    `field:"code=PW,order=7"`
	QuantityAccumulated *float64                `field:"code=PX,decimalPlaces=3,order=8"`
	NumberAssigned      *string                 `field:"code=PY,order=9"`
	DynamicFields       []dynamic.DynamicStruct `field:"code=dynamic"`
}
