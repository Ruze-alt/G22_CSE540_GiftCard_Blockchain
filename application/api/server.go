package api

import (
    "net/http"

    "giftcard/application/api/handlers"
    "giftcard/application/api/middleware"
)

func NewServer(deps handlers.HandlerDeps, frontendDir string) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("POST /cards", handlers.HandleCreateCard(deps))
    mux.HandleFunc("POST /cards/{id}/activate", handlers.HandleActivate(deps))
    mux.HandleFunc("POST /cards/{id}/transfer", handlers.HandleTransfer(deps))
    mux.HandleFunc("POST /cards/{id}/redeem", handlers.HandleRedeem(deps))
    mux.HandleFunc("POST /cards/{id}/suspend", handlers.HandleSuspend(deps))
    mux.HandleFunc("POST /cards/{id}/reactivate", handlers.HandleReactivate(deps))
    mux.HandleFunc("GET /cards/{id}", handlers.HandleGetCard(deps))
    mux.HandleFunc("GET /cards/{id}/balance", handlers.HandleGetBalance(deps))
    mux.HandleFunc("GET /cards/{id}/history", handlers.HandleGetHistory(deps))

    mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

    return middleware.CORS(mux)
}
