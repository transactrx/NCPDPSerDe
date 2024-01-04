package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type Clinical struct {
	SegmentId ncpdp.SegmentId

	DiagnosisCount *int `field:"code=VE,order=2"`
	Diagnoses      []ClinicalDiagnosis

	Measurements []ClinicalMeasurement

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}

type ClinicalDiagnosis struct {
	Qualifier *string `field:"code=WE,order=3"`
	Code      *string `field:"code=DO,order=4"`
}

type ClinicalMeasurement struct {
	Counter   *int       `field:"code=XE,order=5"`
	Date      *time.Time `field:"code=ZE,format=YYYYMMdd,order=6"`
	Time      *time.Time `field:"code=H1,format=HHmm,order=7"`
	Dimension *string    `field:"code=H2,order=8"`
	Unit      *string    `field:"code=H3,order=9"`
	Value     *string    `field:"code=H4,order=10"`
}
