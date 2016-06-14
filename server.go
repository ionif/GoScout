package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func startHandler(w http.ResponseWriter, r *http.Request, title string) {
//	p, err := loadPage(title)
//	if err != nil {
//		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
//		return
//	}
	renderTemplate(w, "start", nil)
}

func enterHandler(w http.ResponseWriter, r *http.Request, title string) {
	name := r.FormValue("name")
	fmt.Printf("enterHandler %s name %s\n", title, name)
	url := "/edit/" + name
	fmt.Printf("redirect to %s\n", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Printf("viewHandler %s\n", title)
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Printf("editHandler %s\n", title)
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Printf("saveHandler %s\n", title)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "start.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view|enter)/([a-zA-Z0-9]+)|/$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("handler: %s\n", r.URL.Path)
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Printf("handler: m %v\n", m)
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", makeHandler(startHandler))
	http.HandleFunc("/enter/", makeHandler(enterHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}