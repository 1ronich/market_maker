package main

import (
	"amm/common"
	"amm/raydium"
	"amm/solana_methods"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
)

const ARMCHE_LP = "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q"
const ARMCHE_mint = "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG"

func convertStrArrayToIntArray(arr []string) []byte {
	intArr := make([]byte, len(arr))
	for i, s := range arr {
		n, _ := strconv.Atoi(s)
		intArr[i] = byte(n)
	}

	return intArr
}


func moain() {
	wallets := []solana.Wallet{
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY1ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY2ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY3ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY4ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY5ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY6ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY7ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY8ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY9ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY1ARRAY"), ",")))},
	}

	var previous_q_price float64
	//var previous_b_price float64
	var balances map[string]Balance = make(map[string]Balance)

	SetBalances(balances, wallets)

	i := 0
	for {
		if i == len(wallets) {
			i = 0
		}

		current_q_price, current_b_price := common.GetPrices(ARMCHE_LP)

		var QtoB string = ""
		var amount uint64

		if i == 0 {
			amount = 0
		} else {
			amount = GetAmountToStabilize(current_q_price, previous_q_price, amount, QtoB) //This will return an amount needed to stabilize to price only if a critical % change is detected, if there is no critical % change then returns amount as 0
		}
		fmt.Println(amount)

		trade(balances, QtoB, current_q_price, current_b_price, amount)

		//previous_b_price = current_b_price
		previous_q_price = current_q_price
		i++

		timeout := rand.Float64() * 10000
		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}
}

func SetBalances(balances map[string]Balance, wallets []solana.Wallet) {
	for _, wallet := range wallets {
		balances[wallet.PrivateKey.String()] = Balance{TokenB: GetVaultBalance(wallet.PublicKey().String(), ARMCHE_mint, "devnet"), TokenQ: float64(solana_methods.GetBalance(wallet.PublicKey().String(), "devnet").Result.Value / 10e9)}
	}
}

func trade(balances map[string]Balance, QtoB string, current_q_price, current_b_price float64, amount uint64) {
	fmt.Println("QtoB before:", QtoB)
	QtoB = BuyOrSell(QtoB) //Randomly decides if to sell or buy or returns as it was if it was already defined
	fmt.Println("QtoB after:", QtoB)

	var available_total float64
	var private_keys []string

	if amount == 0 {
		tries := 0
		for {
			if tries == 5 {
				log.Fatal("BALANCE IS NOT ENOUGH")
			}
			fmt.Println("amount:", amount)

			amounta := AmountForTrade(QtoB, amount, current_b_price) //Randomly returns an amount if it was not already defined
			fmt.Println("amounta:", float64(amounta)/10e9)
			available, pkey := IsAmountAvailable(balances, float64(amounta)/10e9, QtoB) //Checks the Balances map to see if any of the wallets has the enough amount
			if available {
				private_keys = []string{pkey}
				amount = amounta
				break
			} else {
				tries++
				continue
			}

		}
	} else {
		available, prkey := IsAmountAvailable(balances, float64(amount)/10e9, QtoB) //Checks the Balances map to see if any of the wallets has the enough amount
		if !available {

			if QtoB == "true" {
				for pkey, balance := range balances {
					if available_total >= float64(amount)*10e9 {
						break
					}
					available_total += balance.TokenQ
					private_keys = append(private_keys, pkey)
				}
			} else if QtoB == "false" {
				for pkey, balance := range balances {
					if available_total >= float64(amount)*10e9 {
						break
					}
					available_total += balance.TokenB
					private_keys = append(private_keys, pkey)
				}
			}

		} else {
			private_keys = []string{prkey}
		}
	}

	//SimulateSwap(balances, amount, private_keys, QtoB)
	var isQtoB bool
	for _, private_key := range private_keys {
		if QtoB == "true" {
			isQtoB = true
		} else if QtoB == "false" {
			isQtoB = false
		}
		raydium.FullRaydiumDevSwap(amount, private_key, isQtoB, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")
	}
}

type Balance struct {
	TokenQ float64
	TokenB float64
}

func BuyOrSell(QtoB string) string {
	if QtoB == "" {
		if rand.Intn(2) == 0 {
			QtoB = "true"
		} else {
			QtoB = "false"
		}
	}
	return QtoB
}

func AmountForTrade(QtoB string, amount uint64, current_b_price float64) uint64 {
	if amount == 0 {
		if QtoB == "true" {
			amount = uint64(rand.Float64() * 10e7)
		} else if QtoB == "false" {
			amount = uint64((rand.Float64() * current_b_price * 10e7))
		}
	}
	return amount
}

func IsAmountAvailable(balances map[string]Balance, amount float64, QtoB string) (bool, string) {
	for privkey, balance := range balances {
		if QtoB == "true" {
			if balance.TokenQ >= float64(amount)/10e7 {
				return true, privkey
			} else {
				return false, ""
			}
		} else if QtoB == "false" {
			if balance.TokenB >= float64(amount)/10e7 {
				return true, privkey
			} else {
				return false, ""
			}
		}
	}
	return false, ""
}

func SimulateSwap(balances map[string]Balance, amount uint64, pkeys []string, SOLtoARM string) {
	for _, pkey := range pkeys {
		var tokenA string
		var tokenB string

		address := solana.MustPrivateKeyFromBase58(pkey).PublicKey().String()

		if SOLtoARM == "true" {
			tokenA = "SOL"
			tokenB = "ARM"
			balances[address] = Balance{TokenQ: float64(amount) / 10e9}

		} else if SOLtoARM == "false" {
			tokenA = "ARM"
			tokenB = "SOL"
		}

		fmt.Printf("Wallet %s swapped %f %s for %s\n", solana.MustPrivateKeyFromBase58(pkey).PublicKey().String(), float64(amount)/10e9, tokenA, tokenB)
	}

}

func ShouldBuy(percentage_change float64) bool {
	var IsBuy bool
	if percentage_change > 0 {
		//If percent_change is positive price has gone up, meaning that to stabilize the price I have to sell
		IsBuy = false
	} else {
		//If percent_change is negative price has gone down, meaning that to stabilize the price I have to buy
		IsBuy = true
	}
	return IsBuy
}

func CalculateTargetSupply(target_price, current_q_supply, current_b_supply float64) (float64, float64) { //Calculates needed supply to reach X price
	k := float64(current_b_supply * current_q_supply)

	target_b_supply := math.Sqrt(k / target_price)
	target_q_supply := k / target_b_supply

	return target_q_supply, target_b_supply
}

func EXPERIMENTAL() {
	wallets := []solana.Wallet{
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY1ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY2ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY3ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY4ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY5ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY6ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY7ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY8ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY9ARRAY"), ",")))},
		{PrivateKey: (convertStrArrayToIntArray(strings.Split(os.Getenv("PKEY1ARRAY"), ",")))},
	}

	for _, w := range wallets {
		fmt.Println(w.PublicKey().String())
		fmt.Println(MASTER_PRIVKEY)
		sig := Full_Transfer(10000000000, MASTER_PRIVKEY, w.PublicKey().String(), "devnet")
		fmt.Println(sig)
	}
}
	
