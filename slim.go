package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "0.0.0.0:8000", "Listening address and port")

type BytesHandler []byte

func (h BytesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write(h)
}

func main() {
	flag.Parse()
	var handler http.Handler
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal("Error getting working directory: ", err)
		return
	} else {
		if fi, err := os.Stat(cwd); err != nil {
			log.Fatal(fmt.Sprintf("Error stating %s: ", cwd), err)
			return
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				handler = http.FileServer(http.Dir(cwd))
			case mode.IsRegular():
				if bytes, err := ioutil.ReadFile(cwd); err != nil {
					log.Fatal(fmt.Sprintf("Error reading %s: ", cwd), err)
				} else {
					handler = BytesHandler(bytes)
				}
			}
		}
	}
	fmt.Printf("Slim is now listening on %s\n", *addr)
	http.ListenAndServe(*addr, handler)
}
