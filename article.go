package main

import (
	"html/template"
	"io/ioutil"
	"os"

	"github.com/russross/blackfriday"
)

type Article struct {
	Cat   string
	Title string
	Body  []byte
}

func loadArticle(cat, title string) (*Article, error) {
	body, err := ioutil.ReadFile(dataDir + cat + "/" + title)
	if err != nil {
		return nil, err
	}
	return &Article{Cat: cat, Title: title, Body: body}, nil
}

func (a *Article) save() error {
	err := ioutil.WriteFile(dataDir+a.Cat+"/"+a.Title, a.Body, 0600)
	if err != nil {
		err = os.Mkdir(dataDir+a.Cat, 0700)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(dataDir+a.Cat+"/"+a.Title, a.Body, 0600)
}

func (a *Article) HTML() template.HTML {
	markdown := blackfriday.MarkdownCommon(a.Body)
	return template.HTML(markdown)
}
