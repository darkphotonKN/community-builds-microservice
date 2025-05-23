package member

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/models"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/utils/errorutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	Service *MemberService
}

func NewMemberHandler(service *MemberService) *MemberHandler {
	return &MemberHandler{
		Service: service,
	}
}

func (h *MemberHandler) CreateMemberHandler(c *gin.Context) {
	var user models.Member

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with parsing payload as JSON.")})
		return
	}

	err := h.Service.CreateMemberService(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error when attempting to create user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created user."})
}

func (h *MemberHandler) UpdatePasswordMemberHandler(c *gin.Context) {
	var requestData MemberUpdatePasswordRequest
	userId, _ := c.Get("userId")

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with parsing payload as JSON.")})
		return
	}

	err := h.Service.UpdatePasswordMemberService(requestData, userId.(uuid.UUID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error when attempting to update user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created user."})
}

func (h *MemberHandler) UpdateInfoMemberHandler(c *gin.Context) {
	var requestData MemberUpdateInfoRequest
	userId, _ := c.Get("userId")

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with parsing payload as JSON.")})
		return
	}

	err := h.Service.UpdateInfoMemberHandler(requestData, userId.(uuid.UUID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error when attempting to create user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created user."})
}

func (h *MemberHandler) LoginMemberHandler(c *gin.Context) {
	var loginReq MemberLoginRequest

	err := c.ShouldBindJSON(&loginReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when unmarshalling json payload: %s\n", err)})
		return
	}

	userLoginRes, err := h.Service.LoginMemberService(loginReq)

	if errors.Is(err, errorutils.ErrUnauthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": errorutils.ErrUnauthorized.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to login user: %s\n", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully logged in.",
		"result": userLoginRes})
}

func (h *MemberHandler) GetMemberByIdHandler(c *gin.Context) {
	// get id from param
	idParam := c.Param("id")

	// check that its a valid uuid
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		// return to stop flow of function after error response
		return
	}

	user, err := h.Service.GetMemberByIdService(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get user with id %d %s", id, err.Error())})

		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retreived user.",
		"result": user})
}

func (r *MemberRepository) RefreshTokenHandler(c *gin.Context) {
	var refreshTokenReq RefreshTokenRequest

	err := c.ShouldBindJSON(&refreshTokenReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with parsing payload as JSON.")})
		return
	}

}
