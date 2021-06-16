package middlewares

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/internal/tests"
)

// header is a type that helps in testing of request headers.
type header struct {
	key   string
	value string
}

// customLogger is a type helper for testing logging.
type customLogger struct {
	log [][]string
}

func NewCustomLogger() *customLogger {
	log := make([][]string, 0)

	return &customLogger{log}
}

// Log implements logger.Logger interface.
func (cl *customLogger) Log(level, message string) {
	cl.log = append(cl.log, []string{level, message})
}

func TestAuth(t *testing.T) {
	tokenPair, _ := entity.NewTokenPair("1", []string{"test"})
	accessToken := tokenPair.AccessToken.Token

	t.Run("Authenticate middleware test", func(t *testing.T) {
		tt := []struct {
			name                  string
			headers               []header
			message               string
			nextHandlerInvocation bool
			statusCode            int
		}{
			{name: "valid access token", headers: []header{{key: "Authorization", value: "Bearer " + accessToken}}, nextHandlerInvocation: true, statusCode: http.StatusOK},
			{name: "Authorization header is missing", message: errAuthHeaderMissing.Error(), statusCode: http.StatusBadRequest},
			{name: "Authorization header is in wrong format", headers: []header{{key: "Authorization", value: "Bearu"}}, message: errAuthWrongHeaderFormat.Error(), statusCode: http.StatusBadRequest},
			{name: "expired access token", headers: []header{{key: "Authorization", value: "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VyX2lkIjoiMjBjZmU1NTItYTljYi00YmNlLTg1MzctYWU1NjNkYTA3MjFiIiwiUmVmcmVzaF91dWlkIjoiMWM4OGQzZDctMGI2OS00NWY3LTkzMGQtYWViYjg4OTI2MWY5IiwiZXhwIjoxNjA5NjkyNjc0fQ.1V3wgd0cD-u5My1MEb_WoDTPFej5QeMlXG4pTBA4YiZJjcSptc3Pl5MPfEG1k3XF4mx6RAgun3GBowF_DtP0ls"}}, message: "JWT token is not valid: token is expired by", statusCode: http.StatusUnauthorized, nextHandlerInvocation: false},
			{name: "invalid access token", headers: []header{{key: "Authorization", value: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}}, message: "access token is not valid", statusCode: http.StatusUnauthorized, nextHandlerInvocation: false},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})

				handler := Authenticate(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				for _, h := range tc.headers {
					req.Header.Add(h.key, h.value)
				}
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				res := rec.Result()
				defer res.Body.Close()
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Errorf("\t%s\tTest %s:\tCould not read response: %v", tests.Failed, tc.name, err)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to read response", tests.Success, tc.name)

				respBody := strings.TrimSpace(string(b))
				if !strings.Contains(respBody, tc.message) {
					t.Errorf("\t%s\tTest %s:\tWant response body: %v, got response body: %v", tests.Failed, tc.name, tc.message, respBody)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate response body", tests.Success, tc.name)

				if tc.statusCode != res.StatusCode {
					t.Errorf("\t%s\tTest %s:\tWant status code: %d, got status code: %d", tests.Failed, tc.name, tc.statusCode, res.StatusCode)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate status code", tests.Success, tc.name)
			})
		}
	})

	t.Run("Authorize middleware test", func(t *testing.T) {
		tt := []struct {
			name                  string
			role                  string
			context               context.Context
			message               string
			nextHandlerInvocation bool
			statusCode            int
		}{
			{name: "context with  valid JWT access token claims, role inside JWT access token claims matches middleware filter value", role: "user", context: context.WithValue(context.Background(), ClaimsKey, &entity.AccessTokenClaims{User_id: "1", Refresh_uuid: "123", User_roles: []string{"user"}}), nextHandlerInvocation: true, statusCode: http.StatusOK},
			{name: "context with  valid JWT access token claims, role inside JWT access token claims does not matches middleware filter value", role: "admin", message: "you are not authorized for that action", context: context.WithValue(context.Background(), ClaimsKey, &entity.AccessTokenClaims{User_id: "1", Refresh_uuid: "123", User_roles: []string{"user"}}), nextHandlerInvocation: true, statusCode: http.StatusUnauthorized},
			{name: "context with  empty JWT access token claims", message: "you are not authorized for that action", context: context.WithValue(context.Background(), ClaimsKey, &entity.AccessTokenClaims{}), statusCode: http.StatusUnauthorized},
			{name: "empty context", message: errNoClaimsInContext.Error(), context: context.Background(), statusCode: http.StatusBadRequest},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})

				handler := Authorize(tc.role)(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req.WithContext(tc.context))

				res := rec.Result()
				defer res.Body.Close()
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Errorf("\t%s\tTest %s:\tCould not read response: %v", tests.Failed, tc.name, err)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to read response", tests.Success, tc.name)

				respBody := strings.TrimSpace(string(b))
				if !strings.Contains(respBody, tc.message) {
					t.Errorf("\t%s\tTest %s:\tWant response body: %v, got response body: %v", tests.Failed, tc.name, tc.message, respBody)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate response body", tests.Success, tc.name)

				if res.StatusCode != tc.statusCode {
					t.Errorf("\t%s\tTest %s:\tWant status code: %d, got status code: %d", tests.Failed, tc.name, res.StatusCode, tc.statusCode)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate status code", tests.Success, tc.name)
			})
		}
	})

	t.Run("GetJWTClaims function test",
		func(t *testing.T) {
			tt := []struct {
				name    string
				context context.Context
				claims  *entity.AccessTokenClaims
				err     string
			}{
				{name: "valid claims", context: context.WithValue(context.Background(), ClaimsKey, &entity.AccessTokenClaims{User_id: "1234-3124-1234", User_roles: []string{"USER"}, Refresh_uuid: "123-1234-1234"}), claims: &entity.AccessTokenClaims{User_id: "1234-3124-1234", User_roles: []string{"USER"}, Refresh_uuid: "123-1234-1234"}, err: errNoContext.Error()},
				{name: "no context", claims: &entity.AccessTokenClaims{}, err: errNoContext.Error()},
				{name: "empty context", claims: &entity.AccessTokenClaims{}, context: context.Background(), err: errNoClaimsInContext.Error()},
			}
			for _, tc := range tt {
				t.Run(tc.name, func(t *testing.T) {
					claims, err := GetJWTClaims(tc.context)
					if err != nil && tc.err != err.Error() {
						t.Errorf("\t%s\tTest %s:\tWant error : %v, got error: %v", tests.Failed, tc.name, tc.err, err.Error())
					}
					t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate error", tests.Success, tc.name)

					if !reflect.DeepEqual(claims, tc.claims) {
						t.Errorf("\t%s\tTest %s:\tWant claims: %v, got claims: %v", tests.Failed, tc.name, tc.claims, claims)
					}
					t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate claims", tests.Success, tc.name)

				})
			}
		})
}

