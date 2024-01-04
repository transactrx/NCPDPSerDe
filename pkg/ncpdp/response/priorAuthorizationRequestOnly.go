package response

import (
	"github.com/transactrx/NCPDPSerDe/pkg/dynamic"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp"
	responsesegment "github.com/transactrx/NCPDPSerDe/pkg/ncpdp/response/responseSegment"
)

type PriorAuthorizationRequestOnly struct {
	//Header
	Header ncpdp.NcpdpHeader[ncpdp.ResponseHeader] `header:"version=D0,deserializer=ParseNcpdpHeader,serializer=BuildNcpdpHeader,rawField=RawValue"`

	//Shared Segments
	Message         responsesegment.Message `segment:"code=AM20,order=1"`
	DynamicSegments []dynamic.DynamicStruct `segment:"code=dynamic,order=100"`

	//Authorization Groups
	Authorizations []PriorAuthorizationRequestOnlyRecord `group:"max=4"`
}

type PriorAuthorizationRequestOnlyRecord struct {
	Raw                    string                                 `field:"code=rawgroup"`
	Status                 responsesegment.Status                 `segment:"code=AM21,order=1"`
	Claim                  responsesegment.Claim                  `segment:"code=AM22,order=2"`
	PriorAuthorization     responsesegment.PriorAuthorization     `segment:"code=AM26,order=3"`
	CoordinationOfBenefits responsesegment.CoordinationOfBenefits `segment:"code=AM28,order=4"`
	DynamicSegments        []dynamic.DynamicStruct                `segment:"code=dynamic,order=100"`

	//TOOD: Include price update info
}
