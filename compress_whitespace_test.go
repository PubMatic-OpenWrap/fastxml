package fastxml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndentationRemoverWrite(t *testing.T) {
	type args struct {
		input    string
		lastChar byte
	}
	type want struct {
		want         string
		wantError    bool
		wantLastChar byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no_whitespace",
			args: args{
				input:    "<node>text</node>",
				lastChar: 0,
			},
			want: want{
				want:         "<node>text</node>",
				wantError:    false,
				wantLastChar: '>',
			},
		},
		{
			name: "leading_whitespace",
			args: args{
				input:    "   <node>text</node>",
				lastChar: '>',
			},
			want: want{
				want:         "<node>text</node>",
				wantError:    false,
				wantLastChar: '>',
			},
		},
		{
			name: "trailing_whitespace",
			args: args{
				input:    "<node>text</node>   ",
				lastChar: '>',
			},
			want: want{
				want:         "<node>text</node>",
				wantError:    false,
				wantLastChar: '>',
			},
		},
		{
			name: "mixed_whitespace",
			args: args{
				input:    "   <node>   text   </node>   ",
				lastChar: '>',
			},
			want: want{
				want:         "<node>text   </node>",
				wantError:    false,
				wantLastChar: '>',
			},
		},
		{
			name: "nested_elements_with_whitespace",
			args: args{
				input:    "<a>   <b>   <c>text</c>   </b>   </a>",
				lastChar: '>',
			},
			want: want{
				want:         "<a><b><c>text</c></b></a>",
				wantError:    false,
				wantLastChar: '>',
			},
		},
		{
			name: "empty_input",
			args: args{
				input:    "",
				lastChar: 0,
			},
			want: want{
				want:         "",
				wantError:    false,
				wantLastChar: 0,
			},
		},
		{
			name: "only_whitespace",
			args: args{
				input:    "   ",
				lastChar: '>',
			},
			want: want{
				want:         "",
				wantError:    false,
				wantLastChar: '>',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			ir := newCompressWhitespace(buf)
			ir.lastChar = tt.args.lastChar

			_, err := ir.write([]byte(tt.args.input))
			assert.NoError(t, err)
			assert.Equal(t, tt.want.wantLastChar, ir.lastChar)
			got := buf.String()
			assert.Equal(t, tt.want.want, got)
		})
	}
}
