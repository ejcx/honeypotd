package honeypots

type HoneyPot struct {
	Port    string
	Address string
}

type HoneyPotRunnable interface {
	Run(h *HoneyPot) error
}