func TestCheckHeader(t *testing.T) {
	t.Run("CheckHeader middleware", func(t *testing.T) {
		tt := []struct {
			name                  string
			reqHeaders            []header
			wantHeader            header
			message               string
			nextHandlerInvocation bool
			statusCode            int
		}{
			{name: "header test: test", reqHeaders: []header{{key: "test", value: "test"}}, wantHeader: header{key: "test", value: "test"}, nextHandlerInvocation: true, statusCode: http.StatusOK},
			{name: "empty headers", wantHeader: header{key: "content-type", value: "application/json"}, message: "header should have value of", statusCode: http.StatusBadRequest},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})

				handler := CheckHeader(tc.wantHeader.key, tc.wantHeader.value, tc.statusCode)(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				for _, h := range tc.reqHeaders {
					req.Header.Add(h.key, h.value)
				}
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				res := rec.Result()
				defer res.Body.Close()
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Errorf("\t%s\tTest %s:\tCould not read response: %v", tests.Failed, tc.name, err)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to read response", tests.Success, tc.name)

				respBody := strings.TrimSpace(string(b))
				if !strings.Contains(respBody, tc.message) {
					t.Errorf("\t%s\tTest %s:\tWant response body: %v, got response body: %v", tests.Failed, tc.name, tc.message, respBody)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate response body", tests.Success, tc.name)

				if res.StatusCode != tc.statusCode {
					t.Errorf("\t%s\tTest %s:\tWant status code: %d, got status code: %d", tests.Failed, tc.name, res.StatusCode, tc.statusCode)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate status code", tests.Success, tc.name)
			})
		}
	})

	t.Run("CheckJSONHeader middleware", func(t *testing.T) {
		tt := []struct {
			name                  string
			reqHeaders            []header
			message               string
			nextHandlerInvocation bool
			statusCode            int
		}{
			{name: "header content-type: application/json", reqHeaders: []header{{key: "content-type", value: "application/json"}}, nextHandlerInvocation: true, statusCode: http.StatusOK},
			{name: "header content-type: application", reqHeaders: []header{{key: "content-type", value: "application"}}, nextHandlerInvocation: true, statusCode: http.StatusUnsupportedMediaType},
			{name: "empty headers", message: "header should have value of", statusCode: http.StatusUnsupportedMediaType},
		}
		for _, tc := range tt {
			t.Run("CheckJSONHeader middleware"+tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})

				handler := CheckJSONHeader()(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				for _, h := range tc.reqHeaders {
					req.Header.Add(h.key, h.value)
				}
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				res := rec.Result()
				defer res.Body.Close()
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Errorf("\t%s\tTest %s:\tCould not read response: %v", tests.Failed, tc.name, err)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to read response", tests.Success, tc.name)

				respBody := strings.TrimSpace(string(b))
				if !strings.Contains(respBody, tc.message) {
					t.Errorf("\t%s\tTest %s:\tWant response body: %v, got response body: %v", tests.Failed, tc.name, tc.message, respBody)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate response body", tests.Success, tc.name)

				if res.StatusCode != tc.statusCode {
					t.Errorf("\t%s\tTest %s:\tWant status code: %d, got status code: %d", tests.Failed, tc.name, res.StatusCode, tc.statusCode)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive appropriate status code", tests.Success, tc.name)
			})
		}
	})
}

