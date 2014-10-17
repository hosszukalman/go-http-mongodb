package main

import (
	"encoding/json"
	"github.com/goji/httpauth"
	"github.com/julienschmidt/httprouter"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

func deleteProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := mongoCollection.Remove(bson.M{"name": ps.ByName("name")})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result *interface{}
	err := mongoCollection.Find(bson.M{"name": ps.ByName("name")}).One(&result)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		data, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}
	}
}

func saveProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var m map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	changeInfo, err := mongoCollection.Upsert(bson.M{"name": m["name"]}, &m)
	if changeInfo.Updated != 0 {
		log.Println("Updated")
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

var mongoCollection *mgo.Collection

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection = session.DB("gohttptest").C("users")

	router := httprouter.New()

	router.GET("/user/:name", getProfile)
	router.DELETE("/user/:name", deleteProfile)
	router.PUT("/user", saveProfile)

	log.Fatal(http.ListenAndServe(":8080", httpauth.SimpleBasicAuth("user", "pass")(router)))
}
