package fastxml

import (
	"bytes"
)

var (
	whitespace    [256]bool //<space>, \r, \n, \t
	alnum         [256]bool //a-z, A-Z, 0-9
	alpha         [256]bool //a-z, A-Z
	num           [256]bool //0-9
	name          [256]bool //a-z, A-Z, 0-9, _, -
	encodingChars = map[string]byte{
		"amp;":  '&',
		"apos;": '\'',
		"lt;":   '<',
		"gt;":   '>',
		"quot;": '"',
		"#39;":  '\'',
		"#34;":  '"',
		"#xA;":  '\n',
		"#x9;":  '\t',
		"#xD;":  '\r',
	}
)

/*
var xmlReplacerNormal = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"'", "&apos;",
	`"`, "&quot;",
)

var xmlReplacerCanonicalText = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\r", "&#xD;",
)

var xmlReplacerCanonicalAttrVal = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	`"`, "&quot;",
	"\t", "&#x9;",
	"\n", "&#xA;",
	"\r", "&#xD;",
)
*/

func init() {
	whitespace[' '] = true
	whitespace['\r'] = true
	whitespace['\n'] = true
	whitespace['\t'] = true

	//name
	name['_'] = true
	name['-'] = true

	//alnum
	for ch := 'a'; ch <= 'z'; ch++ {
		alnum[ch] = true
		alpha[ch] = true
		name[ch] = true
	}
	for ch := 'A'; ch <= 'Z'; ch++ {
		alnum[ch] = true
		alpha[ch] = true
		name[ch] = true
	}
	for ch := '0'; ch <= '9'; ch++ {
		alnum[ch] = true
		num[ch] = true
		name[ch] = true
	}
}

func _trimCDATA(in []byte, start, end int) (int, int, bool) {
	//`#whitespaces#<![CDATA[ data ]]>#whitespaces#`
	si, ei := _trim(in, start, end)
	//search for <![CDATA[
	found := bytes.HasPrefix(in[si:ei], []byte(cdataStart))
	if found {
		si = si + len(cdataStart)
		ei = ei - len(cdataEnd)
		//if si+len(cdataStart) > ei-len(cdataEnd) {}
		//si, ei = _trim(in, si, ei)
		return si, ei, found
	}
	return start, end, found
}

func skipToken(data []byte, start, end int, front bool) (skipLen int) {
	if end-start < 5 { // Minimum length check for &#xN;
		return 0
	}

	si, ei := start, start+5
	if !front {
		si, ei = end-5, end
	}

	tmp := data[si:ei]
	// Check if the pattern matches &#xN; => &#x[9AD];
	if bytes.HasPrefix(tmp, []byte("&#x")) &&
		tmp[4] == ';' &&
		(tmp[3] == 'A' || tmp[3] == '9' || tmp[3] == 'D') {
		return 5 // Skip the token length of 5
	}

	return 0
}

func _trim(in []byte, start, end int) (int, int) {
	// Remove leading whitespaces and special tokens
	for start < end {
		if whitespace[in[start]] {
			start++
			continue
		}
		if skipLen := skipToken(in, start, end, true); skipLen > 0 {
			start += skipLen
			continue
		}
		break
	}

	// Remove trailing whitespaces and special tokens
	for end > start {
		if whitespace[in[end-1]] {
			end--
			continue
		}
		if skipLen := skipToken(in, start, end, false); skipLen > 0 {
			end -= skipLen
			continue
		}
		break
	}

	return start, end
}

func trimSpace(in string) string {
	si, ei := _trim([]byte(in), 0, len(in))
	return in[si:ei]
}

// escape writes an escaped version of a string to the writer.
func escape[T []byte | string](w Writer, s T) {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '&':
			w.WriteString("&amp;")
		case '<':
			w.WriteString("&lt;")
		case '>':
			w.WriteString("&gt;")
		case '\'':
			w.WriteString("&apos;")
		case '"':
			w.WriteString("&quot;")
		default:
			w.WriteByte(ch)
		}
	}
}

func unescape(w Writer, s []byte) {
	//TODO: use prefix tree for below if these functionality extends

	for i := 0; i < len(s); i++ {
		if s[i] != '&' {
			w.WriteByte(s[i])
			continue
		}

		// Check if the & is followed by a known entity
		found := false
		for key, val := range encodingChars {
			if i+len(key) < len(s) && bytes.HasPrefix(s[i+1:], []byte(key)) {
				w.WriteByte(val)
				i += len(key)
				found = true
				break
			}
		}

		if !found {
			w.WriteByte(s[i])
		}
	}
}

func unescapeBytes(s []byte) []byte {
	b := bytes.Buffer{}
	unescape(&b, s)
	return b.Bytes()
}

func quoteEscape[T []byte | string](w Writer, s T) {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' || ch == '\\' {
			w.WriteByte('\\')
		}
		w.WriteByte(ch)
	}
}

func quoteUnescape[T []byte | string](w Writer, s T) {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '\\' {
			if i+1 < len(s) {
				nextCh := s[i+1]
				if nextCh == '\\' || nextCh == '"' || nextCh == '\'' {
					i++
					ch = nextCh
				}
			}
		}
		w.WriteByte(ch)
	}
}
