package ncpdp

const (
	Empty = ""
)

// Separators
const (
	GROUP   byte = 0x1D
	SEGMENT byte = 0x1E
	FIELD   byte = 0x1C
	ETX     byte = 0x3

	ID_LENGTH = 2
)

// Response Status Codes
const (
	ACCEPTED_STATUS           = "A"
	DUPLICATE_PAID_STATUS     = "D"
	DUPLICATE_REVERSAL_STATUS = "S"
	PAID_STATUS               = "P"
	REJECTED_STATUS           = "R"
)

// Request Segment IDs
const (
	PATIENT_SEGMENT_ID   = "AM01"
	INSURANCE_SEGMENT_ID = "AM04"
	CLAIM_SEGMENT_ID     = "AM07"
)

// Response Segment IDs
const (
	RESPONSE_MESSAGE_SEGMENT_ID = "AM20"
	RESPONSE_STATUS_SEGMENT_ID  = "AM21"
	RESPONSE_PRICING_SEGMENT_ID = "AM23"
)

// Field IDs
const (
	SEGMENT_FIELD_ID = "AM"

	MEDIGAP_ID_FIELD_ID                                  = "2A"
	PATIENT_RESIDENCE_FIELD_ID                           = "4X"
	STATUS_FIELD_ID                                      = "AN"
	FLAT_SALES_TAX_AMOUNT_PAID_FIELD_ID                  = "AW"
	PERCENTAGE_SALES_TAX_AMOUNT_PAID_FIELD_ID            = "AX"
	GROUP_CODE_FIELD_ID                                  = "C1"
	CARDHOLDER_ID_FIELD_ID                               = "C2"
	PLACE_OF_SERVICE_FIELD_ID                            = "C7"
	OTHER_COVERAGE_CODE_FIELD_ID                         = "C8"
	PRESCRIPTION_SERVICE_REFERENCE_NO_FIELD_ID           = "D2"
	FILL_NUMBER_FIELD_ID                                 = "D3"
	DAYS_SUPPLY_FIELD_ID                                 = "D5"
	COMPOUND_CODE_FIELD_ID                               = "D6"
	PRODUCT_SERVICE_ID_FIELD_ID                          = "D7"
	PRODUCT_SERVICE_ID_QUALIFIER_FIELD_ID                = "E1"
	QUANTITY_DISPENSED_FIELD_ID                          = "E7"
	PRESCRIPTION_SERVICE_REFERENCE_NO_QUALIFIER_FIELD_ID = "EM"
	AUTHORIZATION_NUMBER_FIELD_ID                        = "F3"
	MESSAGE_FIELD_ID                                     = "F4"
	PATIENT_PAY_AMOUNT_FIELD_ID                          = "F5"
	INGREDIENT_COST_PAID_FIELD_ID                        = "F6"
	DISPENSING_FEE_PAID_FIELD_ID                         = "F7"
	TOTAL_AMOUNT_PAID_FIELD_ID                           = "F9"
	REJECT_CODE_FIELD_ID                                 = "FB"
	BASIS_OF_REIMBURSEMENT_FIELD_ID                      = "FM"
	ADDITIONAL_MESSAGE_FIELD_ID                          = "FQ"
	MEDICAID_ID_NUMBER_FIELD_ID                          = "N5"
	MEDICAID_AGENCY_NUMBER_FIELD_ID                      = "N6"
	PHARMACY_SERVICE_TYPE_FIELD_ID                       = "U7"
	ADDITIONAL_MESSAGE_QUALIFIER_FIELD_ID                = "UH"
)

// Transaction Codes
const (
	BILLING  = "B1"
	REVERSAL = "B2"
	REBILL   = "B3"
)

// Service/Provider ID Qualfiers
const (
	SERVICE_PROVIDER_QUALFIER_NDC = "03"
)

// Compound Codes
const (
	NOT_A_COMPOUND = "1"
)

// Error messages
const (
	ErrSegmentNotFound = "segment not found: %v"
)
