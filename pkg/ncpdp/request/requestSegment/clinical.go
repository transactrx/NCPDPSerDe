package requestsegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type ClinicalDiagnosis struct {
	Qualifier *string `field:"code=WE"`
	Code      *string `field:"code=DO"`
}

type ClinicalMeasurement struct {
	Counter   *int       `field:"code=XE"`
	Date      *time.Time `field:"code=ZE,format=YYYYMMdd"`
	Time      *time.Time `field:"code=H1,format=HHmm"`
	Dimension *string    `field:"code=H2"`
	Unit      *string    `field:"code=H3"`
	Value     *string    `field:"code=H4"`
}

type Clinical struct {
	SegmentId ncpdp.SegmentId

	DiagnosisCount *int `field:"code=VE"`
	Diagnoses      []ClinicalDiagnosis

	Measurements []ClinicalMeasurement
}
