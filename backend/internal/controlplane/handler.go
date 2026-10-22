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
	accounts hostingaccounts.PersistentRepository
	servers  servers.PersistentRepository
}

func NewHandler(accounts hostingaccounts.PersistentRepository, serverStore servers.PersistentRepository) http.Handler {
	handler := &Handler{accounts: accounts, servers: serverStore}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/accounts", handler.listAccounts)
	mux.HandleFunc("GET /api/v1/servers", handler.listServers)
	mux.HandleFunc("POST /api/v1/servers/registrations", handler.registerServer)
	return mux
}

func (handler *Handler) listAccounts(w http.ResponseWriter, request *http.Request) {
	accounts, err := handler.accounts.ListAll(request.Context())
	if err != nil {
		httpserver.JSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to list hosting accounts"})
		return
	}
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": accounts})
}

func (handler *Handler) listServers(w http.ResponseWriter, request *http.Request) {
	servers, err := handler.servers.ListAll(request.Context())
	if err != nil {
		httpserver.JSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to list managed servers"})
		return
	}
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": servers})
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
