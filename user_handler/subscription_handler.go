package userhandler

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
		return c.Status(fiber.StatusBadRequest).SendString(
			fmt.Sprintf("Validation error: %s", err),
		)
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	db := database.InitiateDataBase()

	tx, err := db.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to start transaction",
		})
	}
	defer tx.Rollback(c.Context())

	var companyID int
	err = tx.QueryRow(c.Context(), `
		INSERT INTO companies (company_name, category, contact_detail, link)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`,
		req.CompanyRequest.CompanyName,
		req.CompanyRequest.Category,
		req.CompanyRequest.ContactDetail,
		req.CompanyRequest.Link,
	).Scan(&companyID)

	if err != nil {
		return err
	}

	var tagID int
	err = tx.QueryRow(c.Context(), `
		INSERT INTO tags (tag_name, color)
		VALUES ($1, $2)
		RETURNING id;
	`,
		req.TagRequest.TagName,
		req.TagRequest.TagColor,
	).Scan(&tagID)

	if err != nil {
		return err
	}

	createdBy := c.Locals("userEmail").(string)

	startDate, err := utils.ParseDate(req.ContractStartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contract_start_date"})
	}
	endDate, err := utils.ParseDate(req.ContractEndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contract_end_date"})
	}

	var subscriptionID int
	err = tx.QueryRow(c.Context(), `
		INSERT INTO subscriptions (
			subscription_name, typ, contract_number, customer_number,
			contract_start_date, contract_end_date,
			cancellation_period, payment_method,
			billing_date, billing_period, price, note,
			company_id, tag_id,created_by,user_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id
	`,
		req.SubscriptionName,
		req.Typ,
		req.ContractNumber,
		req.CustomerNumber,
		startDate,
		endDate,
		req.CancellationPeriod,
		req.PaymentMethod,
		req.BillingDate,
		req.BillingPeriod,
		req.Price,
		req.Note,
		companyID,
		tagID,
		createdBy,
		userID,
	).Scan(&subscriptionID)

	if err != nil {
		return err
	}

	var res data.SubscriptionResponse

	err = tx.QueryRow(c.Context(), `
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
			s.created_at,
			s.created_by,
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
		WHERE s.id = $1
	`, subscriptionID).Scan(
		&res.ID,
		&res.UserID,
		&res.SubscriptionName,
		&res.Typ,
		&res.ContractNumber,
		&res.CustomerNumber,
		&res.ContractStartDate,
		&res.ContractEndDate,
		&res.CancellationPeriod,
		&res.PaymentMethod,
		&res.BillingDate,
		&res.BillingPeriod,
		&res.Price,
		&res.Note,
		&res.CreatedAt,
		&res.CreatedBy,
		&res.Company.ID,
		&res.Company.CompanyName,
		&res.Company.Category,
		&res.Company.ContactDetail,
		&res.Company.Link,
		&res.Tag.ID,
		&res.Tag.TagName,
		&res.Tag.TagColor,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}

	return c.Status(fiber.StatusCreated).JSON(res)
}


func GetSubscriptionsByUserID(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
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

func GetSubscriptionByUserID(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user ID not found",
		})
	}

	id := c.Params("id")

	subscriptionID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `
		SELECT 
			s.id,
			s.user_id,
			s.subscription_name,
			s.typ,
			s.contract_number,
			s.customer_number,
			s.contract_start_date,
			s.contract_end_date,
			s.cancellation_period,
			s.payment_method,
			s.billing_date,
			s.billing_period,
			s.price,
			s.note,
			s.created_by,
			s.updated_by,
			s.deleted_by,
			s.created_at,
			s.updated_at,
			s.deleted_at,
			c.id,
			c.company_name,
			c.category,
			c.contact_detail,
			c.link,
			t.id,
			t.tag_name,
			t.color
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN tags t ON s.tag_id = t.id
		WHERE s.id = $1
          AND s.user_id = $2                            
          AND s.deleted_at IS NULL
	`

	var s data.SubscriptionResponse

	err = database.InitiateDataBase().QueryRow(c.Context(), sqlStatement, subscriptionID, userID).Scan(
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(s)
}

