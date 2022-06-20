package repository

import (
	"math"
	"qchang-test/internal/model"
	"sort"
)

type CashierDeskRepository struct {
	cashierDesk model.CashierDesk
}

func (r *CashierDeskRepository) FillCashierDesk() model.CashierDesk {
	r.cashierDesk = fillCashierDesk()
	return r.cashierDesk
}

func (r *CashierDeskRepository) CheckRemainMoney() (money float64) {
	for value, amount := range r.cashierDesk {
		money += value * float64(amount)
	}

	return
}

func (r *CashierDeskRepository) CalculateChange(paid, priceOfProduct float64) (*model.ChangeMoney, error) {
	if priceOfProduct <= 0 || paid <= 0 {
		return nil, ErrPaidOrPriceIsZero
	}

	var result model.ChangeMoney
	result.AmountBankOrCoinValue = make(model.CashierDesk)
	change := paid - priceOfProduct
	remain := change

	if change < 0 {
		return nil, ErrPriceMoreThanPaid
	}
	if change > r.CheckRemainMoney() {
		return nil, ErrMoneyNotEnough
	}

	keys := make([]float64, 0, len(r.cashierDesk))
	for k := range r.cashierDesk {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, value := range keys {
		amount := r.cashierDesk[value]
		if amount == 0 {
			continue
		}
		if remain >= value {
			bankOrCoinValueAmount := int(math.Floor(remain / value))

			if amount < bankOrCoinValueAmount {
				remainingBankOrCoinValue := bankOrCoinValueAmount - amount
				result.AmountBankOrCoinValue[value] = remainingBankOrCoinValue
				remain = remain - value*float64(remainingBankOrCoinValue)
			} else {
				result.AmountBankOrCoinValue[value] = bankOrCoinValueAmount
				remain = math.Mod(remain, value)
			}

			if remain == 0 {
				break
			}
		}
	}
	if remain != 0 {
		return nil, ErrBankOrCoinNotEnough
	}

	result.Change = change

	return &result, nil
}

func (r *CashierDeskRepository) WithdrawnCashierDesk(transaction model.CashierDesk) error {
	for value, amount := range transaction {
		if r.cashierDesk[value] < amount {
			return ErrTransation
		}
	}

	for value, amount := range transaction {
		r.cashierDesk[value] -= amount
	}
	return nil
}

func (r *CashierDeskRepository) AddCashierDesk(transaction model.CashierDesk) error {
	limit := fillCashierDesk()
	for value, amount := range transaction {
		if amount+r.cashierDesk[value] > limit[value] {
			return ErrTransation
		}
	}

	for value, amount := range transaction {
		r.cashierDesk[value] += amount
	}
	return nil
}

func NewCashierDeskRepository() *CashierDeskRepository {
	return &CashierDeskRepository{
		cashierDesk: fillCashierDesk(),
	}
}

func fillCashierDesk() model.CashierDesk {
	return model.CashierDesk{
		1000: 10,
		500:  20,
		100:  15,
		50:   20,
		20:   30,
		10:   20,
		5:    20,
		1:    20,
		0.25: 50,
	}
}
