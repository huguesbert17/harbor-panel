// Package controlplane exposes the first API surface over domain contracts.
package controlplane

import (
	"encoding/json"
	"net/http"

	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/platform/httpserver"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

type Handler struct {
	accounts *hostingaccounts.MemoryRepository
	servers  *servers.MemoryRepository
}

func NewHandler(accounts *hostingaccounts.MemoryRepository, serverStore *servers.MemoryRepository) http.Handler {
	handler := &Handler{accounts: accounts, servers: serverStore}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/accounts", handler.listAccounts)
	mux.HandleFunc("GET /api/v1/servers", handler.listServers)
	mux.HandleFunc("POST /api/v1/servers/registrations", handler.registerServer)
	return mux
}

func (handler *Handler) listAccounts(w http.ResponseWriter, request *http.Request) {
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": handler.accounts.List(request.Context())})
}

func (handler *Handler) listServers(w http.ResponseWriter, request *http.Request) {
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": handler.servers.List(request.Context())})
}

func (handler *Handler) registerServer(w http.ResponseWriter, request *http.Request) {
	var payload servers.RegistrationRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		httpserver.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	server, err := handler.servers.Register(request.Context(), payload)
	if err != nil {
		httpserver.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	httpserver.JSON(w, http.StatusCreated, map[string]any{"data": server})
}
