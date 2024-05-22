package utils

import (
	"github.com/asaskevich/govalidator"
	"strings"
)

func init() {
	govalidator.CustomTypeTagMap.Set("url", govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
		subject, ok := i.(string)
		if !ok {
			return false
		}
		if !govalidator.IsURL(subject) {
			return false
		}
		return strings.HasPrefix(subject, "http://") || strings.HasPrefix(subject, "https://")
	}))
}
