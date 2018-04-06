package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/archerwq/photo-viewer/conf"
	"github.com/archerwq/photo-viewer/dao"
	. "github.com/archerwq/photo-viewer/pvlog"
	"github.com/archerwq/photo-viewer/server"
)

var configFile = flag.String("config", "etc/photo-viewer.toml", "photo-viewer config file")

func main() {
	PVLog.Println("starting...")
	config := loadConfig()

	err := dao.InitManager(config)
	if err != nil {
		PVLog.Fatal(err.Error())
	}
	defer dao.CleanManager()

	server := server.New(config)
	server.Run()
	defer server.Close()

	wait()
	PVLog.Println("shuting down...")
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
