package adminhandler

import (
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func AdminGetAllSubscriptions(c *fiber.Ctx) error {
	sqlStatement := `
		SELECT 
			s.id,
			s.user_id,
			s.subscription_name,
			COALESCE(s.typ, ''),
			COALESCE(s.contract_number, ''),
			COALESCE(s.customer_number, ''),
			s.contract_start_date,
			s.contract_end_date,
			s.cancellation_period,
			COALESCE(s.payment_method, ''),
			s.billing_date,
			COALESCE(s.billing_period, ''),
			s.price,
			COALESCE(s.note, ''),
			s.created_by,
			s.updated_by,
			s.deleted_by,
			s.created_at,
			s.updated_at,
			s.deleted_at,
			c.id,
			c.company_name,
			COALESCE(c.category, ''),
			COALESCE(c.contact_detail, ''),
			COALESCE(c.link, ''),
			t.id,
			COALESCE(t.tag_name, ''),
			COALESCE(t.color, '')
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN tags t ON s.tag_id = t.id
		WHERE s.deleted_at IS NULL
	`

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement)
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
			&s.UserID,
			&s.SubscriptionName,
			&s.Typ,
			&s.ContractNumber,
			&s.CustomerNumber,
			&s.ContractStartDate,
			&s.ContractEndDate,
			&s.CancellationPeriod,
			&s.PaymentMethod,
			&s.BillingDate,
			&s.BillingPeriod,
			&s.Price,
			&s.Note,
			&s.CreatedBy,
			&s.UpdatedBy,
			&s.DeletedBy,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
			&s.Company.ID,
			&s.Company.CompanyName,
			&s.Company.Category,
			&s.Company.ContactDetail,
			&s.Company.Link,
			&s.Tag.ID,
			&s.Tag.TagName,
			&s.Tag.TagColor,
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

func AdminGetAllSubscriptionByUserId(c *fiber.Ctx) error {

	id := c.Params("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var exists bool
	err = database.InitiateDataBase().QueryRow(c.Context(), `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
    `, userID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	sqlStatement := `
		SELECT 
			s.id,
			s.user_id,
			s.subscription_name,
			COALESCE(s.typ, ''),
			COALESCE(s.contract_number, ''),
			COALESCE(s.customer_number, ''),
			s.contract_start_date,
			s.contract_end_date,
			s.cancellation_period,
			COALESCE(s.payment_method, ''),
			s.billing_date,
			COALESCE(s.billing_period, ''),
			s.price,
			COALESCE(s.note, ''),
			s.created_by,
			s.updated_by,
			s.deleted_by,
			s.created_at,
			s.updated_at,
			s.deleted_at,
			c.id,
			c.company_name,
			COALESCE(c.category, ''),
			COALESCE(c.contact_detail, ''),
			COALESCE(c.link, ''),
			t.id,
			COALESCE(t.tag_name, ''),
			COALESCE(t.color, '')
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN tags t ON s.tag_id = t.id
		WHERE s.user_id = $1
		AND s.deleted_at IS NULL
	`

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement, userID)
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
			&s.UserID,
			&s.SubscriptionName,
			&s.Typ,
			&s.ContractNumber,
			&s.CustomerNumber,
			&s.ContractStartDate,
			&s.ContractEndDate,
			&s.CancellationPeriod,
			&s.PaymentMethod,
			&s.BillingDate,
			&s.BillingPeriod,
			&s.Price,
			&s.Note,
			&s.CreatedBy,
			&s.UpdatedBy,
			&s.DeletedBy,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
			&s.Company.ID,
			&s.Company.CompanyName,
			&s.Company.Category,
			&s.Company.ContactDetail,
			&s.Company.Link,
			&s.Tag.ID,
			&s.Tag.TagName,
			&s.Tag.TagColor,
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
