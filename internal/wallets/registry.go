// Package containts mobile wallet implementations
package wallets

import (
	"github.com/sharkfoursix/go-payup/pkg"

	"github.com/pkg/errors"
)

var (
	registry map[string]pkg.WalletResolver = map[string]pkg.WalletResolver{}
)

func registerMobileWallet(name string, opener pkg.WalletResolver) {
	registry[name] = opener
}

func GetRegisteredMobileWallets() []string {
	list := []string{}
	for l := range registry {
		list = append(list, l)
	}
	return list
}

func Register(name string, resolver pkg.WalletResolver) {
	registry[name] = resolver
}

func New(name, dsn string) (pkg.MobileWallet, error) {
	resolve, found := registry[name]
	if !found {
		return nil, errors.Errorf("unsupported wallet provider `%s`", name)
	}
	return resolve(dsn)
}
