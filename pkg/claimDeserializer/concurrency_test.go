package claimdeserializer

import (
	"sync"
	"testing"

	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

// The serde attribute/definition caches are shared across goroutines, so
// deserialization must be safe and deterministic under concurrent use.
// Run with -race to catch data races.
func Test_DeserializeIsConcurrencySafe(t *testing.T) {
	rawF6 := buildF6FullBillingRequest()

	const goroutines = 10
	const iterations = 25

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				obj := request.Billing{}
				if err := DeserializeType(REQUEST_B1, &obj); err != nil {
					t.Errorf("Failed to deserialize B1: %v", err)
					return
				}
				assertString(t, "Cardholder id under concurrency", obj.Insurance.Cardholder.Id, "POLICYNUMBERTHATISLO")
				if len(obj.Claims) != 1 {
					t.Errorf("Group count mismatch under concurrency. Wanted: 1   Got: %v", len(obj.Claims))
					return
				}

				res, err := Deserialize(rawF6)
				if err != nil {
					t.Errorf("Failed to deserialize F6: %v", err)
					return
				}
				f6, ok := res.(request.Billing)
				if !ok {
					t.Errorf("Expected request.Billing, got: %T", res)
					return
				}
				assertF6FullSampleHeader(t, f6)
				assertF6FullSampleInsurance(t, f6)
				assertF6FullSamplePatient(t, f6)
				if len(f6.Claims) != 1 {
					t.Errorf("F6 group count mismatch under concurrency. Wanted: 1   Got: %v", len(f6.Claims))
					return
				}
				assertF6FullSamplePricing(t, f6)
			}
		}()
	}

	wg.Wait()
}
