package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type LastKnown4Rx struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	IinNumber              *string `field:"code=3E,order=2,sinceVersion=F6"`
	ProcessorControlNumber *string `field:"code=3F,order=3,sinceVersion=F6"`
	GroupId                *string `field:"code=3G,order=4,sinceVersion=F6"`
	CardholderId           *string `field:"code=3H,order=5,sinceVersion=F6"`
	YearOfLastPaidClaim    *string `field:"code=3J,order=6,sinceVersion=F6"`
	MonthOfLastPaidClaim   *string `field:"code=3K,order=7,sinceVersion=F6"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
