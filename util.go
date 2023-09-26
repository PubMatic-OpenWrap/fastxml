package xmlparser

import "bytes"

var whitespace [256]bool

func init() {
	whitespace[' '] = true
	whitespace['\r'] = true
	whitespace['\n'] = true
	whitespace['\t'] = true
}

func _trimCDATA(in []byte, start, end int) (si, ei int) {
	//`#whitespaces#<![CDATA[ data ]]>#whitespaces#`
	si, ei = _trim(in, start, end)
	//search for <![CDATA[
	found := bytes.HasPrefix(in[si:ei], []byte(cdataStart))
	if found {
		si = si + len(cdataStart)
		ei = ei - len(cdataEnd)
		//if si+len(cdataStart) > ei-len(cdataEnd) {}
		return si, ei
		//si, ei = trim(in, si, ei)
	}
	return start, end
}

func _trim(in []byte, start, end int) (int, int) {
	//remove heading whitespaces
	for ; start < end && whitespace[in[start]]; start++ {
	}
	//remove trailing whitespaces
	for ; end > start && whitespace[in[end-1]]; end-- {
	}
	return start, end
}
