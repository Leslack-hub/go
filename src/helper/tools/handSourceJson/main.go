package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"gopkg.in/yaml.v2"
)
var res map[string]interface{}

func main() {
	res = make(map[string]interface{})
	//操作系统指定的路径分隔符
	readdir("/Users/leslack/lsc/test_dir")
	fmt.Println(res)
}

func readdir(p string) {
	//以只读的方式打开目录
	f, err := os.OpenFile(p, os.O_RDONLY, os.ModeDir)
	if err != nil {
		fmt.Println(err.Error())
	}
	//延迟关闭目录
	defer f.Close()
	fileInfo, _ := f.Readdir(-1)
	separator := string(os.PathSeparator)
	_, base := path.Split(p)
	var detail []string
	infos := make(map[string]interface{})
	for _, info := range fileInfo {
		//判断是否是目录
		if info.IsDir() {
			readdir(p+separator+info.Name())
		} else {
			sourceStr := info.Name()
			matched, _ := regexp.MatchString(`.png`, sourceStr)
			if !matched {
				matched, _ = regexp.MatchString(`.yaml`, sourceStr)
				if matched {
					yamlFile := readYaml(p + separator + info.Name())
					for i, k := range yamlFile {
						infos[i] = k
					}
				}
				continue
			}
			if info.Name() == "cover.png" {
				infos["cover"] = info.Name()
			} else if info.Name() == "thumbnail.png" {
				infos["thumbnail"] = info.Name()
			} else {
				detail = append(detail,info.Name())
			}
		}
		infos["detail"] = detail
		res[base] = infos
	}
}

func readYaml(file string) map[string]interface{}{
	var b map[string]interface{}

	yfile, _ := os.Open(file) //test.yaml由下一个例子生成
	defer yfile.Close()

	ydecode := yaml.NewDecoder(yfile)
	ydecode.Decode(&b) //注意这里为指针
	return b
}
