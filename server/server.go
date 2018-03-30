package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/archerwq/photo-viewer/api"
	"github.com/archerwq/photo-viewer/conf"
	. "github.com/archerwq/photo-viewer/pvlog"
	"github.com/gorilla/mux"
)

type Server struct {
	config *conf.Config
	webSrv *http.Server
}

func New(config *conf.Config) *Server {
	s := new(Server)
	s.config = config
	s.webSrv = newWebSrv(config.HTTPAddr, config.Resource)
	return s
}

func newWebSrv(httpAddr, resource string) *http.Server {
	rtr := mux.NewRouter()
	registerURLs(rtr, resource)
	http.Handle("/", rtr)

	return &http.Server{
		Handler:      rtr,
		Addr:         httpAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func registerURLs(rtr *mux.Router, resource string) {
	rtr.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/view.html", resource))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	rtr.HandleFunc("/api/photos", api.QueryPhotos).Methods("GET")

	rtr.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, _ := route.GetPathTemplate()
		qt, _ := route.GetQueriesTemplates()
		p, _ := route.GetPathRegexp()
		qr, _ := route.GetQueriesRegexp()
		m, _ := route.GetMethods()
		log.Println(strings.Join(m, ","), strings.Join(qt, ","), strings.Join(qr, ","), t, p)
		return nil
	})
}

func (s *Server) Run() {
	PVLog.Println("Listening...")
	go func() {
		if err := s.webSrv.ListenAndServe(); err != nil {
			PVLog.Println(err)
		}
	}()
}

func (s *Server) Close() {
	err := s.webSrv.Close()
	if err != nil {
		PVLog.Println("failed to close HTTP Server")
	}
}
