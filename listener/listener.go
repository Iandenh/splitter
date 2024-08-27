package listener

import (
	"fmt"
	"log"
	"net/http"
)

type Listener struct {
	originHostName string
	rewriteHost    bool
	port           int
	upstreams      []string
}

func New(originHostName string, rewriteHost bool, port int, upstreams []string) Listener {
	return Listener{
		originHostName: originHostName,
		rewriteHost:    rewriteHost,
		port:           port,
		upstreams:      upstreams,
	}
}

func (l *Listener) Start() {
	proxyServer := http.NewServeMux()

	proxyServer.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		l.handleRequest(req, w)
	})

	log.Printf("Starting proxy at: http://localhost:%d\n", l.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", l.port), proxyServer); err != nil {
		log.Fatal(err)
	}
}
