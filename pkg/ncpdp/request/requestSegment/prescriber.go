package requestsegment

import "github.com/transactrx/NCPDPSerDe/pkg/ncpdp"

type PrimaryCareProvider struct {
	IdQualifier *string `field:"code=2E"`
	Id          *string `field:"code=DL"`
	LastName    *string `field:"code=4E"`
}

type PrescriberAddress struct {
	Street *string `field:"code=2K"`
	City   *string `field:"code=2M"`
	State  *string `field:"code=2N"`
	Zip    *string `field:"code=2P"`
}

type Prescriber struct {
	SegmentId           ncpdp.SegmentId
	IdQualifier         *string `field:"code=EZ"`
	Id                  *string `field:"code=DB"`
	LastName            *string `field:"code=DR"`
	FirstName           *string `field:"code=2J"`
	Phone               *string `field:"code=PM"`
	Address             PrescriberAddress
	PrimaryCareProvider PrimaryCareProvider
}
