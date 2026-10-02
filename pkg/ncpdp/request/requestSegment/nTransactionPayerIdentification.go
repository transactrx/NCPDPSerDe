package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type NTransactionPayerIdentification struct {
	SegmentId ncpdp.SegmentId `json:"-"`

	PayerIin                    *string `field:"code=9Y,order=2,sinceVersion=F6"`
	PayerProcessorControlNumber *string `field:"code=9Z,order=3,sinceVersion=F6"`
	PayerCardholderId           *string `field:"code=AA,order=4,sinceVersion=F6"`
	PayerGroupId                *string `field:"code=AB,order=5,sinceVersion=F6"`
	PayerAdjudicatedProgramType *string `field:"code=9U,order=6,sinceVersion=F6"`
	TransactionSourceType       *string `field:"code=AC,order=7,sinceVersion=F6"`
	TransactionReconciliationId *string `field:"code=AD,order=8,sinceVersion=F6"`
	TransactionReferenceNumber  *string `field:"code=K5,order=9,sinceVersion=F6"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
