package shenmeci

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	p := filepath.Join("..", "..", "dict", "cedict_1_0_ts_utf-8_mdbg.txt.gz")
	fi, err := os.Stat(p)
	if err != nil {
		b.Skipf("CEDICT file %q not available: %v", p, err)
	}
	if fi.IsDir() {
		b.Skipf("CEDICT file %q not available: wanted file, got directory", p)
	}
	GlobalConfig.CedictPath = p
	for i := 0; i < b.N; i++ {
		LoadCEDICT()
	}
}
