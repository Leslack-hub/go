package main

import (
	"fmt"
	"math"
	"slices"
)

func ladderLength(beginWord string, endWord string, wordList []string) int {
	if !slices.Contains(wordList, endWord) {
		return 0
	}
	step1List, distance := bfs(endWord, append(wordList, beginWord))

	return dfs(beginWord, endWord, 0, step1List, distance)
}

func dfs(begin, end string, level int, step1List map[string][]string, distance map[string]int) int {
	if begin == end {
		return level + 1
	}
	minLevel := math.MaxInt32
	for _, w := range step1List[begin] {
		if distance[begin] == distance[w]+1 {
			minLevel = min(minLevel, dfs(w, end, level+1, step1List, distance))
		}
	}
	if minLevel == math.MaxInt32 {
		return 0
	}
	return minLevel
}

func bfs(endWord string, wordList []string) (map[string][]string, map[string]int) {
	step1List := make(map[string][]string)
	distance := map[string]int{endWord: 0}

	queue := []string{endWord}
	for len(queue) > 0 {
		var nextQueue []string
		for _, word := range queue {
			for _, v := range step1(word, wordList) {
				step1List[word] = append(step1List[word], v)
				if _, ok := distance[v]; !ok {
					distance[v] = distance[word] + 1
					nextQueue = append(nextQueue, v)
				}
			}
		}

		queue = nextQueue
	}
	return step1List, distance
}

func step1(word string, wordList []string) []string {
	var ret []string
	for _, w := range wordList {
		var count int
		for i := 0; i < len(w); i++ {
			if w[i] == word[i] {
				count++
			}
		}
		if count == len(word)-1 {
			ret = append(ret, w)
		}
	}
	return ret
}

func main() {
	fmt.Println(ladderLength("hot", "dog", []string{"hot", "dog"}))
}
