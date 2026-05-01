package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "giftcard/application/core/models"
    "giftcard/application/core/usecase"
    adminservice "giftcard/application/roles/admin/service"
    customerservice "giftcard/application/roles/customer/service"
    retailerservice "giftcard/application/roles/retailer/service"
)

type HandlerDeps struct {
    Issuer   *usecase.IssuerService
    Retailer *retailerservice.Service
    Admin    *adminservice.Service
    Customer *customerservice.Service
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string) {
    lower := strings.ToLower(msg)
    status := http.StatusInternalServerError
    switch {
    case strings.Contains(lower, "does not exist"):
        status = http.StatusNotFound
    case strings.Contains(lower, "already exists"):
        status = http.StatusConflict
    case strings.Contains(lower, "access denied"), strings.Contains(lower, "unauthorized"):
        status = http.StatusForbidden
    case strings.Contains(lower, "invalid"), strings.Contains(lower, "empty"), strings.Contains(lower, "negative"):
        status = http.StatusBadRequest
    }
    writeJSON(w, status, map[string]string{"error": msg})
}

func reqCtx(r *http.Request) (context.Context, context.CancelFunc) {
    return context.WithTimeout(r.Context(), 30*time.Second)
}

// getCard dispatches to the right service based on role.
func getCard(ctx context.Context, cardID, role string, deps HandlerDeps) (*models.GiftCard, error) {
    switch role {
    case "retailer":
        return deps.Retailer.GetGiftCard(ctx, cardID)
    case "customer":
        return deps.Customer.GetGiftCard(ctx, cardID)
    case "admin":
        return deps.Admin.GetGiftCard(ctx, cardID)
    default:
        return deps.Issuer.GetGiftCard(ctx, cardID)
    }
}

func getHistory(ctx context.Context, cardID, role string, deps HandlerDeps) ([]models.GiftCardEvent, error) {
    switch role {
    case "retailer":
        return deps.Retailer.GetGiftCardHistory(ctx, cardID)
    case "customer":
        return deps.Customer.GetGiftCardHistory(ctx, cardID)
    case "admin":
        return deps.Admin.GetGiftCardHistory(ctx, cardID)
    default:
        return deps.Issuer.GetGiftCardHistory(ctx, cardID)
    }
}

func HandleCreateCard(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var body struct {
            CardID   string  `json:"cardID"`
            IssuerID string  `json:"issuerID"`
            Balance  float64 `json:"balance"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            writeError(w, "invalid request body")
            return
        }
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Issuer.IssueGiftCard(ctx, body.CardID, body.IssuerID, body.Balance)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusCreated, map[string]string{"message": msg})
    }
}

func HandleActivate(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Retailer.ActivateGiftCard(ctx, cardID)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"message": msg})
    }
}

func HandleTransfer(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        var body struct {
            NewOwnerID string `json:"newOwnerID"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            writeError(w, "invalid request body")
            return
        }
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Retailer.TransferGiftCard(ctx, cardID, body.NewOwnerID)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"message": msg})
    }
}

func HandleRedeem(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        var body struct {
            Amount float64 `json:"amount"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            writeError(w, "invalid request body")
            return
        }
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Customer.RedeemGiftCard(ctx, cardID, body.Amount)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"message": msg})
    }
}

func HandleSuspend(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        var body struct {
            Reason string `json:"reason"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            writeError(w, "invalid request body")
            return
        }
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Admin.SuspendGiftCard(ctx, cardID, body.Reason)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"message": msg})
    }
}

func HandleReactivate(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        ctx, cancel := reqCtx(r)
        defer cancel()
        msg, err := deps.Admin.ReactivateGiftCard(ctx, cardID)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"message": msg})
    }
}

func HandleGetCard(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        role := r.URL.Query().Get("role")
        ctx, cancel := reqCtx(r)
        defer cancel()
        card, err := getCard(ctx, cardID, role, deps)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, card)
    }
}

func HandleGetBalance(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        role := r.URL.Query().Get("role")
        ctx, cancel := reqCtx(r)
        defer cancel()
        card, err := getCard(ctx, cardID, role, deps)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, map[string]float64{"balance": card.Balance})
    }
}

func HandleGetHistory(deps HandlerDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cardID := r.PathValue("id")
        role := r.URL.Query().Get("role")
        ctx, cancel := reqCtx(r)
        defer cancel()
        events, err := getHistory(ctx, cardID, role, deps)
        if err != nil {
            writeError(w, err.Error())
            return
        }
        writeJSON(w, http.StatusOK, events)
    }
}
