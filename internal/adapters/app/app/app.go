package app

import (
	"fmt"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/response"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/service"
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PersonRouter struct {
	service service.PersonService
}

type IdResponse struct {
	Id int `json:"id,string,omitempty"`
}

func NewPersonRouter(s service.PersonService) *PersonRouter {
	return &PersonRouter{
		service: s,
	}
}

func (r *PersonRouter) GetPerson(c *gin.Context) {
	op := "app.GetPerson"
	idStr := c.Param("id")
	if idStr == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf(" no id in params "))
		log.Print(op, " :no id params")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : invalid id", err))
		log.Print(op, " :invalid id")
		return
	}

	person, err := r.service.GetPerson(c.Request.Context(), id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to get person in service", err))
		log.Print(op, " :failed to get person in service")
		return
	}
	c.JSON(http.StatusOK, person)
}

func (r *PersonRouter) UpdatePerson(c *gin.Context) {
	op := "app.UpdatePerson"
	request := &model.Person{}
	if err := c.ShouldBind(&request); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to update person", err))
		log.Print(op, " :failed to update person")
		return
	}

	err := r.service.UpdatePerson(c.Request.Context(), request)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s: failed to update person in service", err))
		log.Print(op, " :failed to update person in service")
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{"ok"})
}

func (r *PersonRouter) DeletePerson(c *gin.Context) {
	op := "app.DeletePerson"
	request := &IdResponse{}
	if err := c.ShouldBind(&request); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to delete person", err))
		log.Print(op, " :failed to delete person")
		return
	}

	err := r.service.DeletePerson(c.Request.Context(), request.Id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to delete person in service", err))
		log.Print(op, " :failed to delete person")
		return
	}
	c.JSON(http.StatusOK, response.StatusResponse{"ok"})
}

func (r *PersonRouter) GetPersons(c *gin.Context) {
	op := "app.GetPersons"
	data := &model.Person{}
	data.Name = c.Query("name")
	data.Surname = c.Query("surname")
	data.Patronymic = c.Query("patronymic")
	ageStr := c.Query("age")
	if ageStr != "" {
		ageInt, err := strconv.Atoi(ageStr)
		if err != nil {
			response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : invalid age", err))
			log.Print(op, " :invalid age")
			return
		}
		if ageInt < 0 {
			response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : invalid age", err))
			log.Print(op, " :invalid age")
			return
		}
		data.Age = ageInt
	}
	data.Gender = c.Query("gender")
	data.Nationality = c.Query("nationality")

	persons, err := r.service.GetPersons(c.Request.Context(), data)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to get persons", err))
		log.Print(op, " :failed to get persons")
		return
	}
	c.JSON(http.StatusOK, persons)
}

func (r *PersonRouter) AddPerson(c *gin.Context) {
	op := "app.AddPerson"
	var input dto.AddPersonDTO
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s : failed to add person", err))
		log.Print(op, " :failed to add persons")
		return
	}

	id, err := r.service.AddPerson(c.Request.Context(), &input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s: failed to add person to storage", err))
		log.Print(op, " :failed to add persons to storage")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
