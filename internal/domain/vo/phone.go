package vo

import (
	"fmt"
	"strings"
)

type Phone struct {
	Number string
}

func NewPhone(number string) (Phone, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return Phone{}, fmt.Errorf("phone number cannot be empty")
	}
	// Basic normalization map
	if strings.HasPrefix(number, "08") {
		number = "628" + number[2:]
	}
	return Phone{Number: number}, nil
}
