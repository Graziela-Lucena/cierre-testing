package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

type arrange struct {
	db func() map[int]internal.Product
}
type request struct {
	query internal.ProductQuery
}
type response struct {
	p      map[int]internal.Product
	err    error
	errMsg string
}
type testCase struct {
	name     string
	arrange  arrange
	request  request
	response response
}

func TestProductsMap_SearchProducts(t *testing.T) {
	testCases := []testCase{
		{
			name: "Get All Products - Sucesso",
			arrange: arrange{
				db: func() map[int]internal.Product {
					return map[int]internal.Product{
						1: {
							Id: 1,
							ProductAttributes: internal.ProductAttributes{
								Description: "Leite",
								Price:       4.5,
								SellerId:    1,
							},
						},
						2: {
							Id: 2,
							ProductAttributes: internal.ProductAttributes{
								Description: "Cereal",
								Price:       6.0,
								SellerId:    2,
							},
						},
					}
				},
			},
			request: request{
				query: internal.ProductQuery{},
			},
			response: response{
				p: map[int]internal.Product{
					1: {
						Id: 1,
						ProductAttributes: internal.ProductAttributes{
							Description: "Leite",
							Price:       4.5,
							SellerId:    1,
						},
					},
					2: {
						Id: 2,
						ProductAttributes: internal.ProductAttributes{
							Description: "Cereal",
							Price:       6.0,
							SellerId:    2,
						},
					},
				},
				err:    nil,
				errMsg: "",
			},
		},
		{
			name: "Get All Products Lista Vazia - Sucesso",
			arrange: arrange{
				db: func() map[int]internal.Product {
					return make(map[int]internal.Product)
				},
			},
			request: request{
				query: internal.ProductQuery{},
			},
			response: response{
				p:      make(map[int]internal.Product),
				err:    nil,
				errMsg: "",
			},
		},
		{
			name: "Get by Id - Sucesso",
			arrange: arrange{
				db: func() map[int]internal.Product {
					return map[int]internal.Product{
						1: {
							Id: 1,
							ProductAttributes: internal.ProductAttributes{
								Description: "Leite",
								Price:       4.5,
								SellerId:    1,
							},
						},
						2: {
							Id: 2,
							ProductAttributes: internal.ProductAttributes{
								Description: "Cereal",
								Price:       6.0,
								SellerId:    2,
							},
						},
					}
				},
			},
			request: request{
				query: internal.ProductQuery{
					Id: 1,
				},
			},
			response: response{
				p: map[int]internal.Product{
					1: {
						Id: 1,
						ProductAttributes: internal.ProductAttributes{
							Description: "Leite",
							Price:       4.5,
							SellerId:    1,
						},
					},
				},
				err:    nil,
				errMsg: "",
			},
		},
		{
			name: "Get by Id Lista Vazia - Sucesso",
			arrange: arrange{
				db: func() map[int]internal.Product {
					return map[int]internal.Product{
						1: {
							Id: 1,
							ProductAttributes: internal.ProductAttributes{
								Description: "Leite",
								Price:       4.5,
								SellerId:    1,
							},
						},
						2: {
							Id: 2,
							ProductAttributes: internal.ProductAttributes{
								Description: "Cereal",
								Price:       6.0,
								SellerId:    2,
							},
						},
					}
				},
			},
			request: request{
				query: internal.ProductQuery{
					Id: 3,
				},
			},
			response: response{
				p:      map[int]internal.Product{},
				err:    nil,
				errMsg: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.arrange.db()
			rp := repository.NewProductsMap(db)

			p, err := rp.SearchProducts(tc.request.query)

			require.Equal(t, tc.response.p, p)
			require.ErrorIs(t, err, tc.response.err)
			if tc.response.err != nil {
				require.EqualError(t, err, tc.response.errMsg)
			}
		})
	}
}