func TestCors(t *testing.T) {
	t.Run("Cors middleware test", func(t *testing.T) {
		tt := []struct {
			name                  string
			responseHeaders       []header
			requestMethod         string
			nextHandlerInvocation bool
		}{
			{name: "CORS headers", responseHeaders: []header{{key: "Access-Control-Allow-Origin", value: "*"}, {key: "Access-Control-Allow-Methods", value: "GET, POST, OPTIONS"}, {key: "Access-Control-Allow-Headers", value: "Accept, Authorization, Content-Type"}}, nextHandlerInvocation: true},
			{name: "OPTIONS request method", requestMethod: "OPTIONS", nextHandlerInvocation: false},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})

				handler := Cors(nextHandler)
				req := httptest.NewRequest(tc.requestMethod, "/", nil)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				res := rec.Result()

				headers := res.Header
				for _, v := range tc.responseHeaders {
					h := headers.Values(v.key)
					hv := strings.Join(h, " ")
					if hv != v.value {
						t.Errorf("\t%s\tTest %s:\tWant header: %s, value: %s; got header: %s, value: %s", tests.Failed, tc.name, v.key, v.value, v.key, hv)
					}
					t.Logf("\t%s\tTest %s:\tShould be able to receive a appropriate request headers", tests.Success, tc.name)
				}
			})
		}
	})
}
func TestRequestInfo(t *testing.T) {
	t.Run("RequestInfo middleware test", func(t *testing.T) {
		tt := []struct {
			name    string
			message string
		}{
			{name: "check request id"},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if _, err := GetRequestInfo(r.Context()); err != nil {
						t.Errorf("\t%s\tTest %s:\tShould be able to receive a request info from next handler context. Error: %v", tests.Failed, tc.name, err)
					}
				})
				t.Logf("\t%s\tTest %s:\tShould be able to receive a request info from next handler context", tests.Success, tc.name)

				handler := RequestInfo(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)
			})
		}
	})

	t.Run("GetRequestInfo function test", func(t *testing.T) {
		tt := []struct {
			name    string
			context context.Context
			err     string
		}{
			{name: "valid context", context: context.WithValue(context.Background(), RequestKey, &Request{ID: uuid.NewString(), Now: time.Now()})},
			{name: "no context", err: errNoContext.Error()},
			{name: "empty context", context: context.Background(), err: errNoRequestInfoInContext.Error()},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				if _, err := GetRequestInfo(tc.context); err != nil && err.Error() != tc.err {
					t.Errorf("\t%s\tTest %s:\tWant error : %v, got error: %v", tests.Failed, tc.name, tc.err, err.Error())
				}
				t.Logf("\t%s\tTest %s:\tShould be able to receive a request info from context", tests.Success, tc.name)
			})
		}
	})
}

func TestLogger(t *testing.T) {
	t.Run("Logger middleware test", func(t *testing.T) {
		tt := []struct {
			name                  string
			context               context.Context
			nextHandlerInvocation bool
			err                   string
		}{
			{name: "valid context", context: context.WithValue(context.Background(), RequestKey, &Request{ID: uuid.NewString(), Now: time.Now()}), nextHandlerInvocation: true},
			{name: "invalid context", context: context.Background(), nextHandlerInvocation: false},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				cu := NewCustomLogger()

				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !tc.nextHandlerInvocation {
						t.Errorf("\t%s\tTest %s:\tNext handler should be invoked", tests.Failed, tc.name)
					}
					t.Logf("\t%s\tTest %s:\tNext handler should be invoked", tests.Success, tc.name)
				})
				t.Logf("\t%s\tTest %s:\tShould be able to receive a request info from next handler context", tests.Success, tc.name)

				handler := Logger(cu)(nextHandler)
				req := httptest.NewRequest("GET", "/", nil)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req.WithContext(tc.context))
				if len(cu.log) == 0 && tc.nextHandlerInvocation {
					t.Errorf("\t%s\tTest %s:\tShould not be able to log information about request", tests.Failed, tc.name)
				}
				t.Logf("\t%s\tTest %s:\tShould be able to log information about request", tests.Success, tc.name)
			})
		}
	})
}
