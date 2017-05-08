package main

import (
	"encoding/json"
	"eworker/bl"
	"eworker/work"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"

	"golang.org/x/net/netutil"
)

var (
	MaxWorker  int = 2
	MaxQueue   int = 100
	dispatcher *work.Dispatcher
)

func init() {

	// run work dispatch
	// work.Dispatch(MaxWorker, MaxQueue)

	dispatcher = work.NewDispatcher(MaxWorker, MaxQueue)
	dispatcher.Run()
}

func main() {
	// open max cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/payload", payloadHandler)

	// start web server
	log.Println("127.0.0.1:8084 run...")
	l, err := net.Listen("tcp", "127.0.0.1:8084")
	if err != nil {
		log.Printf("ListenError: %v", err)
	}
	defer l.Close()
	l = netutil.LimitListener(l, 1000000) //maximum connection
	http.Serve(l, nil)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	// Read the body into a string for json decoding
	var content = &bl.PayloadCollection{}
	err := json.NewDecoder(io.LimitReader(r.Body, 5000)).Decode(&content)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	// Go through each payload and queue items individually to be posted to S3
	for _, payload := range content.Payloads {
		// let's create a job with the payload
		// Push the work onto the queue.
		dispatcher.Add(payload)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
