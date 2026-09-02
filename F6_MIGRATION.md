# F6 (vF6) Support — Migration Notes

The library now parses and serializes both D.0 and F6 NCPDP telecommunication formats transparently — version is auto-detected from the header. No D0 code paths changed; all existing applications continue to work without modification. Below is what you need to know if you're updating an app to handle F6.

---

## ⚠️ The one breaking-shaped gotcha: `tran.Records[0]` is empty for F6

If you use the **weakly-typed** `NcpdpTransaction[V]` parser (i.e. `NewTransactionRequest` / `NewTransactionResponse` + `ParseNcpdp`), any code that reaches into `tran.Records[0]` directly will **silently misbehave on F6** — and may panic on out-of-range index.

**Why:** F6 transmissions have no group separator (`0x1D`). The body parser only emits a `NcpdpRecord` when it hits `0x1D`, so for F6 every segment lands in `tran.Segments` and `tran.Records` stays empty (`len == 0`).

**What to do — replace direct indexing with the F6-aware helpers:**

```go
// ❌ Will misbehave or panic on F6
seg := tran.Records[0].FindSegment(ncpdp.RESPONSE_STATUS_SEGMENT_ID)
for _, record := range tran.Records {
    seg := record.FindSegment(ncpdp.RESPONSE_STATUS_SEGMENT_ID)
    ...
}

// ✅ Works for both D0 and F6
seg := tran.FindSegmentInRecord(0, ncpdp.RESPONSE_STATUS_SEGMENT_ID)
for i := 0; i < tran.RecordCount(); i++ {
    seg := tran.FindSegmentInRecord(i, ncpdp.RESPONSE_STATUS_SEGMENT_ID)
    ...
}
```

These helpers fall back to the shared-segment slice (`tran.Segments`) and the header's `RecordCount` field when `Records` is empty. Built-in response helpers — `Status()`, `IsPaid()`, `IsRejected()`, `IsStatusOf()`, `GetRejectCodes()`, `GetAdditionalMessages()` — already use the fallbacks internally, so if your app calls those it's already F6-safe.

> **Note:** The **strongly-typed** path (`claimDeserializer.DeserializeRequest` returning a `BillingRequest`, etc.) does populate `Claims[0]` for F6 via `evaluateUngroupedTransaction`. So the gotcha is specific to the weakly-typed `NcpdpTransaction[V]` parser. If you use both paths, audit each separately.

---

## Header (`pkg/ncpdp/header.go`)

- `RequestHeader` carries dual struct tags: `layout` (D0, 56 bytes) and `layoutF6` (F6, 58 bytes). Layout is auto-selected from the first 2 bytes of the raw header — D0 leads with the BIN (digits), F6 leads with `"F6"`.
- F6 field 101-A1 is officially called **IIN** (8 chars) but reuses the existing `Bin` string field — no rename, no new field.
- `ResponseHeader` is identical (31 bytes) for D0 and F6 — no F6 tags needed.
- `FinancialRequestHeader` / `FinancialResponseHeader` are unchanged (no F6 layout).
- New constant: `ncpdp.F6 = "F6"`.

## Transmission shape (vEB+)

- F6 eliminated the group separator (`0x1D`).
- F6 permits exactly **one transaction per transmission**. The serializer returns an error if you try to build with `>1` claim group on an F-prefixed header.
- The serializer auto-omits `0x1D` when the header starts with `"F"` (`useGroupSeparator(version)` in `transaction.go`; `omitGroupSeparators` in `claimSerializer.go`).

## Repeating-field dual mapping (backward compat)

F6 allows certain fields to repeat where D0 had a single occurrence. To preserve backward compatibility (no renames, no retypes) AND avoid data loss, both shapes are kept:

| Segment | D0 scalar (preserved) | F6 slice (new) | Codes |
|---|---|---|---|
| `request.Patient` | `IdQualifier`, `Id` | `Ids []PatientId` | CX, CY |
| `response.Insurance` | `Payer` (struct) | `Payers []Payer` | J7, J8 |

**Read semantics:** When parsing F6 with repeats, the scalar holds the **last** occurrence and the slice holds **all** occurrences. If you care about complete F6 data, read the slice.

**Write semantics:** The serializer (`mergeSliceFieldCodes` in `claimSerializer.go`) suppresses the scalar's codes when the sibling slice is non-empty, so round-trips don't double-emit. Populate either form — don't populate both with the same value, or only the slice will be written.

**Occurrence counts are F6-only and version-gated.** Each group has an F6-added counter — `KR` (Payer/Health Plan ID Count) on `response.Insurance`, `RR` (Patient ID Count) on `request.Patient` / `response.Patient` (AM29). These are `countfor`-derived: the serializer sets them from the slice length automatically (so an omitted or stale count is corrected on build). Because the backing slices are always populated for D0 scalar payloads (the scalar dual-maps into the slice), these counters carry the `sinceVersion=F6` field tag: the serializer omits them entirely for D0 output — even when the count field is set on the struct — since `KR`/`RR` do not exist in D0. For F6, the count is emitted as usual. This is what keeps a D0 deserialize→serialize round-trip from injecting a `KR`/`RR` field the payer never sent.

