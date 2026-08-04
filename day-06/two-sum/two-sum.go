package main

import "fmt"

func TwoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func TwoSumAll(nums []int, target int) [][]int {
	result := make([][]int, 0)

	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				result = append(result, []int{i, j})
			}
		}
	}
	return result
}

func main() {
	nums := []int{2, 4, 7, 1, 5}
	target := 9

	result := TwoSum(nums, target)
	fmt.Println(result)

	allResults := TwoSumAll(nums, target)
	fmt.Println(allResults)
}
