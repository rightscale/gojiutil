// Copyright (c) 2015 RightScale, Inc., see LICENSE

package gojiutil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

var ApplicationJSON = "application/json"

func WriteString(rw http.ResponseWriter, code int, str string) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(code)
	rw.Write([]byte(str))
}

func Printf(rw http.ResponseWriter, code int, message string, args ...interface{}) {
	str := fmt.Sprintf(message, args)
	WriteString(rw, code, str)
}

func WriteJSON(c web.C, rw http.ResponseWriter, code int, obj interface{}) {
	rw.Header().Set("Content-Type", ApplicationJSON+"; charset=utf-8")
	// we could stream the json, but then what do we do with errors?
	//   rw.WriteHeader(code)
	//   err := json.NewEncoder(rw).Encode(obj)
	// instead opt for correctness: render into a buffer and write the buffer
	buf, err := json.Marshal(obj)
	if err == nil {
		rw.WriteHeader(code)
		rw.Write(buf) // we ignore errors here, sigh
	} else {
		c.Env["err"] = err.Error()
		rw.WriteHeader(http.StatusInternalServerError)
		errStr := fmt.Sprintf(`{"err": "Internal Error", "request_id": "%s"}`,
			middleware.GetReqID(c))
		rw.Write([]byte(errStr))
	}
}

// Produce a text/plain error response into the responseWriter and also sets the context to
// reflect the error in a way that the logger groks properly.
// For 500 errors a generic error is returned and the details are only logged.
func ErrorString(c web.C, rw http.ResponseWriter, code int, str string) {
	c.Env["err"] = str
	if code >= 500 {
		errStr := fmt.Sprintf("Internal Error (request ID: %s)", middleware.GetReqID(c))
		http.Error(rw, errStr, code)
	} else {
		http.Error(rw, str, code)
	}
}

// Convenience function to call ErrorString with a format string
func Errorf(c web.C, rw http.ResponseWriter, code int, message string, args ...interface{}) {
	str := fmt.Sprintf(message, args)
	ErrorString(c, rw, code, str)
}

// Convenience function to produce an internal error based on the err argument
func ErrorInternal(c web.C, rw http.ResponseWriter, err error) {
	if err != nil {
		ErrorString(c, rw, 500, err.Error())
	} else {
		ErrorString(c, rw, 500, "nil err passed into gojiutil.ErrorInternal")
	}
}
