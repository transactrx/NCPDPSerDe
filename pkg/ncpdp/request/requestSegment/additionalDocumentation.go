package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type AdditionalDocumentation struct {
	SegmentId                    ncpdp.SegmentId
	TypeId                       *string `field:"code=2Q,order=2"`
	RequestPeriod                AdditionalDocumentationRequestPeriod
	LengthOfNeed                 AdditionalDocumentationLengthOfNeed
	PrescriberSupplierDateSigned *time.Time `field:"code=2T,format=YYYYMMdd,order=8"`
	SupportingDocumentation      *string    `field:"code=2X,order=9"`

	QuestionCount *int `field:"code=2Z,order=10"`
	Questions     []AdditionalDocumentationQuestion

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type AdditionalDocumentationRequestPeriod struct {
	BeginDate   *time.Time `field:"code=2V,format=YYYYMMdd,order=3"`
	RevisedDate *time.Time `field:"code=2W,format=YYYYMMdd,order=4"`
	Status      *string    `field:"code=2U,order=5"`
}

type AdditionalDocumentationLengthOfNeed struct {
	Qualifier *string `field:"code=2S,order=6"`
	Value     *string `field:"code=2R,order=7"`
}

type AdditionalDocumentationQuestion struct {
	Number               *string    `field:"code=4B,order=11"`
	PercentResponse      *float64   `field:"code=4D,decimalPlaces=2,overpunch=true,order=12"`
	ResponseDate         *time.Time `field:"code=4G,format=YYYYMMdd,order=13"`
	DollarAmountResponse *float64   `field:"code=4H,decimalPlaces=2,overpunch=true,order=14"`
	NumericResponse      *int       `field:"code=4J,order=15"`
	AlphanumericResponse *string    `field:"code=4K,order=16"`
}
