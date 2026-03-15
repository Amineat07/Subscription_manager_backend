package data

import "time"

type SubscriptionRequest struct {
	SubscriptionName string    `json:"subscription_name" validate:"required"`
	VendorName       string    `json:"vendor_name" validate:"required"`
	BillingCycle     string    `json:"billing_cycle" validate:"required"`
	RenewaleDate     string    `json:"renewale_date" validate:"required"`
	Categorie        string    `json:"categorie" validate:"required"`
	Status           string    `json:"status" validate:"required"`
	Cost             float32   `json:"cost" validate:"required"`
	Created_at       time.Time `json:"created_at"`
	Updated_at       time.Time `json:"updated_at"`
}

type SubscriptionResponse struct {
	ID               int64      `json:"subscription_id"`
	SubscriptionName string     `json:"subscription_name"`
	VendorName       string     `json:"vendor_name"`
	BillingCycle     string     `json:"billing_cycle"`
	RenewaleDate     string     `json:"renewale_date"`
	Categorie        string     `json:"categorie"`
	Status           string     `json:"status"`
	Cost             float32    `json:"cost"`
	Created_at       time.Time  `json:"created_at"`
	Updated_at       time.Time  `json:"updated_at"`
	Deleted_at       *time.Time `json:"deleted_at"`
}

type SubscriptionRequestUpdate struct {
	SubscriptionName *string   `json:"subscription_name"`
	VendorName       *string   `json:"vendor_name"`
	BillingCycle     *string   `json:"billing_cycle"`
	RenewaleDate     *string   `json:"renewale_date"`
	Categorie        *string   `json:"categorie"`
	Status           *string   `json:"status"`
	Cost             *float32  `json:"cost"`
	Updated_at       time.Time `json:"updated_at"`
}
