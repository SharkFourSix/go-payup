package gopayup_test

import (
	gopayup "github.com/sharkfoursix/go-payup"
	"github.com/sharkfoursix/go-payup/pkg"
	"testing"
)

func TestAirtelMoneyWallet(t *testing.T) {
	var (
		err    error
		ledger pkg.Ledger
		wallet pkg.MobileWallet
		dsn    string
	)

	dsn = ""

	ledger, err = gopayup.NewLedger("sqliteLedger", "")
	if err != nil {
		//t.Fatal(err)
	}

	wallet, err = gopayup.NewMobileWallet("airtelMoney", dsn)
	if err != nil {
		//t.Fatal(err)
	}

	wallet = wallet
	ledger = ledger
	//wallet.VerifyTransaction(nil, "")
	//ledger.VerifyPayment(nil, nil)

	t.Log(gopayup.GetRegisteredLedgers())
	t.Log(gopayup.GetRegisteredWallets())

	//amw.VerifyTransaction()
}
