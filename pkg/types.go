package pkg

import (
	"context"
	"database/sql/driver"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type PaymentStatus int

const (
	PS_NOT_PAID PaymentStatus = iota // Payment has not been completed and is awaiting verification
	PS_PAID                          // Payment has been paid and verified
	PS_ERROR                         // There was an error paying the payment on the provider's side
)

type TransactionStatus int

const (
	TS_PENDING TransactionStatus = iota
	TS_SUCCESS
	TS_FAILED
	TS_EXPIRED
)

var (
	// Transaction not found
	ErrTransactionNotFound = errors.New("transaction not found")

	// Payment not found
	ErrPaymentNotFound = errors.New("payment not found")

	int2PaymentStatusMap map[int]PaymentStatus = map[int]PaymentStatus{
		0: PS_NOT_PAID,
		1: PS_PAID,
		2: PS_ERROR,
	}
)

func (ps PaymentStatus) Value() (driver.Value, error) {
	return int(ps), nil
}

func (ps *PaymentStatus) Scan(value any) (err error) {
	if value == nil {
		*ps = PS_NOT_PAID
		return nil
	}

	set := func(v int) {
		s, found := int2PaymentStatusMap[v]
		if found {
			*ps = s
		} else {
			err = errors.Errorf("invalid payment status ordinal value %d", v)
		}
	}

	switch e := value.(type) {
	case int:
		set(e)
	case int16:
		set(int(e))
	case int8:
		set(int(e))
	case byte:
		set(int(e))
	case int64:
		set(int(e))
	case string:
		v, err := strconv.Atoi(e)
		if err != nil {
			return errors.Wrap(err, "error converting payament status")
		}
		set(v)
	}
	return
}

type Payment struct {
	ID          int64         `prof:"id"`           // Primary ID
	Amount      float64       `prof:"amount"`       // Payment amount
	Description *string       `prof:"description"`  // Description of the payment
	RefID       string        `prof:"ref_id"`       // Payment reference ID, generated by the library. Will be sent to providers
	Msisdn      string        `prof:"msisdn"`       // Payer
	CreatedAt   time.Time     `prof:"created_at"`   // UTC timestamp of when the payment was initially created
	Status      PaymentStatus `prof:"status"`       // Payment status
	TxnDetails  *string       `prof:"txn_details"`  // Text containing provider raw transaction details response
	ValidBefore time.Time     `prof:"valid_before"` // Max validity
	Code        int64         `prof:"payment_code"` // Unique payment code
}

type Transaction interface {
	// Transaction ID by the provider
	ID() string
	// Transaction ID initially specified by the merchant/caller
	RefID() string
	// Returns transaction status, as reported by the provider
	Status() TransactionStatus
	// Amount
	Amount() float64
	// Time the transaction was created. Optional, as some providers may not include this detail
	CreatedAt() *time.Time
}

type MobileWallet interface {
	// Returns transaction details from an MNO.
	//
	// Whether the id is the one created by the merchant or MNO is up to the implementing function.
	//
	// If the function is successful, err will be nil and Transaction will contain transaction details.
	//
	// If the function fails, err will contain an error, and transaction will be nil. In case the transaction
	// was not found, `ErrTransactionNotFound` will be returned.
	VerifyTransaction(ctx context.Context, id string) (Transaction, error)
}

type Ledger interface {
	// Create a payment record
	//
	// 	msisdn
	// Number that will be making this payment
	// 	amount
	// The payment amount
	//
	// The currency will be decided by the wallet implementation
	CreatePayment(ctx context.Context, msisdn string, amount float64) (*Payment, error)

	// Updates a payment record
	UpdatePayment(ctx context.Context, payment *Payment) error

	// Verifies payment given a payment code and transaction ID from a wallet provider
	VerifyPayment(ctx context.Context, paymentCode, transactionID string) (*Payment, error)
}

// Used for instantiating registered wallets
type WalletResolver func(dsn string) (MobileWallet, error)

// Used for instantiating registered ledgers
type LedgerOpener func(dsn string) (Ledger, error)
