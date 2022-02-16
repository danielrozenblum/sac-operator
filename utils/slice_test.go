package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtruct(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				a: []string{"1", "2", "3"},
				b: []string{"2"},
			},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Subtruct(tt.args.a, tt.args.b), "Subtruct(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}
