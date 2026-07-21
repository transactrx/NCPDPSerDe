package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type LastKnown4Rx struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	IinNumber              *string `field:"code=3E,order=2"`
	ProcessorControlNumber *string `field:"code=3F,order=3"`
	GroupId                *string `field:"code=3G,order=4"`
	CardholderId           *string `field:"code=3H,order=5"`
	YearOfLastPaidClaim    *string `field:"code=3J,order=6"`
	MonthOfLastPaidClaim   *string `field:"code=3K,order=7"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
