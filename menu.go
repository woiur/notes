package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"sort"
)

type Menu struct {
	s string
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) HTML() template.HTML {
	return template.HTML(m.s)
}

func (m *Menu) rebuild() {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString(`<nav class="headernav"><ul class="main-nav"><li>
		<a href="/index">Home</a></li><li><a href="/edit/category/title">+</a></li>`)
	dirs, err := content(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range dirs {
		writeStrings(buf, `<li><a href="">`, v, `</a><ul class="articles-nav">`)
		files, err := content(dataDir + "/" + v)
		if err != nil {
			log.Fatal(err)
		}
		for _, u := range files {
			writeStrings(buf, `<li><a href="/view/`, v, "/", u, `">`, u, `</a></li>`)
		}
		buf.WriteString(`</ul></li>`)
	}
	buf.WriteString(`</ul></nav>`)
	m.s = buf.String()
}

func writeStrings(b *bytes.Buffer, ss ...string) {
	for _, v := range ss {
		b.WriteString(v)
	}
}

func content(dir string) ([]string, error) {
	d, err := os.Open(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer d.Close()
	s, err := d.Readdirnames(-1)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sort.Strings(s)
	return s, nil
}
