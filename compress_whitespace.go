package fastxml

type compressWhitespace struct {
	buf      Writer
	lastChar byte
}

func newCompressWhitespace(buf Writer) *compressWhitespace {
	return &compressWhitespace{buf: buf}
}

func (ir *compressWhitespace) WriteString(s string) (int, error) {
	return ir.write([]byte(s))
}

func (ir *compressWhitespace) WriteByte(b byte) error {
	if !whitespace[b] || ir.lastChar != '>' {
		ir.lastChar = b
		return ir.buf.WriteByte(b)
	}
	return nil
}

func (ir *compressWhitespace) Write(p []byte) (int, error) {
	return ir.write(p)
}

func (ir *compressWhitespace) write(b []byte) (int, error) {
	head, tail := 0, 0

	if ir.lastChar == '>' {
		for ; head < len(b) && whitespace[b[head]]; head++ {
			//skip whitespaces
		}
		tail = head
	}

	for head < len(b) {
		if b[head] != '>' || ((head+1) < len(b) && !whitespace[b[head+1]]) {
			head++
			continue
		}
		ir.lastChar = b[head]

		head++ ///pointing to whitespace
		if i, err := ir.buf.Write(b[tail:head]); err != nil {
			return tail + i, err
		}

		for ; head < len(b) && whitespace[b[head]]; head++ {
			//skip whitespaces
		}
		tail = head
	}
	if tail < head {
		if i, err := ir.buf.Write(b[tail:head]); err != nil {
			return tail + i, err
		}
		ir.lastChar = b[head-1]
	}
	return len(b), nil
}
