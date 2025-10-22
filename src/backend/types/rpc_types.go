package types

type SolanaAccountResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			ApiVersion string `json:"apiVersion"`
			Slot       uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Data       [2]string `json:"data"` // ["", "base64"]
			Executable bool      `json:"executable"`
			Lamports   uint64    `json:"lamports"`
			Owner      string    `json:"owner"`
			RentEpoch  uint64    `json:"rentEpoch"`
			Space      uint64    `json:"space"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}
type BlockhashResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Blockhash            string `json:"blockhash"`
			LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}
type GetTokenBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			//ApiVersion string `json:"apiVersion"`
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Amount         string  `json:"amount"`
			Decimals       uint8   `json:"decimals"`
			UiAmount       float64 `json:"uiAmount"`
			UiAmountString string  `json:"uiAmountString"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

type GettProgramAccounts struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`

	Result struct {
		Pubkey  string `json:"pubkey"`
		Account struct {
			Lamports   uint64   `json:"lamports"`
			Owner      string   `json:"owner"`
			Data       []string `json:"data"` // Represents ["<encoded data>", "<encoding>"]
			Executable bool     `json:"executable"`
			RentEpoch  uint64   `json:"rentEpoch"`
			Space      uint64   `json:"space"`
		} `json:"account"`
	} `json:"result"`
}

type TokenAccountsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			APIVersion string `json:"apiVersion"`
			Slot       uint64 `json:"slot"`
		} `json:"context"`
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							IsNative    bool   `json:"isNative"`
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							State       string `json:"state"`
							TokenAmount struct {
								Amount         string  `json:"amount"`
								Decimals       int     `json:"decimals"`
								UIAmount       float64 `json:"uiAmount"`
								UIAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program string `json:"program"`
					Space   int    `json:"space"`
				} `json:"data"`
				Executable bool   `json:"executable"`
				Lamports   uint64 `json:"lamports"`
				Owner      string `json:"owner"`
				RentEpoch  uint64 `json:"rentEpoch"`
				Space      int    `json:"space"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

type SignatureStatusesResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Context struct {
			Slot int `json:"slot"`
		} `json:"context"`
		Value []*struct {
			Slot          int     `json:"slot"`
			Confirmations *int    `json:"confirmations"`
			Err           *string `json:"err"`
			Status        struct {
				Ok interface{} `json:"Ok"`
			} `json:"status"`
			ConfirmationStatus string `json:"confirmationStatus"`
		} `json:"value"`
	} `json:"result"`
}

type GetBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Context struct {
			Slot int `json:"slot"`
		} `json:"context"`
		Value uint64 `json:"value"`
	} `json:"result"`
}
