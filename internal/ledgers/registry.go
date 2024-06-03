// Package containing ledgers
//
// Ledgers act as service providers for various operations such
// as creating, storing, retrieving, validating, and verifying payments.
package ledgers

import (
	"sharkfoursix/go-payup/pkg"

	"github.com/pkg/errors"
)

var (
	ledgerRegistry map[string]pkg.LedgerOpener = map[string]pkg.LedgerOpener{}
)

func registerLedger(name string, opener pkg.LedgerOpener) {
	ledgerRegistry[name] = opener
}

func GetRegisteredLedgers() []string {
	list := []string{}
	for l := range ledgerRegistry {
		list = append(list, l)
	}
	return list
}

func Register(name string, opener pkg.LedgerOpener) {
	ledgerRegistry[name] = opener
}

func New(name, dsn string) (pkg.Ledger, error) {
	o, found := ledgerRegistry[name]
	if !found {
		return nil, errors.Errorf("unsupported ledger `%s`", name)
	}
	return o(dsn)
}
