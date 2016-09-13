/* Server Mux */

package main

import (
  "fmt"
  "log"
  "net/http"
  "html/template"
  "strings"
  "strconv"

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
    action := r.URL.Path[len("/edit/"):]
    if action == "symptoms" {
      value := r.FormValue("body")
      log.Println("Novo sintoma inserido: " + value)
      sintomas = append(sintomas, value)
      http.Redirect(w, r, "/edit/symptoms.html", http.StatusFound)
    } else if action == "signals" {
      fmt.Println(r.FormValue("temperature"))
      fmt.Println(r.FormValue("pressure"))
      http.Redirect(w, r, "/edit/signals.html", http.StatusFound)
    } else if action == "habits" {
      fmt.Println(r.FormValue("fuma"))
      http.Redirect(w, r, "/edit/habits.html", http.StatusFound)
    }
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

  tmpl := r.URL.Path[len("/results/"):]
  t, err := template.ParseFiles("tmpl/" + tmpl)
  if err != nil {
    log.Fatalln(err)
  }

  results := kg.Query(store, sintomas)

  /* fazer um metodo para filtragem */
  var diagnosticos map[string]int
  var resultsSet map[string]bool
  diagnosticos = make(map[string]int)
  resultsSet = make(map[string]bool)

  for i := range results {
    diagnosticos[results[i]]++
  }

  var resposta []string
  var aux string

  for i := range results {
    if resultsSet[results[i]] == false {
      resposta = append(resposta, results[i] + " ocorrencias = " + strconv.Itoa(diagnosticos[results[i]]) + " probabilidade = " + strconv.FormatFloat(float64(diagnosticos[results[i]])/float64(len(results)), 'f', 6, 64))
      resultsSet[results[i]] = true
    }
  }

  aux = strings.Join(resposta, " - ")

  t.Execute(w, aux)
}

func main() {
  store = kg.Init()

  http.HandleFunc("/home/", homeHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/results/", resultsHandler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
  http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img/"))))
  http.ListenAndServe(":8080", nil)
}
