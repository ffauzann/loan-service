package sanitize

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/ffauzann/loan-service/pkg/common/util/str"
)

var sensitiveKeywordsLv1 = []string{
	"password", "secret", "token", "otp", "pin", "passcode", "secretkey",
	"token", "key", "sign", "accesstoken", "refreshtoken", "signature",
	"authorization", "grpcgatewayauthorization", "xsecretkey", "xsignature",
	"publickey", "privatekey",
}

var sensitiveKeywordsLv2 = []string{
	"phone", "phonenumber", "email",
}

var base64Keywords = []string{
	"file", "photo",
}

// Sanitize sanitizes sensitive data before print into logger
func Sanitize(data interface{}) interface{} {
	b, err := json.Marshal(data)
	if err != nil {
		return data
	}

	var obj map[string]interface{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return data
	}

	var deep func(items map[string]interface{}) map[string]interface{}
	deep = func(items map[string]interface{}) map[string]interface{} {
		for k, v := range items {
			if v2, ok := v.(map[string]interface{}); ok {
				items[k] = deep(v2)
			}
			switch genericKey := strings.ToLower(str.RemoveSeparator(k)); {
			case slices.Contains(sensitiveKeywordsLv1, genericKey):
				items[k] = HideData(fmt.Sprintf("%+v", v))
			case slices.Contains(sensitiveKeywordsLv2, genericKey):
				items[k] = HideHalfData(fmt.Sprintf("%+v", v))
			case slices.Contains(base64Keywords, genericKey):
				items[k] = HideBase64Data(fmt.Sprintf("%+v", v))
			}
		}
		return items
	}

	return deep(obj)
}

func HideData(input string) (output string) {
	output = strings.Repeat("*", len(input))
	if len(input) > 20 {
		output = strings.Replace(input, input[5:len(input)-5], strings.Repeat("*", len(input[5:len(input)-5])), 2)
	}
	return
}

func HideBase64Data(input string) (output string) {
	if len(input) > 50 {
		output = input[0:30] + "...." + input[len(input)-10:]
	}
	return
}

func HideHalfData(input string) (output string) {
	l := int(math.Round(float64(len(input) / 3)))
	first := input[:l]
	last := input[l*2:]
	output = first
	for i := 0; i < l; i++ {
		output += "*"
	}
	output += last
	return
}
