package data

import "time"

type SubscriptionRequest struct {
	SubscriptionName string    `json:"subscription_name"`
	VendorName       string    `json:"vendor_name"`
	BillingCycle     string    `json:"billing_cycle"`
	RenewaleDate     string    `json:"renewale_date"`
	Categorie        string    `json:"categorie"`
	Status           string    `json:"status"`
	Cost             float32   `json:"cost"`
	Created_at       time.Time `json:"created_at"`
	Updated_at       time.Time `json:"updated_at"`
}

type SubscriptionResponse struct {
	SubscriptionName string    `json:"subscription_name"`
	VendorName       string    `json:"vendor_name"`
	BillingCycle     string    `json:"billing_cycle"`
	RenewaleDate     string    `json:"renewale_date"`
	Categorie        string    `json:"categorie"`
	Status           string    `json:"status"`
	Cost             float32   `json:"cost"`
}

