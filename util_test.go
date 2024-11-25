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

func Test_trim(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: `empty`, args: ``, want: ``},
		{name: `only_whitespace`, args: `    `, want: ``},
		{name: `heading_whitespace`, args: `  abc`, want: `abc`},
		{name: `trailing_whitespace`, args: `abc  `, want: `abc`},
		{name: `mid_whitespace`, args: `abc  abc`, want: `abc  abc`},
		{name: `all_whitespace`, args: `  abc  abc  `, want: `abc  abc`},
		{name: `newline`, args: "\n\nabc\nabc\n\n", want: "abc\nabc"},
		{name: `tab`, args: "\t\tabc\tabc\t\t", want: "abc\tabc"},
		{name: `cr`, args: "\r\rabc\rabc\r\r", want: "abc\rabc"},
		{name: `&#xA;`, args: "&#xA;&#xA;abc&#xA;abc&#xA;&#xA;", want: "abc&#xA;abc"},
		{name: `&#x9;`, args: "&#x9;&#x9;abc&#x9;abc&#x9;&#x9;", want: "abc&#x9;abc"},
		{name: `&#xD;`, args: "&#xD;&#xD;abc&#xD;abc&#xD;&#xD;", want: "abc&#xD;abc"},
		{name: `mixed`, args: "&#xA; &#x9; \n&#x9;\tabc&#xA; &#x9; \n&#x9;\tabc&#xA; &#x9; \n&#x9;\t", want: "abc&#xA; &#x9; \n&#x9;\tabc"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si, di := _trim([]byte(tt.args), 0, len(tt.args))
			assert.Equal(t, tt.want, string(tt.args[si:di]))
		})
	}
}

func TestSkipToken(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		start    int
		end      int
		front    bool
		expected int
	}{
		{
			name:     "valid_token_at_start_&#xA;",
			data:     []byte("&#xA; rest of the data"),
			start:    0,
			end:      15,
			front:    true,
			expected: 5,
		},
		{
			name:     "valid_token_at_end_&#xD;",
			data:     []byte("some data &#xD;"),
			start:    0,
			end:      15,
			front:    false,
			expected: 5,
		},
		{
			name:     "valid_token_at_start_&#x9;",
			data:     []byte("&#x9; rest of the data"),
			start:    0,
			end:      15,
			front:    true,
			expected: 5,
		},
		{
			name:     "invalid_token",
			data:     []byte("&#xX; rest of the data"),
			start:    0,
			end:      15,
			front:    true,
			expected: 0,
		},
		{
			name:     "not_enough_data_for_start_token",
			data:     []byte("&#x"),
			start:    0,
			end:      3,
			front:    true,
			expected: 0,
		},
		{
			name:     "not_enough_data_for_end_token",
			data:     []byte("&#x"),
			start:    0,
			end:      3,
			front:    false,
			expected: 0,
		},
		{
			name:     "valid_token_in_middle_front=false",
			data:     []byte("data &#x9; more data"),
			start:    0,
			end:      20,
			front:    false,
			expected: 0,
		},
		{
			name:     "valid_token_at_end_unmatched_characters",
			data:     []byte("some data &#xB;"),
			start:    0,
			end:      15,
			front:    false,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipToken(tt.data, tt.start, tt.end, tt.front)
			assert.Equal(t, tt.expected, result)
		})
	}
}
