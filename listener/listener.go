package listener

import (
	"fmt"
	"log"
	"net/http"
)

type Listener struct {
	OriginHostName string
	RewriteHost    bool
	Port           int
}

func (l *Listener) Start() {
	proxyServer := http.NewServeMux()

	proxyServer.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		l.handleRequest(req, w)
	})

	log.Printf("Starting proxy at: http://localhost:%d\n", l.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", l.Port), proxyServer); err != nil {
		log.Fatal(err)
	}
}
