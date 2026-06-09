// consul-api — a small net/http proxy over the official hashicorp/consul/api client.
// Exposes register / health / deregister, mirroring the lesson contract across all 4 languages.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
)

// Handler holds the Consul client shared by every route.
type Handler struct {
	consul *api.Client
}

// registerRequest is the JSON body accepted by POST /consul/register.
type registerRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// writeJSON writes a JSON response with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// register writes the instance into the Consul catalog via the Agent API.
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// The official client maps Go struct fields to Consul's PascalCase wire format.
	reg := &api.AgentServiceRegistration{
		ID:      req.ID,
		Name:    req.Name,
		Address: req.Address,
		Port:    req.Port,
	}
	if err := h.consul.Agent().ServiceRegister(reg); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "registered"})
}

// health lists only the passing instances of a service by logical name.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	service := r.PathValue("service")
	// passingOnly=true makes Consul join catalog + check state and return only live instances.
	entries, _, err := h.consul.Health().Service(service, "", true, nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// deregister removes exactly one instance from the catalog by its instance ID.
func (h *Handler) deregister(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Deregister by exact instance ID, not by service name.
	if err := h.consul.Agent().ServiceDeregister(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "deregistered"})
}

func main() {
	cfg := api.DefaultConfig()
	// CONSUL_BASE_URL is host:port form; strip the scheme for the official client Address field.
	if base := os.Getenv("CONSUL_BASE_URL"); base != "" {
		cfg.Address = stripScheme(base)
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	h := &Handler{consul: client}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /consul/register", h.register)
	mux.HandleFunc("GET /consul/health/{service}", h.health)
	mux.HandleFunc("POST /consul/deregister/{id}", h.deregister)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := "0.0.0.0:" + port
	log.Printf("consul-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// stripScheme removes a leading http:// or https:// so the client receives host:port.
func stripScheme(url string) string {
	for _, prefix := range []string{"http://", "https://"} {
		if len(url) >= len(prefix) && url[:len(prefix)] == prefix {
			return url[len(prefix):]
		}
	}
	return url
}
