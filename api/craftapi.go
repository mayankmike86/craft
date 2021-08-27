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

func InsertMydataConfig(db database.MongoUtilsInterfce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside InsertCraftConfig", db)
		craftConfig := models.CraftConfiguration{}

		body, err := ioutil.ReadAll(r.Body)
		fmt.Print("inside body", body)
		if err != nil {
			fmt.Print("inside 1")
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		err = json.Unmarshal(body, &craftConfig)
		if err != nil {
			fmt.Print("inside 2")

			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		res, err := db.Insert(craftConfig)
		if err != nil {
			fmt.Print("inside 3")
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
		fmt.Print("inside InsertCraftConfig", db)
		//craftConfig := models.CraftConfiguration{}

		reqVars := mux.Vars(r)
		timestamp := reqVars["timestamp"]
		key := reqVars["key"]
		value := reqVars["value"]

		fmt.Println("inside timestamp: ", timestamp)
		fmt.Println("inside key: ", key)
		fmt.Println("inside value: ", value)

		var params models.Params
		timestamp64, err := strconv.ParseInt(timestamp, 10, 64)
		params.Timestamp = timestamp64
		if err != nil {
			fmt.Print("inside 1")
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}

		params.Key = key
		params.Value = value
		IDs, err := db.FetchMyData(params)

		// err = json.Unmarshal(body, &craftConfig)
		// if err != nil {
		// 	fmt.Print("inside 2")

		// 	w.WriteHeader(http.StatusBadRequest)
		// 	setResponseHeaders(w)
		// 	json.NewEncoder(w).Encode(err.Error())
		// }

		// res, err := db.Insert(craftConfig)
		if err != nil {
			fmt.Print("inside 3")
			w.WriteHeader(http.StatusBadRequest)
			setResponseHeaders(w)
			json.NewEncoder(w).Encode(err.Error())
		}
		setResponseHeaders(w)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(IDs)

	}
}

func setResponseHeaders(respWriter http.ResponseWriter) {
	respWriter.Header().Set("content-type", "application/scim+json")
}
func RegisterRoutes(router *mux.Router, db database.MongoUtilsInterfce) {

	router.Methods("GET").Path("/health").HandlerFunc(handleHealthCheck)

	router.HandleFunc("/mydata", InsertMydataConfig(db)).Methods("POST")
	router.HandleFunc("/mydata/{timestamp}/{key}/{value}", GetMydataConfig(db)).Methods("GET")

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
