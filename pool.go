package fastxml

import (
	"bytes"
	"sync"
)

var (
	bufferedPool *sync.Pool
)

func init() {
	bufferedPool = newBufferedPool()
}

// newBufferedPool returns the singleton instance of bufferedPool
func newBufferedPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

func getBuffer() *bytes.Buffer {
	buf := bufferedPool.Get().(*bytes.Buffer)
	buf.Reset() // Ensure the buffer is clean before use
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	buf.Reset() // Reset before putting it back
	bufferedPool.Put(buf)
}
