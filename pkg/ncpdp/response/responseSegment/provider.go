package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Provider struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	DataSourceOfInvalidProviderDetermination             *string `field:"code=ZV,order=2,sinceVersion=F6"`
	StateCodeForDataSourceOfInvalidProviderDetermination *string `field:"code=ZZ,order=3,sinceVersion=F6"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
