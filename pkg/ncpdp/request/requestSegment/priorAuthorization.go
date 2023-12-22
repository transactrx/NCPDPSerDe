package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PriorAuthorizationRequestAndBilling struct {
	SegmentId                ncpdp.SegmentId
	RequestType              *string    `field:"code=PA"`
	RequestPeriodBeginDate   *time.Time `field:"code=PB,format=YYYYMMdd"`
	RequestPeriodEndDate     *time.Time `field:"code=PC,format=YYYYMMdd"`
	BasisOfRequest           *string    `field:"code=PD"`
	AuthorizedRepresentative AuthorizedRepresentative
	NumberAssigned           *string `field:"code=PY"`
	Number                   *string `field:"code=F3"`
	SupportingDocumentation  *string `field:"code=PP"`
}

type AuthorizedRepresentative struct {
	FirstName                       *string `field:"code=PE"`
	LastName                        *string `field:"code=PF"`
	AuthorizedRepresentativeAddress AuthorizedRepresentativeAddress
}

type AuthorizedRepresentativeAddress struct {
	Street *string `field:"code=PG"`
	City   *string `field:"code=PH"`
	State  *string `field:"code=PJ"`
	Zip    *string `field:"code=PK"`
}
