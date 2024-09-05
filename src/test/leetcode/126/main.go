package main

import "fmt"

// bfs+dfs(如果是双向bfs，效果会更好)
func findLadders(beginWord string, endWord string, wordList []string) [][]string {
	dict := make(map[string]bool, 0)
	for _, v := range wordList {
		dict[v] = true
	}
	if !dict[endWord] {
		return [][]string{}
	}
	dict[beginWord] = true
	graph := make(map[string][]string)
	distance := bfs(endWord, dict, graph)
	res := make([][]string, 0)
	dfs(beginWord, endWord, &res, []string{}, distance, graph)
	return res
}

// 回溯实现方式一：（个人偏好这个，更符合模板）
func dfs(beginWord string, endWord string, res *[][]string, path []string, distance map[string]int, graph map[string][]string) {
	if beginWord == endWord {
		path = append(path, beginWord)
		tmp := make([]string, len(path))
		copy(tmp, path)
		*res = append(*res, tmp)
		path = path[:len(path)-1]
		return
	}
	for _, v := range graph[beginWord] {
		if distance[beginWord] == distance[v]+1 {
			path = append(path, beginWord)
			dfs(v, endWord, res, path, distance, graph)
			path = path[:len(path)-1]
		}
	}
}

// 从终点出发，进行bfs，计算每一个点到达终点的距离
func bfs(endWord string, dict map[string]bool, graph map[string][]string) map[string]int {
	var distance = make(map[string]int, 0)
	var queue = make([]string, 0)
	queue = append(queue, endWord)
	distance[endWord] = 0 //初始值
	for len(queue) != 0 {
		cursize := len(queue)
		for i := 0; i < cursize; i++ {
			word := queue[0]
			queue = queue[1:]
			for _, v := range expand1(word, dict) {
				graph[v] = append(graph[v], word)
				if _, ok := distance[v]; !ok {
					distance[v] = distance[word] + 1
					queue = append(queue, v)
				}
			}
		}
	}
	return distance
}

// 获得邻接点
func expand(word string, dict map[string]bool) []string {
	expansion := make([]string, 0) //保存word的邻接点
	//从word的每一位开始
	chs := []rune(word)
	for i := 0; i < len(word); i++ {
		tmp := chs[i] //保存当前位，方便后序进行复位
		for c := 'a'; c <= 'z'; c++ {
			//如果一样则直接跳过，之所以用tmp，是因为chs[i]在变
			if tmp == c {
				continue
			}
			chs[i] = c
			newstr := string(chs)
			//新单词在dict中不存在，则跳过
			if dict[newstr] {
				expansion = append(expansion, newstr)
			}
		}
		chs[i] = tmp //单词位复位
	}
	return expansion
}

//func main() {
//	strings := []string{"dot", "dog", "lot", "log"}
//	fmt.Println(expand1("hot", strings))
//	fmt.Println(expand("hot", map[string]bool{
//		"dot": true,
//		"dog": true,
//		"lot": true,
//		"log": true,
//	}))
//}

func expand1(word string, dict map[string]bool) []string {
	var ret []string
	var arr []string
	for k := range dict {
		arr = append(arr, k)
	}
	for _, v := range arr {
		var i int
		var count int
		for i < len(word) {
			if word[i] == v[i] {
				count++
			}
			i++
		}
		if count == len(word)-1 {
			ret = append(ret, v)
		}
	}
	return ret
}

func differByOneChar(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	diffCount := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			diffCount++
			if diffCount > 1 {
				return false
			}
		}
	}

	return diffCount == 1
}

func main() {
	fmt.Println(findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}))
}
