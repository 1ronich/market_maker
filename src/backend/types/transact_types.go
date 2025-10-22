package types

import "github.com/gagliardetto/solana-go"

type Atastruct struct {
	PayerPrivateKey    solana.PrivateKey
	LPmints            [2]solana.PublicKey
	PoolSourceToken    solana.PublicKey
	PoolDestToken      solana.PublicKey
	UserSourceTokenAcc solana.PublicKey
	UserDestTokenAcc   solana.PublicKey
}

type SwapInfo struct {
	PayerPkey   solana.PrivateKey
	PayerPubkey solana.PublicKey

	Pool PoolInfo

	AmountIn     uint64
	MinAmountOut uint64
	Slippage     uint64

	RPC string
}

type PoolInfo struct {
	Address solana.PublicKey
	LPMintA solana.PublicKey
	LPMintB solana.PublicKey

	AmmProgram    solana.PublicKey // AMM protocol program ID (e.g., Raydium)
	AmmPool       solana.PublicKey // AMM pool state account
	AmmAuthority  solana.PublicKey // AMM pool authority
	AmmOpenOrders solana.PublicKey // AMM pool open orders account (on Serum)
	AmmCoinVault  solana.PublicKey // AMM pool's base token vault
	AmmPcVault    solana.PublicKey // AMM pool's quote token vault

	MarketProgram     solana.PublicKey // Serum DEX program ID
	Market            solana.PublicKey // Serum market account
	MarketBids        solana.PublicKey // Serum market bids orderbook
	MarketAsks        solana.PublicKey // Serum market asks orderbook
	MarketEventQueue  solana.PublicKey // Serum market event queue
	MarketCoinVault   solana.PublicKey // Serum market's base token vault
	MarketPcVault     solana.PublicKey // Serum market's quote token vault
	MarketVaultSigner solana.PublicKey // Serum market vault authority

	UserTokenSource      solana.PublicKey // User's source token account
	UserTokenDestination solana.PublicKey // User's destination token account
	UserSourceOwner      solana.PublicKey // User's wallet (signer)
}

type Transf struct {
	Sender    solana.PrivateKey
	Receiver  solana.PublicKey
	Authority solana.PublicKey
	Amount    uint64
}
