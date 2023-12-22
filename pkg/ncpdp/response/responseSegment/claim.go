package responsesegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Claim struct {
	SegmentId                                ncpdp.SegmentId
	PrescriptionServiceReference             PrescriptionServiceReference
	PreferredProductCount                    *int `field:"code=9F"`
	PreferredProducts                        []PreferredProduct
	MedicaidSubrogationInternalControlNumber *string `field:"code=N4"`
}

type PrescriptionServiceReference struct {
	Qualifier *string `field:"code=EM"`
	Number    *string `field:"code=D2"`
}

type PreferredProduct struct {
	IdQualifier        *string  `field:"code=AP"`
	Id                 *string  `field:"code=AR"`
	Incentive          *float64 `field:"code=AS,decimalPlaces=2,overpunch=true"`
	CostShareIncentive *float64 `field:"code=AT,decimalPlaces=2,overpunch=true"`
	Description        *string  `field:"code=AU"`
}
