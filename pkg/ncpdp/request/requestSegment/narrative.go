package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Narrative struct {
	SegmentId ncpdp.SegmentId
	Message   *string `field:"code=BM"`
}
