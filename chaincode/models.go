package main

type GiftCard struct {
	CardID          string         `json:"cardId"`
	OwnerID         string         `json:"ownerId"`
	OwnerMSP        string         `json:"ownerMsp"`
	IssuerID        string         `json:"issuerId"`
	IssuerMSP       string         `json:"issuerMsp"`
	RetailerID      string         `json:"retailerId"`
	RetailerMSP     string         `json:"retailerMsp"`
	Balance         float64        `json:"balance"`
	OriginalBalance float64        `json:"originalBalance"`
	Status          GiftCardStatus `json:"status"`
	CreatedAt       string         `json:"createdAt"`
	ActivatedAt     string         `json:"activatedAt"`
	LastUpdatedAt   string         `json:"lastUpdatedAt"`
}

type GiftCardEvent struct {
	EventID   string    `json:"eventId"`
	CardID    string    `json:"cardId"`
	EventType EventType `json:"eventType"`
	ActorID   string    `json:"actorId"`
	ActorMSP  string    `json:"actorMsp"`
	ActorRole string    `json:"actorRole"`
	Timestamp string    `json:"timestamp"`
	TxID      string    `json:"txId"`
	Details   string    `json:"details"`
}

type ClientIdentityInfo struct {
	ClientID string `json:"clientId"`
	MSPID    string `json:"mspId"`
	Role     string `json:"role,omitempty"` // for now its optional since we're only using MSP (Org1 vs Org2)
}
