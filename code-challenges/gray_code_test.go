package graycode

import (
	"fmt"
	"testing"
)

func TestGrayCode(t *testing.T) {
	tests := []struct {
		n int
	}{
		{
			n: 2,
		},
		{
			n: 1,
		},
		{
			n: 3,
		},
	}

	for _, tt := range tests {
		ans := grayCode(tt.n)
		fmt.Printf("%d\n", len(ans))
		for i := 0; i < len(ans); i++ {
			fmt.Printf("grayCode(%d): %d\n", tt.n, ans[i])
		}
	}
}
