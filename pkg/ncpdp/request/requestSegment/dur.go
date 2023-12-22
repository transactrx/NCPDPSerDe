package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type DurCoAgent struct {
	Qualifier *string `field:"code=J9"`
	Id        *string `field:"code=H6"`
}

type DurItem struct {
	Counter                 *int    `field:"code=7E"`
	ReasonForServiceCode    *string `field:"code=E4"`
	ProfessionalServiceCode *string `field:"code=E5"`
	ResultOfServiceCode     *string `field:"code=E6"`
	LevelOfEffort           *string `field:"code=8E"`
	CoAgent                 DurCoAgent
}

type Dur struct {
	SegmentId ncpdp.SegmentId

	Items []DurItem
}
