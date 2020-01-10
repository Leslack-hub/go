package perslset

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"testing"
)

func TestItemSaver(t *testing.T) {
	profile := map[string]string{
		"Age": "11",
		"name": "lishanchuan",
	}
	id, err := save(profile)
	if err!= nil {
		panic(err)
	}

	client, err := elastic.NewClient(elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	resp, err := client.Get().Index("dating_profile").Type("zhenai").Id(id).Do(context.Background())

	if err != nil {
		panic(err)
	}

	var actual model.Profile
	err = json.Unmarshal(*resp.Source, &actual)
	if err != nil {
		panic(err)
	}

}

