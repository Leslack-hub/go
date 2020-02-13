package main

import (
	"fmt"
	"sort"
	"strings"
)

func groupAnagrams(strs []string) [][]string {
	if len(strs) <= 1 {
		return [][]string{}
	}
	var res [][]string
	typeList := make(map[string][]string)
	for i := 0; i < len(strs); i++ {
		var tempList []string
		for j := 0; j < len(strs[i]); j++ {
			tempList = append(tempList, string(strs[i][j]))
		}
		sort.Strings(tempList)
		key := strings.Join(tempList, "")
		val, ok := typeList[key]
		if !ok {
			typeList[key] = []string{strs[i]}
		} else {
			typeList[key] = append(val, strs[i])
		}
	}
	for _, v := range typeList {
		res = append(res, v)
	}
	return res
}

func judgeType(location map[int][]byte, str []byte) (map[int][]byte, int) {
	var isAppend bool
	var num int
	for k, v := range location {
		numTemp := 0
		for i := range v {
			if inArray(str, v[i]) {
				numTemp++
			}
		}
		if numTemp == len(v) {
			num = k
			isAppend = true
			break
		}
		num++
	}

	if !isAppend {
		num++
		location[num] = str
	}

	return location, num
}

func inArray(array []byte, val byte) bool {
	var inArray bool
	for i := range array {
		if array[i] == val {
			inArray = true
			break
		}
	}
	return inArray
}
func main() {
	fmt.Println(groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"}))
}
