package tests

import (
	"net/http"

	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/vproductgrp"
	"github.com/testvergecloud/testApi/business/web"
)

func vproductQuery200(sd seedData) []tableData {
	total := len(sd.admins[1].products) + len(sd.users[1].products)

	allProducts := toAppVProducts(sd.admins[1].User, sd.admins[1].products)
	allProducts = append(allProducts, toAppVProducts(sd.users[1].User, sd.users[1].products)...)

	table := []tableData{
		{
			name:       "basic",
			url:        "/v1/vproducts?page=1&rows=10&orderBy=product_id,DESC",
			token:      sd.admins[1].token,
			statusCode: http.StatusOK,
			method:     http.MethodGet,
			resp:       &web.PageDocument[vproductgrp.AppProduct]{},
			expResp: &web.PageDocument[vproductgrp.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       allProducts,
			},
			cmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*web.PageDocument[vproductgrp.AppProduct])
				exp := y.(*web.PageDocument[vproductgrp.AppProduct])

				var found int
				for _, r := range resp.Items {
					for _, e := range exp.Items {
						if e.ID == r.ID && e.UserName == r.UserName {
							found++
							break
						}
					}
				}

				if found != total {
					return "number of expected products didn't match"
				}

				return ""
			},
		},
	}

	return table
}
