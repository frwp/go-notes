package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/RianWardanaPutra/notes-v1/httputil"
	"github.com/RianWardanaPutra/notes-v1/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (c *Controller) ListUsers(ctx *gin.Context) {
	var users []model.User
	// url?find=''
	var q = ctx.Query("find")
	if q == "" {
		c.db.Order("name").Find(&users)
		ctx.JSON(http.StatusOK, users)
		return
	}

	c.db.Where("name LIKE ?", "%"+q+"%").Find(&users)
	ctx.JSON(http.StatusOK, users)
}

func (c *Controller) AddUser(ctx *gin.Context) {
	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	var newUser model.User
	json.Unmarshal(reqBody, &newUser)

	if name := newUser.Name; name == "" {
		er := errors.New("name cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	if email := newUser.Email; email == "" {
		er := errors.New("email cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	var user model.User

	c.db.First(&user, model.User{
		Email: newUser.Email,
	})

	if user.ID != 0 {
		er := errors.New("records with this email already exists")
		httputil.NewError(ctx, 409, er)
		return
	}

	c.db.Create(&newUser)

	// model.SetGormTags(&newUser)
	ctx.JSON(http.StatusCreated, newUser)
}

func (c *Controller) FindUserById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, 400, er)
		return
	}

	var user model.User
	c.db.First(&user, id)
	if user.ID == 0 {
		httputil.NewError(ctx, 404, gorm.ErrRecordNotFound)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) UpdateUserById(ctx *gin.Context) {
	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, http.StatusBadRequest, er)
		return
	}

	var newUser model.User
	json.Unmarshal(reqBody, &newUser)

	var oldUser model.User
	c.db.First(&oldUser, id)
	if oldUser.ID == 0 {
		httputil.NewError(ctx, 404, gorm.ErrRecordNotFound)
		return
	}

	c.db.Model(&oldUser).Updates(&newUser)
	ctx.JSON(http.StatusOK, oldUser)
}

func (c *Controller) DeleteUserById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, http.StatusBadRequest, er)
		return
	}

	var user model.User
	c.db.First(&user, id)
	if user.ID == 0 {
		httputil.NewError(ctx, 404, gorm.ErrRecordNotFound)
		return
	}

	c.db.Delete(&user)
	ctx.JSON(http.StatusOK, user)
}
