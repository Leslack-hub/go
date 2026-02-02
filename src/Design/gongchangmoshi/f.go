package factory

import "fmt"

type SplitterFactory interface {
	CreateSplitter() ISplitter
}

type ISplitter interface {
	split()
}

type BinarySplitter struct {
}

func (b *BinarySplitter) split() {
	fmt.Println("binary 文件分解")
}
func (this *BinarySplitter) CreateSplitter() ISplitter {
	return &BinarySplitter{}
}

type TextSplitter struct {
}

func (this *TextSplitter) split() {
	fmt.Println("text 文件分解")
}
func (this *TextSplitter) CreateSplitter() ISplitter {
	return &TextSplitter{}
}

type MainForm struct {
	factory SplitterFactory
}

func (this MainForm) Create() {
	splitter := this.factory.CreateSplitter()
	splitter.split()
}
