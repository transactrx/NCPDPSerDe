package claimserializer

import (
	"fmt"
	"sync"
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// Serialization shares cached field-code sets and segment definitions across
// goroutines, so concurrent serialization of the same object must be safe and
// always produce identical output. Run with -race to catch data races.
func Test_SerializeIsConcurrencySafe(t *testing.T) {
	shared := request.Billing{}
	if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &shared); err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	baseline, err := Serialize(&shared)
	if err != nil {
		t.Fatalf("Failed to serialize baseline: %v", err)
	}

	const goroutines = 10
	const iterations = 25

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				if !serializeSharedConcurrently(t, &shared, baseline) || !roundTripLocalConcurrently(t, baseline) {
					return
				}
			}
		}()
	}

	wg.Wait()
}

// Serialize the shared object (concurrent reads of one struct).
func serializeSharedConcurrently(t *testing.T, shared *request.Billing, baseline string) bool {
	t.Helper()

	got, err := Serialize(shared)
	if err != nil {
		t.Errorf("Failed to serialize shared object: %v", err)
		return false
	}

	if got != baseline {
		t.Errorf("Serialized output mismatch under concurrency.\nWanted: %q\nGot:    %q", baseline, got)
		return false
	}

	return true
}

// Full round trip with a goroutine-local object.
func roundTripLocalConcurrently(t *testing.T, baseline string) bool {
	t.Helper()

	local := request.Billing{}
	if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &local); err != nil {
		t.Errorf("Failed to deserialize local object: %v", err)
		return false
	}

	got, err := Serialize(&local)
	if err != nil {
		t.Errorf("Failed to serialize local object: %v", err)
		return false
	}

	if got != baseline {
		t.Errorf("Local round-trip output mismatch under concurrency.\nWanted: %q\nGot:    %q", baseline, got)
		return false
	}

	return true
}

// Where Test_SerializeIsConcurrencySafe exercises one shared F6 object, this
// test hammers the full NCPDP -> model -> JSON -> model -> NCPDP chain from
// many goroutines over DISTINCT transaction shapes at once, so the shared
// reflection caches (typeOwnsCode, field codes, segment/field attributes) see
// interleaved first-populations for different types. Every result must be
// byte-identical to the sequential baseline: a race or cache
// cross-contamination shows up as a mismatch, and `go test -race` pins the
// access pattern itself.
func Test_ConcurrentRoundTripsMatchSequentialBaseline(t *testing.T) {
	samples := []string{
		REQUEST_B1, REQUEST_B2, REQUEST_B3, REQUEST_E1,
		RESPONSE_B2, RESPONSE_B3, RESPONSE_E1,
	}

	baseline := make([]string, len(samples))
	for i, raw := range samples {
		out, err := fullRoundTrip(raw)
		if err != nil {
			t.Fatalf("sample %d baseline round trip failed: %v", i, err)
		}
		baseline[i] = out
	}

	const workers = 8
	const iterations = 50

	errs := make(chan error, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for iter := 0; iter < iterations; iter++ {
				for i, raw := range samples {
					out, err := fullRoundTrip(raw)
					if err != nil {
						errs <- fmt.Errorf("worker %d iter %d sample %d: %w", worker, iter, i, err)
						return
					}
					if out != baseline[i] {
						errs <- fmt.Errorf("worker %d iter %d sample %d: output diverged from sequential baseline", worker, iter, i)
						return
					}
				}
			}
		}(w)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

func fullRoundTrip(raw string) (string, error) {
	obj, err := claimdeserializer.Deserialize(raw)
	if err != nil {
		return "", err
	}
	return jsonRoundTripSerialize(obj)
}
