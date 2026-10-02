package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	IdCount *int `field:"code=RR,order=2,countfor=Ids,sinceVersion=F6"`
	Ids     []PatientId

	FirstName     *string                 `field:"code=CA,order=5"`
	LastName      *string                 `field:"code=CB,order=6"`
	BirthDate     *time.Time              `field:"code=C4,format=YYYYMMdd,order=7"`
	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PatientId struct {
	Qualifier *string `field:"code=CX,order=3,sinceVersion=F6"`
	Id        *string `field:"code=CY,order=4,sinceVersion=F6"`
}
