package vproductgrp

import (
	"errors"
	"net/http"

	"github.com/testvergecloud/testApi/business/core/views/vproduct"
	"github.com/testvergecloud/testApi/business/web/v1/order"
	"github.com/testvergecloud/testApi/foundation/validate"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByProductID = "product_id"
		orderByUserID    = "user_id"
		orderByName      = "name"
		orderByCost      = "cost"
		orderByQuantity  = "quantity"
		orderByUserName  = "user_name"
	)

	orderByFields := map[string]string{
		orderByProductID: vproduct.OrderByProductID,
		orderByUserID:    vproduct.OrderByUserID,
		orderByName:      vproduct.OrderByName,
		orderByCost:      vproduct.OrderByCost,
		orderByQuantity:  vproduct.OrderByQuantity,
		orderByUserName:  vproduct.OrderByUserName,
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
