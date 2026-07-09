package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Claim struct {
	SegmentId                          ncpdp.SegmentId
	PrescriptionServiceReference       PrescriptionServiceReference
	ProductService                     ProductService
	AssociatedPrescriptionService      AssociatedPrescriptionService
	ProcedureModifierCodeCount         *int `field:"code=SE,order=8"`
	ProcedureModifierCodes             []ProcedureModifierCode
	QuantityDispensed                  *float64   `field:"code=E7,decimalPlaces=3,order=10"`
	FillNumber                         *string    `field:"code=D3,order=11"`
	DaysSupply                         *int       `field:"code=D5,order=12"`
	CompoundCode                       *string    `field:"code=D6,order=13"`
	DispenseAsWritten                  *string    `field:"code=D8,order=14"`
	DateWritten                        *time.Time `field:"code=DE,format=YYYYMMdd,order=15"`
	RefillsAuthorized                  *int       `field:"code=DF,order=16"`
	OriginCode                         *string    `field:"code=DJ,order=17"`
	SubmissionClarificationCodeCount   *int       `field:"code=NX,order=18"`
	SubmissionClarificationCodes       []SubmissionClarificationCode
	QuantityPrescribed                 *float64 `field:"code=ET,decimalPlaces=3,order=20"`
	OtherCoverageCode                  *string  `field:"code=C8,order=21"`
	SpecialPackagingIndicator          *string  `field:"code=DT,order=22"`
	OriginallyPrescribedProductService OriginallyPrescribedProductService
	AlternateId                        *string `field:"code=CW,order=26"`
	ScheduledPrescriptionId            *string `field:"code=EK,order=27"`
	UnitOfMeasure                      *string `field:"code=28,order=28"`
	LevelOfService                     *string `field:"code=DI,order=29"`
	PriorAuthorization                 ClaimPriorAuthorization
	IntermediaryAuthorization          IntermediaryAuthorization
	DispensingStatus                   *string  `field:"code=HD,order=34"`
	IntendedQuantityDispensed          *float64 `field:"code=HF,decimalPlaces=3,order=35"`
	IntendedDaysSupply                 *int     `field:"code=HG,order=36"`
	DelayReasonCode                    *string  `field:"code=NV,order=37"`
	TransactionReferenceNumber         *string  `field:"code=K5,order=38"`
	PatientAssignmentIndicator         *string  `field:"code=MT,order=39"`
	RouteOfAdministration              *string  `field:"code=E2,order=40"`
	CompoundType                       *string  `field:"code=G1,order=41"`
	MedicaidSubrogationControlNumber   *string  `field:"code=N4,order=42"`
	PharmacyServiceType                *string  `field:"code=U7,order=43"`
	ReconciliationId                   *string  `field:"code=34,order=44"`
	SubmissionTypeCodeCount            *int     `field:"code=KZ,order=49"`
	SubmissionTypeCodes                []SubmissionTypeCode
	MultipleOrderGroupId               *string  `field:"code=M3,order=51"`
	MultipleOrderGroupReasonCode       *string  `field:"code=M4,order=52"`
	TotalPrescribedQuantityRemaining   *float64 `field:"code=KW,decimalPlaces=3,order=53"`
	PreparationEnvironmentType         *string  `field:"code=KU,order=54"`
	PreparationEnvironmentEventCode    *string  `field:"code=KT,order=55"`
	TimeOfService                      *string  `field:"code=Y6,order=56"`
	SalesTransactionId                 *string  `field:"code=ZF,order=57"`
	ReportedAdjudicatedProgramType     *string  `field:"code=ZS,order=58"`
	OriginalManufacturerProduct        OriginalManufacturerProduct
	LTPACDispenseFrequency             *int                    `field:"code=KK,order=61"`
	LTPACBillingMethodology            *string                 `field:"code=KH,order=62"`
	NumberOfLTPACDispensingEvents      *int                    `field:"code=KM,order=63"`
	DoNotDispenseBeforeDate            *time.Time              `field:"code=K9,format=YYYYMMdd,order=64"`
	DynamicFields                      []dynamic.DynamicStruct `field:"code=dynamic"`
}

type PrescriptionServiceReference struct {
	Qualifier *string `field:"code=EM,order=2"`
	Number    *string `field:"code=D2,order=3"`
}

type AssociatedPrescriptionService struct {
	Number                   *string    `field:"code=EN,order=6"`
	Date                     *time.Time `field:"code=EP,format=YYYYMMdd,order=7"`
	ReferenceNumberQualifier *string    `field:"code=XZ,order=45"`
	ProviderIdQualifier      *string    `field:"code=XX,order=46"`
	ProviderId               *string    `field:"code=XY,order=47"`
	FillNumber               *string    `field:"code=X0,order=48"`
}

type ProductService struct {
	Qualifier *string `field:"code=E1,order=4"`
	Id        *string `field:"code=D7,order=5"`
}

type OriginallyPrescribedProductService struct {
	Qualifier *string  `field:"code=EJ,order=23"`
	Code      *string  `field:"code=EA,order=24"`
	Quantity  *float64 `field:"code=EB,decimalPlaces=3,order=25"`
}

type ClaimPriorAuthorization struct {
	TypeCode        *string `field:"code=EU,order=30"`
	NumberSubmitted *string `field:"code=EV,order=31"`
}

type IntermediaryAuthorization struct {
	TypeId *string `field:"code=EW,order=32"`
	Id     *string `field:"code=EX,order=33"`
}

type ProcedureModifierCode struct {
	Code *string `field:"code=ER,order=9"`
}

type SubmissionClarificationCode struct {
	Code *string `field:"code=DK,order=19"`
}

type SubmissionTypeCode struct {
	Code *string `field:"code=K8,order=50"`
}

type OriginalManufacturerProduct struct {
	Qualifier *string `field:"code=4P,order=59"`
	Id        *string `field:"code=4N,order=60"`
}
