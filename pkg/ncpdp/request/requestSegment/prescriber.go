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
	IdAssociatedState   *string                 `field:"code=ZK,order=14"`
	IdAssociatedCountry *string                 `field:"code=3B,order=15"`
	PhoneExtension      *string                 `field:"code=7T,order=16"`
	MiddleName          *string                 `field:"code=0F,order=17"`
	AlternateId         PrescriberAlternateId
	DeaNumber           *string                 `field:"code=KV,order=25"`
	PlaceOfService      *string                 `field:"code=RG,order=26"`
	DynamicFields       []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PrimaryCareProvider struct {
	IdQualifier *string `field:"code=2E,order=6"`
	Id          *string `field:"code=DL,order=7"`
	LastName    *string `field:"code=4E,order=8"`
}

type PrescriberAddress struct {
	Street      *string `field:"code=2K,order=10"`
	City        *string `field:"code=2M,order=11"`
	State       *string `field:"code=2N,order=12"`
	Zip         *string `field:"code=2P,order=13"`
	StreetLine1 *string `field:"code=7U,order=18"`
	StreetLine2 *string `field:"code=7V,order=19"`
	CountryCode *string `field:"code=3C,order=20"`
}

type PrescriberAlternateId struct {
	Qualifier         *string `field:"code=ZM,order=21"`
	Id                *string `field:"code=ZP,order=22"`
	AssociatedState   *string `field:"code=ZQ,order=23"`
	AssociatedCountry *string `field:"code=3A,order=24"`
}
