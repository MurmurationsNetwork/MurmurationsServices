package service

import (
	"reflect"
	"testing"
)

func TestSplitEscapedComma(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "No comma",
			input:  "apple",
			expect: []string{"apple"},
		},
		{
			name:   "Basic Split",
			input:  "apple,banana",
			expect: []string{"apple", "banana"},
		},
		{
			name:   "First character is comma",
			input:  ",apple,banana",
			expect: []string{"apple", "banana"},
		},
		{
			name:   "Escaped comma",
			input:  "Banana\\, Inc.",
			expect: []string{"Banana, Inc."},
		},
		{
			name:   "Complicated Conditions",
			input:  "\\,apple\\,banana,cherry",
			expect: []string{",apple,banana", "cherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitEscapedComma(tt.input); !reflect.DeepEqual(
				got,
				tt.expect,
			) {
				t.Errorf(
					"splitEscapedComma(%v) = %v, want %v",
					tt.input,
					got,
					tt.expect,
				)
			}
		})
	}
}
