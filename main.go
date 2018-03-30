package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/archerwq/photo-viewer/conf"
	. "github.com/archerwq/photo-viewer/pvlog"
	"github.com/archerwq/photo-viewer/server"
)

var configFile = flag.String("config", "etc/photo-viewer.toml", "photo-viewer config file")

func main() {
	PVLog.Println("starting...")
	config := loadConfig()

	server := server.New(config)
	server.Run()

	wait()

	fmt.Println("shuting down...")
	server.Close()
	PVLog.Println("stopped!")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Path[len("/echo/"):]
	fmt.Fprint(w, m)
}

func loadConfig() *conf.Config {
	flag.Parse()
	config, err := conf.Load(*configFile)
	if err != nil {
		PVLog.Fatal(err.Error())
	}
	PVLog.Println(*config)
	return config
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
