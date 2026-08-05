package main

import "fmt"

// nums[2, 7, 11] target= 9

func TwoSumMap(nums []int, target int) []int {

	repo := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		need := target - nums[i]

		index, ok := repo[need]
		if ok {
			return []int{index, i}
		}
		repo[nums[i]] = i
	}
	return nil
}

func main() {
	nums := []int{2, 5, 6, 7, 11}
	target := 9

	fmt.Println(TwoSumMap(nums, target))
}
