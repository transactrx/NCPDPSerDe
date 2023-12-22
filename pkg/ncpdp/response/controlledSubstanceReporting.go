package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type ControlledSubstanceReporting struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,parseMethod=ParseNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message responsesegment.Message `segment:"code=AM20"`

	//Claim Groups
	Claims []ControlledSubstanceReportingRecord `group:"max=4"`
}

type ControlledSubstanceReportingRecord struct {
	Raw    string                 `field:"code=rawgroup"`
	Status responsesegment.Status `segment:"code=AM21"`
	Claim  responsesegment.Claim  `segment:"code=AM22"`

	//TOOD: Include price update info
}
