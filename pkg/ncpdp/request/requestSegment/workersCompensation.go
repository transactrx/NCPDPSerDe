package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type EmployerAddress struct {
	Street *string `field:"code=CG"`
	City   *string `field:"code=CH"`
	State  *string `field:"code=CI"`
	Zip    *string `field:"code=CJ"`
}

type Employer struct {
	Name        *string `field:"code=CF"`
	Address     EmployerAddress
	Phone       *string `field:"code=CK"`
	ContactName *string `field:"code=CL"`
}

type PayTo struct {
	Qualifier    *string `field:"code=TS"`
	Id           *string `field:"code=TT"`
	Name         *string `field:"code=TU"`
	PayToAddress PayToAddress
}

type PayToAddress struct {
	Street *string `field:"code=TV"`
	City   *string `field:"code=TW"`
	State  *string `field:"code=TX"`
	Zip    *string `field:"code=TY"`
}

type GenericEquivalentProduct struct {
	Qualifier *string `field:"code=TZ"`
	Id        *string `field:"code=UA"`
}

type WorkersCompensation struct {
	SegmentId                  ncpdp.SegmentId
	InjuryDate                 *time.Time `field:"code=DY,format=YYYYMMdd"`
	Employer                   Employer
	CarrierId                  *string `field:"code=CR"`
	ClaimReferenceId           *string `field:"code=DZ"`
	BillingEntityTypeIndicator *string `field:"code=TR"`
	PayTo                      PayTo
	GenericEquivalentProduct   GenericEquivalentProduct
}
