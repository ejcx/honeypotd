package honeypots

import (
	"fmt"
	"log"
	"net/http"
)

type HoneyPot struct {
	Port    string
	Address string
}

type HoneyPotRunnable interface {
	Run(h *HoneyPot) error
}

type HTTPPot struct {
}

func (p *HTTPPot) Run(h *HoneyPot) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { return })
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", h.Address, h.Port), nil))
	return nil
}