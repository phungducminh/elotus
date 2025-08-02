package code_challenges

func findLength_Original(nums1 []int, nums2 []int) int {
	// 1,2,3,2,1
	// 3,2,1,4,7
	// dp[i][j] is maxmimum length of subarray end of nums1[i] and nums2[j]
	// if nums1[i] == nums2[j] -> dp[i][j] = max(dp[i][j], dp[i-1][j-1] + 1)
	n, m := len(nums1), len(nums2)
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, m+1)
	}

	ans := 0
	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = max(dp[i][j], dp[i-1][j-1] + 1)
				ans = max(ans, dp[i][j])
			} else {
				dp[i][j] = 0
			}
		}
	}
	return ans
}

func findLength(nums1 []int, nums2 []int) int {
	n, m := len(nums1), len(nums2)
   	prv := make([]int, m+1) 
	cur := make([]int, m+1)

	ans := 0
	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			if nums1[i-1] == nums2[j-1] {
				cur[j] = max(cur[j], prv[j-1] + 1)
				ans = max(ans, cur[j])
			} else {
				cur[j] = 0
			}
		}
		cur, prv = prv, cur
	}
	return ans
}
