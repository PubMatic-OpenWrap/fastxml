package fastxml

import (
	"testing"
)

func TestXPath_AddAndGet(t *testing.T) {
	root := &xpath{data: "/", childs: make(map[string]*xpath)}

	paths := [][]string{
		{"a", "b", "c"},
		{"a", "b", "d"},
		{"a", "e"},
	}

	for _, p := range paths {
		root.add(p)
	}

	tests := []struct {
		path       []string
		expected   string
		shouldFind bool
	}{
		{[]string{"a", "b", "c"}, "c", true},
		{[]string{"a", "b", "d"}, "d", true},
		{[]string{"a", "e"}, "e", true},
		{[]string{"a", "x"}, "", false},
	}

	for _, test := range tests {
		node := root
		for _, key := range test.path {
			node = node.get(key)
			if node == nil {
				break
			}
		}
		if test.shouldFind && node == nil {
			t.Errorf("Expected to find path %v, but did not", test.path)
		} else if !test.shouldFind && node != nil {
			t.Errorf("Expected not to find path %v, but found %v", test.path, node.data)
		} else if test.shouldFind && node.data != test.expected {
			t.Errorf("Expected to find path %v with data %v, but found %v", test.path, test.expected, node.data)
		}
	}
}

func TestGetXPath(t *testing.T) {
	paths := [][]string{
		{"a", "b", "c"},
		{"a", "b", "d"},
		{"a", "e"},
	}

	xpath := GetXPath(paths)

	tests := []struct {
		path       []string
		expected   string
		shouldFind bool
	}{
		{[]string{"a", "b", "c"}, "c", true},
		{[]string{"a", "b", "d"}, "d", true},
		{[]string{"a", "e"}, "e", true},
		{[]string{"a", "x"}, "", false},
	}

	for _, test := range tests {
		node := xpath
		for _, key := range test.path {
			node = node.get(key)
			if node == nil {
				break
			}
		}
		if test.shouldFind && node == nil {
			t.Errorf("Expected to find path %v, but did not", test.path)
		} else if !test.shouldFind && node != nil {
			t.Errorf("Expected not to find path %v, but found %v", test.path, node.data)
		} else if test.shouldFind && node.data != test.expected {
			t.Errorf("Expected to find path %v with data %v, but found %v", test.path, test.expected, node.data)
		}
	}
}
