package utils

import (
	"github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"time"
	"unsafe"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

var (
	DefaultTimeLocation = time.Local
)

func RegisterTimeSerializer(timeFormat string, location *time.Location) {
	// 添加time.Time类型的序列化器
	encodeFunc := func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		stream.WriteString((*(*time.Time)(ptr)).Format(timeFormat))
	}
	isEmptyFunc := func(ptr unsafe.Pointer) bool {
		return (*(*time.Time)(ptr)).IsZero()
	}
	now := time.Now()
	jsoniter.RegisterTypeEncoderFunc(reflect2.TypeOf(now).String(), encodeFunc, isEmptyFunc)

	// 添加time.Time类型的反序列表器
	decodeFunc := func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if iter.WhatIsNext() != jsoniter.StringValue {
			iter.ReportError("parse time.Time", "invalid time.Time value of format \""+timeFormat+"\"")
			return
		}
		t, err := time.ParseInLocation(timeFormat, iter.ReadString(), location)
		if err != nil {
			iter.ReportError("parse time.Time", "invalid time.Time value of format \""+timeFormat+"\"")
			return
		}
		*(*time.Time)(ptr) = t
	}
	jsoniter.RegisterTypeDecoderFunc(reflect2.TypeOf(now).String(), decodeFunc)

	// 添加*time.Time类型的序列化器
	pointerEncodeFunc := func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		pTime := *(**time.Time)(ptr)
		if pTime == nil {
			stream.WriteNil()
		} else {
			stream.WriteString((*pTime).Format(timeFormat))
		}
	}
	isEmptyPointerFunc := func(ptr unsafe.Pointer) bool {
		pTime := *(**time.Time)(ptr)
		return pTime == nil
	}
	jsoniter.RegisterTypeEncoderFunc(reflect2.TypeOf(&now).String(), pointerEncodeFunc, isEmptyPointerFunc)

	// 添加*time.Time类型的反序列化器
	pointerDecodeFunc := func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		whatIsValue := iter.WhatIsNext()

		if whatIsValue == jsoniter.NilValue {
			iter.Read()
			*(**time.Time)(ptr) = nil
			return
		} else if whatIsValue == jsoniter.StringValue {
			t, err := time.ParseInLocation(timeFormat, iter.ReadString(), location)
			if err == nil {
				*(**time.Time)(ptr) = &t
				return
			}
		}
		if iter.WhatIsNext() != jsoniter.StringValue {
			iter.ReportError("parse time.Time", "invalid time.Time value of format \""+timeFormat+"\"")
			return
		}
	}
	jsoniter.RegisterTypeDecoderFunc(reflect2.TypeOf(&now).String(), pointerDecodeFunc)
}
