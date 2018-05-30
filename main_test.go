package chomper

import (
	"testing"
)

func BenchmarkGet(b *testing.B) {
	bs := New(b.N)

	for i := 0; i < b.N; i++ {
		bs.Get(i)
	}
}
