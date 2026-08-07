package main

import (
	"errors"
	"fmt"
)

func FindSecondLargest(nums []int) (int, error) {
	if len(nums) < 2 {
		return 0, errors.New("not enough numbers")
	}
	largest := nums[0]
	secondLargest := 0
	hasSecond := false

	for _, current := range nums {
		if current > largest {
			secondLargest = largest
			largest = current
			hasSecond = true
		} else if current < largest && (!hasSecond || current > secondLargest) {
			secondLargest = current
			hasSecond = true
		}
	}

	if !hasSecond {
		return 0, errors.New("second largest not found")
	}
	return secondLargest, nil

}

func main() {
	nums := []int{1, 2, 5, -1, 20, 99, 48}

	result, err := FindSecondLargest(nums)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("second largest number:", result)
}
