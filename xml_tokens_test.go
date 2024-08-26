package fastxml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTokenType(t *testing.T) {
	type args struct {
		in    string
		index int
	}
	tests := []struct {
		name string
		args []args
		want xmlTokenType
	}{
		{
			name: `start_tag`,
			args: []args{
				{in: `<test/>`, index: 1},
			},
			want: startXMLToken,
		},
		{
			name: `end_tag`,
			args: []args{
				{in: `</test>`, index: 1},
			},
			want: endXMLToken,
		},
		{
			name: `processing_tag`,
			args: []args{
				{in: `<?xml version = "1.0" encoding = "UTF-8" standalone = "no" ?>`, index: 1},
			},
			want: processingXMLToken,
		},
		{
			name: `comments_tag`,
			args: []args{
				{in: `<!-- commented code -->`, index: 1},
			},
			want: commentsXMLToken,
		},
		{
			name: `cdata_tag`,
			args: []args{
				{in: `<![CDATA[test]]>`, index: 1},
			},
			want: cdataXMLToken,
		},
		{
			name: `doctype_tag`,
			args: []args{
				{in: `<!DOCTYPE list SYSTEM "example.dtd">`, index: 1},
			},
			want: doctypeXMLToken,
		},
		{
			name: `invalid_tag`,
			args: []args{
				{in: ``, index: 0},
				{in: `<test/>`, index: 10},
				{in: `<123test/>`, index: 1},
				{in: `< test/>`, index: 1},
				{in: `<! -- test -->`, index: 1},
				{in: `<! [CDATA[test]]>`, index: 1},
				{in: `<![ CDATA[test]]>`, index: 1},
				{in: `<![CDATAtest]]>`, index: 1},
				{in: `<! DOCTYPE list SYSTEM "example.dtd">`, index: 1},
			},
			want: unknownXMLToken,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.args {
				got := getTokenType([]byte(arg.in), arg.index)
				assert.Equal(t, tt.want, got, fmt.Sprintf("[input] args:%s, index:%v", arg.in, arg.index))
			}
		})
	}
}

