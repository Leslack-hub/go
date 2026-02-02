package factory

import "testing"

func TestTextSplitter_split(t *testing.T) {
	form := MainForm{&TextSplitter{}}
	form.Create()
}