package handler

import (
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func AddSubscription(c *fiber.Ctx) error {
	var req data.SubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	sqlstatment := `INSERT INTO subscriptions (subscription_name,vendor_name,cost,billing_cycle,renewale_date,catagorie,status,
	created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := database.InitiateDataBase().Exec(
		c.Context(),
		sqlstatment,
		req.SubscriptionName,
		req.VendorName,
		req.Cost,
		req.BillingCycle,
		req.RenewaleDate,
		req.Categorie,
		req.Status,
		req.Created_at,
		req.Updated_at,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := data.SubscriptionResponse{
		SubscriptionName: req.SubscriptionName,
		VendorName:       req.VendorName,
		Cost:             req.Cost,
		BillingCycle:     req.BillingCycle,
		RenewaleDate:     req.RenewaleDate,
		Categorie:        req.Categorie,
		Status:           req.Status,
	}

	return c.Status(fiber.StatusCreated).JSON(&response)
}
