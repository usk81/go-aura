package router

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoute_Build(t *testing.T) {
	type fields struct {
		Middlewares []func(http.Handler) http.Handler
		Endpoints   []EndpointPattern
		SubRoutes   []SubRoute
	}
	type args struct {
		r *chi.Mux
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Middlewares: nil,
				Endpoints:   nil,
				SubRoutes:   nil,
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "simple",
			fields: fields{
				Middlewares: nil,
				Endpoints: []EndpointPattern{
					{
						Pattern: "/test",
						Endpoints: map[string]Endpoint{
							http.MethodGet: {
								Handler: mockHandlerFunc,
							},
							http.MethodHead: {
								Handler: mockHandlerFunc,
							},
							http.MethodPost: {
								Handler: mockHandlerFunc,
							},
							http.MethodPut: {
								Handler: mockHandlerFunc,
							},
							http.MethodPatch: {
								Handler: mockHandlerFunc,
							},
							http.MethodDelete: {
								Handler: mockHandlerFunc,
							},
							http.MethodConnect: {
								Handler: mockHandlerFunc,
							},
							http.MethodOptions: {
								Handler: mockHandlerFunc,
							},
							http.MethodTrace: {
								Handler: mockHandlerFunc,
							},
						},
					},
				},
				SubRoutes: []SubRoute{},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "rest",
			fields: fields{
				Middlewares: nil,
				Endpoints: []EndpointPattern{
					{
						Pattern: "/",
						Endpoints: map[string]Endpoint{
							http.MethodGet: {
								Handler: mockHandlerFunc,
							},
							http.MethodPost: {
								Handler: mockHandlerFunc,
							},
						},
					},
					{
						Pattern: "/:id",
						Endpoints: map[string]Endpoint{
							http.MethodGet: {
								Handler: mockHandlerFunc,
							},
							http.MethodPut: {
								Handler: mockHandlerFunc,
							},
							http.MethodDelete: {
								Handler: mockHandlerFunc,
							},
						},
					},
				},
				SubRoutes: []SubRoute{},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "subroute",
			fields: fields{
				Middlewares: nil,
				Endpoints:   nil,
				SubRoutes: []SubRoute{
					{
						Pattern: "/user",
						Route: Route{
							Middlewares: nil,
							Endpoints: []EndpointPattern{
								{
									Pattern: "/",
									Endpoints: map[string]Endpoint{
										http.MethodGet: {
											Handler: mockHandlerFunc,
										},
										http.MethodPost: {
											Handler: mockHandlerFunc,
										},
									},
								},
								{
									Pattern: "/:id",
									Endpoints: map[string]Endpoint{
										http.MethodGet: {
											Handler: mockHandlerFunc,
										},
										http.MethodPut: {
											Handler: mockHandlerFunc,
										},
										http.MethodDelete: {
											Handler: mockHandlerFunc,
										},
									},
								},
							},
							SubRoutes: nil,
						},
					},
				},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "subroute",
			fields: fields{
				Middlewares: nil,
				Endpoints:   nil,
				SubRoutes: []SubRoute{
					{
						Pattern: "/product",
						Route: Route{
							Middlewares: nil,
							Endpoints: []EndpointPattern{
								{
									Pattern: "/",
									Endpoints: map[string]Endpoint{
										http.MethodGet: {
											Handler: mockHandlerFunc,
										},
										http.MethodPost: {
											Handler: mockHandlerFunc,
										},
									},
								},
							},
							SubRoutes: []SubRoute{
								{
									Pattern: "/:id",
									Route: Route{
										Middlewares: nil,
										Endpoints: []EndpointPattern{
											{
												Pattern: "/",
												Endpoints: map[string]Endpoint{
													http.MethodGet: {
														Handler: mockHandlerFunc,
													},
													http.MethodPut: {
														Handler: mockHandlerFunc,
													},
													http.MethodDelete: {
														Handler: mockHandlerFunc,
													},
												},
											},
											{
												Pattern: "/test",
												Endpoints: map[string]Endpoint{
													http.MethodGet: {
														Handler: mockHandlerFunc,
													},
												},
											},
										},
										SubRoutes: nil,
									},
								},
							},
						},
					},
				},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "all",
			fields: fields{
				Middlewares: nil,
				Endpoints: []EndpointPattern{
					{
						Pattern: "/all",
						Endpoints: map[string]Endpoint{
							"All": {
								Handler: mockHandlerFunc,
							},
						},
					},
				},
				SubRoutes: []SubRoute{},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: false,
		},
		{
			name: "invalid_method",
			fields: fields{
				Middlewares: nil,
				Endpoints: []EndpointPattern{
					{
						Pattern: "/test",
						Endpoints: map[string]Endpoint{
							"foobar": {
								Handler: mockHandlerFunc,
							},
						},
					},
				},
				SubRoutes: []SubRoute{},
			},
			args: args{
				r: chi.NewRouter(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &Route{
				Middlewares: tt.fields.Middlewares,
				Endpoints:   tt.fields.Endpoints,
				SubRoutes:   tt.fields.SubRoutes,
			}
			if err := rt.Build(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Route.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
			routeDebug(tt.args.r)
		})
	}
}

func routeDebug(r *chi.Mux) (err error) {
	err = chi.Walk(r, func() chi.WalkFunc {
		return func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			route = strings.Replace(route, "/*/", "/", -1)
			fmt.Printf("[%s] %s || handler: %s\n", method, route, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
			return nil
		}
	}())
	return
}

func mockHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
