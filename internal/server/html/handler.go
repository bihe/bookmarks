package html

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/bihe/commons-go/cookies"
	"github.com/bihe/commons-go/errors"
	"github.com/bihe/commons-go/handler"
)

const templateDir = "templates"
const errorTemplate = "error.tmpl"

// TemplateHandler is used to display errors using HTML templates
type TemplateHandler struct {
	handler.Handler
	CookieSettings cookies.Settings
	Version        string
	Build          string
	BasePath       string
}

// HandleError returns a HTML template showing errors
func (h *TemplateHandler) HandleError(w http.ResponseWriter, r *http.Request) error {
	cookie := cookies.NewAppCookie(h.CookieSettings)
	var (
		msg       string
		err       string
		isError   bool
		isMessage bool
		init      sync.Once
		tmpl      *template.Template
		e         error
	)

	init.Do(func() {
		path := filepath.Join(h.BasePath, templateDir, errorTemplate)
		tmpl, e = template.ParseFiles(path)
	})
	if e != nil {
		return e
	}

	// read (flash)
	err = cookie.Get(errors.FlashKeyError, r)
	if err != "" {
		isError = true
	}
	msg = cookie.Get(errors.FlashKeyInfo, r)
	if msg != "" {
		isMessage = true
	}

	// clear (flash)
	cookie.Del(errors.FlashKeyError, w)
	cookie.Del(errors.FlashKeyInfo, w)

	data := map[string]interface{}{
		"year":      time.Now().Year(),
		"appname":   "bookmarks.binggl.net",
		"version":   fmt.Sprintf("%s-%s", h.Version, h.Build),
		"isError":   isError,
		"error":     err,
		"isMessage": isMessage,
		"msg":       msg,
	}

	return tmpl.Execute(w, data)
}
