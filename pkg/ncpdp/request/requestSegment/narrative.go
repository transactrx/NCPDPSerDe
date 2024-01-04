package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Narrative struct {
	SegmentId     ncpdp.SegmentId
	Message       *string                 `field:"code=BM,order=2"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
