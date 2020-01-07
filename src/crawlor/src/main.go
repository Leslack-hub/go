package src

import (
	"fmt"
	"regexp"
	"./src/fetcher/fetcher"
)

func main() {
	printCityList(all)
}

func printCityList(contents []byte) {
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" [^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(contents, -1)
	var info = make(map[string]string)
	for _, m := range matches {
		info[string(m[2])] = string(m[1])
	}
	fmt.Println("len:", len(matches))
}
