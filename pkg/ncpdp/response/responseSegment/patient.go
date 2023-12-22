package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId ncpdp.SegmentId
	FirstName *string    `field:"code=CA"`
	LastName  *string    `field:"code=CB"`
	BirthDate *time.Time `field:"code=C4,format=YYYYMMdd"`
}
