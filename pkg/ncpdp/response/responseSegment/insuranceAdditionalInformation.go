package responsesegment

import (
	"time"

	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type InsuranceAdditionalInformation struct {
	SegmentId      ncpdp.SegmentId `json:"-"`
	MedicarePartD  MedicarePartD
	ContractNumber *string `field:"code=U1,order=4"`
	FormularyId    *string `field:"code=FF,order=5"`
	PlanName       *string `field:"code=96,order=9,sinceVersion=F6"`
}

type MedicarePartD struct {
	MedicarePartDCoverageCode        *string                 `field:"code=UR,order=2"`
	LICSLevel                        *string                 `field:"code=UQ,order=3"`
	BenefitId                        *string                 `field:"code=U6,order=6"`
	NextMedicarePartDEffectiveDate   *time.Time              `field:"code=US,format=YYYYMMdd,order=7"`
	NextMedicarePartDTerminationDate *time.Time              `field:"code=UT,format=YYYYMMdd,order=8"`
	DynamicFields                    []dynamic.DynamicStruct `field:"code=dynamic"`
}
