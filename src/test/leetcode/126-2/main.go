package main

import "fmt"

func findLadders(beginWord string, endWord string, wordList []string) [][]string {
	var exists bool
	for _, word := range wordList {
		if word == endWord {
			exists = true
			break
		}
	}
	if !exists {
		return [][]string{}
	}
	graph, distance := bfs(endWord, append(wordList, beginWord))

	var ret [][]string
	var path []string
	dfs(beginWord, endWord, path, &ret, distance, graph)
	return ret
}

func dfs(begin, end string, path []string, ret *[][]string, distance map[string]int, graph map[string][]string) {
	if begin == end {
		*ret = append(*ret, append(append([]string{}, path...), end))
		return
	}
	for _, word := range graph[begin] {
		if distance[begin] == distance[word]+1 {
			path = append(path, begin)
			dfs(word, end, path, ret, distance, graph)
			path = path[:len(path)-1]
		}
	}
}

func bfs(endWord string, wordList []string) (map[string][]string, map[string]int) {
	step1List := make(map[string][]string)
	distance := map[string]int{endWord: 0}

	queue := []string{endWord}
	for len(queue) > 0 {
		var nextQueue []string
		for _, word := range queue {
			for _, step1Word := range step1(word, wordList) {
				step1List[word] = append(step1List[word], step1Word)
				if _, ok := distance[step1Word]; !ok {
					distance[step1Word] = distance[word] + 1
					nextQueue = append(nextQueue, step1Word)
				}
			}
		}
		queue = nextQueue
	}
	return step1List, distance
}

func step1(word string, wordList []string) []string {
	var ret []string
	for _, v := range wordList {
		var count int
		for i := 0; i < len(word); i++ {
			if word[i] == v[i] {
				count++
			}
		}
		if count == len(word)-1 {
			ret = append(ret, v)
		}
	}
	return ret
}

func main() {
	fmt.Println(findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}))
}
