package domain

import (
	"errors"
	"github.com/google/uuid"
)

type AccountStatus int

const (
	AccountStatusActive AccountStatus = iota
	AccountStatusBlocked
	AccountStatusClosed
)

type Account struct {
	id      uuid.UUID
	ownerID uuid.UUID
	balance Money
	status  AccountStatus
}

func NewAccount(ownerID uuid.UUID, currency string) (*Account, error) {
	id := uuid.New()
	
	balance, err := NewMoney(0, currency)
	if err != nil {
		return nil, err
	}
	
	status := AccountStatusActive
	
	return &Account{id, ownerID, balance, status}, nil
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) OwnerID() uuid.UUID {
	return a.ownerID
}

func (a *Account) Balance() Money {
	return a.balance
}

func (a *Account) Status() AccountStatus {
	return a.status
}

func (a *Account) Deposit(amount Money) error {
	if a.status != AccountStatusActive {
		return errors.New("Can't deposit at inactive account!")
	}
	if a.balance.Currency() != amount.Currency() {
		return errors.New("Currency must be equals!")
	}
	if amount.IsZero() {
		return errors.New("Meaningless add 0!")
	}
	var err error
	a.balance, err = a.balance.Add(amount)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) Withdraw(amount Money) error {
	if a.status != AccountStatusActive {
		return errors.New("Can't withdraw at inactive account!")
	}
	if a.balance.Currency() != amount.Currency() {
		return errors.New("Currency must be equals!")
	}
	if amount.Amount() == 0 {
		return errors.New("Meaningless subtract 0!")
	}
	var err error
	a.balance, err = a.balance.Subtract(amount)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) Block() error {
	if a.status != AccountStatusActive {	
		return errors.New("Can't block inactive account!")
	}
	a.status = AccountStatusBlocked
	return nil
}

func (a *Account) Close() error {	
	if a.status != AccountStatusActive {	
		return errors.New("Cannot close inactive account!")
	}
	if !a.balance.IsZero() {
		return errors.New("Cannot close account with non-zero balance!")
    }
	a.status = AccountStatusClosed
	return nil
}
