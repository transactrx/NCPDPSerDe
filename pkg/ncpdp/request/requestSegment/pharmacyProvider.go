package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type PharmacyProvider struct {
	SegmentId   ncpdp.SegmentId
	IdQualifier *string `field:"code=EY"`
	Id          *string `field:"code=E9"`
}
