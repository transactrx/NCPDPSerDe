package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type FacilityAddress struct {
	Street *string `field:"code=3U"`
	City   *string `field:"code=5J"`
	State  *string `field:"code=3V"`
	Zip    *string `field:"code=6D"`
}

type Facility struct {
	SegmentId ncpdp.SegmentId
	Id        *string `field:"code=8C"`
	Name      *string `field:"code=3Q"`
	Address   FacilityAddress
}
