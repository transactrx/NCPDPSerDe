package requestsegment

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
)

type NTransactionPayerIdentification struct {
	SegmentId ncpdp.SegmentId

	PayerIin                    *string `field:"code=9Y,order=2"`
	PayerProcessorControlNumber *string `field:"code=9Z,order=3"`
	PayerCardholderId           *string `field:"code=AA,order=4"`
	PayerGroupId                *string `field:"code=AB,order=5"`
	PayerAdjudicatedProgramType *string `field:"code=9U,order=6"`
	TransactionSourceType       *string `field:"code=AC,order=7"`
	TransactionReconciliationId *string `field:"code=AD,order=8"`
	TransactionReferenceNumber  *string `field:"code=K5,order=9"`

	DynamicFields []dynamic.DynamicStruct `field:"code=dynamic"`
}
