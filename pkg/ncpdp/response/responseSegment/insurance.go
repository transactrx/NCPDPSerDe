package responsesegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Insurance struct {
	SegmentId              ncpdp.SegmentId `json:"-"`
	GroupId                *string         `field:"code=C1,order=2"`
	PlanId                 *string         `field:"code=FO,order=3"`
	NetworkReimbursementId *string         `field:"code=2F,order=4"`
	PayerIdCount           *int            `field:"code=KR,order=5,countfor=Payers,sinceVersion=F6"`
	Payer                  Payer
	// Payers captures every J7/J8 occurrence when the F6 payer ID repeats; the
	// singular Payer struct is kept for D0 backward compatibility.
	Payers        []Payer `field:"sinceVersion=F6"`
	Medicaid      Medicaid
	CardholderId  *string                 `field:"code=C2,order=10"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type Medicaid struct {
	Id           *string `field:"code=N5,order=8"`
	AgencyNumber *string `field:"code=N6,order=9"`
}

type Payer struct {
	Qualifier *string `field:"code=J7,order=6"`
	Id        *string `field:"code=J8,order=7"`
}
