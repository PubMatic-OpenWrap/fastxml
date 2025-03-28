package fastxml

import (
	"bytes"
	"sync"
)

var (
	bufferedPool      *sync.Pool
	xmlOperationsPool *sync.Pool
)

func init() {
	bufferedPool = newBufferedPool()
	xmlOperationsPool = newXMLOperationsPool()
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

//--------------------------------------------------------------------------------------------

// newXMLOperationsPool returns the singleton instance of bufferedPool
func newXMLOperationsPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return make([]xmlOperation, 0, 10)
		},
	}
}

func getXMLOperations() []xmlOperation {
	ops := xmlOperationsPool.Get().([]xmlOperation)
	return ops[:0]
}

func putXMLOperations(ops []xmlOperation) {
	ops = ops[:0] // Reset before putting it back
	xmlOperationsPool.Put(ops[:])
}
