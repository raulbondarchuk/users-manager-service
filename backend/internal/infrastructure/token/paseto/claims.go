package paseto

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PasetoClaims — typed fields that you want to store in the token.
type PasetoClaims struct {
	Username    string `json:"username"`
	CompanyID   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	Roles       string `json:"roles"`
	// IsPrimary     bool   `json:"isPrimary"`
	OwnerUsername string `json:"ownerUsername"`

	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func capitalizeKey(key string) string {
	if len(key) == 0 {
		return key
	}
	return strings.ToUpper(key[:1]) + key[1:]
}

func (c *PasetoClaims) SetFromMap(data map[string]interface{}) {
	v := reflect.ValueOf(c).Elem()
	for key, value := range data {
		// Capitalize the key to match the field name
		fieldName := capitalizeKey(key)
		field := v.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			val := reflect.ValueOf(value)
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			} else if field.Kind() == reflect.Int && val.Kind() == reflect.String {
				// Example: if the number is saved as a string
				if num, err := strconv.Atoi(val.String()); err == nil {
					field.SetInt(int64(num))
				}
			} else if field.Kind() == reflect.Bool && val.Kind() == reflect.String {
				if b, err := strconv.ParseBool(val.String()); err == nil {
					field.SetBool(b)
				}
			}
		}
	}
}
