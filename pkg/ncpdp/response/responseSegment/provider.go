package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Provider struct {
	SegmentId ncpdp.SegmentId

	DataSourceOfInvalidProviderDetermination             *string `field:"code=ZV,order=2"`
	StateCodeForDataSourceOfInvalidProviderDetermination *string `field:"code=ZZ,order=3"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
