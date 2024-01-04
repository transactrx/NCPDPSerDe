package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId      ncpdp.SegmentId
	IdQualifier    *string    `field:"code=CX,order=2"`
	Id             *string    `field:"code=CY,order=3"`
	BirthDate      *time.Time `field:"code=C4,format=YYYYMMdd,order=4"`
	GenderCode     *string    `field:"code=C5,order=5"`
	FirstName      *string    `field:"code=CA,order=6"`
	LastName       *string    `field:"code=CB,order=7"`
	Address        PatientAddress
	Phone          *string                 `field:"code=CQ,order=12"`
	PlaceOfService *string                 `field:"code=C7,order=13"`
	EmployerId     *string                 `field:"code=CZ,order=14"`
	SmokerCode     *string                 `field:"code=1C,order=15"`
	Pregnant       *string                 `field:"code=2C,order=16"`
	Email          *string                 `field:"code=HN,order=17"`
	Residence      *string                 `field:"code=4X,order=18"`
	DynamicFields  []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PatientAddress struct {
	Street *string `field:"code=CM,order=8"`
	City   *string `field:"code=CN,order=9"`
	State  *string `field:"code=CO,order=10"`
	Zip    *string `field:"code=CP,order=11"`
}
