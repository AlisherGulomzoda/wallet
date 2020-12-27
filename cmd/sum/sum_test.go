package sum

import (
	"testing"
)

func BenchmarkRegular(b *testing.B) {
	want := int64(4_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Regular()
		b.StopTimer()
		if want != result {
			b.Errorf("invalide result: got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

func BenchmarkConcurently(b *testing.B) {
	want := int64(4_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Concurently()
		b.StopTimer()
		if want != result {
			b.Errorf("invalide result: got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}