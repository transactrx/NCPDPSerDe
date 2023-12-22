package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type DataCollection struct {
	SegmentId ncpdp.SegmentId

	InventoryCategoryCode                       *string    `field:"code=&A"`
	PatientId                                   *string    `field:"code=&B"` //Not de-identified
	PayerSequence                               *string    `field:"code=&C"`
	SubstitutionType                            *string    `field:"code=&D"`
	DispensingType                              *string    `field:"code=&E"`
	CompanyNumber                               *string    `field:"code=&F"`
	PayPlan                                     *string    `field:"code=&G"`
	PatientLocationCode                         *string    `field:"code=&H"`
	ClaimIndicator                              *string    `field:"code=&I"`
	Facility                                    *string    `field:"code=&J"`
	TransactionTime                             *time.Time `field:"code=&K,format=HHmmss"`
	MedicarePartDPriceCode                      *string    `field:"code=&L"`
	TestTransaction                             *string    `field:"code=&M"`
	PharmacyNpi                                 *string    `field:"code=&N"`
	CycleFillListIndictor                       *string    `field:"code=&T"`
	AcquisitionCost                             *float64   `field:"code=&U,decimalPlaces=2"`
	SellingPrice                                *float64   `field:"code=&V,decimalPlaces=2"`
	PaymentType                                 *string    `field:"code=&W"`
	NursingHome                                 *string    `field:"code=&X"`
	PatientCodeDeidentified                     *string    `field:"code=&Y"`
	PrescriptionStopDate                        *time.Time `field:"code=&Z,format=YYYYMMdd"`
	PriceCodeType                               *string    `field:"code=#A"`
	PatientLanguage                             *string    `field:"code=#B"`
	PrescriberDea                               *string    `field:"code=#C"`
	SecondaryTernaryIndicator                   *string    `field:"code=#D"`
	DecimalQuantity                             *float64   `field:"code=#E"`
	PrescriberStateId                           *string    `field:"code=#F"`
	PrescribedNdc                               *string    `field:"code=#G"`
	ElectronicRxTransactionReferenceNumber      *string    `field:"code=#H"` //UIB-030-01
	ElectronicRxInitiatorControlReferenceNumber *string    `field:"code=#I"` //UIH-030-01
	PharmacyDea                                 *string    `field:"code=#K"`
	PrescriberNpi                               *string    `field:"code=#L"`
	DrugAWP                                     *float64   `field:"code=#M,decimalPlaces=2"`
	ClassificationCode                          *string    `field:"code=#N"`
	PHIParticipation                            *string    `field:"code=#O"`
	MarketingParticipation                      *string    `field:"code=#P"`
	PatientSocialSecurityNumber                 *string    `field:"code=#Q"`
	PatientDriversLicenseNumber                 *string    `field:"code=#R"`
	PriorAuthorizationProcessed                 *string    `field:"code=#S"`
	NPPADea                                     *string    `field:"code=#T"`
	NPPANpi                                     *string    `field:"code=#U"`
	NPPAStateId                                 *string    `field:"code=#V"`
	DrugCode                                    *string    `field:"code=#W"`
	PharmacistLastName                          *string    `field:"code=#X"`
	PharmacistFirstName                         *string    `field:"code=#Y"`
	PharmacistInitials                          *string    `field:"code=#Z"`
	PharmacistNpi                               *string    `field:"code=!A"`
	PharmacistStateLicenseNumber                *string    `field:"code=!B"`
	PrescriberDeaSuffix                         *string    `field:"code=!C"`
	PrescriberXDea                              *string    `field:"code=!D"`
	TreatmentType                               *string    `field:"code=!E"`
	PrescribedQuantity                          *float64   `field:"code=!F,decimalPlaces=3"`
	CdtDiagnosisCode                            *string    `field:"code=!G"`
	PartialFillCounter                          *string    `field:"code=!H"`
	PrescriptionSerialNumber                    *string    `field:"code=!I"`
	OwnerLastName                               *string    `field:"code=!J"`
	OwnerFirstName                              *string    `field:"code=!K"`
	OwnerBirthDate                              *time.Time `field:"code=!L,format=YYYYMMdd"`
	PrescriberStateLicenseQualifier             *string    `field:"code=!M"`
	NPPAStateLicenseQualifier                   *string    `field:"code=!N"`
}
