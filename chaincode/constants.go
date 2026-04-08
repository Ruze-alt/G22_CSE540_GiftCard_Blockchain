package main

type GiftCardStatus string
type ParticipantRole string
type EventType string

const (
	StatusCreated           GiftCardStatus = "CREATED"
	StatusTransferred       GiftCardStatus = "TRANSFERRED"
	StatusActivated         GiftCardStatus = "ACTIVATED"
	StatusPartiallyRedeemed GiftCardStatus = "PARTIALLY_REDEEMED"
	StatusRedeemed          GiftCardStatus = "REDEEMED"
	StatusSuspended         GiftCardStatus = "SUSPENDED"
)

const (
	RoleIssuer   ParticipantRole = "ISSUER"
	RoleRetailer ParticipantRole = "RETAILER"
	RoleCustomer ParticipantRole = "CUSTOMER"
	RoleAdmin    ParticipantRole = "ADMIN"
)

const (
	EventTypeCardCreated     EventType = "CARD_CREATED"
	EventTypeCardTransferred EventType = "CARD_TRANSFERRED"
	EventTypeCardActivated   EventType = "CARD_ACTIVATED"
	EventTypeCardRedeemed    EventType = "CARD_REDEEMED"
	EventTypeCardSuspended   EventType = "CARD_SUSPENDED"
	EventTypeCardReactivated EventType = "CARD_REACTIVATED"
)

const (
	Org1MSP = "Org1MSP"
	Org2MSP = "Org2MSP"
)

const EventObjectType string = "event"
