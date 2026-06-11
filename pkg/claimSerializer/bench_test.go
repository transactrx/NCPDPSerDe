package claimserializer

import (
	"testing"

	claimdeserializer "github.com/transactrx/NCPDPSerDe/pkg/claimDeserializer"
	"github.com/transactrx/NCPDPSerDe/pkg/ncpdp/request"
)

func Benchmark_SerializeB1(b *testing.B) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(REQUEST_B1, &obj); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Serialize(&obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_SerializeF6(b *testing.B) {
	obj := request.Billing{}
	if err := claimdeserializer.DeserializeType(buildF6BillingRequest(), &obj); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Serialize(&obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}
