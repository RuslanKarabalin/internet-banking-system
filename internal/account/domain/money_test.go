package domain

import "testing"

func TestNewMoney_Success(t *testing.T) {
	money, err := NewMoney(10050, "USD")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if money.Amount() != 10050 {
		t.Errorf("expected amount 10050, got %d", money.Amount())
	}

	if money.Currency() != "USD" {
		t.Errorf("expected currency USD, got %s", money.Currency())
	}
}

func TestNewMoney_NegativeAmount(t *testing.T) {
	_, err := NewMoney(-10050, "USD")

	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}

func TestNewMoney_EmptyCurrency(t *testing.T) {
	_, err := NewMoney(10050, "")

	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}

func TestMoney_Add_SameCurrency(t *testing.T) {
	// given
	money100, _ := NewMoney(10000, "USD")
	money50, _ := NewMoney(5000, "USD")

	// when
	money150, err150 := money100.Add(money50)

	// then
	if err150 != nil {
		t.Errorf("expected no error, got %v", err150)
	}

	if money150.Amount() != 15000 {
		t.Errorf("expected amount 15000, got %d", money150.Amount())
	}

	if money150.Currency() != "USD" {
		t.Errorf("expected currency USD, got %s", money150.Currency())
	}
}

func TestMoney_Add_DifferentCurrencies(t *testing.T) {
	// given
	usd, _ := NewMoney(10000, "USD")
	eur, _ := NewMoney(5000, "EUR")

	// when
	_, errUsdPlusEur := usd.Add(eur)

	// then
	if errUsdPlusEur == nil {
		t.Errorf("expected error, got %v", errUsdPlusEur)
	}
}

func TestMoney_Subtract_Success(t *testing.T) {
	// given
	money100, _ := NewMoney(10000, "USD")
	money30, _ := NewMoney(3000, "USD")

	// when
	money70, err70 := money100.Subtract(money30)

	// then
	if err70 != nil {
		t.Errorf("expected no error, got %v", err70)
	}

	if money70.Amount() != 7000 {
		t.Errorf("expected amount 7000, got %d", money70.Amount())
	}

	if money70.Currency() != "USD" {
		t.Errorf("expected currency USD, got %s", money70.Currency())
	}
}

func TestMoney_Subtract_InsufficientFunds(t *testing.T) {
	// given
	money50, _ := NewMoney(5000, "USD")
	money100, _ := NewMoney(10000, "USD")

	// when
	_, errMinus50 := money50.Subtract(money100)

	// then
	if errMinus50 == nil {
		t.Errorf("expected error, got %v", errMinus50)
	}
}
