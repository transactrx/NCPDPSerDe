package ncpdp

import (
	"log"
	"strconv"
	"strings"
)

type TransmissionData struct {
	TransmissionId   string
	PreEditId        string
	TransmissionType string
	FinalResponse    bool
	RouteCode        string
	RoutingAddress   string
	Request          *NcpdpTransaction[RequestHeader]
	Response         *NcpdpTransaction[ResponseHeader]
}

type Transmissions []TransmissionData

func NewTransmissionList() *Transmissions {
	tr := make(Transmissions, 0, 6)
	return &tr
}

// Log transmission event
func (td *TransmissionData) LogEvent(logger *log.Logger, format, reason string) {
	if td == nil {
		return
	}

	logger.Printf(
		format,
		td.TransmissionId,
		td.Request.Header.Value.Bin,
		td.Request.Header.Value.Pcn,
		td.Request.GetGroupCode(),
		td.Request.Header.Value.TransactionCode,
		td.Request.Header.Value.RecordCount,
		td.Response.Status(),
		td.Response.GetRejectCodes(),
		reason)
}

func (td *TransmissionData) GetRequestHash(recordIndex int) string {
	if td == nil {
		return Empty
	}

	if recordIndex >= len(td.Request.Records) {
		return Empty
	}

	fillNumber := td.Request.FindFirstField(CLAIM_SEGMENT_ID, FILL_NUMBER_FIELD_ID, recordIndex).GetIntOrDefault(0)

	builder := strings.Builder{}
	builder.WriteString(td.Request.Header.Value.ServiceProviderId)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.Header.Value.ServiceProviderIdQualifier)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.Header.Value.DateOfService)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.Header.Value.Bin)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.Header.Value.Pcn)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.GetGroupCode())
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.Header.Value.TransactionCode)
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.FindFirstField(CLAIM_SEGMENT_ID, PRESCRIPTION_SERVICE_REFERENCE_NO_FIELD_ID, recordIndex).GetString())
	builder.WriteString(Pipe)
	builder.WriteString(td.Request.FindFirstField(CLAIM_SEGMENT_ID, PRODUCT_SERVICE_ID_FIELD_ID, recordIndex).GetString())
	builder.WriteString(Pipe)
	builder.WriteString(strconv.Itoa(fillNumber))

	return builder.String()
}
