package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Insurance struct {
	SegmentId              ncpdp.SegmentId
	GroupId                *string `field:"code=C1,order=2"`
	PlanId                 *string `field:"code=FO,order=3"`
	NetworkReimbursementId *string `field:"code=2F,order=4"`
	Payer                  Payer
	Medicaid               Medicaid
	CardholderId           *string                 `field:"code=C2,order=9"`
	DynamicFields          []dynamic.DynamicStruct `field:"code=dynamic"`
}

type Medicaid struct {
	Id           *string `field:"code=N5,order=7"`
	AgencyNumber *string `field:"code=N6,order=8"`
}

type Payer struct {
	Qualifier *string `field:"code=J7,order=5"`
	Id        *string `field:"code=J8,order=6"`
}
