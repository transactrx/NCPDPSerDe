package claimdeserializer

import "testing"

func Benchmark_DeserializeB1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Deserialize(REQUEST_B1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_DeserializeF6Full(b *testing.B) {
	raw := buildF6FullBillingRequest()
	for i := 0; i < b.N; i++ {
		_, err := Deserialize(raw)
		if err != nil {
			b.Fatal(err)
		}
	}
}
