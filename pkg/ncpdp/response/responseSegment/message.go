package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Message struct {
	SegmentId     ncpdp.SegmentId
	Message       *string                 `field:"code=F4,order=2"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
