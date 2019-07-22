package shenmeci

import (
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadCEDICT()
	}
}

func init() {
	LoadConfig()
	ValidateConfig()
	LoadCEDICT()
}
