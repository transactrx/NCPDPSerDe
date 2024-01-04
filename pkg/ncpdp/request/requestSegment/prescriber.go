package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Prescriber struct {
	SegmentId           ncpdp.SegmentId
	IdQualifier         *string `field:"code=EZ,order=2"`
	Id                  *string `field:"code=DB,order=3"`
	LastName            *string `field:"code=DR,order=4"`
	FirstName           *string `field:"code=2J,order=9"`
	Phone               *string `field:"code=PM,order=5"`
	Address             PrescriberAddress
	PrimaryCareProvider PrimaryCareProvider
	DynamicFields       []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PrimaryCareProvider struct {
	IdQualifier *string `field:"code=2E,order=6"`
	Id          *string `field:"code=DL,order=7"`
	LastName    *string `field:"code=4E,order=8"`
}

type PrescriberAddress struct {
	Street *string `field:"code=2K,order=10"`
	City   *string `field:"code=2M,order=11"`
	State  *string `field:"code=2N,order=12"`
	Zip    *string `field:"code=2P,order=13"`
}
