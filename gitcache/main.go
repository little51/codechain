package main

import (
	"flag"
	"log"
	"net/http"
)

var g_Basedir string
var port string

func main() {
	//log params
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("TRACE: ")
	//flag params
	flag.StringVar(&g_Basedir, "b", "f:/temp1", "默认为f:/temp1")
	flag.StringVar(&port, "p", "5000", "端口号，默认为5000")
	flag.Parse()
	log.Printf("basedir:%v port:%v", g_Basedir, port)
	//cron
	Cron()
	//listen
	http.HandleFunc("/", RequestHandler(g_Basedir))
	address := ":" + port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	} else {
		log.Printf("ListenAndServer: %s", address)
	}
}
