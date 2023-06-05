package admin

import (
	"fmt"
	"net/http"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"point-of-sale/utils/dto"
	"point-of-sale/utils/gen"
	"point-of-sale/utils/res"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetMembership(c echo.Context) error {
	var (
		page        int
		limit       = 10
		offset      int
		total       int64
		memberships []*model.Membership
	)

	temp := c.QueryParam("page")

	if temp == "" {
		return c.JSON(http.StatusBadRequest, "required paramter `page`")
	}

	page, err := strconv.Atoi(temp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "page must be integer")
	}

	offset = (page - 1) * limit

	if err := config.Db.Table("membership").Offset(offset).Find(&memberships).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.Table("membership").Model(&model.Membership{}).Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(total),
	}
	response := res.Responsedata(200, "Success", "Membership list", memberships, pages)

	return c.JSON(http.StatusOK, response)
}

func AddMembership(c echo.Context) error {
	request := dto.AddMembershipRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	birthDay, err := time.Parse("2006-01-02", request.BirthDay)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid BirthDay format")
	}

	memberCode := fmt.Sprintf("%s-%d", gen.RandomStrGen(), gen.RandomIntGen())
	membership := model.Membership{
		MemberCode: memberCode,
		Name:       request.Name,
		Email:      request.Email,
		Phone:      request.Phone,
		BirthDay:   birthDay,
		Level:      "Bronze",
		Point:      0,
		CreatedAt:  time.Now(),
	}

	if err := config.Db.Table("membership").Create(&membership).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(201, "Success", "Membership created", membership)

	return c.JSON(http.StatusOK, response)
}

func AddPoint(c echo.Context) error {
	var (
		membership model.Membership
		point      = 0
	)

	request := dto.AddPointRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if request.TotalTransaction < 0 {
		return c.JSON(http.StatusBadRequest, "total transaction cannot be smaller than 0")
	} else if request.TotalTransaction < 50001 {
		point += 10
	} else if request.TotalTransaction < 100001 {
		point += 20
	} else if request.TotalTransaction < 150001 {
		point += 30
	} else {
		point += 40
	}

	fmt.Println(request.ID)
	if err := config.Db.Table("membership").First(&membership, request.ID).Error; err != nil {
		return err
	}

	membership.Point += point

	if membership.Point < 0 {
		return c.JSON(http.StatusBadRequest, "point cannot be smaller than 0")
	} else if membership.Point < 1000 {
		membership.Level = "Bronze"
	} else if membership.Point < 2000 {
		membership.Level = "Silver"
	} else {
		membership.Level = "Gold"
	}

	if err := config.Db.Table("membership").Updates(&membership).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(201, "Success", "Membership edited", membership)

	return c.JSON(http.StatusOK, response)
}

func EditMembership(c echo.Context) error {
	request := dto.EditMembershipRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	birthDay, err := time.Parse("2006-01-02", request.BirthDay)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid BirthDay format")
	}

	membership := model.Membership{
		ID:       intID,
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		BirthDay: birthDay,
	}

	if err := config.Db.Table("membership").Updates(&membership).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.Table("membership").First(&membership, intID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	response := res.Response(200, "Success", "Membership edited", membership)

	return c.JSON(http.StatusOK, response)
}

func DeleteMembership(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := config.Db.Table("membership").Where("id = ?", intID).Delete(&model.Membership{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := res.Response(200, "Success", "Membership deleted", "")

	return c.JSON(http.StatusOK, response)
}
