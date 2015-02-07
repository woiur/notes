package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

type Page struct {
	*Menu
	*Article
}

var validator = regexp.MustCompile(`^[\pL0-9]+$`)

var templates = template.Must(template.ParseFiles(
	tmplDir+"edit.html",
	tmplDir+"view.html",
	tmplDir+"index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/index", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//	fmt.Println(menu)
	err := templates.ExecuteTemplate(w, "index.html", menu)
	if err != nil {
		println("aDF")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseParams(ps httprouter.Params) (string, string, error) {
	cat := ps.ByName("category")
	title := ps.ByName("title")
	if !validator.MatchString(cat) || !validator.MatchString(title) {
		return "", "", errors.New("Wrong params")
	}
	return cat, title, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, title, err := parseParams(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	art, err := loadArticle(cat, title)
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	renderTemplate(w, "view", &Page{menu, art})
}

func editHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, title, err := parseParams(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	art, err := loadArticle(cat, title)
	if err != nil {
		art = &Article{Cat: cat, Title: title}
	}
	renderTemplate(w, "edit", &Page{menu, art})
}

func saveHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, title, err := parseParams(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	c := r.FormValue("cat")
	t := r.FormValue("title")
	//Remove old article if the category or the title was changed
	if cat != c || title != t {
		rm(cat, title)
	}
	art := &Article{Cat: c, Title: t, Body: []byte(r.FormValue("body"))}
	err = art.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	menu.rebuild()
	http.Redirect(w, r, "/view/"+art.Cat+"/"+art.Title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cat, title, err := parseParams(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	rm(cat, title)
	menu.rebuild()
	http.Redirect(w, r, "/index", http.StatusFound)
}

func rm(d, f string) {
	err := os.Remove(dataDir + d + "/" + f)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.Remove(dataDir + d) //Remove dir if empty
	if err != nil {
		log.Println(err)
	}
}
