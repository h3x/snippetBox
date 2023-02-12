package main

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "errors"
     
    "github.com/h3x/snippetbox/internal/models"
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        app.notFound(w)
		return
	}
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    ts,err := template.ParseFiles(files...)
    if err != nil {
        app.errorLog.Println(err.Error())
        app.serverError(w, err)
        return
    }
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        app.serverError(w, err)
    }
}


func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
        app.notFound(w)
		return
	}

    snip,err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }
    fmt.Fprintf(w, "%+v", snip)
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w,http.StatusMethodNotAllowed)
		return
	}

    // tmp
    title := "O snail"
    content := "O snail\nClimb mount fuji\nBut slowly, slowly"
    expires := 7

    id,err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }
    
    http.Redirect(w,r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
