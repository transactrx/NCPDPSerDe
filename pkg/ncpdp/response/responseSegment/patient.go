package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId     ncpdp.SegmentId
	FirstName     *string                 `field:"code=CA,order=2"`
	LastName      *string                 `field:"code=CB,order=3"`
	BirthDate     *time.Time              `field:"code=C4,format=YYYYMMdd,order=4"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
