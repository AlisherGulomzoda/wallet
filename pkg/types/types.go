package types

type Money int64

type PaymentCategoty string

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
	Category  PaymentCategoty
	Status    PaymentStatus
}

type Phone string

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Messenger interface {
	Send(message string) (ok bool)
	Receive() (message string, ok bool)
}

type Telegram struct {

}

func (t *Telegram) Send(message string) bool {
	return true
}

func (t *Telegram) Receive() (message string, ok bool) {
	return "", true
}

type error interface {
	Error() string
}

type Error string

func (e Error) Error() string {
	return string(e)
}

type Favorite struct {
	ID string
	AccountID int64
	Name string
	Amount Money
	Category PaymentCategoty
}