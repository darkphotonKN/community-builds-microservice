package payment

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	Client PaymentClient
}

func NewHandler(client PaymentClient) *PaymentHandler {
	return &PaymentHandler{
		Client: client,
	}
}

/*
* @Summary Create Subscription
* @Description Create a new subscription for a user with a specified price ID
 */
func (h *PaymentHandler) CreateSubscriptionHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	email, _ := c.Get("email")

	fmt.Printf("\nCreating subscription with userId: %s and email: %s\n\n", userIdStr, email)

	var createSubScriptionReq CreateSubScriptionRequest

	if err := c.ShouldBindJSON(&createSubScriptionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	grpcPayload := &pb.CreateSubscriptionRequest{
		Name:        createSubScriptionReq.Name,
		Description: createSubScriptionReq.Description,
		Price:       createSubScriptionReq.Price,
		MemberId:    userIdStr.(string),
		Email:       email.(string),
	}

	response, err := h.Client.CreateSubscription(c.Request.Context(), grpcPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Create subscription", "result": response})
}

/*
* @Summary Create Customer
* @Description Create a new customer in the payment system
 */
func (h *PaymentHandler) CreateCustomerHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	email, _ := c.Get("email")

	fmt.Printf("\nCreating customer with userId: %s and email: %s\n\n", userIdStr, email)

	grpcPayload := &pb.CreateCustomerRequest{
		MemberId: userIdStr.(string),
		Email:    email.(string),
	}

	response, err := h.Client.CreateCustomer(c.Request.Context(), grpcPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Create customer", "result": response})
}
