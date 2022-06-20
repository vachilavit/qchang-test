package repository

import "errors"

var (
	ErrPaidOrPriceIsZero   = errors.New("paid or price is zero")
	ErrPriceMoreThanPaid   = errors.New("price of product more than paid")
	ErrMoneyNotEnough      = errors.New("money in cashier not enough")
	ErrBankOrCoinNotEnough = errors.New("bank or coin not enough")
	ErrTransation          = errors.New("transation fail")
)
