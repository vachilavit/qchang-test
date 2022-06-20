package model

type CashierDesk map[float64]int // float64 = value, int = amount

type ChangeMoney struct {
	Change                float64
	AmountBankOrCoinValue CashierDesk
}
