package xmlparser

import (
	"fmt"
)

type xmlAttribute struct {
	key, value xmlTagIndex
}

func (a xmlAttribute) Key(in []byte) []byte {
	return in[a.key.si:a.key.ei]
}

func (a xmlAttribute) Value(in []byte) []byte {
	return in[a.value.si:a.value.ei]
}

func (a xmlAttribute) String(in []byte) string {
	return fmt.Sprintf("%s:%s", in[a.key.si:a.key.ei], in[a.value.si:a.value.ei])
}

func parseAttributes(in []byte, si, ei int) (attributes []xmlAttribute) {
	found := true
	for found {
		var attr xmlAttribute

		//parsing key
		attr.key.si, attr.key.ei, found = _parseKey(in, si, ei)
		if found {
			//parsing = separator
			i := attr.key.ei
			for ; i < ei && in[i] != '='; i = i + 1 {
			}
			//parsing value
			attr.value.si, attr.value.ei, found = _parseValue(in, i+1, ei)
		}
		if found {
			attributes = append(attributes, attr)
			si = attr.value.ei + 1
		}
	}
	return
}

func _parseKey(in []byte, si, ei int) (int, int, bool) {
	len := ei
	for ; si < len && whitespace[in[si]]; si = si + 1 {
	}
	for ei = si; ei < len && in[ei] != '=' && !whitespace[in[ei]]; ei = ei + 1 {
	}
	return si, ei, (ei < len && si != ei)
}

func _parseValue(in []byte, si, ei int) (int, int, bool) {
	len := ei
	for ; si < len && whitespace[in[si]]; si = si + 1 {
	}
	if si < len {
		if !(in[si] == '\'' || in[si] == '"') {
			return 0, 0, false
		}

		quote := in[si]
		for ei = si + 1; ei < len && in[ei] != quote; ei = ei + 1 {
		}
	}
	return si + 1, ei, (ei < len)
}
