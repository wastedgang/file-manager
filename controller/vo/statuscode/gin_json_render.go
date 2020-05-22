package statuscode

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type JSONRender struct {
	Data interface{}
}

var (
	jsonContentType = []string{"application/json; charset=utf-8"}
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Render (JSONRender) writes data with custom ContentType.
func (r JSONRender) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSONRender) writes JSONRender ContentType.
func (r JSONRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := jsoniter.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

