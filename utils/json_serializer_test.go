package utils

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
	"time"
)

func TestRegisterTimeSerializer(t *testing.T) {
	RegisterTimeSerializer(DefaultTimeFormat, DefaultTimeLocation)

	var testTime, zeroTime time.Time
	testTime = time.Date(2020, time.May, 23, 1, 11, 0, 0, time.Now().Location())
	timeString := testTime.Format(DefaultTimeFormat)
	zeroTimeString := zeroTime.Format(DefaultTimeFormat)
	var obj1 struct {
		Time             time.Time  `json:"time"`
		EmptyTime        time.Time  `json:"empty_time"`
		TimePointer      *time.Time `json:"time_pointer"`
		EmptyTimePointer *time.Time `json:"empty_time_pointer"`
	}
	obj1.Time = testTime
	obj1.TimePointer = &testTime
	jsonBytes, err := jsoniter.Marshal(obj1)
	if err != nil {
		t.Error(err)
		return
	}

	jsonMap := make(map[string]*string)
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		t.Error(err)
		return
	}
	if *jsonMap["time"] != timeString || *jsonMap["empty_time"] != zeroTimeString ||
		*jsonMap["time_pointer"] != timeString || jsonMap["empty_time_pointer"] != nil {
		t.Fail()
		return
	}

	var decodeTestObj struct {
		Time             time.Time  `json:"time"`
		TimePointer      *time.Time `json:"time_pointer"`
		EmptyTimePointer *time.Time `json:"empty_time_pointer"`
	}

	jsonString := fmt.Sprintf(`{"time":"%s","time_pointer":"%s","empty_time_pointer":null}`, timeString, timeString)
	err = jsoniter.UnmarshalFromString(jsonString, &decodeTestObj)
	if err != nil {
		t.Error(err)
		return
	}
	if !decodeTestObj.Time.Equal(testTime) || !decodeTestObj.TimePointer.Equal(testTime) || decodeTestObj.EmptyTimePointer != nil {
		t.Fail()
		return
	}
}
