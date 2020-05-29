package perslset

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Save: got item"+"#%d: %v", itemCount, item)
			_, err := save(item)
			if err != nil {
				log.Printf("Item Saver err saving Item %v: %v", item, err)
			}

		}
	}()

	return out
}

func save(item interface{}) (id string, err error) {
	// 内网没有办法sniff 需要关闭
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return "", err
	}

	response, err := client.Index().Index("dating_profile").Type("zhenai").BodyJson(item).Do(context.Background())
	if err != nil {
		return "", err
	}

	fmt.Println(response)
	// fmt.Printf("%+v",response) 打印结构体
	return response.Id, nil
}
