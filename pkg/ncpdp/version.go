package ncpdp

// versionRank orders the NCPDP telecommunication versions this library
// supports. Higher rank = newer version. Supporting a new version REQUIRES an
// entry here — plus header layout tags (see resolveLayoutTagName) and
// sinceVersion field scoping where applicable. A version absent from this
// table ranks below every known version, which preserves the library's
// historical behavior of treating unrecognized headers as D0-era.
var versionRank = map[string]int{
	D0: 10,
	F6: 20,
}

// KnownVersion reports whether the version is one this library models.
func KnownVersion(version string) bool {
	_, known := versionRank[version]
	return known
}

// VersionAtLeast reports whether version is the same as or newer than minimum.
// An unknown version ranks oldest; an unknown minimum can never be satisfied,
// so a field scoped to a version this build does not know is always omitted.
func VersionAtLeast(version, minimum string) bool {
	minRank, known := versionRank[minimum]
	if !known {
		return false
	}
	return versionRank[version] >= minRank
}

// OmitsGroupSeparator reports whether the version eliminated the 0x1D group
// separator. vEB removed it; F6 is the first such version this library
// models, and every later version inherits the change.
func OmitsGroupSeparator(version string) bool {
	return VersionAtLeast(version, F6)
}

// HeaderLeadsWithVersion reports whether the version's header layout places
// the version code in the first two bytes (the vEB+ redesign; older request
// headers lead with the 6-digit BIN instead).
func HeaderLeadsWithVersion(version string) bool {
	return VersionAtLeast(version, F6)
}

// DetectTransmissionVersion extracts the NCPDP version from a raw
// transmission or header. vEB+ headers (and D0 response headers) lead with
// the version code; D0-era request headers carry it after the 6-digit BIN.
// Returns Empty when no known version is found — callers treat that as
// oldest-known (D0-era) behavior for backward compatibility.
func DetectTransmissionVersion(raw string) string {
	if len(raw) >= 2 && KnownVersion(raw[:2]) {
		return raw[:2]
	}
	if len(raw) >= 8 && KnownVersion(raw[6:8]) {
		return raw[6:8]
	}
	return Empty
}