The `sinceVersion` tag is the general mechanism for this: any field scoped to a newer version is skipped when serializing an older transmission. Version ordering lives in one place — the rank table in `pkg/ncpdp/version.go` (`VersionAtLeast`, `OmitsGroupSeparator`, `HeaderLeadsWithVersion`) — so a field scoped to F6 is emitted for F6 and every later version. Supporting a future NCPDP version means adding one rank entry there (plus header layout tags and any new `sinceVersion` scoping).

**Every F6-only field now carries the tag.** All fields marked `" - F6"` in the JSON schema have `sinceVersion=F6` in their `field` tag, so consumers that reflect over the structs (rule engines, field catalogs, doc generators) can distinguish D0 fields from F6 additions programmatically, and a hand-populated F6-only field can never leak into D0 output. Two tag placements exist:

- **Code-bearing fields** (leaves and `countfor` counters): the serializer omits them from D0 output entirely.
- **Code-less tags on repeating-group slices** (`request.Patient.Ids`, `request.Pricing.RegulatoryFees`, `response.Insurance.Payers` — i.e. `field:"sinceVersion=F6"` with no `code=`): metadata only. The serializer does NOT skip these for D0, because dual-mapped groups carry D0 data — the scalar's codes (CX/CY, J7/J8) are written from the slice. F6-only leaves inside such groups carry their own leaf tag.

## Strongly-typed path (`pkg/claimDeserializer`, `pkg/claimSerializer`)

- `DeserializeRequest` infers the transaction code from byte 2 (F6) vs byte 8 (D0) automatically.
- For F6 raw with no `0x1D`, `evaluateUngroupedTransaction` locates the single claim group by scanning for the first segment ID belonging to the group struct; that segment plus everything after it becomes `Claims[0]`. Everything before becomes shared segments.
- Round-trip is symmetric — F6 in → F6 out, D0 in → D0 out — no caller-side switching needed.

## New F6-only fields

The JSON schema (`pkg/ncpdp/schemas/ncpdp-schemas.json`) gained **161 F6-only fields** across Patient, Prescriber, Claim, Pricing, Coordination of Benefits, and other segments. All are tagged `" - F6"` in their description and are **optional pointers** — D0 transactions leave them nil. Existing D0 fields are untouched, so existing schema consumers keep working.

## Useful new transaction helpers

- `tran.Version() string` — reads the version from the parsed header.
- `tran.RecordCount() int` — record count, falling back to the header's count when `Records` is empty (F6).
- `tran.FindSegmentInRecord(recordIdx int, segId string) *NcpdpSegment` — the F6-safe replacement for `tran.Records[i].FindSegment(...)`. With `recordIdx == 0`, falls back to shared segments when `Records` is empty.

## What did NOT change

- Public API signatures — no breaking changes.
- D0 parsing/serialization output — byte-for-byte identical.
- All pre-existing fields — same names, same types, same positions.
- Field codes — D0's existing codes (BIN, AM segments, etc.) are unchanged.

## Migration checklist for existing apps

1. **Search your codebase for `.Records[0]`**, `.Records[i]`, and `range tran.Records` patterns on `NcpdpTransaction[V]`. Replace with `FindSegmentInRecord` / `RecordCount()` loops.
2. **If you read patient IDs or insurance payers**, decide whether you need the F6 repeating data — switch to `Patient.Ids` / `Insurance.Payers` slices, or accept that the scalar holds only the last occurrence for F6 multi-ID payloads.
3. **If you build transmissions manually**, you can keep building D0 the same way. For F6, set the header version to `"F6"` and the BIN/IIN field width to 8 — the serializer handles the rest.
4. **If you guard transaction count**, F6 → 1 transaction per transmission, enforced at serialize time.
5. **Test against a real F6 sample** — `D:\POW\NCPDP Files\F6 Format\RawSample_F6.txt` and the spec PDF in the same folder are the references the library was tested against.

## Reference tests to read

If you want canonical examples of the parsed shape:

- `pkg/ncpdp/header_test.go` — F6 header round-trip
- `pkg/claimDeserializer/claimDeserializer_test.go` — `Test_CanParseF6BillingRequest`, `Test_CanParseF6FullFieldSample`, `Test_F6UnknownSegmentsRespectGroupBoundary`
- `pkg/claimSerializer/claimSerializer_test.go` — `Test_CanRoundTripF6BillingRequest`, `Test_F6SerializeRejectsMultipleClaimGroups`
- `pkg/ncpdp/response_test.go` — `TestResponseHelpersAcrossVersions` (proves the helpers work for both D0 records-shaped and F6 segments-shaped responses)
