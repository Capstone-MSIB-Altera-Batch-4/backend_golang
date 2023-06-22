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
		response := res.Response(http.StatusBadRequest, "error", "required parameter page", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	page, err := strconv.Atoi(temp)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "page must be integer", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	offset = (page - 1) * limit

	if err := config.Db.Table("memberships").Offset(offset).Limit(limit).Find(&memberships).Error; err != nil {
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := config.Db.Table("memberships").Model(&model.Membership{}).Count(&total).Error; err != nil {
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
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
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Periksa jika salah satu data kosong
	if request.Name == "" || request.Email == "" || request.BirthDay == "" {
		response := res.Response(http.StatusBadRequest, "error", "Missing required data", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	birthDay, err := time.Parse("2006-01-02", request.BirthDay)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "Invalid BirthDay format", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	memberCode := fmt.Sprintf("%s-%d", gen.RandomStrGen(), gen.RandomIntGen())
	membership := model.Membership{
		MemberCode: memberCode,
		Name:       request.Name,
		Email:      request.Email,
		Phone:      int64(request.Phone),
		BirthDay:   birthDay,
		Level:      "Bronze",
		Point:      0,
		CreatedAt:  time.Now(),
	}

	if err := config.Db.Table("memberships").Create(&membership).Error; err != nil {
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
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
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	if request.TotalTransaction < 0 {
		response := res.Response(http.StatusBadRequest, "error", "total transaction cannot be smaller than 0", nil)
		return c.JSON(http.StatusBadRequest, response)
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
	if err := config.Db.Table("memberships").First(&membership, request.ID).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
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

	if err := config.Db.Table("memberships").Updates(&membership).Error; err != nil {
		response := res.Response(201, "error", "Membership not found", membership)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := res.Response(201, "Success", "Membership edited", membership)
	return c.JSON(http.StatusOK, response)
}

func EditMembership(c echo.Context) error {
	request := dto.EditMembershipRequest{}
	if err := c.Bind(&request); err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	birthDay, err := time.Parse("2006-01-02", request.BirthDay)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "Invalid BirthDay format", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	membership := model.Membership{
		ID:       intID,
		Name:     request.Name,
		Email:    request.Email,
		Phone:    int64(request.Phone),
		BirthDay: birthDay,
	}

	if err := config.Db.Table("memberships").Updates(&membership).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if err := config.Db.Table("memberships").First(&membership, intID).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
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

	if err := config.Db.Table("memberships").Where("id = ?", intID).Delete(&model.Membership{}).Error; err != nil {
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := res.Response(200, "Success", "Membership deleted", nil)
	return c.JSON(http.StatusOK, response)
}

func SearchMembership(c echo.Context) error {
	var (
		page        int
		limit       = 10
		offset      int
		total       int64
		memberships []*model.Membership
	)

	temp := c.QueryParam("page")

	if temp == "" {
		response := res.Response(http.StatusBadRequest, "error", "required parameter page", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	page, err := strconv.Atoi(temp)
	if err != nil {
		response := res.Response(http.StatusBadRequest, "error", "page must be integer", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	offset = (page - 1) * limit

	keyword := c.QueryParam("member_code")

	if keyword == "" {
		return c.JSON(http.StatusBadRequest, "required parameter `member_code`")
	}

	if err := config.Db.Where("member_code LIKE ?", "%"+keyword+"%").Offset(offset).Limit(limit).Find(&memberships).Error; err != nil {
		response := res.Response(http.StatusInternalServerError, "error", err.Error(), nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	if err := config.Db.Model(&model.Membership{}).Where("member_code LIKE ?", "%"+keyword+"%").Count(&total).Error; err != nil {
		response := res.Response(http.StatusBadRequest, "error", err.Error(), nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	pages := res.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: int(total),
	}
	response := res.Responsedata(http.StatusOK, "success", "successfully retrieved data", memberships, pages)

	return c.JSON(http.StatusOK, response)
}
