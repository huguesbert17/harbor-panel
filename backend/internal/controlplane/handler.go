// Package controlplane exposes the first API surface over domain contracts.
package controlplane

import (
	"encoding/json"
	"net/http"

	"github.com/harbor-panel/harbor-panel/internal/domains"
	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/platform/httpserver"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

type Handler struct {
	accounts hostingaccounts.PersistentRepository
	servers  servers.PersistentRepository
	domains  domains.Service
}

func NewHandler(accounts hostingaccounts.PersistentRepository, serverStore servers.PersistentRepository, domainService domains.Service) http.Handler {
	handler := &Handler{accounts: accounts, servers: serverStore, domains: domainService}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/accounts", handler.listAccounts)
	mux.HandleFunc("GET /api/v1/servers", handler.listServers)
	mux.HandleFunc("POST /api/v1/servers/registrations", handler.registerServer)
	mux.HandleFunc("GET /api/v1/domains", handler.listDomains)
	mux.HandleFunc("POST /api/v1/domains", handler.createDomain)
	mux.HandleFunc("GET /api/v1/dns/zones", handler.listZones)
	mux.HandleFunc("POST /api/v1/dns/zones", handler.createZone)
	mux.HandleFunc("GET /api/v1/dns/zones/{zoneID}/records", handler.listRecords)
	mux.HandleFunc("POST /api/v1/dns/zones/{zoneID}/records", handler.createRecord)
	return mux
}

func (handler *Handler) listDomains(w http.ResponseWriter, request *http.Request) {
	items, err := handler.domains.ListDomains(request.Context(), request.URL.Query().Get("account_id"))
	if err != nil {
		httpserver.JSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to list domains"})
		return
	}
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (handler *Handler) createDomain(w http.ResponseWriter, request *http.Request) {
	var command domains.CreateDomainCommand
	if err := json.NewDecoder(request.Body).Decode(&command); err != nil {
		httpserver.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	item, err := handler.domains.CreateDomain(request.Context(), command)
	if err != nil {
		httpserver.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	httpserver.JSON(w, http.StatusCreated, map[string]any{"data": item})
}

func (handler *Handler) listZones(w http.ResponseWriter, request *http.Request) {
	items, err := handler.domains.ListZones(request.Context(), request.URL.Query().Get("account_id"))
	if err != nil {
		httpserver.JSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to list DNS zones"})
		return
	}
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (handler *Handler) createZone(w http.ResponseWriter, request *http.Request) {
	var command domains.CreateZoneCommand
	if err := json.NewDecoder(request.Body).Decode(&command); err != nil {
		httpserver.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	item, err := handler.domains.CreateZone(request.Context(), command)
	if err != nil {
		httpserver.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	httpserver.JSON(w, http.StatusCreated, map[string]any{"data": item})
}

func (handler *Handler) listRecords(w http.ResponseWriter, request *http.Request) {
	items, err := handler.domains.ListRecords(request.Context(), request.PathValue("zoneID"))
	if err != nil {
		httpserver.JSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to list DNS records"})
		return
	}
	httpserver.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (handler *Handler) createRecord(w http.ResponseWriter, request *http.Request) {
	var command domains.CreateDNSRecordCommand
	if err := json.NewDecoder(request.Body).Decode(&command); err != nil {
		httpserver.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	item, err := handler.domains.CreateRecord(request.Context(), request.PathValue("zoneID"), command)
	if err != nil {
		httpserver.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	httpserver.JSON(w, http.StatusCreated, map[string]any{"data": item})
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
