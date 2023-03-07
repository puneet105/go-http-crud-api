package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)
type Apple struct{
	ID int  				`json:"id"`
	Kind string 			`json:"kind"`
	BatchID int 			`json:"batchId"`
	StorageLocation string 	`json:"storageLocation"`
}

var apples []Apple

type AppleHandler struct {}

func NewAppleHandler() *AppleHandler {
	return &AppleHandler{}
}

func (h *AppleHandler) Router() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/apples" && r.Method == http.MethodPost {
			h.createApple(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/apples/") && r.Method == http.MethodDelete {
			h.deleteApple(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/apples/") && r.Method == http.MethodGet {
			h.getApple(w, r)
		} else if r.URL.Path == "/apple-stats" && r.Method == http.MethodGet {
			h.countApple(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (h *AppleHandler) createApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json")
	var apple Apple
	_ = json.NewDecoder(r.Body).Decode(&apple)
	apples = append(apples, apple)
	log.Println(apples)
	//json.NewEncoder(w).Encode(apple)
	fmt.Fprint(w,http.StatusOK)
}

func (h *AppleHandler) deleteApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json")
	rawID := strings.TrimPrefix(r.URL.Path, "/apples/")
	appleID, err := strconv.Atoi(rawID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for index, item := range apples{
		if item.ID == appleID{
			apples = append(apples[:index], apples[index+1:]...)
			break
		}
	}
	fmt.Fprint(w,http.StatusOK)
}

func (h *AppleHandler) getApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	rawID := strings.TrimPrefix(r.URL.Path, "/apples/")
	appleID, err := strconv.Atoi(rawID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, item := range apples{
		if item.ID == appleID{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func (h *AppleHandler) countApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	appleStatsResult := map[string]int{}
	for _, item := range apples{
		_, ok := appleStatsResult[item.Kind]
		if ok{
			appleStatsResult[item.Kind] = appleStatsResult[item.Kind] + 1
		}else{
			appleStatsResult[item.Kind] = 1
		}
	}
	json.NewEncoder(w).Encode(appleStatsResult)
}
func main(){
	appleHandler := NewAppleHandler()
	http.ListenAndServe(":8080", appleHandler.Router())
}



/*
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type AppleHandler struct {
}

func NewAppleHandler() *AppleHandler {
	return &AppleHandler{}
}

type Apple struct{
	ID int  //`json:"id"`
	Kind string //`json:"kind"`
	BatchID int //`json:"batchId"`
	StorageLocation string //`json:"storageLocation"`
}

var apples []Apple
// Router creates the following routes for the API
//
// GET /apples/ID
// DELETE /apples/ID
// POST /apples
// GET /apple-stats
func (h *AppleHandler) Router() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/apples" && r.Method == http.MethodPost {
			h.CreateApple(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/apples/") && r.Method == http.MethodDelete {
			h.DeleteApple(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/apples/") && r.Method == http.MethodGet {
			h.GetApple(w, r)
		} else if r.URL.Path == "/apple-stats" && r.Method == http.MethodGet {
			h.CountApples(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// CreateApple create a new Apple and stores it in the memory storage.
func (h *AppleHandler) CreateApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json")
	var apple Apple
	_ = json.NewDecoder(r.Body).Decode(&apple)
	apples = append(apples, apple)
	log.Println(apples)
	//json.NewEncoder(w).Encode(apple)
	fmt.Fprint(w,http.StatusOK)
}

// DeleteApple deletes an existing Apple from the memory storage.
func (h *AppleHandler) DeleteApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json")
	rawID := strings.TrimPrefix(r.URL.Path, "/apples/")
	appleID, err := strconv.Atoi(rawID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for index, item := range apples{
		if item.ID == appleID{
			apples = append(apples[:index], apples[index+1:]...)
			break
		}
	}
	fmt.Fprint(w,http.StatusOK)
	//json.NewEncoder(w).Encode(apples)
}

// GetApple gets an existing Apple from the memory storage.
func (h *AppleHandler) GetApple(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	rawID := strings.TrimPrefix(r.URL.Path, "/apples/")
	appleID, err := strconv.Atoi(rawID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, item := range apples{
		if item.ID == appleID{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// CountApples get counts for each kind of apples.
func (h *AppleHandler) CountApples(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	appleStatsResult := map[string]int{}
	for _, item := range apples{
		_, ok := appleStatsResult[item.Kind]
		if ok{
			appleStatsResult[item.Kind] = appleStatsResult[item.Kind] + 1
		}else{
			appleStatsResult[item.Kind] = 1
		}
	}
	json.NewEncoder(w).Encode(appleStatsResult)
}

func main() {
	appleHandler := NewAppleHandler()

	ts := httptest.NewServer(appleHandler.Router())
	defer ts.Close()

	runTests(ts)
}

func runTests(ts *httptest.Server) {
	// Create an Apple
	appleBody := "{\"ID\":1,\"Kind\":\"Golden\",\"BatchID\":42,\"StorageLocation\":\"NYC\"}"
	body := &bytes.Buffer{}
	body.WriteString(appleBody)
	res, err := http.Post(fmt.Sprintf("%s/apples", ts.URL), "text/json", body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("invalid status code %d", res.StatusCode)
	}

	// Get created apple
	res, err = http.Get(fmt.Sprintf("%s/apples/1", ts.URL))
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("invalid status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if strings.TrimSpace(string(b)) != appleBody {
		log.Fatalf("expected another result for apple 1: got %s, expected %s", string(b), appleBody)
	}

	// Create some apples
	apples := []string{
		"{\"ID\":2,\"Kind\":\"Granny smith\",\"BatchID\":43,\"StorageLocation\":\"NYC\"}",
		"{\"ID\":3,\"Kind\":\"Golden\",\"BatchID\":42,\"StorageLocation\":\"PAR\"}",
		"{\"ID\":4,\"Kind\":\"Golden\",\"BatchID\":44,\"StorageLocation\":\"NYC\"}",
		"{\"ID\":5,\"Kind\":\"Golden\",\"BatchID\":44,\"StorageLocation\":\"NYC\"}",
		"{\"ID\":6,\"Kind\":\"Granny smith\",\"BatchID\":42,\"StorageLocation\":\"NYC\"}",
		"{\"ID\":7,\"Kind\":\"Granny smith\",\"BatchID\":42,\"StorageLocation\":\"NYC\"}",
	}
	for _, a := range apples {
		body = &bytes.Buffer{}
		body.WriteString(a)
		res, err = http.Post(fmt.Sprintf("%s/apples", ts.URL), "text/json", body)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			log.Fatalf("invalid status code %d", res.StatusCode)
		}
	}

	// Get apple stats
	appleStatsResult := map[string]int{}
	res, err = http.Get(fmt.Sprintf("%s/apple-stats", ts.URL))
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("invalid status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&appleStatsResult); err != nil {
		log.Fatalf("could not decode apple stats: %s", err)
	}
	if appleStatsResult["Granny smith"] != 3 || appleStatsResult["Golden"] != 4 {
		log.Fatalf("invalid apple stats")
	}

	// Delete an apple
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/apples/1", ts.URL), nil)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("invalid status code %d", res.StatusCode)
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("invalid status code %d", res.StatusCode)
	}

	log.Println("all good!")
}
*/
