package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Facility struct {
	SegmentId     ncpdp.SegmentId `json:"-"`
	IdQualifier   *string         `field:"code=3Z,order=2"`
	Id            *string         `field:"code=8C,order=3"`
	Name          *string         `field:"code=3Q,order=4"`
	Address       FacilityAddress
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type FacilityAddress struct {
	Street      *string `field:"code=3U,order=5"`
	StreetLine1 *string `field:"code=7M,order=6"`
	StreetLine2 *string `field:"code=7N,order=7"`
	City        *string `field:"code=5J,order=8"`
	State       *string `field:"code=3V,order=9"`
	Zip         *string `field:"code=6D,order=10"`
	CountryCode *string `field:"code=1X,order=11"`
}
