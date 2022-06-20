package test

import (
	"qchang-test/internal/model"
	"qchang-test/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCashierDeskRepositoryFillCashierDesk(t *testing.T) {
	cashierDeskRepository := repository.NewCashierDeskRepository()
	assert.Equal(t, model.CashierDesk{
		1000: 10,
		500:  20,
		100:  15,
		50:   20,
		20:   30,
		10:   20,
		5:    20,
		1:    20,
		0.25: 50,
	}, cashierDeskRepository.FillCashierDesk())
}

func TestCashierDeskRepositoryCheckRemainMoney(t *testing.T) {
	cashierDeskRepository := repository.NewCashierDeskRepository()
	assert.Equal(t, 23432.5, cashierDeskRepository.CheckRemainMoney())
}

func TestCashierDeskRepositorWithdrawnCashierDesk(t *testing.T) {
	t.Parallel()
	t.Run("normal", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		err := cashierDeskRepository.WithdrawnCashierDesk(model.CashierDesk{
			1000: 9,
		})
		assert.NoError(t, err)
		assert.Equal(t, 14432.5, cashierDeskRepository.CheckRemainMoney())
	})

	t.Run("transation fail", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		err := cashierDeskRepository.WithdrawnCashierDesk(model.CashierDesk{
			1000: 11,
		})
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrTransation, err)
		}
		assert.Equal(t, 23432.5, cashierDeskRepository.CheckRemainMoney())
	})
}

func TestCashierDeskRepositorAddCashierDesk(t *testing.T) {
	t.Parallel()
	t.Run("normal", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		err := cashierDeskRepository.WithdrawnCashierDesk(model.CashierDesk{
			1000: 1,
		})
		assert.NoError(t, err)
		err = cashierDeskRepository.AddCashierDesk(model.CashierDesk{
			1000: 1,
		})
		assert.NoError(t, err)
		assert.Equal(t, 23432.5, cashierDeskRepository.CheckRemainMoney())
	})

	t.Run("transation fail", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		err := cashierDeskRepository.AddCashierDesk(model.CashierDesk{
			1000: 1,
		})
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrTransation, err)
		}
		assert.Equal(t, 23432.5, cashierDeskRepository.CheckRemainMoney())
	})
}

func TestCashierDeskRepositoryCalculateChange(t *testing.T) {
	t.Parallel()
	t.Run("normal", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		changeMoney, err := cashierDeskRepository.CalculateChange(3000, 1111.50)
		assert.NoError(t, err)
		assert.Equal(t, &model.ChangeMoney{
			Change: 1888.5,
			AmountBankOrCoinValue: model.CashierDesk{
				1000: 1,
				500:  1,
				100:  3,
				50:   1,
				20:   1,
				10:   1,
				5:    1,
				1:    3,
				0.25: 2,
			},
		}, changeMoney)
	})

	t.Run("paid or price is zero", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		changeMoney, err := cashierDeskRepository.CalculateChange(0, 0)
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrPaidOrPriceIsZero, err)
		}
		assert.Nil(t, changeMoney)
	})

	t.Run("price of product more than paid", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		changeMoney, err := cashierDeskRepository.CalculateChange(1, 1000)
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrPriceMoreThanPaid, err)
		}
		assert.Nil(t, changeMoney)
	})

	t.Run("money in cashier not enough", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		changeMoney, err := cashierDeskRepository.CalculateChange(cashierDeskRepository.CheckRemainMoney()+2, 1)
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrMoneyNotEnough, err)
		}
		assert.Nil(t, changeMoney)
	})

	t.Run("bank or coin not enough", func(t *testing.T) {
		cashierDeskRepository := repository.NewCashierDeskRepository()
		err := cashierDeskRepository.WithdrawnCashierDesk(model.CashierDesk{
			0.25: 48,
		})
		assert.NoError(t, err)

		changeMoney, err := cashierDeskRepository.CalculateChange(2, 1.25)
		if assert.Error(t, err) {
			assert.Equal(t, repository.ErrBankOrCoinNotEnough, err)
		}
		assert.Nil(t, changeMoney)
	})
}
