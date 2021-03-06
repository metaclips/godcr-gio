package ui

import (
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	loading        bool // true if the window is in the middle of an operation that cannot be stopped
	dialog         bool // true if the window dialog modal is open
	renamingWallet bool // true if the wallets-page is renaming a wallet
}

// updateStates changes the wallet state based on the received update
func (win *Window) updateStates(update interface{}) {
	switch e := update.(type) {
	case wallet.MultiWalletInfo:
		*win.walletInfo = e
		if len(win.outputs.tabs) != win.walletInfo.LoadedWallets {
			win.reloadTabs()
		}
		win.states.loading = false
		return
	case *wallet.Transactions:
		win.walletTransactions = e
		return
	case wallet.CreatedSeed:
		win.current = win.WalletsPage
		win.states.dialog = false
	case wallet.Restored:
		win.resetSeeds()
		win.current = win.WalletsPage
		win.states.dialog = false
	case wallet.DeletedWallet:
		win.selected = 0
		win.states.dialog = false
	case wallet.AddedAccount:
		win.states.dialog = false
	}
	win.states.loading = true
	win.wallet.GetMultiWalletInfo()
	win.wallet.GetAllTransactions(0, 10, 0)

	log.Debugf("Updated with multiwallet info: %+v\n and window state %+v", win.walletInfo, win.states)
}

func (win *Window) reloadTabs() {
	win.outputs.tabs = make([]decredmaterial.TabItem, win.walletInfo.LoadedWallets)
	for i := range win.outputs.tabs {
		win.outputs.tabs[i] = decredmaterial.TabItem{
			Button: win.theme.Button(win.walletInfo.Wallets[i].Name),
		}
	}
	win.tabs.SetTabs(win.outputs.tabs)
}