func Test_getTokenEndIndex(t *testing.T) {
	type args struct {
		in         string
		startIndex int
		ttype      xmlTokenType
	}
	type want struct {
		index  int
		inline bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: `start_token_empty_string`,
			args: args{in: ``, ttype: startXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `start_token_first_char`,
			args: args{in: `>`, ttype: startXMLToken},
			want: want{index: 1, inline: false},
		},
		{
			name: `start_token_valid`,
			args: args{in: `<test>dummy_text`, ttype: startXMLToken},
			want: want{index: 6, inline: false},
		},
		{
			name: `start_token_inline`,
			args: args{in: `<test/>dummy_text`, ttype: startXMLToken},
			want: want{index: 7, inline: true},
		},
		{
			name: `start_token_not_found`,
			args: args{in: `<test`, ttype: startXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `start_token_empty_inline`,
			args: args{in: `/>`, ttype: startXMLToken},
			want: want{index: 2, inline: true},
		},
		{
			name: `end_token_empty`,
			args: args{in: ``, ttype: endXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `end_token_first_char`,
			args: args{in: `>`, ttype: endXMLToken},
			want: want{index: 1, inline: false},
		},
		{
			name: `end_token_valid`,
			args: args{in: `</test>dummy_text`, ttype: endXMLToken},
			want: want{index: 7, inline: false},
		},
		{
			name: `end_token_not_found`,
			args: args{in: `</test`, ttype: endXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `preprocessing_token_empty`,
			args: args{in: ``, ttype: processingXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `preprocessing_token_first_char`,
			args: args{in: `?>`, ttype: processingXMLToken},
			want: want{index: 2, inline: false},
		},
		{
			name: `preprocessing_token_valid`,
			args: args{in: `<? test ?>dummy_test`, ttype: processingXMLToken},
			want: want{index: 10, inline: false},
		},
		{
			name: `preprocessing_token_not_found`,
			args: args{in: `<? test `, ttype: processingXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `preprocessing_token_missing_questionmark`,
			args: args{in: `<? test >`, ttype: processingXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_empty`,
			args: args{in: ``, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_first_char`,
			args: args{in: `-->`, ttype: commentsXMLToken},
			want: want{index: 3, inline: false},
		},
		{
			name: `comments_token_first_char_missing_dashes`,
			args: args{in: `>`, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_first_char_missing_one_dash`,
			args: args{in: `->`, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_valid`,
			args: args{in: `<!-- test -->dummy_text`, ttype: commentsXMLToken},
			want: want{index: 13, inline: false},
		},
		{
			name: `comments_token_not_found`,
			args: args{in: `<!-- test`, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_contains_gt`,
			args: args{in: `<!-- test > -->dummy_text`, ttype: commentsXMLToken},
			want: want{index: 15, inline: false},
		},
		{
			name: `comments_token_missing_both_dashes`,
			args: args{in: `<!-- test >`, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `comments_token_missing_first_dash`,
			args: args{in: `<!-- test ->`, ttype: commentsXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `cdata_token_empty`,
			args: args{in: ``, ttype: cdataXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `cdata_token_first_char`,
			args: args{in: `]]>`, ttype: cdataXMLToken},
			want: want{index: 3, inline: false},
		},
		{
			name: `cdata_token_first_char_missing_one_bracket`,
			args: args{in: `]>`, ttype: cdataXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `cdata_token_first_missing_all_brackets`,
			args: args{in: `>`, ttype: cdataXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `cdata_token_valid`,
			args: args{in: `<![CDATA[ test ]]>`, ttype: cdataXMLToken},
			want: want{index: 18, inline: false},
		},
		{
			name: `cdata_token_not_found`,
			args: args{in: `<![CDATA[ test ]]`, ttype: cdataXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `doctype_token_empty`,
			args: args{in: ``, ttype: doctypeXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `doctype_token_valid`,
			args: args{in: `<!DOCTYPE test>`, ttype: doctypeXMLToken},
			want: want{index: 15, inline: false},
		},
		{
			name: `doctype_token_with_bracket`,
			args: args{in: `<!DOCTYPE test[<!ELEMENT note>]>`, ttype: doctypeXMLToken},
			want: want{index: 32, inline: false},
		},
		{
			name: `doctype_token_with_multiple_bracket`,
			args: args{in: `<!DOCTYPE test[<!ELEMENT note [<!ELEMENT note1>]>]>`, ttype: doctypeXMLToken},
			want: want{index: 49, inline: false},
		},
		{
			name: `unknown_token_empty`,
			args: args{in: ``, ttype: unknownXMLToken},
			want: want{index: -1, inline: false},
		},
		{
			name: `unknown_token_valid`,
			args: args{in: `test>`, ttype: unknownXMLToken},
			want: want{index: 5, inline: false},
		},
		{
			name: `unknown_token_not_found`,
			args: args{in: `test`, ttype: unknownXMLToken},
			want: want{index: -1, inline: false},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, inline := getTokenEndIndex([]byte(tt.args.in), tt.args.startIndex, tt.args.ttype)
			assert.Equal(t, tt.want.index, index)
			assert.Equal(t, tt.want.inline, inline)
		})
	}
}

func Test_getTokenNameIndex(t *testing.T) {
	type args struct {
		in         string
		startIndex int
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: `empty_string`,
			args: args{in: ``, startIndex: 0},
			want: ``,
		},
		{
			name: `end_tag_with_gt`,
			args: args{in: `<test>`, startIndex: 1},
			want: `test`,
		},
		{
			name: `end_tag_with_space`,
			args: args{in: `<test k="v">`, startIndex: 1},
			want: `test`,
		},
		{
			name: `end_tag_inline`,
			args: args{in: `<test/>`, startIndex: 1},
			want: `test`,
		},
		{
			name: `namespace`,
			args: args{in: `<xn:test/>`, startIndex: 1},
			want: `test`,
		},
		{
			name: `invalid_multiple_namespace_separator`,
			args: args{in: `<xn1:xn2:test/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `invalid_special_character_start`,
			args: args{in: `<#test/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `invalid_special_character_mid`,
			args: args{in: `<te#st/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `invalid_special_character_last`,
			args: args{in: `<test#/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `invalid_start_with_number`,
			args: args{in: `<1test/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `invalid_start_with_dash`,
			args: args{in: `<-test/>`, startIndex: 1},
			want: ``,
		},
		{
			name: `contains_dash`,
			args: args{in: `<test-1/>`, startIndex: 1},
			want: `test-1`,
		},
		{
			name: `contains_dash`,
			args: args{in: `<test-1/>`, startIndex: 1},
			want: `test-1`,
		},
		{
			name: `start_with_underscore`,
			args: args{in: `<_test/>`, startIndex: 1},
			want: `_test`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si, ei := getTokenNameIndex([]byte(tt.args.in), tt.args.startIndex)
			assert.Equal(t, tt.want, string(tt.args.in[si:ei]))
		})
	}
}
