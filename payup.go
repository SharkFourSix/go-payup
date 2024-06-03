package gopayup

import (
	"github.com/sharkfoursix/go-payup/internal/ledgers"
	"github.com/sharkfoursix/go-payup/internal/wallets"
	"github.com/sharkfoursix/go-payup/pkg"
)

func NewMobileWallet(name, dsn string) (pkg.MobileWallet, error) {
	return wallets.New(name, dsn)
}

// Creates a ledger
func NewLedger(name, dsn string) (pkg.Ledger, error) {
	return ledgers.New(name, dsn)
}

func GetRegisteredLedgers() []string {
	return ledgers.GetRegisteredLedgers()
}

func GetRegisteredWallets() []string {
	return wallets.GetRegisteredMobileWallets()
}

// provides a top level method of registering a wallet
func RegisterWallet(name string, resolver func(dsn string) (pkg.MobileWallet, error)) {
	wallets.Register(name, resolver)
}

func RegisterLedger(name string, opener func(dsn string) (pkg.Ledger, error)) {
	ledgers.Register(name, opener)
}
