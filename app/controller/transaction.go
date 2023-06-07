package controller

import (
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
		// 1. create order
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 2. create order items
		var orderItems []model.OrderItems
		var totalAmount int
		for _, item := range request.Items {
			product := model.Product{}
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return err
			}
			subtotal := item.Quantity * product.Price
			orderItem := model.OrderItems{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Subtotal:  subtotal,
				Note:      item.Note,
			}
			if err := tx.FirstOrCreate(&orderItem, model.OrderItems{
				OrderID:   order.ID,
				ProductID: item.ProductID,
			}).Error; err != nil {
				return err
			}
			orderItem.Products = product
			orderItems = append(orderItems, orderItem)
			totalAmount += subtotal
		}
		order.Items = orderItems

		// 3. Create transaction
		service := model.Service{}
		user := c.Get("user").(model.User)
		if err := tx.First(&service).Order("id DESC").Limit(1).Error; err != nil {
			return err
		}
		serviceCharge := float64(service.Service) / 100.0
		transaction := model.Transaction{
			OrderID:    order.ID,
			Status:     "Paid",
			Payment:    request.Payment,
			Amount:     totalAmount + int(float64(totalAmount)*serviceCharge),
			Service:    service.Service,
			MemberCode: request.MemberCode,
			UserID:     user.ID,
		}
		order.Transaction = transaction
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// 4. Kalkulasi dan update point member
		totalAmountForPoints := transaction.Amount
		if totalAmountForPoints > 0 {
			member := model.Membership{}
			if err := tx.Where("member_code = ?", request.MemberCode).First(&member).Error; err == nil {
				points := 0
				if totalAmountForPoints <= 50000 {
					points = 10
				} else if totalAmountForPoints <= 100000 {
					points = 20
				} else if totalAmountForPoints <= 150000 {
					points = 30
				} else if totalAmountForPoints <= 200000 {
					points = 40
				} else {
					// Kelipatan
					points = (totalAmountForPoints / 10000) * 10
				}

				member.Point += points
				if member.Point >= 100 && member.Point <= 1999 {
					member.Level = "bronze"
				} else if member.Point >= 2000 && member.Point <= 4999 {
					member.Level = "silver"
				} else if member.Point >= 5000 {
					member.Level = "gold"
				}

				if err := tx.Save(&member).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", "failed to create order", err.Error())
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Transform response
	transformedOrder := res.TransformOrderResponse(order)
	response := res.Response(http.StatusCreated, "success", "success create order", transformedOrder)
	return c.JSON(http.StatusCreated, response)
}

//func RequestPayment(c echo.Context) error {
//	request := dto.CreateOrderRequest{}
//	if err := c.Bind(&request); err != nil {
//		response := res.Response(http.StatusBadRequest, "error", "failed input data", err.Error())
//		return c.JSON(http.StatusBadRequest, response)
//	}
//
//	if request.OrderOption != "dine in" && request.OrderOption != "take away" {
//		response := res.Response(http.StatusBadRequest, "error", "failed input data", "order option only 'take away' or 'dine in'")
//		return c.JSON(http.StatusBadRequest, response)
//	}
//
//	orderCount := generator.GetOrderCount()
//	today := time.Now().Format("02012006")
//	orderCode := generator.GenerateOrderCode(orderCount, today)
//
//	order := model.Order{
//		OrderCode:   orderCode,
//		Name:        request.Name,
//		OrderOption: request.OrderOption,
//		NumberTable: request.TableNumber,
//	}
//
//	err := config.Db.Transaction(func(tx *gorm.DB) error {
//		// 1. create order
//		if err := tx.Create(&order).Error; err != nil {
//			return err
//		}
//
//		// 2. create order items
//		var orderItems []model.OrderItems
//		var totalAmount float64
//		for _, item := range request.Items {
//			product := model.Product{}
//			if err := tx.First(&product, item.ProductID).Error; err != nil {
//				return err
//			}
//			subtotal := item.Quantity * product.Price
//			orderItem := model.OrderItems{
//				OrderID:   order.ID,
//				ProductID: item.ProductID,
//				Quantity:  item.Quantity,
//				Subtotal:  subtotal,
//				Note:      item.Note,
//			}
//			if err := tx.FirstOrCreate(&orderItem, model.OrderItems{
//				OrderID:   order.ID,
//				ProductID: item.ProductID,
//			}).Error; err != nil {
//				return err
//			}
//			orderItem.Products = product // Set product pada orderItem
//			orderItems = append(orderItems, orderItem)
//			totalAmount += float64(subtotal) // Ubah tipe data subtotal menjadi float64
//		}
//		order.Items = orderItems
//
//		// 3. Create transaction
//		service := model.Service{}
//		user := c.Get("user").(model.User)
//		if err := tx.First(&service).Order("id DESC").Limit(1).Error; err != nil {
//			return err
//		}
//		serviceCharge := float64(service.Service) / 100.0
//		transaction := model.Transaction{
//			OrderID: order.ID,
//			Status:  "paid",
//			Payment: request.Payment,
//			Amount:  int(totalAmount + (totalAmount * serviceCharge)),
//			Service: service.Service,
//			UserID:  user.ID,
//		}
//		order.Transaction = transaction
//		if err := tx.Create(&transaction).Error; err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		response := res.Response(http.StatusInternalServerError, "error", "failed to create order", err.Error())
//		return c.JSON(http.StatusInternalServerError, response)
//	}
//
//	// Transform response
//	transformedOrder := res.TransformOrderResponse(order)
//	response := res.Response(http.StatusCreated, "success", "success create order", transformedOrder)
//	return c.JSON(http.StatusCreated, response)
//}
