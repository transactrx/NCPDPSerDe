package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Dur struct {
	SegmentId ncpdp.SegmentId

	Items []DurItem

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type DurCoAgent struct {
	Qualifier *string `field:"code=J9,order=7"`
	Id        *string `field:"code=H6,order=8"`
}

type DurItem struct {
	Counter                 *int    `field:"code=7E,order=2"`
	ReasonForServiceCode    *string `field:"code=E4,order=3"`
	ProfessionalServiceCode *string `field:"code=E5,order=4"`
	ResultOfServiceCode     *string `field:"code=E6,order=5"`
	LevelOfEffort           *string `field:"code=8E,order=6"`
	CoAgent                 DurCoAgent
}
