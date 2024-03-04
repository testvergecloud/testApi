package payment

import (
	"go-starter/config"
	"go-starter/internal/validator"
	"go-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
)

var (
	FreePlan         string = "freePlan"
	ProfessionalPlan string = "professionalPlan"
	PremiumPlan      string = "premiumPlan"
)

type CreateCheckoutSessionDto struct {
	PriceId    string `json:"priceId" binding:"required"`
	CustomerId string `json:"customerId" binding:"required"`
}

type PaymentController interface {
	FindAllSubscriptionPlans(*gin.Context)
	CreateCheckoutSession(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type paymentController struct {
	config *config.Config
}

type SubscriptionPlan struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ProductId string  `json:"productId"`
	PriceId   string  `json:"priceId"`
}

func NewPaymentController(c *config.Config) PaymentController {
	stripe.Key = c.StripeSecretKey
	return &paymentController{config: c}
}

func (sc *paymentController) FindAllSubscriptionPlans(ctx *gin.Context) {
	var plans []SubscriptionPlan

	params := &stripe.PriceListParams{
		LookupKeys: []*string{&FreePlan, &ProfessionalPlan, &PremiumPlan},
	}

	i := price.List(params)

	for i.Next() {
		switch i.Price().LookupKey {
		case FreePlan:
			plans = append(plans,
				SubscriptionPlan{
					Name:      "Free Plan",
					Price:     i.Price().UnitAmountDecimal / 100,
					ProductId: i.Price().Product.ID,
					PriceId:   i.Price().ID})
		case ProfessionalPlan:
			plans = append(plans,
				SubscriptionPlan{
					Name:      "Professional Plan",
					Price:     i.Price().UnitAmountDecimal / 100,
					ProductId: i.Price().Product.ID,
					PriceId:   i.Price().ID})
		case PremiumPlan:
			plans = append(plans,
				SubscriptionPlan{
					Name:      "Premium Plan",
					Price:     i.Price().UnitAmountDecimal / 100,
					ProductId: i.Price().Product.ID,
					PriceId:   i.Price().ID})

		}
	}

	response.Success(http.StatusOK, ctx, plans)
}

func (sc *paymentController) CreateCheckoutSession(ctx *gin.Context) {
	var dto CreateCheckoutSessionDto
	if err := ctx.BindJSON(&dto); err != nil {
		validator.HandleError(ctx, err)
		return
	}

	checkoutParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(dto.PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		Customer:   stripe.String(dto.CustomerId),
		SuccessURL: stripe.String(sc.config.PanelAddress + "/payment/result?success=true&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(sc.config.PanelAddress + "/payment/result?canceled=true"),
	}

	s, err := session.New(checkoutParams)
	if err != nil {
		response.Error(http.StatusInternalServerError, ctx, []response.ApiError{{Code: 999, Message: err.Error()}}, nil)
	}

	// var resp = make(map[string]string)
	// resp["url"] = s.URL

	resp := struct {
		Url string
	}{Url: s.URL}

	response.Success(http.StatusSeeOther, ctx, resp)
}

func (sc *paymentController) RegisterRoutes(rg *gin.RouterGroup) {
	paymentRoutes := rg.Group("/payment")
	paymentRoutes.GET("/plans", sc.FindAllSubscriptionPlans)
	paymentRoutes.POST("/", sc.CreateCheckoutSession)
}
