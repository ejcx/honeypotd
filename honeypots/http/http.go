package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ejcx/honeypotd/honeypots"
	"github.com/ejcx/honeypotd/notification/twilio"
)

type HTTPPot struct {
}

func (p *HTTPPot) Run(h *honeypots.HoneyPot) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := twilio.Notify(fmt.Sprintf("Access from %s\n", r.RemoteAddr))
		if err != nil {
			log.Println("Error notifying", err)
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", h.Address, h.Port), nil))
	return nil
}
