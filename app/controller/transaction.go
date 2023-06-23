package controller

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	generator "point-of-sale/utils/gen"
	"point-of-sale/utils/res"
	"time"
)

func RequestPayment(c echo.Context) error {
	request := dto.CreateOrderRequest{}
	if err := c.Bind(&request); err != nil {
		response := res.Response(http.StatusBadRequest, "error", "failed input data", err.Error())
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := generator.ValidateData(&request); err != nil {
		response := res.Response(http.StatusBadRequest, "error", "failed input data", err)
		return c.JSON(http.StatusBadRequest, response)
	}

	if request.OrderOption != "DINE_IN" && request.OrderOption != "TAKE_AWAY" {
		response := res.Response(http.StatusBadRequest, "error", "failed input data", "order option only 'DINE_IN' or 'TAKE_AWAY'")
		return c.JSON(http.StatusBadRequest, response)
	}

	user := c.Get("user").(model.User)

	// Process membership points
	if request.MemberCode != "" {
		member, err := getMemberByCode(request.MemberCode)
		if err != nil {
			response := res.Response(http.StatusBadRequest, "error", "invalid member code", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		order, err := createOrder(&request, user, member)
		if err != nil {
			response := res.Response(http.StatusInternalServerError, "error", "failed to create order", err.Error())
			return c.JSON(http.StatusInternalServerError, response)
		}

		// Transform response
		transformedOrder := res.TransformOrderResponse(*order)
		response := res.Response(http.StatusCreated, "success", "success create order", transformedOrder)
		return c.JSON(http.StatusCreated, response)
	}

	// No membership code provided, proceed without processing membership points
	order, err := createOrder(&request, user, nil)
	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", "failed to create order", err.Error())
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Transform response
	transformedOrder := res.TransformOrderResponse(*order)
	response := res.Response(http.StatusCreated, "success", "success create order", transformedOrder)
	return c.JSON(http.StatusCreated, response)
}

func createOrder(request *dto.CreateOrderRequest, user model.User, member *model.Membership) (*model.Order, error) {
	orderCount := generator.GetOrderCount()
	today := time.Now().Format("02012006")
	orderCode := generator.GenerateOrderCode(orderCount, today)

	order := model.Order{
		OrderCode:   orderCode,
		Name:        request.Name,
		OrderOption: request.OrderOption,
		NumberTable: request.TableNumber,
	}

	err := config.Db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Create order items
		var orderItems []model.OrderItems
		var totalAmount int
		for _, item := range request.Items {
			product := model.Product{}
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("product with ID %d not found", item.ProductID)
			}

			subtotal := item.Quantity * product.Price
			orderItem := model.OrderItems{
				OrderID:     order.ID,
				ProductName: product.Name,
				Quantity:    item.Quantity,
				Subtotal:    subtotal,
				Note:        item.Note,
			}

			if item.Quantity > product.Quantity {
				return fmt.Errorf("quantity exceeded available stock for product with ID %d", item.ProductID)
			}

			product.Quantity -= item.Quantity
			if err := tx.Model(&product).UpdateColumn("quantity", product.Quantity).Error; err != nil {
				return err
			}

			if err := tx.FirstOrCreate(&orderItem, model.OrderItems{
				OrderID:     order.ID,
				ProductName: product.Name,
			}).Error; err != nil {
				return err
			}

			orderItems = append(orderItems, orderItem)
			totalAmount += subtotal
		}

		order.Items = orderItems

		// Create transaction
		service := model.Service{}
		if err := tx.Order("id DESC").Limit(1).First(&service).Error; err != nil {
			return err
		}
		serviceCharge := float64(service.Service) / 100.0
		transaction := model.Transaction{
			OrderID:    order.ID,
			Status:     "PAID",
			Payment:    request.Payment,
			Amount:     int(float64(totalAmount) + (float64(totalAmount) * serviceCharge)),
			Service:    service.Service,
			MemberCode: request.MemberCode,
			UserID:     user.ID,
		}
		order.Transaction = transaction

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process membership points if a member is provided
	if member != nil {
		if err := processMembershipPoints(&order, member); err != nil {
			return nil, err
		}
	}

	return &order, nil
}

func processMembershipPoints(order *model.Order, member *model.Membership) error {
	totalAmountForPoints := order.Transaction.Amount
	if totalAmountForPoints > 0 {
		points := calculateMemberPoints(totalAmountForPoints)
		member.Point += points
		member.Level = getMemberLevel(member.Point)

		if err := config.Db.Save(&member).Error; err != nil {
			return err
		}
	}

	return nil
}

func getMemberByCode(memberCode string) (*model.Membership, error) {
	member := model.Membership{}
	err := config.Db.Where("member_code = ?", memberCode).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("member code not found")
		}
		return nil, err
	}

	return &member, nil
}

func calculateMemberPoints(totalAmount int) int {
	points := 0
	if totalAmount <= 50000 {
		points = 10
	} else if totalAmount <= 100000 {
		points = 20
	} else if totalAmount <= 150000 {
		points = 30
	} else if totalAmount <= 200000 {
		points = 40
	} else {
		// Multiple of 10,000
		points = (totalAmount / 10000) * 10
	}
	return points
}

func getMemberLevel(points int) string {
	if points >= 100 && points <= 1999 {
		return "bronze"
	} else if points >= 2000 && points <= 4999 {
		return "silver"
	} else if points >= 5000 {
		return "gold"
	}
	return ""
}
