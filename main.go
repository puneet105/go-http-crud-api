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