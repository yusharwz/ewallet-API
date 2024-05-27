package json

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	jsonResponse struct {
		Code    string      `json:"responseCode"`
		Message string      `json:"responseMessage,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	jsonResponseWithPaging struct {
		Code    string      `json:"responseCode"`
		Message string      `json:"responseMessage"`
		Data    interface{} `json:"data,omitempty"`
		Paging  interface{} `json:"paging,omitempty"`
	}

	jsonErrorResponse struct {
		Code    string `json:"responseCode"`
		Message string `json:"responseMessage"`
		Error   string `json:"error,omitempty"`
	}

	ValidationField struct {
		FieldName string `json:"field"`
		Message   string `json:"message"`
	}

	jsonBadRequestResponse struct {
		Code             string            `json:"responseCode"`
		Message          string            `json:"responseMessage"`
		ErrorDescription []ValidationField `json:"error_description,omitempty"`
	}

	PagingInfo struct {
		Page      string `json:"page,omitempty"`
		TotalData string `json:"totalData,omitempty"`
	}
)

func NewResponSuccesPaging(c *gin.Context, result interface{}, message, servisCode, responCode string, page, totalData string) {
	c.JSON(http.StatusOK, jsonResponseWithPaging{
		Code:    "200" + servisCode + responCode,
		Message: message,
		Data:    result,
		Paging: PagingInfo{
			Page:      page,
			TotalData: totalData,
		},
	})
}

func NewResponSucces(c *gin.Context, result interface{}, message, servisCode, responCode string) {
	c.JSON(http.StatusOK, jsonResponse{
		Code:    "200" + servisCode + responCode,
		Message: message,
		Data:    result,
	})
}

func NewResponBadRequest(c *gin.Context, validationField []ValidationField, message, serviceCode, errorCode string) {
	c.JSON(http.StatusBadRequest, jsonBadRequestResponse{
		Code:             "400" + serviceCode + errorCode,
		Message:          message,
		ErrorDescription: validationField,
	})
}

func NewResponseError(c *gin.Context, err, serviceCode, errorCode string) {
	log.Error().Msg(err)
	c.JSON(http.StatusInternalServerError, jsonErrorResponse{
		Code:    "500" + serviceCode + errorCode,
		Message: "Internal Server Error",
		Error:   err,
	})
}

func NewResponseForbidden(c *gin.Context, message, serviceCode, errorCode string) {
	c.JSON(http.StatusForbidden, jsonErrorResponse{
		Code:    "403" + serviceCode + errorCode,
		Message: message,
	})
}

func NewResponseUnauthorized(c *gin.Context, message, serviceCode, errorCode string) {
	c.JSON(http.StatusUnauthorized, jsonErrorResponse{
		Code:    "401" + serviceCode + errorCode,
		Message: message,
	})
}
