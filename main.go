package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"	
)

const (
	tmplDir = "./tmpl/"
	dataDir = "./data/"
)

var menu *Menu

func main() {
	menu = NewMenu()
	menu.rebuild()

	rtr := httprouter.New()
	rtr.GET("/", rootHandler)
	rtr.GET("/index/", indexHandler)
	rtr.GET("/view/:category/:title", viewHandler)
	rtr.GET("/edit/:category/:title", editHandler)
	rtr.POST("/save/:category/:title", saveHandler)
	rtr.GET("/delete/:category/:title", deleteHandler)
	rtr.ServeFiles("/static/*filepath", http.Dir("./static/"))

	http.ListenAndServe(":8090", rtr)
}
