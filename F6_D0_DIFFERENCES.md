# Complete F6 vs D.0 Differences

**Sources:** NCPDP Telecommunication Implementation Guide vF6 (`D:\POW\NCPDP Files\F6 Format\Telecommunication-Imp-Guide-VF6.pdf`), the official F6 Data Dictionary (`Data Dictionary\NCPDP_F6_Data_Dictionary.xlsx`), and the team's change-tracking workbook (`vF6 Planning.xlsx`, sheets "F6 Request Changes" / "F6 Response Changes"). Cross-verified against this library's struct definitions — every spec-new field below is implemented in NCPDPSerDe (verified programmatically; zero gaps).

Format notation is NCPDP picture format: `x(20)` = alphanumeric max 20, `9(8)` = numeric 8, `s9(6)v99` = signed/overpunch 6+2 decimals, `ECS` = extended character set.

---

## 1. Transmission-structure changes

- **Group separator (`0x1D`) eliminated.** F6 transmissions contain no group separator; all segments follow the header directly.
- **One transaction per transmission.** Transaction Count (109-A9) is fixed at `1` in F6 (D.0 allowed up to 4).
- **Version value:** 102-A2 = `F6`, and it moves to the first position of the header (F6 transmissions lead with `"F6"` instead of BIN digits — that's how version auto-detection works).
- **Segment renames:** Response Coordination of Benefits (AM28) → **Response Other Payers Segment**. The "Preferred Product" field family in Response Claim is renamed **"Formulary Alternative"**.
- **Coupon segment (AM09) is eliminated** — all three of its fields are removed (§4).

## 2. Header changes — request header 56 → 58 bytes

| Field | D.0 position/format | F6 position/format | Change |
|---|---|---|---|
| 102-A2 Version/Release | 6–7 `x(2)` | **0–1** `x(2)` | Moved to first position |
| 103-A3 Transaction Code | 8–9 `x(2)` | **2–3** `x(2)` | Moved to second position |
| 101-A1 BIN → **IIN Number** | 0–5 `9(6)` | 4–11 `9(8)` | **Renamed and lengthened 6 → 8** |
| 104-A4 Processor Control Number | 10–19 `x(10)` | 12–21 | Shifted +2 |
| 109-A9 Transaction Count | 20 `x(1)` | 22 | Shifted +2; value fixed at 1 |
| 202-B2 Service Provider ID Qualifier | 21–22 `x(2)` | 23–24 | Shifted +2 |
| 201-B1 Service Provider ID | 23–37 `x(15)` | 25–39 | Shifted +2 |
| 401-D1 Date of Service | 38–45 `9(8)` | 40–47 | Shifted +2 |
| 110-AK Software Vendor/Certification ID | 46–55 `x(10)` | 48–57 | Shifted +2 |

The **response header is unchanged** (31 bytes, same layout).

## 3. Field length / format changes (existing fields)

### Request side

| Field | Name | Segment | D.0 | F6 |
|---|---|---|---|---|
| 101-A1 | IIN Number (was BIN Number) | Header | `9(6)` | `9(8)` |
| 350-HN | Patient E-Mail Address | Patient (AM01) | `x(80)` | `ECS(80)` |
| 444-E9 | Provider ID | Pharmacy Provider (AM02) | `x(15)` | `x(35)` |
| 427-DR | Prescriber Last Name | Prescriber (AM03) | `x(15)` | `x(35)` |
| 364-2J | Prescriber First Name | Prescriber (AM03) | `x(12)` | `x(35)` |
| 411-DB | Prescriber ID | Prescriber (AM03) | `x(15)` | `x(35)` |
| 470-4E | Primary Care Provider Last Name | Prescriber (AM03) | `x(15)` | `x(35)` |
| 421-DL | Primary Care Provider ID | Prescriber (AM03) | `x(15)` | `x(35)` |
| 340-7C | Other Payer ID | COB/Other Payments (AM05) | `x(10)` | `x(30)` |
| 352-NQ | Other Payer-Patient Responsibility Amount | COB/Other Payments (AM05) | `s9(8)v99` | `s9(9)v99` |
| 431-DV | Other Payer Amount Paid | COB/Other Payments (AM05) | `s9(6)v99` | `s9(9)v99` |
| 120-TU | Pay To Name | Workers' Compensation (AM06) | `x(20)` | `x(70)` |
| 126-UA | Generic Equivalent Product ID | Workers' Compensation (AM06) | `x(19)` | `x(40)` |
| 462-EV | Prior Authorization ID Submitted (was PA Number Submitted) | Claim (AM07) | `9(11)` | `x(35)` |
| 420-DK | Submission Clarification Code | Claim (AM07) | `9(2)` | `9(3)` |
| 429-DT | Special Packaging Indicator | Claim (AM07) | `9(1)` | `9(2)` |
| 454-EK | Scheduled Prescription ID Number | Claim (AM07) | `x(12)` | `x(30)` |
| 407-D7 | Product/Service ID | Claim (AM07) | `x(19)` | `x(40)` |
| 445-EA | Originally Prescribed Product/Service Code | Claim (AM07) | `x(19)` | `x(40)` |
| 476-H6 | DUR Co-Agent ID | DUR/PPS (AM08) | `x(19)` | `x(40)` |
| 450-EF | Compound Dosage Form Description Code | Compound (AM10) | `x(2)` | `x(15)` |
| 448-ED | Compound Ingredient Quantity | Compound (AM10) | `9(7)v999` | `9(7)v9999999` |
| 489-TE | Compound Product ID | Compound (AM10) | `x(19)` | `x(40)` |
| 449-EE | Compound Ingredient Drug Cost | Compound (AM10) | `s9(6)v99` | `s9(9)v99` |
| 433-DX | Patient Pay Amount Reported (was Patient Paid Amount Submitted) | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 481-HA | Regulatory Fee Amount Submitted (was Flat Sales Tax Amount Submitted) | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 482-GE | Percentage Tax Amount Submitted (was Percentage Sales Tax…) | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 483-HE | Percentage Tax Rate Submitted (was Percentage Sales Tax Rate…) | Pricing (AM11) | `s9(3)v4` | `9(3)v9999` |
| 480-H9 | Other Amount Claim Submitted | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 412-DC | Dispensing Fee Submitted | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 438-E3 | Incentive Amount Submitted | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 409-D9 | Ingredient Cost Submitted | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 426-DQ | Usual And Customary Charge | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 430-DU | Gross Amount Due | Pricing (AM11) | `s9(6)v99` | `s9(9)v99` |
| 381-4H | Question Dollar Amount Response | Additional Documentation (AM14) | `s9(9)v99` | `s9(12)v99` |
| 385-3Q | Facility Name | Facility (AM15) | `x(30)` | `x(70)` |

(Pattern: dollar fields widen to `s9(9)v99`, provider names/IDs widen to `x(35)`, product IDs widen to `x(40)`.)

### Response side

| Field | Name | Segment | D.0 | F6 |
|---|---|---|---|---|
| 553-AR | Formulary Alternative ID (was Preferred Product ID) | Response Claim (AM22) | `x(19)` | `x(40)` |
| 554-AS | Formulary Alternative Incentive (was Preferred Product Incentive) | Response Claim (AM22) | `s9(6)v99` | `s9(9)v99` |
| 555-AT | Formulary Alternative Estimated Patient Cost Share (was Preferred Product Cost Share Incentive) | Response Claim (AM22) | `s9(6)v99` | `s9(9)v99` |
| 558-AW | Regulatory Fee Amount Paid (was Flat Sales Tax Amount Paid) | Response Pricing (AM23) | `s9(6)v99` | `s9(9)v99` |
| 559-AX | Percentage Tax Amount Paid (was Percentage Sales Tax Amount Paid) | Response Pricing (AM23) | `s9(6)v99` | `s9(9)v99` |
| 560-AY | Percentage Tax Rate Paid (was Percentage Sales Tax Rate Paid) | Response Pricing (AM23) | `s9(3)v4` | `9(3)v9999` |
| 575-EQ | Patient Percentage Tax Amount (was Patient Sales Tax Amount) | Response Pricing (AM23) | `s9(6)v99` | `s9(9)v99` |

## 4. Field removals

### Request side (22 fields)

| Field | Name | Segment | Replacement (if any) |
|---|---|---|---|
| 314-CE | Home Plan | Insurance (AM04) | — |
| 990-MG | Other Payer BIN Number | Insurance (AM04) | — |
| 991-MH | Other Payer Processor Control Number | Insurance (AM04) | — |
| 356-NU | Other Payer Cardholder ID | Insurance (AM04) | — |
| 992-MJ | Other Payer Group ID | Insurance (AM04) | — |
| 322-CM | Patient Street Address | Patient (AM01) | 7A / 7B address lines |
| 365-2K | Prescriber Street Address | Prescriber (AM03) | 7U / 7V address lines |
| 463-EW | Intermediary Authorization Type ID | Claim (AM07) | Response Intermediary segment 8S |
| 464-EX | Intermediary Authorization ID | Claim (AM07) | Response Intermediary segment 8T |
| 880-K5 | Transaction Reference Number | Claim (AM07) | Moved to new AM38 segment |
| 114-N4 | Medicaid Subrogation ICN/TCN | Claim (AM07) | — |
| 485-KE | Coupon Type | Coupon (AM09) | Segment eliminated |
| 486-ME | Coupon Number | Coupon (AM09) | Segment eliminated |
| 487-NE | Coupon Value Amount | Coupon (AM09) | Segment eliminated |
| 392-MU | Benefit Stage Count | COB/Other Payments (AM05) | 9W Benefit Stage Indicator Count |
| 393-MV | Benefit Stage Qualifier | COB/Other Payments (AM05) | 9X Benefit Stage Indicator |
| 394-MW | Benefit Stage Amount | COB/Other Payments (AM05) | — |
| 316-CG | Employer Street Address | Workers' Compensation (AM06) | 8D / 7G address lines |
| 121-TV | Pay To Street Address | Workers' Compensation (AM06) | 7R / 7S address lines |
| 113-N3 | Medicaid Paid Amount | Pricing (AM11) | — |
| 498-PG | Authorized Rep. Street Address | Prior Authorization (AM12) | 7D / 8B address lines |
| 386-3U | Facility Street Address | Facility (AM15) | 7M / 7N address lines |

(Also removed at name level: 321-CL Employer Contact Name, split into 7H first / 7J last.)

### Response side (21 fields)

| Field | Name | Segment | Replacement (if any) |
|---|---|---|---|
| 138-UQ | CMS LIS Level | Response Insurance Additional Info (AM27) | KF LIS Level in new AM39 segment |
| 987-MA | URL | Response Status (AM21) | BC Help Desk Contact Information |
| 549-7F | Help Desk Phone Number Qualifier | Response Status (AM21) | BF Help Desk Contact Info Qualifier |
| 550-8F | Help Desk Phone Number | Response Status (AM21) | BC Help Desk Contact Information |
| 114-N4 | Medicaid Subrogation ICN/TCN | Response Claim (AM22) | KY Subrogation Requestor's Reconciliation ID |
| 523-FN | Amount Attributed to Sales Tax | Response Pricing (AM23) | Patient Pay Components (KQ/KN) |
| 517-FH | Amount Applied to Periodic Deductible | Response Pricing (AM23) | Patient Pay Components |
| 518-FI | Amount of Copay | Response Pricing (AM23) | Patient Pay Components |
| 520-FK | Amount Exceeding Periodic Benefit Maximum | Response Pricing (AM23) | Patient Pay Components |
| 571-NZ | Amount Attributed to Processor Fee | Response Pricing (AM23) | Patient Pay Components |
| 574-2Y | Plan Sales Tax Amount | Response Pricing (AM23) | — |
| 129-UD | Health Plan-Funded Assistance Amount | Response Pricing (AM23) | Patient Pay Components |
| 133-UJ | Amount Attributed to Provider Network Selection | Response Pricing (AM23) | Patient Pay Components |
| 134-UK | Amount Attributed to Product Selection/Brand Drug | Response Pricing (AM23) | Patient Pay Components |
| 135-UM | Amount Attributed to Product Selection/Non-Preferred Formulary | Response Pricing (AM23) | Patient Pay Components |
| 136-UN | Amount Attributed to Product Selection/Brand Non-Preferred | Response Pricing (AM23) | Patient Pay Components |
| 137-UP | Amount Attributed to Coverage Gap | Response Pricing (AM23) | Patient Pay Components |
| 392-MU | Benefit Stage Count | Response Pricing (AM23) | 9W Benefit Stage Indicator Count |
| 393-MV | Benefit Stage Qualifier | Response Pricing (AM23) | 9X Benefit Stage Indicator |
| 394-MW | Benefit Stage Amount | Response Pricing (AM23) | — |
| 570-NS | DUR Additional Text | Response DUR/PPS (AM24) | Replaced by 17 discrete DUR fields (§6) |

> **Library note:** NCPDPSerDe keeps all removed fields in its structs so D.0 traffic still parses; on F6 transactions they are simply never populated.

## 5. Field renames (same field number, new name)

### Request side

| Field | D.0 name | F6 name | Segment |
|---|---|---|---|
| 101-A1 | BIN Number | IIN Number | Header |
| 323-CN | Patient City | Patient City Address | Patient (AM01) |
| 324-CO | Patient State or Province | Patient State or Province Address | Patient (AM01) |
| 326-CQ | Patient Phone Number | Patient Telephone Number | Patient (AM01) |
| 462-EV | Prior Authorization Number Submitted | Prior Authorization ID Submitted | Claim (AM07) |
| 498-PM | Prescriber Phone Number | Prescriber Telephone Number | Prescriber (AM03) |
| 368-2P | Prescriber ZIP/Postal Zone | Prescriber ZIP/Postal Code | Prescriber (AM03) |
| 320-CK | Employer Phone Number | Employer Telephone Number | Workers' Comp (AM06) |
| 124-TY | Pay To ZIP/Postal Zone | Pay To ZIP/Postal Code | Workers' Comp (AM06) |
| 498-PF | Authorized Rep. Last Name | Authorized Representative Last Name | Prior Auth (AM12) |
| 498-PH | Authorized Rep. City | Authorized Representative City | Prior Auth (AM12) |
| 498-PJ | Authorized Rep. State/Province | Authorized Representative State/Province | Prior Auth (AM12) |
| 498-PK | Authorized Rep. Zip/Postal Code | Authorized Representative ZIP/Postal Code | Prior Auth (AM12) |
| 498-PY | Prior Authorization Number - Assigned | Prior Authorization ID Assigned | Prior Auth (AM12) |
| 433-DX | Patient Paid Amount Submitted | Patient Pay Amount Reported | Pricing (AM11) |
| 481-HA | Flat Sales Tax Amount Submitted | Regulatory Fee Amount Submitted | Pricing (AM11) |
| 482-GE | Percentage Sales Tax Amount Submitted | Percentage Tax Amount Submitted | Pricing (AM11) |
| 483-HE | Percentage Sales Tax Rate Submitted | Percentage Tax Rate Submitted | Pricing (AM11) |
| 484-JE | Percentage Sales Tax Basis Submitted | Percentage Tax Basis Submitted | Pricing (AM11) |
| 389-6D | Facility ZIP/Postal Zone | Facility ZIP/Postal Code | Facility (AM15) |

### Response side

| Field | D.0 name | F6 name | Segment |
|---|---|---|---|
| 568-J7 | Payer ID Qualifier | Payer/Health Plan ID Qualifier | Response Insurance (AM25) |
| 569-J8 | Payer ID | Payer/Health Plan ID | Response Insurance (AM25) |
| 551-9F | Preferred Product Count | Formulary Alternative Product Count | Response Claim (AM22) |
| 552-AP | Preferred Product ID Qualifier | Formulary Alternative ID Qualifier | Response Claim (AM22) |
| 553-AR | Preferred Product ID | Formulary Alternative ID | Response Claim (AM22) |
| 554-AS | Preferred Product Incentive | Formulary Alternative Incentive | Response Claim (AM22) |
| 555-AT | Preferred Product Cost Share Incentive | Formulary Alternative Estimated Patient Cost Share | Response Claim (AM22) |
| 556-AU | Preferred Product Description | Formulary Alternative Description | Response Claim (AM22) |
| 557-AV | Tax Exempt Indicator | Percentage Tax Exempt Indicator | Response Pricing (AM23) |
| 558-AW | Flat Sales Tax Amount Paid | Regulatory Fee Amount Paid | Response Pricing (AM23) |
| 559-AX | Percentage Sales Tax Amount Paid | Percentage Tax Amount Paid | Response Pricing (AM23) |
| 560-AY | Percentage Sales Tax Rate Paid | Percentage Tax Rate Paid | Response Pricing (AM23) |
| 561-AZ | Percentage Sales Tax Basis Paid | Percentage Tax Basis Paid | Response Pricing (AM23) |
| 575-EQ | Patient Sales Tax Amount | Patient Percentage Tax Amount | Response Pricing (AM23) |
| 531-FV | Quantity of Previous Fill | Quantity of Previous Dispensing | Response DUR/PPS (AM24) |
| 530-FU | Previous Date of Fill | Previous Date of Service | Response DUR/PPS (AM24) |
| 498-PY | Prior Authorization Number - Assigned | Prior Authorization ID Assigned | Response Prior Auth (AM26) |
| 127-UB | Other Payer Help Desk Phone Number | Other Payer Help Desk Telephone Number | Response Other Payers (AM28) |

## 6. Field additions to existing segments

### Request side

**Patient (AM01)** — plus CX/CY become repeating (max 9, counted by RR)

| Field | Name | Format | Notes |
|---|---|---|---|
| 618-RR | Patient ID Count | `9(1)` | max 9 CX/CY pairs |
| A43-1K | Patient Country Code | `x(2)` | |
| B38-1Y | Patient ID Associated Country Code | `x(2)` | |
| A22-YR | Patient ID Associated State/Province Address | `x(2)` | |
| E09-0C | Patient Middle Name | `x(25)` | |
| E10-0D | Patient Name Prefix | `x(10)` | |
| E11-0E | Patient Name Suffix | `x(10)` | |
| B08-7A | Patient Street Address Line 1 | `x(40)` | replaces 322-CM |
| B09-7B | Patient Street Address Line 2 | `x(40)` | replaces 322-CM |
| E06-S8 | Species | `x(9)` | |
| A45-1R | Veterinary Use Indicator | `x(1)` | |

**Pharmacy Provider (AM02)**

| Field | Name | Format |
|---|---|---|
| B96-4A | Provider First Name | `x(35)` |
| B97-4M | Provider Last Name | `x(35)` |

**Prescriber (AM03)**

| Field | Name | Format |
|---|---|---|
| E12-0F | Prescriber Middle Name | `x(25)` |
| B27-7U | Prescriber Street Address Line 1 (replaces 365-2K) | `x(40)` |
| B28-7V | Prescriber Street Address Line 2 | `x(40)` |
| B42-3C | Prescriber Country Code | `x(2)` |
| A24-ZK | Prescriber ID Associated State/Province Address | `x(2)` |
| B41-3B | Prescriber ID Associated Country Code | `x(2)` |
| B26-7T | Prescriber Telephone Number Extension | `9(8)` |
| A25-ZM | Prescriber Alternate ID Qualifier | `x(2)` |
| A26-ZP | Prescriber Alternate ID | `x(35)` |
| A27-ZQ | Prescriber Alternate ID Associated State/Province Address | `x(2)` |
| B40-3A | Prescriber Alternate ID Associated Country Code | `x(2)` |
| D01-KV | Prescriber DEA Number | `x(15)` |
| D57-RG | Prescriber Place of Service | `9(2)` |

**Claim (AM07)**

| Field | Name | Format | Notes |
|---|---|---|---|
| 579-XX | Associated Prescription/Service Provider ID Qualifier | `x(2)` | |
| 580-XY | Associated Prescription/Service Provider ID | `x(35)` | |
| 581-XZ | Associated Prescription/Service Reference Number Qualifier | `x(2)` | |
| 582-X0 | Associated Prescription/Service Fill Number | `9(2)` | |
| B98-34 | Reconciliation ID | `x(30)` | |
| D16-KZ | Submission Type Code Count | `9(1)` | max 5 |
| D17-K8 | Submission Type Code | `x(2)` | repeating |
| D21-M3 | Multiple Prescription/Service Order Group ID | `x(35)` | |
| D22-M4 | Multiple Prescription/Service Order Group Reason Code | `x(2)` | |
| D02-KW | Total Prescribed Quantity Remaining | `9(7)v999` | |
| C99-KU | Preparation Environment Type | `x(2)` | |
| C98-KT | Preparation Environment Event Code | `9(3)` | |
| 678-Y6 | Time of Service | `9(6)` | |
| 681-ZF | Sales Transaction ID | `x(80)` | |
| A29-ZS | Reported Adjudicated Program Type | `x(3)` | |
| C02-4P | Original Manufacturer Product ID Qualifier | `x(2)` | |
| C01-4N | Original Manufacturer Product ID | `x(40)` | |
| C91-KK | LTPAC Dispense Frequency | `9(2)` | |
| C90-KH | LTPAC Billing Methodology | `9(1)` | |
| C92-KM | Number of LTPAC Dispensing Events | `9(3)` | |
| D18-K9 | Do Not Dispense Before Date | `9(8)` | |

**Coordination of Benefits / Other Payments (AM05)**

| Field | Name | Format | Notes |
|---|---|---|---|
| C47-9T | Other Payer Adjudicated Program Type | `x(3)` | |
| C49-9V | Other Payer Reconciliation ID | `x(30)` | |
| D51-P7 | Other Payer Percentage Tax Exempt Indicator | `x(1)` | |
| D53-P9 | Other Payer Regulatory Fee Type Count | `9(1)` | max 3 |
| D63-RN | Other Payer Regulatory Fee Type Code | `x(2)` | repeating |
| D52-P8 | Other Payer Regulatory Fee Exempt Indicator | `x(1)` | repeating |
| C50-9W | Benefit Stage Indicator Count | `9(1)` | max 4; replaces 392-MU |
| C51-9X | Benefit Stage Indicator | `x(2)` | repeating; replaces 393-MV |

**Workers' Compensation (AM06)**

| Field | Name | Format | Notes |
|---|---|---|---|
| B15-8D | Employer Street Address Line 1 | `x(40)` | replaces 316-CG |
| B16-7G | Employer Street Address Line 2 | `x(40)` | replaces 316-CG |
| B35-1V | Employer Country Code | `x(2)` | |
| B19-7K | Employer Telephone Number Extension | `9(8)` | |
| B17-7H | Employer Contact First Name | `x(35)` | replaces 321-CL |
| B18-7J | Employer Contact Last Name | `x(35)` | replaces 321-CL |
| B24-7R | Pay To Street Address Line 1 | `x(40)` | replaces 121-TV |
| B25-7S | Pay To Street Address Line 2 | `x(40)` | replaces 121-TV |
| B39-1Z | Pay To Country Code | `x(2)` | |

**Compound (AM10)**

| Field | Name | Format |
|---|---|---|
| C60-AG | Compound Level of Complexity | `x(2)` |

**Pricing (AM11)**

| Field | Name | Format | Notes |
|---|---|---|---|
| D60-RK | Regulatory Fee Count | `9(1)` | max 3 |
| D61-RL | Regulatory Fee Type Code | `x(2)` | repeating |
| D14-KX | Subrogation Amount Requested | `s9(9)v99` | |

**Prior Authorization (AM12)**

| Field | Name | Format | Notes |
|---|---|---|---|
| B13-7D | Authorized Representative Street Address Line 1 | `x(40)` | replaces 498-PG |
| B14-8B | Authorized Representative Street Address Line 2 | `x(40)` | replaces 498-PG |
| B34-1U | Authorized Representative Country Code | `x(2)` | |

**Facility (AM15)**

| Field | Name | Format | Notes |
|---|---|---|---|
| B95-3Z | Facility ID Qualifier | `x(2)` | |
| B20-7M | Facility Street Address Line 1 | `x(40)` | replaces 386-3U |
| B21-7N | Facility Street Address Line 2 | `x(40)` | replaces 386-3U |
| B37-1X | Facility Country Code | `x(2)` | |

### Response side

**Response Insurance (AM25)** — J7/J8 become repeating (max 3, counted by KR)

| Field | Name | Format | Notes |
|---|---|---|---|
| C96-KR | Payer/Health Plan ID Count | `9(1)` | max 3 J7/J8 pairs |

**Response Insurance Additional Information (AM27)**

| Field | Name | Format |
|---|---|---|
| 600-96 | Plan Name | `x(70)` |

**Response Patient (AM29)** — patient ID group is new here

| Field | Name | Format | Notes |
|---|---|---|---|
| 618-RR | Patient ID Count | `9(1)` | max 9 |
| 331-CX | Patient ID Qualifier | `x(2)` | repeating |
| 332-CY | Patient ID | `x(20)` | repeating |

**Response Status (AM21)**

| Field | Name | Format | Notes |
|---|---|---|---|
| B98-34 | Reconciliation ID | `x(30)` | |
| C72-BH | Help Desk Support Type Count | `9(2)` | max 5 |
| C71-BG | Help Desk Support Type | `x(2)` | repeating |
| C67-BB | Help Desk Business Unit Type Count | `9(2)` | max 5 |
| C66-BA | Help Desk Business Unit Type | `x(2)` | repeating |
| C70-BF | Help Desk Contact Information Qualifier | `x(2)` | replaces 549-7F |
| C68-BC | Help Desk Contact Information | `x(255)` | replaces 550-8F and 987-MA |
| C69-BD | Help Desk Contact Information Extension | `x(8)` | |
| A28-ZR | Adjudicated Program Type | `x(3)` | |
| B04-BT | Next Available Fill Date | `9(8)` | |

**Response Claim (AM22)**

| Field | Name | Format | Notes |
|---|---|---|---|
| E89-ZO | Formulary Alternative Effective Date | `9(8)` | in Formulary Alternative group |
| D42-PV | Formulary Alternative Plan Benefit Tier | `x(2)` | " |
| D43-PZ | Formulary Alternative Reason Code | `x(2)` | " |
| D44-P0 | Formulary Alternative Required Therapy Indicator Count | `x(2)` | max 6 |
| D45-P1 | Formulary Alternative Required Therapy Indicator | `x(2)` | repeating |
| D46-P2 | Formulary Alternative Required Therapy Time Period Qualifier | `x(2)` | |
| D47-P3 | Formulary Alternative Required Therapy Time Period Duration | `x(3)` | |
| D48-P4 | Formulary Alternative Required Therapy Time Period Start Date | `9(8)` | |
| D49-P5 | Formulary Alternative Required Therapy Time Period End Date | `9(8)` | |
| D54-RC | Plan Benefit Override Indicator | `x(6)` | |
| D55-RD | Plan Benefit Override Value Count | `9(1)` | max 9 |
| D56-RF | Plan Benefit Override Value | `x(10)` | repeating |
| 931-F8 | Maximum Age Qualifier | `x(1)` | |
| 932-GA | Maximum Age | `9(3)` | |
| 943-GQ | Minimum Age Qualifier | `x(1)` | |
| 944-GR | Minimum Age | `9(3)` | |
| D20-M2 | Minimum Amount Qualifier | `x(3)` | |
| D19-M1 | Minimum Amount | `9(10)v999` | |
| 934-GC | Maximum Amount Qualifier | `x(2)` | |
| 933-GB | Maximum Amount | `9(10)v999` | |
| 935-GF | Maximum Amount Time Period | `x(2)` | |
| 936-GG | Maximum Amount Time Period Start Date | `9(8)` | |
| 937-GH | Maximum Amount Time Period End Date | `9(8)` | |
| 938-GJ | Maximum Amount Time Period Units | `9(4)` | |
| D25-M7 | Remaining Amount Qualifier | `x(3)` | |
| D24-M6 | Remaining Amount | `9(10)v999` | |
| C59-AF | Benefit Type Opportunity Count | `9(1)` | max 9 |
| C58-AE | Benefit Type Opportunity | `x(2)` | repeating |
| D15-KY | Subrogation Requestor's Reconciliation ID | `x(30)` | replaces 114-N4 |

**Response Pricing (AM23)**

| Field | Name | Format | Notes |
|---|---|---|---|
| C94-KP | Patient Pay Component Count | `9(2)` | max 25 |
| C95-KQ | Patient Pay Component Qualifier | `x(2)` | repeating; qualifier values replace the 12 removed "Amount Attributed to…" fields |
| C93-KN | Patient Pay Component Amount | `s9(9)v99` | repeating |
| D60-RK | Regulatory Fee Count | `9(1)` | |
| D61-RL | Regulatory Fee Type Code | `x(2)` | repeating |
| D62-RM | Regulatory Fee Exempt Indicator | `x(1)` | repeating |
| C50-9W | Benefit Stage Indicator Count | `9(1)` | replaces 392-MU |
| C51-9X | Benefit Stage Indicator | `x(2)` | repeating; replaces 393-MV |
| D65-RS | Patient Regulatory Fee Amount | `s9(9)v99` | |
| B33-6G | Professional Service Fee Contracted/Reimbursable Amount | `s9(9)v99` | |

**Response DUR/PPS (AM24)** — all replacing free-text 570-NS DUR Additional Text

| Field | Name | Format |
|---|---|---|
| 475-J9 | DUR/DUE Co-Agent ID Qualifier | `x(2)` |
| 476-H6 | DUR/DUE Co-Agent ID | `x(40)` |
| E93-ZC | DUR/DUE Co-Agent Description | `x(40)` |
| F01-Z4 | Other Prescriber ID Qualifier | `x(2)` |
| F02-Z3 | Other Prescriber ID | `x(35)` |
| E99-Z5 | Other Prescriber Last Name | `x(35)` |
| F03-Z2 | Other Prescriber Telephone Number | `9(10)` |
| E94-ZA | Unit of Measure for Previous Dispensed Quantity | `x(3)` |
| F05-Z0 | DUR/DUE Compound Product ID Qualifier | `x(2)` |
| F04-Z1 | DUR/DUE Compound Product ID | `x(40)` |
| F06-YO | DUR/DUE Maximum Daily Dose Quantity | `9(7)v999` |
| F07-YL | DUR/DUE Maximum Daily Dose Unit of Measure | `x(3)` |
| F08-YJ | DUR/DUE Minimum Daily Dose Quantity | `9(7)v999` |
| F09-YI | DUR/DUE Minimum Daily Dose Unit of Measure | `x(3)` |

(Request DUR/PPS also gains Other Pharmacy fields Z6–Z9 per the library/schema; the response side additionally carries Other Pharmacy ID/Name/Telephone via the same DUR item group.)

**Response Other Payers (AM28, renamed from Response Coordination of Benefits)**

| Field | Name | Format |
|---|---|---|
| 339-6C | Other Payer ID Qualifier | `x(2)` |
| D41-PQ | Other Payer Relationship Type | `x(3)` |
| D50-P6 | Other Payer Benefit Classification | `x(10)` |
| C47-9T | Other Payer Adjudicated Program Type | `x(3)` |
| D23-M5 | Other Payer Name | `x(70)` |
| B23-7Q | Other Payer Help Desk Telephone Number Extension | `9(8)` |

**Response Prior Authorization (AM26)**

| Field | Name | Format |
|---|---|---|
| 498-PT | Prior Authorization Expiration Date | `9(8)` |

## 7. New segments (F6-only) with full field definitions

### Request: Intermediary Segment (AM19)
Used on: Billing, Rebill, Prior Authorization, Service Billing, Service Rebill.

| Field | Name | Format | Notes |
|---|---|---|---|
| B44-8G | Intermediary ID Count | `9(1)` | max 9 |
| B45-8H | Intermediary ID Type Code | `x(2)` | repeating group |
| B46-8J | Intermediary ID Type Entity | `x(2)` | " |
| B47-8K | Intermediary ID Qualifier | `x(2)` | " |
| B48-8M | Intermediary ID | `x(20)` | " |
| B49-8N | Intermediary ID State/Province Address | `x(2)` | " |
| B55-8U | Intermediary ID Country Code | `x(2)` | " |

### Request: Last Known 4Rx Segment (AM37)
Used on: Eligibility (E1).

| Field | Name | Format |
|---|---|---|
| B56-3E | Last Known IIN Number | `9(8)` |
| B57-3F | Last Known Processor Control Number | `x(10)` |
| B58-3G | Last Known Group ID | `x(15)` |
| B59-3H | Last Known Cardholder ID | `x(20)` |
| B60-3J | Year of Last Paid Claim | `9(4)` |
| B61-3K | Month of Last Paid Claim | `9(2)` |

### Request: N Transaction Payer Identification Segment (AM38)
Used on: Information Reporting (N1/N2/N3).

| Field | Name | Format |
|---|---|---|
| C52-9Y | N Payer IIN | `9(8)` |
| C53-9Z | N Payer Processor Control Number | `x(10)` |
| C54-AA | N Payer Cardholder ID | `x(20)` |
| C55-AB | N Payer Group ID | `x(15)` |
| C48-9U | N Payer Adjudicated Program Type | `x(3)` |
| C56-AC | N Transaction Source Type | `9(2)` |
| C57-AD | N Transaction Reconciliation ID | `x(30)` |
| 880-K5 | Transaction Reference Number | `x(30)` (moved here from D.0 Claim segment) |

### Response: Intermediary Segment (AM36)
Used on: Billing, Prior Authorization, Rebill, Service Billing, Service Rebill responses.

| Field | Name | Format | Notes |
|---|---|---|---|
| B52-8R | Response Intermediary Authorization Count | `9(1)` | |
| B53-8S | Response Intermediary Authorization Type ID | `9(2)` | repeating group |
| B54-8T | Response Intermediary Authorization ID | `x(20)` | " |
| B51-8Q | Intermediary Message | `x(200)` | repeating |
| C86-KC | Intermediary Help Desk Support Type Count | `9(2)` | |
| C85-KB | Intermediary Help Desk Support Type | `x(2)` | repeating |
| C81-G9 | Intermediary Help Desk Business Unit Type Count | `9(2)` | |
| C80-G8 | Intermediary Help Desk Business Unit Type | `x(2)` | repeating group |
| C84-KA | Intermediary Help Desk Contact Information Qualifier | `x(2)` | " |
| C82-JP | Intermediary Help Desk Contact Information | `x(255)` | " |
| C83-JR | Intermediary Help Desk Contact Information Extension | `x(8)` | " |

### Response: Other Related Benefit Detail Segment (AM39)
Used on: Billing, Eligibility, Predetermination of Benefits, Prior Authorization, Prior Authorization Inquiry, Rebill, Service Billing, Service Rebill responses.

| Field | Name | Format | Notes |
|---|---|---|---|
| C97-KS | Plan Type | `x(4)` | |
| C88-KF | LIS Level | `9(1)` | replaces 138-UQ CMS LIS Level |
| C87-KD | LIS Effective Date | `9(8)` | |
| C89-KG | LIS Termination Date | `9(8)` | |
| C61-AH | Disability Effective Date | `9(8)` | |
| C63-A5 | ESRD Indicator | `x(2)` | |
| C62-AJ | ESRD Effective Date | `9(8)` | |
| C64-A6 | ESRD Termination Date | `9(8)` | |
| C76-G4 | Hospice Effective Date | `9(8)` | |
| C79-G7 | Hospice Termination Date | `9(8)` | |
| C78-G6 | Hospice Provider Number | `x(15)` | |
| C77-G5 | Hospice Facility Name | `x(70)` | |
| C65-A8 | Hospice Telephone Number | `9(10)` | |
| C73-BJ | Institutional Indicator | `x(2)` | |
| C74-BK | Institutional Effective Date | `9(8)` | |
| C75-GD | Institutional Termination Date | `9(8)` | |
| D26-M8 | Other Benefit Count | `9(1)` | starts repeating group |
| D40-PN | Other Benefit Type Code | `x(2)` | group member |
| D34-MZ | Other Benefit Effective Date | `9(8)` | " |
| D39-NN | Other Benefit Termination Date | `9(8)` | " |
| 729-TA | State/Province Address | `x(2)` | " |
| D38-N9 | Other Benefit Type ID | `x(15)` | " |
| D35-N1 | Other Benefit Facility Name | `x(70)` | " |
| D36-N7 | Other Benefit Facility Telephone Number | `9(10)` | " |
| D37-N8 | Other Benefit Detail Information Count | `9(1)` | starts second group |
| D32-MS | Other Benefit Detail Information Indicator | `x(2)` | group member |
| D28-MM | Other Benefit Detail Information Effective Date | `9(8)` | " |
| D33-MX | Other Benefit Detail Information Termination Date | `9(8)` | " |
| D31-MR | Other Benefit Detail Information Provider Number | `x(15)` | " |
| D29-MN | Other Benefit Detail Information Facility Name | `x(70)` | " |
| D30-MP | Other Benefit Detail Information Facility Telephone Number | `9(10)` | " |

### Response: Provider Segment (AM40)
Used on: Billing, Predetermination of Benefits, Prior Authorization, Prior Authorization Inquiry, Prior Authorization Request Only, Rebill, Service Billing, Service Rebill responses.

| Field | Name | Format |
|---|---|---|
| E87-ZV | Data Source of Invalid Provider Determination | `x(2)` |
| E88-ZZ | State Code for Data Source of Invalid Provider Determination | `x(2)` |

## 8. Repeating-count (max occurrence) values noted for F6

Submission Clarification Code 5 · Submission Type Code 5 · Procedure Modifier Code 10 · Diagnosis Code 5 · Clinical Information Counter 5 · DUR/PPS Code Counter 9 · COB Other Payments 3 · Other Payer Amount Paid 9 · Other Payer-Patient Responsibility 25 · Compound Ingredients 25 · Compound Ingredient Modifier 10 · Question Number/Letter 50 · Patient IDs 9 · Payer/Health Plan IDs 3 · Regulatory Fees 3 · Benefit Stage Indicators 4 · Patient Pay Components 25 · Intermediary IDs 9.

## 9. Value-level changes (out of scope here, tracked separately)

F6 also adds/removes **code values** on many existing fields (e.g., new qualifier values, Patient Pay Component qualifiers replacing the removed itemized amount fields). These are tracked in the "New Field Values" sheet of `vF6 Planning.xlsx` (~560 entries) and the External Code List.

## 10. NCPDPSerDe implementation notes

- **All spec-new F6 fields above are implemented** (verified: every "New" field code in the planning workbook exists in the library's struct tags).
- **Removed/renamed D.0 fields are retained** under their D.0 names so D.0 traffic and existing consumers are unaffected; on F6 they're simply absent (renames are the same field number/code on the wire, so parsing is unaffected).
- Repeating patient IDs (CX/CY) and payer IDs (J7/J8) are dual-mapped: D.0-compatible scalars (last occurrence) plus F6 slices `Patient.Ids` / `Insurance.Payers` (all occurrences).
- Body fields are parsed as variable-length delimited values, so the §3 length changes affect validation/truncation concerns in consuming apps, not this library's parsing — with two exceptions the library handles explicitly, both version-resolved from the transmission being parsed/built:
  - **448-ED Compound Ingredient Quantity** changed its implied decimal count (3 → 7). The struct tag carries `decimalPlaces=3,decimalPlacesSince=F6:7`, so D0 divides by 10³ and F6 by 10⁷.
  - **483-HE / 560-AY Percentage Tax Rate** became unsigned in F6. The tags carry `overpunch=true,unsignedSince=F6`: parsing stays lenient either way, but F6 serialization emits plain digits with no overpunch character while D0 keeps the signed form.
  - **379-4D Question Percent Response** is unsigned `9(3)v99` in the F6 dictionary and carries the same `unsignedSince=F6` treatment. (The D0 dictionary CSV's format column also says unsigned and no D0→F6 change is tracked for it, so the library's historical D0 overpunch on 4D is likely a long-standing quirk; D0 output is left byte-identical to production behavior pending verification.)
- See `F6_MIGRATION.md` for app-migration guidance (weakly-typed `Records[0]` gotcha, helper methods, serializer rules).
