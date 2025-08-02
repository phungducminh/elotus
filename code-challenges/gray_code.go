package graycode

func grayCode(n int) []int {
	ans := make([]int, 0)
	set := make(map[int]struct{})
	solve(n, &ans, set, 0, 0)
	return ans
}

func solve(n int, list *[]int, set map[int]struct{}, code int, i int) bool {
	if i == (1 << n) - 1 {
		for bit := 0; bit < n; bit++ {
			if (1 << bit) == code {
				*list = append(*list, code)
				set[code] = struct{}{}
				return true
			}
		}
		return false
	}
	*list = append(*list, code)
	set[code] = struct{}{}
	for bit := 0; bit < n; bit++ {
		// change 1 bit 0->1 or 1->0
		nextCode := code ^ (1 << bit)
		if _, ok := set[nextCode]; !ok {
			if solve(n, list, set, nextCode, i+1) {
				return true
			}
		}
	}

	delete(set, code)
	*list = (*list)[:len(*list)-1]
	return false
}
