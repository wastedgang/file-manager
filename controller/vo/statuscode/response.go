package statuscode

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Response) SetCode(code string) *Response {
	return &Response{Code: code, Message: r.Message, Data: r.Data}
}

func (r *Response) SetMessage(message string) *Response {
	return &Response{Code: r.Code, Message: message, Data: r.Data}
}

func (r *Response) SetData(data interface{}) *Response {
	return &Response{Code: r.Code, Message: r.Message, Data: data}
}

func (r *Response) AddField(fieldName string, fieldData interface{}) *Response {
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		data = make(map[string]interface{})
	}
	data[fieldName] = fieldData
	return &Response{Code: r.Code, Message: r.Message, Data: data}
}
