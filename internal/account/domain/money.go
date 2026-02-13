package domain

import "errors"

type Money struct {
	amount   int64
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("Amount cannot be negative!")
	}
	if len(currency) == 0 {
		return Money{}, errors.New("Currency can't be empty!")
	}
	return Money{amount, currency}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("Currency must be equals!")
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("Currency must be equals!")
	}
	if m.amount-other.amount < 0 {
		return Money{}, errors.New("Amount cannot be negative!")
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) Equals(other Money) bool {
	return m.currency == other.currency && m.amount == other.amount
}
