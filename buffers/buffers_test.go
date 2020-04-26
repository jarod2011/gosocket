package buffers

import (
	"io"
	"testing"
)

func TestBuffers(t *testing.T) {
	buf := New()
	b := buf.Get()
	n, err := b.Write([]byte{0x11, 0x22, 0x33, 0x44})
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("err write len %d", 4)
	}
	buf.Put(b)
	b1 := buf.Get()
	if b1.Len() == 4 {
		t.Fatalf("buf len %v", b1.Len())
	}
}

func BenchmarkBuffers_Get(b *testing.B) {
	buf := New()
	for i := 1; i < b.N; i++ {
		m := buf.Get()
		buf.Put(m)
	}
}
