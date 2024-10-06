package banyan

import (
	"flag"
	"net/http"

	"banyan/config"
	"banyan/log"
)

func Init() {
	flag.Parse()
	log.Setup()
	config.Configuration.Load()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}
