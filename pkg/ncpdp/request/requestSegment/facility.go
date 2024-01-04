package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Facility struct {
	SegmentId     ncpdp.SegmentId
	Id            *string `field:"code=8C,order=2"`
	Name          *string `field:"code=3Q,order=3"`
	Address       FacilityAddress
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type FacilityAddress struct {
	Street *string `field:"code=3U,order=4"`
	City   *string `field:"code=5J,order=5"`
	State  *string `field:"code=3V,order=6"`
	Zip    *string `field:"code=6D,order=7"`
}
