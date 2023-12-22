package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type InsuranceAdditionalInformation struct {
	SegmentId      ncpdp.SegmentId
	MedicarePartD  MedicarePartD
	ContractNumber *string `field:"code=U1"`
	FormularyId    *string `field:"code=FF"`
}

type MedicarePartD struct {
	MedicarePartDCoverageCode        *string    `field:"code=UR"`
	LICSLevel                        *string    `field:"code=UQ"`
	BenefitId                        *string    `field:"code=U6"`
	NextMedicarePartDEffectiveDate   *time.Time `field:"code=US,format=YYYYMMdd"`
	NextMedicarePartDTerminationDate *time.Time `field:"code=UT,format=YYYYMMdd"`
}
