package fastxml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unescape(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: `empty`, args: ``, want: ``},
		{name: `no_escape`, args: `abcdefg 01234 AABCDEF`, want: `abcdefg 01234 AABCDEF`},
		{name: `all_escape`, args: `&lt;&quot;&apos;&amp;&apos;&quot;&gt;`, want: `<"'&'">`},
		{name: `end_with_&`, args: `test&`, want: `test&`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			unescape(&buf, []byte(tt.args))
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
