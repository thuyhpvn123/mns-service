package controller

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	// "io"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/meta-node-blockchain/meta-node-mns/internal/request"
	"github.com/meta-node-blockchain/meta-node-mns/internal/service"
	"github.com/meta-node-blockchain/meta-node-mns/internal/usecase"
	"github.com/meta-node-blockchain/meta-node-mns/pkg/api"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"errors"
)

type Controller interface {
	GetNames(c *gin.Context)
	CheckExpire()
	GetOwnerByName(c *gin.Context)
	VerifyDomain(c *gin.Context)
}
type controller struct {
	nameUsecase usecase.NameUsecase
	servs       service.SendTransactionService
}

func NewController(
	nameUsecase usecase.NameUsecase,
	servs service.SendTransactionService,
) Controller {
	return &controller{
		nameUsecase,
		servs,
	}
}

func (h *controller) GetNames(c *gin.Context) {
	var queryParam request.GetNamesRequest
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		// Handle error
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, err, c)
		return
	}
	validate := validator.New()

	err := validate.Struct(&queryParam)
	if err != nil {
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, err, c)
		return
	}
	var names []string
	names, err = h.nameUsecase.GetNamesByOwnerAdd(queryParam.Owner)
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	result := gin.H{
		"message": "successful request",
		"data":    names,
	}
	api.ResponseWithStatusAndData(http.StatusOK, result, c)
}
func (h *controller) CheckExpire() {
	_, err := h.nameUsecase.CheckExpire()
	if err != nil {
		logger.Error("Err when CheckExpire")
	}
}

func (h *controller) VerifyDomain(c *gin.Context) {
	var queryParam request.VerifyDomainRequest
	if err := c.ShouldBindJSON(&queryParam); err != nil {
		// Handle error
		api.ResponseWithErrorAndMessage(http.StatusBadRequest, err, c)
		return
	}
	domain := queryParam.Domain
	owner := queryParam.Owner
	label := queryParam.Label

	// Construct the URL with the domain name
	baseURL := "http://" + domain + "/verify.txt"
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	// Create the request URL with the domain as a query parameter
	// queryParams := apiURL.Query()
	// queryParams.Set("domain", domain)
	// apiURL.RawQuery = queryParams.Encode()

	// Create the HTTP request
	resp, err := http.Get(apiURL.String())
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// body, _ := io.ReadAll(resp.Body)
		errMessage := fmt.Sprintf("Received non-OK HTTP status: %s.", resp.Status)
		api.ResponseWithError(errors.New(errMessage), c)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	// Convert the body to string
	bodyStr := string(body)
	// Trim any whitespace (including \n) from bodyStr
	bodyStr = strings.TrimSpace(bodyStr)
	// Convert the string to uint64
	code, err := strconv.ParseUint(bodyStr, 10, 64)
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	fmt.Println("code:", code)
	kq, err := h.servs.CheckDomain(domain, code, owner, label)
	if err != nil {
		api.ResponseWithError(err, c)
		return
	}
	if kq == true {
		result, err := h.servs.VerifyDomain(domain, code, owner, label)
		if err != nil {
			api.ResponseWithError(err, c)
			return
		}
		if result == true {
			message := gin.H{
				"message": "successful request",
				"data":    true,
			}
			api.ResponseWithStatusAndData(http.StatusOK, message, c)
		} else {
			message := gin.H{
				"message": "successful request",
				"data":    false,
			}
			api.ResponseWithStatusAndData(http.StatusOK, message, c)
		}
	} else {
		message := gin.H{
			"message": "successful request",
			"data":    false,
		}
		api.ResponseWithStatusAndData(http.StatusOK, message, c)
	}

}

func (h *controller) GetOwnerByName(c *gin.Context) {
	owner, _ := h.nameUsecase.GetOwnerByName(c.Param("domain") + ".mtd")
	c.String(http.StatusOK, hex.EncodeToString(common.HexToAddress(owner).Bytes()))
}
