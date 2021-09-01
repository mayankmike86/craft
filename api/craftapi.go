package api

import (
	"craft/database"
	"craft/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	limitCount = 500
)

func InsertMydataConfig(db database.MongoUtilsInterfce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside InsertCraftConfig")
		craftConfig := models.CraftConfiguration{}

		body, err := ioutil.ReadAll(r.Body)
		fmt.Print("inside body", body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		err = json.Unmarshal(body, &craftConfig)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		res, err := db.Insert(craftConfig)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}
		setResponseHeaders(w)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)

	}
}

func GetMydataConfig(db database.MongoUtilsInterfce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside InsertCraftConfig")

		reqVars := mux.Vars(r)
		timestamp := reqVars["timestamp"]
		key := reqVars["key"]
		value := reqVars["value"]

		limit := r.URL.Query().Get("limit")
		limitInt := limitCount
		if limit != "" {
			limitInt, _ = strconv.Atoi(limit)
		}
		var params models.Params
		timestamp64, err := strconv.ParseInt(timestamp, 10, 64)
		params.Timestamp = timestamp64
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		params.Key = key
		params.Value = value
		params.Limit = limitInt

		IDs, err := db.FetchMyData(params)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}
		setResponseHeaders(w)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(IDs)

	}
}

func UpdateMydataConfig(db database.MongoUtilsInterfce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside InsertCraftConfig")
		craftConfig := models.CraftConfiguration{}
		reqVars := mux.Vars(r)
		id := reqVars["id"]
		body, err := ioutil.ReadAll(r.Body)
		fmt.Print("inside body", body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		err = json.Unmarshal(body, &craftConfig)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}
		var keys []string
		for _, attributes := range craftConfig.Attributes {
			keys = append(keys, attributes.Key)
		}
		res, msg, err := db.UpdateMydataConfig(id, keys, craftConfig)
		fmt.Println("msg: ", msg)
		if err != nil && msg == "Bad Request" {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(res)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}
		setResponseHeaders(w)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)

	}
}

func setResponseHeaders(respWriter http.ResponseWriter) {
	respWriter.Header().Set("content-type", "application/scim+json")
}
func RegisterRoutes(router *mux.Router, db database.MongoUtilsInterfce) {

	router.Methods("GET").Path("/health").HandlerFunc(handleHealthCheck)

	router.HandleFunc("/mydata", InsertMydataConfig(db)).Methods("POST")
	router.HandleFunc("/mydata/{timestamp}/{key}/{value}", GetMydataConfig(db)).Methods("GET")
	router.HandleFunc("/mydata/{id}", UpdateMydataConfig(db)).Methods("PATCH")

}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
func handleHealthCheck(respWriter http.ResponseWriter, _ *http.Request) {
	fmt.Println("health is ok")

	var data = []byte("Server is up and running")
	writeJSONResponse(respWriter, http.StatusOK, data)

}
