package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func memoryAlocationBenchmark(w http.ResponseWriter, req *http.Request) {
	cont, err := strconv.Atoi(mux.Vars(req)["iterations"])
	if err != nil {
		cont = 10000
	}
	text := ""
	for index := 0; index < cont; index++ {
		str := "Memory test. Save data in variables. Number: " + strconv.Itoa(index) + " \n"
		text = text + str
	}
	fmt.Fprintf(w, "Memory test. Executed a loop with %s iterations!\n", strconv.Itoa(cont))
}

func fileReadBenchmark(w http.ResponseWriter, req *http.Request) {
	cont, err := strconv.Atoi(mux.Vars(req)["iterations"])
	fmt.Println(err)
	if err != nil {
		cont = 10000
	}
	text := ""
	for index := 0; index < cont; index++ {
		content, err := ioutil.ReadFile("files/test.txt")
		check(err)
		text = string(content)
	}
	fmt.Fprintf(w, "Reading file test. Executed a loop with %s iterations!\nFile content:\n%s", strconv.Itoa(cont), text)
}

func fileWriteBenchmark(w http.ResponseWriter, req *http.Request) {
	cont, err := strconv.Atoi(mux.Vars(req)["iterations"])
	if err != nil {
		cont = 10000
	}
	dat, err := ioutil.ReadFile("files/test.txt")
	check(err)
	id := uuid.New()
	os.MkdirAll("tmp/"+id.String(), os.ModePerm)
	for index := 0; index < cont; index++ {
		err = ioutil.WriteFile("tmp/"+id.String()+"/test"+strconv.Itoa(index)+".txt", dat, 0644)
		check(err)
	}
	os.RemoveAll("tmp/" + id.String())
	fmt.Fprintf(w, "Writing file test. RequestID: "+id.String()+"Executed a loop with %s iterations!\n", strconv.Itoa(cont))
}

func pingHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "Pong!")

}

func interServicesBenchmark(w http.ResponseWriter, req *http.Request) {
	url, ok := req.URL.Query()["url"]
	if !ok || len(url[0]) < 1 {
		fmt.Println("Url Param 'url' is missing")
		return
	}
	log.Println(url[0])
	request, err := http.NewRequest("GET", url[0], nil)
	check(err)
	fmt.Fprintf(w, "Remote Test! "+request.URL.String())
}

func allOperationsBenchmark(w http.ResponseWriter, req *http.Request) {
	memoryAlocationBenchmark(w, req)
	fileReadBenchmark(w, req)
	fileWriteBenchmark(w, req)
}

func main() {
	r := mux.NewRouter()
	r.Use(apmgorilla.Middleware())
	r.HandleFunc("/memoryAlocationBenchmark", memoryAlocationBenchmark)
	r.HandleFunc("/memoryAlocationBenchmark/{iterations}", memoryAlocationBenchmark)
	r.HandleFunc("/fileReadBenchmark", fileReadBenchmark)
	r.HandleFunc("/fileReadBenchmark/{iterations}", fileReadBenchmark)
	r.HandleFunc("/fileWriteBenchmark", fileWriteBenchmark)
	r.HandleFunc("/fileWriteBenchmark/{iterations}", fileWriteBenchmark)
	r.HandleFunc("/ping", pingHandler)
	r.HandleFunc("/interServicesBenchmark", interServicesBenchmark)
	r.HandleFunc("/allOperations", allOperationsBenchmark)
	fmt.Println("Server is up and listening on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", r))
}
