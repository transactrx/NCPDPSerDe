package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type AdditionalDocumentationRequestPeriod struct {
	BeginDate   *time.Time `field:"code=2V,format=YYYYMMdd"`
	RevisedDate *time.Time `field:"code=2W,format=YYYYMMdd"`
	Status      *string    `field:"code=2U"`
}

type AdditionalDocumentationLengthOfNeed struct {
	Qualifier *string `field:"code=2S"`
	Value     *string `field:"code=2R"`
}

type AdditionalDocumentationQuestion struct {
	Number               *string    `field:"code=4B"`
	PercentResponse      *float64   `field:"code=4D,decimalPlaces=2,overpunch=true"`
	ResponseDate         *time.Time `field:"code=4G,format=YYYYMMdd"`
	DollarAmountResponse *float64   `field:"code=4H,decimalPlaces=2,overpunch=true"`
	NumericResponse      *int       `field:"code=4J"`
	AlphanumericResponse *string    `field:"code=4K"`
}

type AdditionalDocumentation struct {
	SegmentId                    ncpdp.SegmentId
	TypeId                       *string `field:"code=2Q"`
	RequestPeriod                AdditionalDocumentationRequestPeriod
	LengthOfNeed                 AdditionalDocumentationLengthOfNeed
	PrescriberSupplierDateSigned *time.Time `field:"code=2T,format=YYYYMMdd"`
	SupportingDocumentation      *string    `field:"code=2X"`

	QuestionCount *int `field:"code=2Z"`
	Questions     []AdditionalDocumentationQuestion
}
