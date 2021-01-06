package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"text/template"
	"time"
)

type emily struct {
	CM struct {
		Healthy int `json:"healthy"`
		Total   int `json:"total"`
	} `json:"cm"`
	Web struct {
		Community bool `json:"community"`
		Store     bool `json:"store"`
		API       bool `json:"api"`
	}
}

type tmpl struct {
	HealthyPct int
	Store      bool
	Community  bool
	API        bool
}

var (
	em        emily
	lastFetch int64
	file      *template.Template
	hpc       *http.Client
)

func updateEm() (*emily, error) {
	if hpc == nil {
		hpc = &http.Client{}
	}

	res, err := hpc.Get(os.Getenv("EMILY"))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bx, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	em := &emily{}
	err = json.Unmarshal(bx, em)
	if err != nil {
		return nil, err
	}

	return em, nil
}

func main() {
	if os.Getenv("EMILY") == "" {
		panic("A path to an Emily server must be specified as the EMILY environment variable.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3201"
	}

	bx, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(fmt.Errorf("Could not open template file for reading: %v", err))
	}

	file, err = template.New("maria").Parse(string(bx))
	if err != nil {
		panic(fmt.Errorf("Could not parse template: %v", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Unix()-lastFetch > 15 {
			pEm, err := updateEm()
			if err != nil {
				w.WriteHeader(500)
				w.Header().Add("Content-Type", "text/plain")
				w.Write([]byte("Could not conntect to Emily"))

				return
			}

			lastFetch = time.Now().Unix()
			em = *pEm
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")

		file.Execute(w, tmpl{
			HealthyPct: int(math.Floor((float64(em.CM.Healthy) / float64(em.CM.Total)) * 100)),
			Store:      em.Web.Store,
			Community:  em.Web.Store,
			API:        em.Web.API,
		})
	})

	panic(http.ListenAndServe(":"+port, nil))
}
