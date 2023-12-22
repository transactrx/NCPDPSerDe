package responsesegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type Insurance struct {
	SegmentId              ncpdp.SegmentId
	GroupId                *string `field:"code=C1"`
	PlanId                 *string `field:"code=FO"`
	NetworkReimbursementId *string `field:"code=2F"`
	Payer                  Payer
	Medicaid               Medicaid
	CardholderId           *string `field:"code=C2"`
}

type Medicaid struct {
	Id           *string `field:"code=N5"`
	AgencyNumber *string `field:"code=N6"`
}

type Payer struct {
	Qualifier *string `field:"code=J7"`
	Id        *string `field:"code=J8"`
}
