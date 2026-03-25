package lib

import (
	"html"
	"time"

	"github.com/lib/pq"
)

func NewBoolPointer(value bool) *bool {
	x := value
	return &x
}

func BoolPtrToBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func StringPtrToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func StringToStringPtr(value string) *string {
	return &value
}

func CurrentTimeUTCPtr() *time.Time {
	currentTimeUTC := time.Now().UTC()
	return &currentTimeUTC
}

func StringArrayRemove(slice pq.StringArray, index int64) pq.StringArray {
	return append(slice[:index], slice[index+1:]...)
}

func StringArrayRemoveItem(slice pq.StringArray, item string) pq.StringArray {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func GetStringPtr(data string) *string {
	v := data
	return &v
}

func Sanitize(input string) string {
	return html.EscapeString(input)
}
