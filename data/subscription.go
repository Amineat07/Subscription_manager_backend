package data

import "time"

type SubscriptionRequest struct {
	CompanyRequest     CompanyRequest `json:"company"`
	TagRequest         TagRequest     `json:"tag"`
	ContractStartDate  *string        `json:"contract_start_date"`
	ContractEndDate    *string        `json:"contract_end_date"`
	SubscriptionName   string         `json:"subscription_name" validate:"required"`
	Typ                *string        `json:"typ"`
	ContractNumber     string         `json:"contract_number"`
	CustomerNumber     string         `json:"customer_number"`
	PaymentMethod      string         `json:"payment_method" validate:"required"`
	BillingPeriod      string         `json:"billing_period" validate:"required"`
	Note               string         `json:"note"`
	CancellationPeriod *int64         `json:"cancellation_period"`
	BillingDate        *int64         `json:"billing_date" validate:"required"`
	Price              float64        `json:"price" validate:"required"`
}

type CompanyRequest struct {
	CompanyName   string `json:"company_name" validate:"required"`
	Category      string `json:"category"`
	ContactDetail string `json:"contact_detail"`
	Link          string `json:"link"`
}

type TagRequest struct {
	TagName  string `json:"tag_name"`
	TagColor string `json:"tag_color"`
}

type SubscriptionResponse struct {
	ID                 int64           `json:"id"`
	ContractStartDate  *time.Time      `json:"contract_start_date"`
	ContractEndDate    *time.Time      `json:"contract_end_date"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          *time.Time      `json:"updated_at"`
	DeletedAt          *time.Time      `json:"deleted_at"`
	Company            CompanyResponse `json:"company"`
	Tag                TagResponse     `json:"tag"`
	SubscriptionName   string          `json:"subscription_name"`
	Typ                string          `json:"typ"`
	ContractNumber     string          `json:"contract_number"`
	CustomerNumber     string          `json:"customer_number"`
	PaymentMethod      string          `json:"payment_method"`
	BillingPeriod      string          `json:"billing_period"`
	Note               string          `json:"note"`
	CreatedBy          *string         `json:"created_by"`
	UpdatedBy          *string         `json:"updated_by"`
	DeletedBy          *string         `json:"deleted_by"`
	CancellationPeriod *int64          `json:"cancellation_period"`
	BillingDate        int64           `json:"billing_date"`
	Price              float64         `json:"price"`
}

type CompanyResponse struct {
	ID            int64  `json:"id"`
	CompanyName   string `json:"company_name"`
	Category      string `json:"category"`
	ContactDetail string `json:"contact_detail"`
	Link          string `json:"link"`
}

type TagResponse struct {
	ID       int64  `json:"id"`
	TagName  string `json:"tag_name"`
	TagColor string `json:"tag_color"`
}

type SubscriptionRequestUpdate struct {
	ContractStartDate       *time.Time               `json:"contract_start_date"`
	ContractEndDate         *time.Time               `json:"contract_end_date"`
	CompanyRequestUpdate    *CompanyRequestUpdate    `json:"company"`
	TagRequestUpdateRequest *TagRequestUpdateRequest `json:"tag"`
	SubscriptionName        *string                  `json:"subscription_name"`
	Typ                     *string                  `json:"typ"`
	ContractNumber          *string                  `json:"contract_number"`
	CustomerNumber          *string                  `json:"customer_number"`
	PaymentMethod           *string                  `json:"payment_method"`
	BillingPeriod           *string                  `json:"billing_period"`
	CreatedBy               string                   `json:"created_by"`
	Note                    *string                  `json:"note"`
	CancellationPeriod      *int64                   `json:"cancellation_period"`
	BillingDate             *int64                   `json:"billing_date"`
	Price                   *float64                 `json:"price"`
}

type CompanyRequestUpdate struct {
	CompanyName   *string `json:"company_name"`
	Category      *string `json:"category"`
	ContactDetail *string `json:"contact_detail"`
	Link          *string `json:"link"`
}

type TagRequestUpdateRequest struct {
	TagName  *string `json:"tag_name"`
	TagColor *string `json:"tag_color"`
}
