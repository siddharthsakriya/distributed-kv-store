package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func Handler(kv *KVServer) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		val, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		respond(w, kv, Command{Action: "PUT", Key: r.PathValue("key"), Val: val})
	})

	mux.HandleFunc("GET /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, kv, Command{Action: "GET", Key: r.PathValue("key")})
	})

	mux.HandleFunc("DELETE /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		respond(w, kv, Command{Action: "DELETE", Key: r.PathValue("key")})
	})

	return mux
}

func respond(w http.ResponseWriter, kv *KVServer, c Command) {
	data, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := kv.Submit(data)

	switch {
	case errors.Is(err, ErrNotLeader):
		http.Error(w, "not leader", http.StatusServiceUnavailable)
	case errors.Is(err, ErrTimeout):
		http.Error(w, "timeout", http.StatusGatewayTimeout)
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		w.Write(result)
	}
}
