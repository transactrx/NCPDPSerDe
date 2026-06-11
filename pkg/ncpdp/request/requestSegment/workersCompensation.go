package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type WorkersCompensation struct {
	SegmentId                  ncpdp.SegmentId
	InjuryDate                 *time.Time `field:"code=DY,format=YYYYMMdd,order=2"`
	Employer                   Employer
	CarrierId                  *string `field:"code=CR,order=10"`
	ClaimReferenceId           *string `field:"code=DZ,order=11"`
	BillingEntityTypeIndicator *string `field:"code=TR,order=12"`
	PayTo                      PayTo
	GenericEquivalentProduct   GenericEquivalentProduct
	DynamicFields              []dynamic.DynamicStruct `field:"code=dynamic"`
}

type EmployerAddress struct {
	Street      *string `field:"code=CG,order=4"`
	StreetLine1 *string `field:"code=8D,order=22"`
	StreetLine2 *string `field:"code=7G,order=23"`
	City        *string `field:"code=CH,order=5"`
	State       *string `field:"code=CI,order=6"`
	Zip         *string `field:"code=CJ,order=7"`
	CountryCode *string `field:"code=1V,order=24"`
}

type Employer struct {
	Name             *string `field:"code=CF,order=3"`
	Address          EmployerAddress
	Phone            *string `field:"code=CK,order=8"`
	PhoneExtension   *string `field:"code=7K,order=25"`
	ContactName      *string `field:"code=CL,order=9"`
	ContactFirstName *string `field:"code=7H,order=26"`
	ContactLastName  *string `field:"code=7J,order=27"`
}

type PayTo struct {
	Qualifier    *string `field:"code=TS,order=13"`
	Id           *string `field:"code=TT,order=14"`
	Name         *string `field:"code=TU,order=15"`
	PayToAddress PayToAddress
}

type PayToAddress struct {
	Street      *string `field:"code=TV,order=16"`
	StreetLine1 *string `field:"code=7R,order=28"`
	StreetLine2 *string `field:"code=7S,order=29"`
	City        *string `field:"code=TW,order=17"`
	State       *string `field:"code=TX,order=18"`
	Zip         *string `field:"code=TY,order=19"`
	CountryCode *string `field:"code=1Z,order=30"`
}

type GenericEquivalentProduct struct {
	Qualifier *string `field:"code=TZ,order=20"`
	Id        *string `field:"code=UA,order=21"`
}
