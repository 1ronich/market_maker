package main

import (
	"amm/common"
	"amm/raydium"
	"amm/solana_methods"
	"amm/utils"
	"os"

	"fmt"
	"log"
	"math/rand"

	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"
)

var MASTER_PUBKEY = "3enQJKSDb3doDF5xHXiLvWniAcKGJK36SNLQywkADCdQ"
var MASTER_PRIVKEY string

func init() {
	err := godotenv.Load("/home/maester/Desktop/Projects/amm/src/backend/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	MASTER_PRIVKEY = os.Getenv("MASTER_PRIVATEKEY")
}

func GetVaultBalance(pubkey, mint, cluster string) float64 {
	vaultAcc, _, _ := solana.FindAssociatedTokenAddress(solana.MustPublicKeyFromBase58(pubkey), solana.MustPublicKeyFromBase58(mint)) //a is wsol

	res := solana_methods.GetTokenAccountBalance(vaultAcc.String(), cluster)

	return res.Result.Value.UiAmount
}

func maain() {
	//RealPhase(100000000000000000000)
	//CreateWallets(10, true)
	//LoadWallets()

}

func findPubkey() {

	var accounts []solana.Wallet
	utils.ReadJSONFile(&accounts, "./wallets.json")

	for _, account := range accounts {
		balance := solana_methods.GetBalance(account.PublicKey().String(), "devnet").Result.Value
		fmt.Println(balance)
		fmt.Println(account.PrivateKey.String())
		if balance > 1000000 {
			sig := Full_Transfer(solana_methods.GetBalance(account.PublicKey().String(), "devnet").Result.Value-5000, account.PrivateKey.String(), "3enQJKSDb3doDF5xHXiLvWniAcKGJK36SNLQywkADCdQ", "devnet")
			fmt.Println(sig)
		}
	}
}

func EstimatedFee(swap_amount float64) float64 {
	return (swap_amount / 100) * 25
}

func RealPhase(max_price float64) {
	fmt.Println(MASTER_PRIVKEY)
	fmt.Println(MASTER_PUBKEY)

	var multiplier float64
	wallets := LoadWallets()

	var ARM_wallets []map[*solana.Wallet]uint64
	for {
		var last_b_supply float64
		var last_q_supply float64
		var err error
		for {

			last_q_supply, last_b_supply, err = common.GetSupp("97p3acpjRH9kZJ189fRHzVzmWiTSsSoGZeCwFNmHDQhf", "devnet", "v4")
			if err != nil {
				fmt.Println(fmt.Errorf("error getting supply: %v\n", err))
				continue
			}
			fmt.Println("q supply:", last_q_supply)
			fmt.Println("b supply:", last_b_supply)
			if last_b_supply == 0 || last_q_supply == 0 {
				time.Sleep(200 * time.Millisecond)
				continue
			} else {
				break
			}
		}

		current_q_price, current_b_price := common.CalculatePrices(last_b_supply, last_q_supply)
		fmt.Println("current_b_price:", current_b_price)
		fmt.Println("current_q_price:", current_q_price)

		multiplier_rand := rand.Intn(2)
		if multiplier_rand%100 == 0 {
			multiplier = 3
		} else if multiplier_rand%9 == 0 {
			multiplier = 1.6
		} else {
			multiplier = 1.2
		}

		wallet := wallets[rand.Intn(11)]

		sell_or_buy := rand.Intn(2)
		if sell_or_buy == 0 { //buy
			fmt.Println("Swapping Q (selling SOL for ARM)")

			buy_amount := GetBuyAmount(multiplier)
			RealBuy(wallet, buy_amount)

			ARM_wallets = append(ARM_wallets, map[*solana.Wallet]uint64{&wallet: buy_amount})
		} else { //sell
			fmt.Println("Swapping B (selling ARM for SOL)")

			sell_amount := GetSellAmount(multiplier, current_b_price)
			RealSell(wallet, sell_amount)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func GetBuyAmount(multiplier float64) uint64 {
	var in_amount float64
	for {
		in_amount = (rand.Float64() * (multiplier))

		if in_amount == 0 {
			continue
		} else {
			break
		}
	}
	amount := uint64(in_amount * 10e8)

	return amount
}

func GetSellAmount(multiplier, token_price_in_b float64) uint64 {
	in_amount := (rand.Float64() * (multiplier)) * token_price_in_b

	amount := uint64(in_amount * 10e8)

	return amount
}

func RealSell(wallet solana.Wallet, amount uint64) {
	balance := GetVaultBalance(wallet.PublicKey().String(), "E725QEK1AeemfLSTnVcNZTaXabtwbzJ5JXeTDo7fN5Ca", "devnet")
	fmt.Println("balance:", balance)

	if balance < float64(amount) {

		siga := Full_Transfer(amount, MASTER_PRIVKEY, wallet.PublicKey().String(), "devnet")
		fmt.Println("Inner transf from main:", siga)
	}

	raydium.FullRaydiumDevSwap(amount, MASTER_PRIVKEY, false, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q") //"4atvsLG6qoARYpxmL4PmoziHmydEoP9RTpVUFbfQ1BF82ZQ3tm27Lofsq1THgezYt1UUkFYVxUUbK1YpqQHGZJVH", false)
}

// 291431920747
// 609769516514
func Full_Transfer(amount uint64, private_key, receiver, cluster string) string {
	transfer := solana_methods.Transfer(amount, solana.MustPrivateKeyFromBase58(private_key), solana.MustPublicKeyFromBase58(receiver))
	fmt.Println("amount:", amount)
	signature, _ := solana_methods.SignAndSendTransaction(transfer, solana.MustPrivateKeyFromBase58(private_key).PublicKey(), solana.MustPrivateKeyFromBase58(private_key), cluster)
	return signature.String()
}

func RealBuy(wallet solana.Wallet, amount uint64) {

	//privatekey := wallet.PrivateKey
	//pubkey := wallet.PublicKey()

	balance := GetVaultBalance("3enQJKSDb3doDF5xHXiLvWniAcKGJK36SNLQywkADCdQ", "So11111111111111111111111111111111111111112", "devnet")
	if balance < float64(amount) {
		lamports := amount
		siga, _ := solana_methods.SignAndSendTransaction(solana_methods.Transfer(lamports, solana.MustPrivateKeyFromBase58(MASTER_PRIVKEY), wallet.PublicKey()), solana.MustPublicKeyFromBase58(MASTER_PUBKEY), solana.MustPrivateKeyFromBase58(MASTER_PRIVKEY), "devnet")
		fmt.Println("Inner transf from main:", siga)
	}

	fmt.Println("buy amount:", amount)

	raydium.FullRaydiumDevSwap(amount, wallet.PrivateKey.String(), true, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q") //"4atvsLG6qoARYpxmL4PmoziHmydEoP9RTpVUFbfQ1BF82ZQ3tm27Lofsq1THgezYt1UUkFYVxUUbK1YpqQHGZJVH", true)
}

func CalculatePriceImpact(current_price, old_price float64) float64 { //Price impact = current price
	return old_price - current_price/old_price
}

func CreateWallets(amount int, save_to_file bool) {
	var wallets []solana.Wallet
	for i := 0; i < amount; i++ {
		account := solana.NewWallet()
		wallets = append(wallets, *account)

	}
	fmt.Println(wallets)
	fmt.Println(len(wallets))

	if save_to_file {
		utils.WriteToJSONFile(wallets, "wallets.json")
	}
}

func LoadWallets() []solana.Wallet {
	var wallets []solana.Wallet // MyWalletType implements solana.Wallet
	utils.ReadJSONFile(&wallets, "wallets.json")

	return wallets
}
