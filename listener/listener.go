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
	log.Printf("Starting proxy at: http://localhost:%d\n", l.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", l.port), l); err != nil {
		log.Fatal(err)
	}
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.handleRequest(w, r)
}
