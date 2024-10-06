package node

import (
	"banyan/config"
	"banyan/log"
	"banyan/message"
	"io"
	"net/http"
	"net/url"
)

// http request header names
const (
	HTTPClientID  = "Id"
	HTTPCommandID = "Cid"
)

// serve serves the http REST API request from clients
func (n *node) http() {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", n.handleQuery)

	// http string should be in form of ":8080"
	ip, err := url.Parse(config.Configuration.HTTPAddrs[n.id])
	if err != nil {
		log.Fatal("http url parse error: ", err)
	}
	port := ":" + ip.Port()
	n.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Info("http server starting on ", port)
	log.Fatal(n.server.ListenAndServe())
}

func (n *node) handleQuery(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var query message.Query
	query.C = make(chan message.QueryReply)
	n.TxChan <- query
	reply := <-query.C
	_, err := io.WriteString(w, reply.Info)
	if err != nil {
		log.Error(err)
	}
}
