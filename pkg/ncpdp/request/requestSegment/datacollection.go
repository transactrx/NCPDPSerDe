package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type DataCollection struct {
	SegmentId                                   ncpdp.SegmentId
	InventoryCategoryCode                       *string                 `field:"code=&A,order=2"`
	PatientId                                   *string                 `field:"code=&B,order=3"` //Not de-identified
	PayerSequence                               *string                 `field:"code=&C,order=4"`
	SubstitutionType                            *string                 `field:"code=&D,order=5"`
	DispensingType                              *string                 `field:"code=&E,order=6"`
	CompanyNumber                               *string                 `field:"code=&F,order=7"`
	PayPlan                                     *string                 `field:"code=&G,order=8"`
	PatientLocationCode                         *string                 `field:"code=&H,order=9"`
	ClaimIndicator                              *string                 `field:"code=&I,order=10"`
	Facility                                    *string                 `field:"code=&J,order=11"`
	TransactionTime                             *time.Time              `field:"code=&K,format=HHmmss,order=12"`
	MedicarePartDPriceCode                      *string                 `field:"code=&L,order=13"`
	TestTransaction                             *string                 `field:"code=&M,order=14"`
	PharmacyNpi                                 *string                 `field:"code=&N,order=15"`
	CycleFillListIndictor                       *string                 `field:"code=&T,order=16"`
	AcquisitionCost                             *float64                `field:"code=&U,decimalPlaces=2,order=17"`
	SellingPrice                                *float64                `field:"code=&V,decimalPlaces=2,order=18"`
	PaymentType                                 *string                 `field:"code=&W,order=19"`
	NursingHome                                 *string                 `field:"code=&X,order=20"`
	PatientCodeDeidentified                     *string                 `field:"code=&Y,order=21"`
	PrescriptionStopDate                        *time.Time              `field:"code=&Z,format=YYYYMMdd,order=22"`
	PriceCodeType                               *string                 `field:"code=#A,order=23"`
	PatientLanguage                             *string                 `field:"code=#B,order=24"`
	PrescriberDea                               *string                 `field:"code=#C,order=25"`
	SecondaryTernaryIndicator                   *string                 `field:"code=#D,order=26"`
	DecimalQuantity                             *float64                `field:"code=#E,order=27"`
	PrescriberStateId                           *string                 `field:"code=#F,order=28"`
	PrescribedNdc                               *string                 `field:"code=#G,order=29"`
	ElectronicRxTransactionReferenceNumber      *string                 `field:"code=#H,order=30"` //UIB-030-01
	ElectronicRxInitiatorControlReferenceNumber *string                 `field:"code=#I,order=31"` //UIH-030-01
	PharmacyDea                                 *string                 `field:"code=#K,order=32"`
	PrescriberNpi                               *string                 `field:"code=#L,order=33"`
	DrugAWP                                     *float64                `field:"code=#M,decimalPlaces=2,order=34"`
	ClassificationCode                          *string                 `field:"code=#N,order=35"`
	PHIParticipation                            *string                 `field:"code=#O,order=36"`
	MarketingParticipation                      *string                 `field:"code=#P,order=37"`
	PatientSocialSecurityNumber                 *string                 `field:"code=#Q,order=38"`
	PatientDriversLicenseNumber                 *string                 `field:"code=#R,order=39"`
	PriorAuthorizationProcessed                 *string                 `field:"code=#S,order=40"`
	NPPADea                                     *string                 `field:"code=#T,order=41"`
	NPPANpi                                     *string                 `field:"code=#U,order=42"`
	NPPAStateId                                 *string                 `field:"code=#V,order=43"`
	DrugCode                                    *string                 `field:"code=#W,order=44"`
	PharmacistLastName                          *string                 `field:"code=#X,order=45"`
	PharmacistFirstName                         *string                 `field:"code=#Y,order=46"`
	PharmacistInitials                          *string                 `field:"code=#Z,order=47"`
	PharmacistNpi                               *string                 `field:"code=!A,order=48"`
	PharmacistStateLicenseNumber                *string                 `field:"code=!B,order=49"`
	PrescriberDeaSuffix                         *string                 `field:"code=!C,order=50"`
	PrescriberXDea                              *string                 `field:"code=!D,order=51"`
	TreatmentType                               *string                 `field:"code=!E,order=52"`
	PrescribedQuantity                          *float64                `field:"code=!F,decimalPlaces=3,order=53"`
	CdtDiagnosisCode                            *string                 `field:"code=!G,order=54"`
	PartialFillCounter                          *string                 `field:"code=!H,order=55"`
	PrescriptionSerialNumber                    *string                 `field:"code=!I,order=56"`
	OwnerLastName                               *string                 `field:"code=!J,order=57"`
	OwnerFirstName                              *string                 `field:"code=!K,order=58"`
	OwnerBirthDate                              *time.Time              `field:"code=!L,format=YYYYMMdd,order=59"`
	PrescriberStateLicenseQualifier             *string                 `field:"code=!M,order=60"`
	NPPAStateLicenseQualifier                   *string                 `field:"code=!N,order=61"`
	AnimalOwnerSex                              *string                 `field:"code=!O,order=62"`
	DeliveryStreetAddress                       *string                 `field:"code=!P,order=63"`
	DeliveryCity                                *string                 `field:"code=!Q,order=64"`
	DeliveryState                               *string                 `field:"code=!R,order=65"`
	DeliveryZip                                 *string                 `field:"code=!S,order=66"`
	AllowHighDollarAwpUpdate                    *string                 `field:"code=!T,order=67"`
	DynamicFields                               []dynamic.DynamicStruct `field:"code=dynamic"`
}
