package productgrp

import (
	"errors"
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/web/order"
	"github.com/testvergecloud/testApi/foundation/validate"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByProductID = "product_id"
		orderByUserID    = "user_id"
		orderByName      = "name"
		orderByCost      = "cost"
		orderByQuantity  = "quantity"
	)

	orderByFields := map[string]string{
		orderByProductID: product.OrderByProductID,
		orderByName:      product.OrderByName,
		orderByCost:      product.OrderByCost,
		orderByQuantity:  product.OrderByQuantity,
		orderByUserID:    product.OrderByUserID,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByProductID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
