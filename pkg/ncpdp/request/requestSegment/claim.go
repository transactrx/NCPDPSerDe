package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type PrescriptionServiceReference struct {
	Qualifier *string `field:"code=EM"`
	Number    *string `field:"code=D2"`
}

type AssociatedPrescriptionService struct {
	Number *string    `field:"code=EN"`
	Date   *time.Time `field:"code=EP,format=YYYYMMdd"`
}

type ProductService struct {
	Qualifier *string `field:"code=E1"`
	Id        *string `field:"code=D7"`
}

type OriginallyPrescribedProductService struct {
	Qualifier *string  `field:"code=EJ"`
	Code      *string  `field:"code=EA"`
	Quantity  *float64 `field:"code=EB,decimalPlaces=3"`
}

type ClaimPriorAuthorization struct {
	TypeCode        *string `field:"code=EU"`
	NumberSubmitted *string `field:"code=EV"`
}

type IntermediaryAuthorization struct {
	TypeId *string `field:"code=EW"`
	Id     *string `field:"code=EX"`
}

type ProcedureModifierCode struct {
	Code *string `field:"code=ER"`
}

type SubmissionClarificationCode struct {
	Code *string `field:"code=DK"`
}

type Claim struct {
	SegmentId                          ncpdp.SegmentId
	PrescriptionServiceReference       PrescriptionServiceReference
	ProductService                     ProductService
	AssociatedPrescriptionService      AssociatedPrescriptionService
	ProcedureModifierCodeCount         *int `field:"code=SE"`
	ProcedureModifierCodes             []ProcedureModifierCode
	QuantityDispensed                  *float64   `field:"code=E7,decimalPlaces=3"`
	FillNumber                         *string    `field:"code=D3"`
	DaysSupply                         *int       `field:"code=D5"`
	CompoundCode                       *string    `field:"code=D6"`
	DispenseAsWritten                  *string    `field:"code=D8"`
	DateWritten                        *time.Time `field:"code=DE,format=YYYYMMdd"`
	RefillsAuthorized                  *int       `field:"code=DF"`
	OriginCode                         *string    `field:"code=DJ"`
	SubmissionClarificationCodeCount   *int       `field:"code=NX"`
	SubmissionClarificationCodes       []SubmissionClarificationCode
	QuantityPrescribed                 *float64 `field:"code=ET,decimalPlaces=3"`
	OtherCoverageCode                  *string  `field:"code=C8"`
	SpecialPackagingIndicator          *string  `field:"code=DT"`
	OriginallyPrescribedProductService OriginallyPrescribedProductService
	AlternateId                        *string `field:"code=CW"`
	ScheduledPrescriptionId            *string `field:"code=EK"`
	UnitOfMeasure                      *string `field:"code=28"`
	LevelOfService                     *string `field:"code=DI"`
	PriorAuthorization                 ClaimPriorAuthorization
	IntermediaryAuthorization          IntermediaryAuthorization
	DispensingStatus                   *string  `field:"code=HD"`
	IntendedQuantityDispensed          *float64 `field:"code=HF,decimalPlaces=3"`
	IntendedDaysSupply                 *int     `field:"code=HG"`
	DelayReasonCode                    *string  `field:"code=NV"`
	TransactionReferenceNumber         *string  `field:"code=K5"`
	PatientAssignmentIndicator         *string  `field:"code=MT"`
	RouteOfAdministration              *string  `field:"code=E2"`
	CompoundType                       *string  `field:"code=G1"`
	MedicaidSubrogationControlNumber   *string  `field:"code=N4"`
	PharmacyServiceType                *string  `field:"code=U7"`
}
