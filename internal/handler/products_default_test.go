package handler_test

import (
	"app/internal"
	"app/internal/handler"
	"app/internal/repository"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type arrange struct {
	repoMock func() *repository.ProductsMock
}
type request struct {
	request  func() *http.Request
	response *httptest.ResponseRecorder
}
type response struct {
	code    int
	body    string
	headers http.Header
}
type testCase struct {
	name     string
	arrange  arrange
	request  request
	response response
}

func TestProduct_GetProducts(t *testing.T) {

	testCases := []testCase{
		{
			name: "Get bem sucedido",
			arrange: arrange{
				repoMock: func() *repository.ProductsMock {
					rpMock := repository.NewProductsMock()
					rpMock.FuncSearchProducts = func(query internal.ProductQuery) (map[int]internal.Product, error) {
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
									SellerId:    1,
								},
							},
						}, nil
					}
					return rpMock
				},
			},
			request: request{
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					return r
				},
				response: httptest.NewRecorder(),
			},
			response: response{
				code: http.StatusOK,
				body: `
					{"message": "success", "data": {
						"1": {"id": 1, "description": "Leite", "price": 4.5, "seller_id": 1},
						"2": {"id": 2, "description": "Cereal", "price": 6.0, "seller_id": 1}
					}}
				`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
		{
			name: "Get bem sucedido - sem produtos",
			arrange: arrange{
				repoMock: func() *repository.ProductsMock {
					rpMock := repository.NewProductsMock()
					rpMock.FuncSearchProducts = func(query internal.ProductQuery) (map[int]internal.Product, error) {
						return map[int]internal.Product{}, nil
					}
					return rpMock
				},
			},
			request: request{
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					return r
				},
				response: httptest.NewRecorder(),
			},
			response: response{
				code: http.StatusOK,
				body: `
					{"message": "success", "data": {}}
				`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
		{
			name: "Get sem sucesso",
			arrange: arrange{
				repoMock: func() *repository.ProductsMock {
					rpMock := repository.NewProductsMock()
					rpMock.FuncSearchProducts = func(query internal.ProductQuery) (map[int]internal.Product, error) {
						return nil, errors.New("error: repository")
					}
					return rpMock
				},
			},
			request: request{
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					return r
				},
				response: httptest.NewRecorder(),
			},
			response: response{
				code: http.StatusInternalServerError,
				body: fmt.Sprintf(
					`{"status": "%s", "message": "%s"}`,
					http.StatusText(http.StatusInternalServerError),
					"internal error",
				),
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
		{
			name: "Get by id bem sucedido",
			arrange: arrange{
				repoMock: func() *repository.ProductsMock {
					rpMock := repository.NewProductsMock()
					rpMock.FuncSearchProducts = func(query internal.ProductQuery) (map[int]internal.Product, error) {
						return map[int]internal.Product{
							1: {
								Id: 1,
								ProductAttributes: internal.ProductAttributes{
									Description: "Leite",
									Price:       4.5,
									SellerId:    1,
								},
							},
						}, nil
					}
					return rpMock
				},
			},
			request: request{
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					q := r.URL.Query()
					q.Add("id", "1")
					r.URL.RawQuery = q.Encode()
					return r
				},
				response: httptest.NewRecorder(),
			},
			response: response{
				code: http.StatusOK,
				body: `
					{"message": "success", "data": {
						"1": {"id": 1, "description": "Leite", "price": 4.5, "seller_id": 1}
					}}
				`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
		{
			name: "Get by id bem sucedido - sem produtos",
			arrange: arrange{
				repoMock: func() *repository.ProductsMock {
					rpMock := repository.NewProductsMock()
					rpMock.FuncSearchProducts = func(query internal.ProductQuery) (map[int]internal.Product, error) {
						return map[int]internal.Product{}, nil
					}
					return rpMock
				},
			},
			request: request{
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					q := r.URL.Query()
					q.Add("id", "1")
					r.URL.RawQuery = q.Encode()
					return r
				},
				response: httptest.NewRecorder(),
			},
			response: response{
				code: http.StatusOK,
				body: `
					{"message": "success", "data": {}}
				`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
	}

	for _, testC := range testCases {
		t.Run(testC.name, func(t *testing.T) {
			repoMock := testC.arrange.repoMock()
			hd := handler.NewProductsDefault(repoMock)
			hdFunc := hd.GetProducts()

			hdFunc(testC.request.response, testC.request.request())

			require.Equal(t, testC.response.code, testC.request.response.Code)
			require.Equal(t, testC.response.headers, testC.request.response.Header())
			require.JSONEq(t, testC.response.body, testC.request.response.Body.String())
		})
	}
}
