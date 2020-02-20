package html

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/handler"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestErrorPage(t *testing.T) {
	r := chi.NewRouter()

	cookieSettings := cookies.Settings{
		Path:   "/",
		Domain: "localhost",
	}

	errHandler := &TemplateHandler{
		Handler: handler.Handler{
			ErrRep: &errors.ErrorReporter{
				CookieSettings: cookieSettings,
				ErrorPath:      "error",
			},
		},
		Version:        "1",
		Build:          "2",
		CookieSettings: cookieSettings,
		BasePath:       "../../../",
	}

	r.Get("/error", errHandler.Call(errHandler.HandleError))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	if body == "" {
		t.Errorf("could not get the template!")
	}
}
