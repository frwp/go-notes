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

func (c *Controller) ListNotes(ctx *gin.Context) {
	var notes []model.Note
	// url?find=''
	var q = ctx.Query("find")
	if q == "" {
		c.db.Order("updated_at").Find(&notes)
		ctx.JSON(http.StatusOK, notes)
		return
	}

	c.db.Where("title LIKE ?", "%"+q+"%").Find(&notes)

	// model.SetGormTags(&notes)
	ctx.JSON(http.StatusOK, notes)
}

func (c *Controller) AddNote(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, 400, er)
		return
	}

	var user model.User
	c.db.First(&user, id)
	if user.ID == 0 {
		httputil.NewError(ctx, 404, errors.New("user not found"))
		return
	}

	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	var newNote model.Note
	json.Unmarshal(reqBody, &newNote)
	newNote.UserId = id

	if title := newNote.Title; title == "" {
		er := errors.New("title cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	if description := newNote.Description; description == "" {
		er := errors.New("description cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	var note model.Note

	c.db.First(&note, model.Note{
		Title: newNote.Title,
	})

	if note.ID != 0 {
		er := errors.New("records with this title already exist")
		httputil.NewError(ctx, 409, er)
		return
	}

	c.db.Create(&newNote)

	ctx.JSON(http.StatusCreated, newNote)
}

func (c *Controller) FindNotesByUserId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, 400, er)
		return
	}

	var user model.User
	c.db.First(&user, id)

	if user.ID == 0 {
		httputil.NewError(ctx, 404, errors.New("user not found"))
		return
	}

	var notes []model.Note

	c.db.Order("created_at desc").Model(&user).Related(&notes)

	ctx.JSON(http.StatusOK, Returns{
		ID:    user.ID,
		Notes: notes,
	})
}

func (c *Controller) FindNoteById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request")
		httputil.NewError(ctx, 400, er)
		return
	}

	var note model.Note
	c.db.First(&note, id)

	if note.ID == 0 {
		httputil.NewError(ctx, 404, gorm.ErrRecordNotFound)
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func (c *Controller) UpdateNoteByIdAndUID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		er := errors.New("bad request: bad user id")
		httputil.NewError(ctx, 400, er)
		return
	}

	var user model.User
	c.db.First(&user, id)

	if user.ID == 0 {
		httputil.NewError(ctx, 404, gorm.ErrRecordNotFound)
		return
	}

	var note model.Note
	noteId, noteErr := strconv.Atoi(ctx.Param("note_id"))
	if noteErr != nil {
		er := errors.New("bad request: bad note id")
		httputil.NewError(ctx, 400, er)
		return
	}

	c.db.First(&note, noteId)

	if note.ID == 0 {
		httputil.NewError(ctx, 404, errors.New("note not found"))
		return
	}

	if note.UserId != id {
		httputil.NewError(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	var newNote model.Note
	json.Unmarshal(reqBody, &newNote)

	if title := newNote.Title; title == "" {
		er := errors.New("title cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	if description := newNote.Description; description == "" {
		er := errors.New("description cannot be empty")
		httputil.NewError(ctx, 400, er)
		return
	}

	var oldNote model.Note

	c.db.First(&oldNote, model.Note{
		Title: newNote.Title,
	})

	if oldNote.UserId != id {
		er := errors.New("records with this title already exist")
		httputil.NewError(ctx, 409, er)
		return

	}

	c.db.Model(&oldNote).Updates(&newNote)

	ctx.JSON(http.StatusOK, oldNote)
}
