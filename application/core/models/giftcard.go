package models

type GiftCard struct {
    CardID          string  `json:"cardID"`
    OwnerID         string  `json:"ownerID"`
    OwnerMSP        string  `json:"ownerMSP"`
    IssuerID        string  `json:"issuerID"`
    IssuerMSP       string  `json:"issuerMSP"`
    RetailerID      string  `json:"retailerID,omitempty"`
    RetailerMSP     string  `json:"retailerMSP,omitempty"`
    Balance         float64 `json:"balance"`
    OriginalBalance float64 `json:"originalBalance"`
    Status          string  `json:"status"`
    CreatedAt       string  `json:"createdAt"`
    ActivatedAt     string  `json:"activatedAt,omitempty"`
    LastUpdatedAt   string  `json:"lastUpdatedAt"`
}

type GiftCardEvent struct {
    EventID       string `json:"eventID"`
    CardID        string `json:"cardID"`
    EventType     string `json:"eventType"`
    ActorID       string `json:"actorID"`
    ActorMSP      string `json:"actorMSP"`
    ActorRole     string `json:"actorRole"`
    Description   string `json:"description"`
    TransactionID string `json:"transactionID"`
    Timestamp     string `json:"timestamp"`
}
