package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Coupon struct {
	SegmentId ncpdp.SegmentId

	Type   *string  `field:"code=KE"`
	Number *string  `field:"code=ME"`
	Amount *float64 `field:"code=NE,decimalPlaces=2,overpunch=true"`
}
