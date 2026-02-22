package vo

import "fmt"

type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64, currency string) Money {
	if currency == "" {
		currency = "IDR"
	}
	return Money{Amount: amount, Currency: currency}
}

func (m Money) Formatted() string {
	return fmt.Sprintf("%s %.2f", m.Currency, m.Amount)
}
