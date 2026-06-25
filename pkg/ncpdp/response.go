package ncpdp

import (
	"fmt"
	"strings"
)

func (tran *NcpdpTransaction[ResponseHeader]) Status() string {
	if tran == nil {
		return Empty
	}

	// FindSegmentInRecord falls back to shared segments when Records is empty
	// (F6 transmissions have no group separator, so the response status segment
	// lives in tran.Segments rather than tran.Records[0]).
	seg := tran.FindSegmentInRecord(0, RESPONSE_STATUS_SEGMENT_ID)
	if seg == nil {
		return Empty
	}

	field := seg.FindFirstField(STATUS_FIELD_ID)
	if field == nil {
		return Empty
	}

	return field.Value
}

func (tran *NcpdpTransaction[ResponseHeader]) IsStatusOf(status string) bool {
	if tran == nil {
		return false
	}

	return tran.Status() == status
}

func (tran *NcpdpTransaction[ResponseHeader]) IsPaid() bool {
	if tran == nil {
		return false
	}

	return tran.IsStatusOf(PAID_STATUS) || tran.IsStatusOf(DUPLICATE_PAID_STATUS)
}

func (tran *NcpdpTransaction[ResponseHeader]) IsRejected() bool {
	if tran == nil {
		return false
	}

	return tran.IsStatusOf(REJECTED_STATUS)
}

func (tran *NcpdpTransaction[ResponseHeader]) GetRejectCodes() []string {
	codes := []string{}

	if tran == nil {
		return codes
	}

	// Iterate via RecordCount so F6 responses (no group separator, segments live
	// at the transaction level) are read through FindSegmentInRecord's fallback.
	for i := 0; i < tran.RecordCount(); i++ {
		seg := tran.FindSegmentInRecord(i, RESPONSE_STATUS_SEGMENT_ID)
		if seg == nil {
			continue
		}

		for _, field := range seg.FindAllFields(REJECT_CODE_FIELD_ID) {
			codes = append(codes, strings.TrimSpace(field.Value))
		}
	}

	return codes
}

func (tran *NcpdpTransaction[ResponseHeader]) GetAdditionalMessages() map[string]string {
	messages := make(map[string]string)

	if tran == nil {
		return messages
	}

	for i := 0; i < tran.RecordCount(); i++ {
		seg := tran.FindSegmentInRecord(i, RESPONSE_STATUS_SEGMENT_ID)
		if seg == nil {
			continue
		}

		qfrFields := seg.FindAllFields(ADDITIONAL_MESSAGE_QUALIFIER_FIELD_ID)
		msgFields := seg.FindAllFields(ADDITIONAL_MESSAGE_FIELD_ID)

		for j := 0; j < len(msgFields); j++ {
			qfr := Empty
			msg := msgFields[j].GetString()

			if j < len(qfrFields) {
				qfr = qfrFields[j].GetString()
			}

			if qfr == Empty {
				qfr = fmt.Sprintf("%v", j)
			}

			messages[qfr] = msg
		}
	}

	return messages
}
