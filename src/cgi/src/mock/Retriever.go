package mock

type Retriever struct {
	Content string
}

func (r Retriever) Get(Url string) string {
	return Url
}
