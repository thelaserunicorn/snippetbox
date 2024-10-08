package main
import (
  "fmt"
  "html/template"
  "log"
  "net/http"
  "strconv"
)

func home(w http.ResponseWriter, r *http.Request){
  w.Header().Add("Server", "GO")

  files := []string{
    "./ui/html/base.tmpl",
    "./ui/html/partials/nav.tmpl",
    "./ui/html/pages/home.tmpl",
  }
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Print(err.Error())
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  err = ts.ExecuteTemplate(w,"base", nil)
  if err != nil {
    log.Print(err.Error())
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

func snippetView(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil || id < 1 {
    http.NotFound(w, r)
    return
  }
  fmt.Fprintf(w, "SNIPPET %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("FORM FOR CREATING A SNIPPET"))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusCreated)
  w.Write([]byte("SAVE A NEW SNIPPET"))
}
