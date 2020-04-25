package main

import (
	"fmt"
	"testing"
)

func Test_dfs(t *testing.T) {
	type args struct {
		left  int
		right int
		idx   int
		bytes []byte
		res   *[]string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "generateParenthesis",
			args: args{
				left:  3,
				right: 2,
				idx:   2,
				bytes: nil,
				res:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(generateParenthesis(tt.args.left))
		})
	}
}