func UpdateSubscriptionByUserID(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userEmail, ok := c.Locals("userEmail").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	id := c.Params("id")
	subscriptionID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var req data.SubscriptionRequestUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	db := database.InitiateDataBase()
	tx, err := db.Begin(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to start transaction",
		})
	}
	defer tx.Rollback(c.Context())

	var companyID, tagID *int64
	err = tx.QueryRow(c.Context(), `
        SELECT company_id, tag_id 
        FROM subscriptions 
        WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
    `, subscriptionID, userID).Scan(&companyID, &tagID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	if req.CompanyRequestUpdate != nil && companyID != nil {
		_, err = tx.Exec(c.Context(), `
            UPDATE companies
            SET company_name   = COALESCE(NULLIF($1, ''), company_name),
                category       = COALESCE(NULLIF($2, ''), category),
                contact_detail = COALESCE(NULLIF($3, ''), contact_detail),
                link           = COALESCE(NULLIF($4, ''), link),
                updated_at     = NOW()
            WHERE id = $5 AND owner_id = $6
        `,
			req.CompanyRequestUpdate.CompanyName,
			req.CompanyRequestUpdate.Category,
			req.CompanyRequestUpdate.ContactDetail,
			req.CompanyRequestUpdate.Link,
			*companyID,
			userID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if req.TagRequestUpdateRequest != nil && tagID != nil {
		_, err = tx.Exec(c.Context(), `
            UPDATE tags
            SET tag_name = COALESCE(NULLIF($1, ''), tag_name),
                color    = COALESCE(NULLIF($2, ''), color)
            WHERE id = $3 AND user_id = $4
        `,
			req.TagRequestUpdateRequest.TagName,
			req.TagRequestUpdateRequest.TagColor,
			*tagID,
			userID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	var updated data.SubscriptionResponse
	err = tx.QueryRow(c.Context(), `
        UPDATE subscriptions
        SET subscription_name   = COALESCE(NULLIF($1,  ''), subscription_name),
            typ                 = COALESCE(NULLIF($2,  ''), typ),
            contract_number     = COALESCE(NULLIF($3,  ''), contract_number),
            customer_number     = COALESCE(NULLIF($4,  ''), customer_number),
            contract_start_date = COALESCE($5,  contract_start_date),
            contract_end_date   = COALESCE($6,  contract_end_date),
            cancellation_period = COALESCE($7,  cancellation_period),
            payment_method      = COALESCE(NULLIF($8,  ''), payment_method),
            billing_date        = COALESCE($9,  billing_date),
            billing_period      = COALESCE(NULLIF($10, ''), billing_period),
            price               = COALESCE($11, price),
            note                = COALESCE(NULLIF($12, ''), note),
            updated_at          = NOW(),
            updated_by          = $13
        WHERE id = $14
          AND user_id = $15
          AND deleted_at IS NULL
        RETURNING id,
                  subscription_name,
                  typ,
                  contract_number,
                  customer_number,
                  contract_start_date,
                  contract_end_date,
                  cancellation_period,
                  payment_method,
                  billing_date,
                  billing_period,
                  price,
                  note,
                  created_at,
                  updated_at,
                  deleted_at,
                  created_by,
                  updated_by,
                  deleted_by
    `,
		req.SubscriptionName,
		req.Typ,
		req.ContractNumber,
		req.CustomerNumber,
		req.ContractStartDate,
		req.ContractEndDate,
		req.CancellationPeriod,
		req.PaymentMethod,
		req.BillingDate,
		req.BillingPeriod,
		req.Price,
		req.Note,
		userEmail,
		subscriptionID,
		userID,
	).Scan(
		&updated.ID,
		&updated.SubscriptionName,
		&updated.Typ,
		&updated.ContractNumber,
		&updated.CustomerNumber,
		&updated.ContractStartDate,
		&updated.ContractEndDate,
		&updated.CancellationPeriod,
		&updated.PaymentMethod,
		&updated.BillingDate,
		&updated.BillingPeriod,
		&updated.Price,
		&updated.Note,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
		&updated.CreatedBy,
		&updated.UpdatedBy,
		&updated.DeletedBy,
	)
	if err != nil {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error()})
		}
	}

	if companyID != nil {
		err = tx.QueryRow(c.Context(), `
            SELECT id, company_name, category, contact_detail, link
            FROM companies WHERE id = $1
        `, *companyID).Scan(
			&updated.Company.ID,
			&updated.Company.CompanyName,
			&updated.Company.Category,
			&updated.Company.ContactDetail,
			&updated.Company.Link,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if tagID != nil {
		err = tx.QueryRow(c.Context(), `
            SELECT id, tag_name, color
            FROM tags WHERE id = $1
        `, *tagID).Scan(
			&updated.Tag.ID,
			&updated.Tag.TagName,
			&updated.Tag.TagColor,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "subscription updated",
		"subscription": updated,
	})
}

func DeleteSubscriptionByUserID(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userEmail, ok := c.Locals("userEmail").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	subscriptionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `
        UPDATE subscriptions
        SET deleted_at = NOW(),
            deleted_by = $1
        WHERE id = $2
          AND user_id = $3
          AND deleted_at IS NULL
    `

	result, err := database.InitiateDataBase().Exec(
		c.Context(),
		sqlStatement,
		userEmail,      
		subscriptionID, 
		userID,        
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "subscription not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "subscription deleted successfully",
	})
}
