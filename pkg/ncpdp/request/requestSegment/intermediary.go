package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Intermediary struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	IdCount *int `field:"code=8G,order=2,countfor=Ids,sinceVersion=F6"`
	Ids     []IntermediaryId

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type IntermediaryId struct {
	TypeCode             *string `field:"code=8H,order=3,sinceVersion=F6"`
	TypeEntity           *string `field:"code=8J,order=4,sinceVersion=F6"`
	Qualifier            *string `field:"code=8K,order=5,sinceVersion=F6"`
	Id                   *string `field:"code=8M,order=6,sinceVersion=F6"`
	StateProvinceAddress *string `field:"code=8N,order=7,sinceVersion=F6"`
	CountryCode          *string `field:"code=8U,order=8,sinceVersion=F6"`
}
