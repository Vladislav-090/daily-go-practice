package main

import "fmt"

func ContainsDuplicate(nums []int) bool {
	repo := make(map[int]bool)

	for _, num := range nums {
		_, ok := repo[num]

		if ok {
			return true
		}
		repo[num] = true
	}
	return false
}

func main() {
	nums := []int{1, 2, 4, 5, 6, 1}

	result := ContainsDuplicate(nums)

	fmt.Println("contains duplicate:", result)
}
