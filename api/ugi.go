package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/archerwq/photo-viewer/dao"
	"github.com/gorilla/mux"
)

var PlaceHodler = 1

// GET /api/ugis/{sha1}
func GetUGI(w http.ResponseWriter, r *http.Request) {
	sha1 := mux.Vars(r)["sha1"]
	ugi, err := dao.Manager.UGIDao.Get(sha1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(map[string]interface{}{
		"success": true,
		"ugi":     ugi,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type putRequest struct {
	Tag   string `json: "tag"`
	Story string `json: "story"`
}

// PUT /api/ugis/{sha1}
func PutUGI(w http.ResponseWriter, r *http.Request) {
	sha1 := mux.Vars(r)["sha1"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var pr putRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if pr.Tag != "" {
		err = dao.Manager.UGIDao.AddTag(sha1, pr.Tag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if pr.Story != "" {
		err = dao.Manager.UGIDao.AddStory(sha1, pr.Story)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	js, err := json.Marshal(map[string]interface{}{"success": true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
