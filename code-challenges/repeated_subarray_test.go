package code_challenges

import (
	// "fmt"
	"testing"
)

func TestFindLength01(t *testing.T) {
	testFindLength(t, findLength)	
}

func TestFindLength02(t *testing.T) {
	testFindLength(t, findLength_Original)	
}

func testFindLength(t *testing.T, testFunc func([]int, []int)int) {
	tests := []struct {
		nums1  []int
		nums2  []int
		expect int
	}{
		{
			nums1:  []int{1, 2, 3, 2, 1},
			nums2:  []int{3, 2, 1, 4, 7},
			expect: 3,
		},
		{
			nums1:  []int{0, 0, 0, 0, 0},
			nums2:  []int{0, 0, 0, 0, 0},
			expect: 5,
		},
	}
	for _, tt := range tests {
		actual := testFunc(tt.nums1, tt.nums2)
		// fmt.Printf("findLength %d, %d\n", actual, tt.expect)
		if actual != tt.expect {
			t.Errorf("findLength(%v, %v), expect=%d, actual=%d", tt.nums1, tt.nums2, tt.expect, actual)
		}
	}
}
