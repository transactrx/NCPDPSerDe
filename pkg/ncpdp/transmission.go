package ncpdp

import "log"

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
