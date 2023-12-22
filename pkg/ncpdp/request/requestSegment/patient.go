package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId      ncpdp.SegmentId
	IdQualifier    *string    `field:"code=CX"`
	Id             *string    `field:"code=CY"`
	BirthDate      *time.Time `field:"code=C4,format=YYYYMMdd"`
	GenderCode     *string    `field:"code=C5"`
	FirstName      *string    `field:"code=CA"`
	LastName       *string    `field:"code=CB"`
	Address        PatientAddress
	Phone          *string `field:"code=CQ"`
	PlaceOfService *string `field:"code=C7"`
	EmployerId     *string `field:"code=CZ"`
	SmokerCode     *string `field:"code=1C"`
	Pregnant       *string `field:"code=2C"`
	Email          *string `field:"code=HN"`
	Residence      *string `field:"code=4X"`
}

type PatientAddress struct {
	Street *string `field:"code=CM"`
	City   *string `field:"code=CN"`
	State  *string `field:"code=CO"`
	Zip    *string `field:"code=CP"`
}
