package strutils

import (
	"reflect"
	"testing"
)

// nolint
func TestStringSplit(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantRets []string
	}{
		{"", args{s: ""}, []string{""}},
		{"", args{s: "    "}, []string{""}},
		{"", args{s: "ab cd"}, []string{"ab", "cd"}},
		{"", args{s: "ab  cd"}, []string{"ab", "cd"}},
		{"", args{s: " ab  cd "}, []string{"ab", "cd"}},
		{"", args{s: "a 'b c' d"}, []string{"a", "b c", "d"}},
		{"", args{s: "a \"b c\" d"}, []string{"a", "b c", "d"}},
		{"", args{s: `a "b c" d`}, []string{"a", "b c", "d"}},
		{"", args{s: `a "b 'c' c" 'a "x x" b' d`}, []string{"a", "b 'c' c", "a \"x x\" b", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRets := StringSplit(tt.args.s); !reflect.DeepEqual(gotRets, tt.wantRets) {
				t.Errorf("StringSplit() = %v, want %v", gotRets, tt.wantRets)
			}
		})
	}
}
