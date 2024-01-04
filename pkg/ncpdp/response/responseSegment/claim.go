package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Claim struct {
	SegmentId                                ncpdp.SegmentId
	PrescriptionServiceReference             PrescriptionServiceReference
	PreferredProductCount                    *int `field:"code=9F,order=4"`
	PreferredProducts                        []PreferredProduct
	MedicaidSubrogationInternalControlNumber *string                 `field:"code=N4,order=10"`
	DynamicFields                            []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PrescriptionServiceReference struct {
	Qualifier *string `field:"code=EM,order=2"`
	Number    *string `field:"code=D2,order=3"`
}

type PreferredProduct struct {
	IdQualifier        *string  `field:"code=AP,order=5"`
	Id                 *string  `field:"code=AR,order=6"`
	Incentive          *float64 `field:"code=AS,decimalPlaces=2,overpunch=true,order=7"`
	CostShareIncentive *float64 `field:"code=AT,decimalPlaces=2,overpunch=true,order=8"`
	Description        *string  `field:"code=AU,order=9"`
}
