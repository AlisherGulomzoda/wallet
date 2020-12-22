package types

type Money int64

type PaymrntCategoty string

type PaymentStatus string

const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymrntCategoty
	Status    PaymentStatus
}

type Phone string

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
