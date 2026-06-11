package claimserializer

import (
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
				// Serialize the shared object (concurrent reads of one struct)
				got, err := Serialize(&shared)
				if err != nil {
					t.Errorf("Failed to serialize shared object: %v", err)
					return
				}
				if got != baseline {
					t.Errorf("Serialized output mismatch under concurrency.\nWanted: %q\nGot:    %q", baseline, got)
					return
				}

				// Full round trip with a goroutine-local object
				local := request.Billing{}
				if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &local); err != nil {
					t.Errorf("Failed to deserialize local object: %v", err)
					return
				}
				got, err = Serialize(&local)
				if err != nil {
					t.Errorf("Failed to serialize local object: %v", err)
					return
				}
				if got != baseline {
					t.Errorf("Local round-trip output mismatch under concurrency.\nWanted: %q\nGot:    %q", baseline, got)
					return
				}
			}
		}()
	}

	wg.Wait()
}
