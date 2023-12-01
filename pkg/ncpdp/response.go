package ncpdp

import (
	"fmt"
)

func (tran *NcpdpTransaction[ResponseHeader]) Status() string {
	seg := tran.Records[0].FindSegment(RESPONSE_STATUS_SEGMENT_ID)
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
	return tran.Status() == status
}

func (tran *NcpdpTransaction[ResponseHeader]) IsPaid() bool {
	return tran.IsStatusOf(PAID_STATUS) || tran.IsStatusOf(DUPLICATE_PAID_STATUS)
}

func (tran *NcpdpTransaction[ResponseHeader]) IsRejected() bool {
	return tran.IsStatusOf(REJECTED_STATUS)
}

func (tran *NcpdpTransaction[ResponseHeader]) GetRejectCodes() []string {
	codes := []string{}

	for _, record := range tran.Records {
		seg := record.FindSegment(RESPONSE_STATUS_SEGMENT_ID)

		if seg != nil {
			segFields := seg.FindAllFields(REJECT_CODE_FIELD_ID)
			for _, field := range segFields {
				codes = append(codes, field.Value)
			}
		}
	}

	return codes
}

func (tran *NcpdpTransaction[ResponseHeader]) GetAdditionalMessages() map[string]string {
	messages := make(map[string]string)

	for _, record := range tran.Records {
		seg := record.FindSegment(RESPONSE_STATUS_SEGMENT_ID)

		if seg != nil {
			qfrFields := seg.FindAllFields(ADDITIONAL_MESSAGE_QUALIFIER_FIELD_ID)
			msgFields := seg.FindAllFields(ADDITIONAL_MESSAGE_FIELD_ID)

			for i := 0; i < len(msgFields); i++ {
				qfr := Empty
				msg := msgFields[i].GetString()

				if i < len(qfrFields) {
					qfr = qfrFields[i].GetString()
				}

				if qfr == Empty {
					qfr = fmt.Sprintf("%v", i)
				}

				messages[qfr] = msg
			}
		}
	}

	return messages
}
