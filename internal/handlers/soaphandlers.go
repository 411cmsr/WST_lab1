package handlers

import (

	"WST_lab1_server_new1/internal/database/postgres"

	"WST_lab1_server_new1/internal/models"
	"bytes"

	"encoding/xml"

	"fmt"
	"io"
	"net/http"


	"github.com/gin-gonic/gin"

)

/*
Структура обработчика для разделения логики обработки запросов от доступа к данным
*/
type StorageHandler struct {
	Storage *postgres.Storage
}



// Обработчик SOAP запросов
func (sh *StorageHandler) SOAPHandler(c *gin.Context) {

	var envelope models.Envelope

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading request body")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := xml.Unmarshal(body, &envelope); err != nil {
		fmt.Println("Error decoding XML:", err)
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	fmt.Printf("Decoded Envelope: %+v\n", envelope)

	switch {
	case envelope.Body.SearchPerson != nil:
		sh.searchPersonHandler(c, envelope.Body.SearchPerson)
	default:
		fmt.Println("Unsupported action")
		c.String(http.StatusBadRequest, "Unsupported action")
		return
	}
}


// Метод поиска записей по запросу
func (h *StorageHandler) searchPersonHandler(c *gin.Context, request *models.SearchPersonRequest) {

	persons, err := h.Storage.PersonRepository.SearchPerson(request.Query)
	if err != nil {
		return
	}

	if len(persons) == 0 {
		fmt.Println("No persons found.")
		return
	} else {
		fmt.Printf("Found persons: %+v\n", persons)
	}

	// Формируем результат в формате SOAP
	response := models.SearchPersonResponse{
		Persons: persons,
	}
	fmt.Printf("Response: %+v\n", response)
	c.XML(http.StatusOK, response)
}
