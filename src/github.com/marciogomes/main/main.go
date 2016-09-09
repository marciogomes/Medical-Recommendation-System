package main

import (
  "fmt"
  "log"
  "net/http"
  "html/template"

  "github.com/marciogomes/kg"
  "github.com/cayleygraph/cayley"
)

var store *cayley.Handle
var sintomas []string

func homeHandler(w http.ResponseWriter, r *http.Request) {
  tmpl := r.URL.Path[len("/home/"):]
  t, err := template.ParseFiles("tmpl/" + tmpl)
  if err != nil {
    log.Fatalln(err)
  }
  t.Execute(w, nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    tmpl := r.URL.Path[len("/edit/"):]
    t, err := template.ParseFiles("tmpl/" + tmpl)
    if err != nil {
      log.Fatalln(err)
    }
    t.Execute(w, nil)
  } else {
    value := r.FormValue("body")
    fmt.Println("Novo sintoma inserido: " + value)
    sintomas = append(sintomas, value)
    http.Redirect(w, r, "/edit/edit.html", http.StatusFound)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println(sintomas)
  tmpl := r.URL.Path[len("/view/"):]
  t, err := template.ParseFiles("tmpl/" + tmpl)
  if err != nil {
    log.Fatalln(err)
  }
  t.Execute(w, sintomas)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println(sintomas)
  tmpl := r.URL.Path[len("/results/"):]
  t, err := template.ParseFiles("tmpl/" + tmpl)
  if err != nil {
    log.Fatalln(err)
  }

  doencas := kg.Query(store, sintomas)

  t.Execute(w, doencas)
}

func main() {
  store = kg.Init()

  http.HandleFunc("/home/", homeHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/results/", resultsHandler)
  http.ListenAndServe(":8080", nil)
}
