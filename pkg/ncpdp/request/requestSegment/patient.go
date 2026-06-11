package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Patient struct {
	SegmentId              ncpdp.SegmentId
	IdCount                *int    `field:"code=RR,order=2"`
	IdQualifier            *string `field:"code=CX,order=3"`
	Id                     *string `field:"code=CY,order=4"`
	// Ids captures every CX/CY occurrence when the F6 patient ID repeats; the
	// scalar IdQualifier/Id fields are kept for D0 backward compatibility.
	Ids                    []PatientId
	BirthDate              *time.Time `field:"code=C4,format=YYYYMMdd,order=5"`
	GenderCode             *string    `field:"code=C5,order=6"`
	FirstName              *string    `field:"code=CA,order=7"`
	MiddleName             *string    `field:"code=0C,order=8"`
	LastName               *string    `field:"code=CB,order=9"`
	NameSuffix             *string    `field:"code=0E,order=10"`
	NamePrefix             *string    `field:"code=0D,order=11"`
	Address                PatientAddress
	Phone                  *string                 `field:"code=CQ,order=19"`
	PlaceOfService         *string                 `field:"code=C7,order=20"`
	EmployerId             *string                 `field:"code=CZ,order=21"`
	SmokerCode             *string                 `field:"code=1C,order=22"`
	Pregnant               *string                 `field:"code=2C,order=23"`
	Email                  *string                 `field:"code=HN,order=24"`
	Residence              *string                 `field:"code=4X,order=25"`
	IdAssociatedState      *string                 `field:"code=YR,order=26"`
	IdAssociatedCountry    *string                 `field:"code=1Y,order=27"`
	VeterinaryUseIndicator *string                 `field:"code=1R,order=28"`
	Species                *string                 `field:"code=S8,order=29"`
	DynamicFields          []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PatientId struct {
	Qualifier *string `field:"code=CX,order=3"`
	Id        *string `field:"code=CY,order=4"`
}

type PatientAddress struct {
	Street      *string `field:"code=CM,order=12"`
	StreetLine1 *string `field:"code=7A,order=13"`
	StreetLine2 *string `field:"code=7B,order=14"`
	City        *string `field:"code=CN,order=15"`
	State       *string `field:"code=CO,order=16"`
	Zip         *string `field:"code=CP,order=17"`
	CountryCode *string `field:"code=1K,order=18"`
}
