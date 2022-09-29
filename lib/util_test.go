package lib

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	condition := func(n int) bool {
		return n%2 == 0
	}

	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"empty", args{[]int{}}, []int{}},
		{"filled", args{[]int{1, 2, 3, 4, 5, 6, 7}}, []int{2, 4, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.input, condition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
