package handler

import (
	"fmt"
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
)

func AddSubscription(c *fiber.Ctx) error {
	var req data.SubscriptionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
	}

	sqlStatement := `
	INSERT INTO public.subscriptions
	(subscription_name, vendor_name, cost, billing_cycle, renewale_date, categorie, status, created_at, updated_at)
	VALUES ($1,$2,ROUND($3::numeric,2),$4,$5,$6,$7,NOW(),NOW())
	RETURNING id, subscription_name, vendor_name, cost, billing_cycle, renewale_date, categorie, status, created_at, updated_at
	`

	var inserted data.SubscriptionResponse
	err := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		req.SubscriptionName,
		req.VendorName,
		req.Cost,
		req.BillingCycle,
		req.RenewaleDate,
		req.Categorie,
		req.Status,
	).Scan(
		&inserted.ID,
		&inserted.SubscriptionName,
		&inserted.VendorName,
		&inserted.Cost,
		&inserted.BillingCycle,
		&inserted.RenewaleDate,
		&inserted.Categorie,
		&inserted.Status,
		&inserted.Created_at,
		&inserted.Updated_at,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(inserted)
}

func GetSubscriptions(c *fiber.Ctx) error {

	sqlstatement := `SELECT id, subscription_name, vendor_name, billing_cycle, renewale_date, categorie, status, cost , created_at, 
	updated_at, deleted_at FROM public.subscriptions WHERE deleted_at IS NULL`

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlstatement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	subs := []data.SubscriptionResponse{}

	for rows.Next() {
		var s data.SubscriptionResponse

		err := rows.Scan(
			&s.ID,
			&s.SubscriptionName,
			&s.VendorName,
			&s.BillingCycle,
			&s.RenewaleDate,
			&s.Categorie,
			&s.Status,
			&s.Cost,
			&s.Created_at,
			&s.Updated_at,
			&s.Deleted_at,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		subs = append(subs, s)
	}

	return c.Status(fiber.StatusOK).JSON(subs)
}

func GetSubscription(c *fiber.Ctx) error {

	id := c.Params("id")

	subscriptionId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlstatement := `SELECT id, subscription_name, vendor_name, billing_cycle, renewale_date, categorie, status, cost , created_at, 
	updated_at, deleted_at FROM public.subscriptions WHERE id=$1 AND deleted_at IS NULL`

	var s data.SubscriptionResponse

	err = database.InitiateDataBase().QueryRow(c.Context(), sqlstatement, subscriptionId).Scan(
		&s.ID,
		&s.SubscriptionName,
		&s.VendorName,
		&s.BillingCycle,
		&s.RenewaleDate,
		&s.Categorie,
		&s.Status,
		&s.Cost,
		&s.Created_at,
		&s.Updated_at,
		&s.Deleted_at,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(s)

}

func UpdateSubscription(c *fiber.Ctx) error {

	id := c.Params("id")

	subscriptionId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var req data.SubscriptionRequestUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	sqlstatement := `
	UPDATE public.subscriptions
	SET subscription_name = COALESCE($1, subscription_name),
	    vendor_name       = COALESCE($2, vendor_name),
	    cost              = COALESCE($3, cost),
	    billing_cycle     = COALESCE($4, billing_cycle),
	    renewale_date     = COALESCE($5, renewale_date),
	    categorie         = COALESCE($6, categorie),
	    status            = COALESCE($7, status),
	    updated_at        = NOW()
	WHERE id = $8
	RETURNING id, subscription_name, vendor_name, cost, billing_cycle, renewale_date, categorie, status, created_at, updated_at, deleted_at
	`

	var updated data.SubscriptionResponse
	err = database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlstatement,
		req.SubscriptionName,
		req.VendorName,
		req.Cost,
		req.BillingCycle,
		req.RenewaleDate,
		req.Categorie,
		req.Status,
		subscriptionId,
	).Scan(
		&updated.ID,
		&updated.SubscriptionName,
		&updated.VendorName,
		&updated.Cost,
		&updated.BillingCycle,
		&updated.RenewaleDate,
		&updated.Categorie,
		&updated.Status,
		&updated.Created_at,
		&updated.Updated_at,
		&updated.Deleted_at,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "subscription updated",
		"subscription": updated,
	})
}

func DeleteSubscription(c *fiber.Ctx) error {
	id := c.Params("id")

	subscriptionId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlstatement := `UPDATE public.subscriptions SET deleted_at = NOW() WHERE id = $1`

	subs, err := database.InitiateDataBase().Exec(c.Context(), sqlstatement, subscriptionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if subs.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "subscription deleted successfully",
	})

}
