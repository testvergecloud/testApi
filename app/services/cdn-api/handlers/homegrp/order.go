package homegrp

import (
	"errors"
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/home"
	"github.com/testvergecloud/testApi/business/web/order"
	"github.com/testvergecloud/testApi/foundation/validate"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByID     = "home_id"
		orderByType   = "type"
		orderByUserID = "user_id"
	)

	orderByFields := map[string]string{
		orderByID:     home.OrderByID,
		orderByType:   home.OrderByType,
		orderByUserID: home.OrderByUserID,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
