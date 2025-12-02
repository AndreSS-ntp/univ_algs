package ioutil

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadIntsFromFile(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var nums []int

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}

			nums = append(nums, n)
		}
	}

	return nums, nil
}
