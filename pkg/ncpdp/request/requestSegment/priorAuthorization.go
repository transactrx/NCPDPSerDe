package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PriorAuthorizationRequestAndBilling struct {
	SegmentId                ncpdp.SegmentId `json:"-"`
	RequestType              *string         `field:"code=PA,order=2"`
	RequestPeriodBeginDate   *time.Time      `field:"code=PB,format=YYYYMMdd,order=3"`
	RequestPeriodEndDate     *time.Time      `field:"code=PC,format=YYYYMMdd,order=4"`
	BasisOfRequest           *string         `field:"code=PD,order=5"`
	AuthorizedRepresentative AuthorizedRepresentative
	NumberAssigned           *string                 `field:"code=PY,order=12"`
	Number                   *string                 `field:"code=F3,order=13"`
	SupportingDocumentation  *string                 `field:"code=PP,order=14"`
	DynamicFields            []dynamic.DynamicStruct `field:"code=dynamic"`
}

type AuthorizedRepresentative struct {
	FirstName                       *string `field:"code=PE,order=6"`
	LastName                        *string `field:"code=PF,order=7"`
	AuthorizedRepresentativeAddress AuthorizedRepresentativeAddress
}

type AuthorizedRepresentativeAddress struct {
	Street      *string `field:"code=PG,order=8"`
	StreetLine1 *string `field:"code=7D,order=15,sinceVersion=F6"`
	StreetLine2 *string `field:"code=8B,order=16,sinceVersion=F6"`
	City        *string `field:"code=PH,order=9"`
	State       *string `field:"code=PJ,order=10"`
	Zip         *string `field:"code=PK,order=11"`
	CountryCode *string `field:"code=1U,order=17,sinceVersion=F6"`
}
