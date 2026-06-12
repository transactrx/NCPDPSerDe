package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PharmacyProvider struct {
	SegmentId     ncpdp.SegmentId
	IdQualifier   *string                 `field:"code=EY,order=2"`
	Id            *string                 `field:"code=E9,order=3"`
	FirstName     *string                 `field:"code=4A,order=4"`
	LastName      *string                 `field:"code=4M,order=5"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
