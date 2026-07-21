package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Coupon struct {
	SegmentId     ncpdp.SegmentId         `json:"-"`
	Type          *string                 `field:"code=KE,order=2"`
	Number        *string                 `field:"code=ME,order=3"`
	Amount        *float64                `field:"code=NE,decimalPlaces=2,overpunch=true,order=4"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
