package ping

import "net/http"

type Ping struct{}

func (p *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func NewPing() *Ping {
	return &Ping{}
}
