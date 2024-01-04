package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Coupon struct {
	SegmentId     ncpdp.SegmentId
	Type          *string                 `field:"code=KE,order=1"`
	Number        *string                 `field:"code=ME,order=2"`
	Amount        *float64                `field:"code=NE,decimalPlaces=2,overpunch=true,order=3"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
