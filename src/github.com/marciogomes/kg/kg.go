package kg

import (
  "fmt"
  "log"
  "os"

  "github.com/cayleygraph/cayley"
  "github.com/cayleygraph/cayley/quad"
  "github.com/cayleygraph/cayley/quad/cquads"
)

func Init() *cayley.Handle {
  store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

  file, err := os.Open("data/quads.nq")
	if(err != nil) {
		log.Fatalln(err)
	}

  decoder := cquads.NewDecoder(file)
	for {
		aux, err := decoder.Unmarshal()
		if err != nil {
			break
		}
		store.AddQuad(aux)
		//fmt.Println(aux)
	}

  return store
}

func Query(store *cayley.Handle, sintomas []string) []string {
  var doencas []string
  // fazendo o caminho

  for i := range sintomas {
    p := cayley.StartPath(store, quad.IRI(sintomas[i])).In(quad.IRI("http://health-lifesci.schema.org/cause"))

    err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
  		nativeValue := quad.NativeOf(value)
      doencas = append(doencas, fmt.Sprint(nativeValue))
  	})
  	if err != nil {
  		log.Fatalln(err)
  	}
  }
  return doencas
}
