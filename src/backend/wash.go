package main

import (
	"amm/common"
	"amm/raydium"
	"amm/solana_methods"
	"amm/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gagliardetto/solana-go"
)

/*
wallets := []solana.Wallet{
		{PrivateKey: solana.PrivateKey{8, 22, 119, 198, 58, 1, 217, 249, 134, 121, 151, 215, 23, 57, 26, 31, 42, 151, 197, 13, 199, 54, 198, 9, 86, 242, 162, 233, 8, 239, 104, 146, 222, 244, 215, 196, 165, 253, 93, 233, 200, 240, 187, 83, 7, 138, 70, 148, 228, 162, 81, 31, 39, 203, 20, 227, 255, 155, 240, 198, 206, 12, 11, 122}},          //G1L3rMsvEjvPUY8iGmxSm9wJkNWNxDApXonjHw8YQTBf
		{PrivateKey: solana.PrivateKey{106, 44, 177, 134, 249, 176, 138, 15, 83, 143, 120, 198, 208, 183, 153, 94, 232, 238, 116, 126, 17, 125, 40, 218, 158, 210, 83, 161, 178, 224, 117, 114, 190, 198, 90, 163, 79, 224, 149, 247, 243, 101, 143, 173, 52, 100, 121, 246, 230, 118, 9, 208, 84, 251, 58, 105, 72, 147, 95, 49, 22, 69, 5, 112}},  //DqhtVv6nR6umivMdEKV4NDLJ9PjCSLC35od4FCAmqP4B
		{PrivateKey: solana.PrivateKey{110, 36, 180, 222, 59, 212, 17, 101, 172, 172, 231, 239, 193, 50, 156, 69, 196, 42, 227, 92, 100, 141, 17, 108, 121, 47, 179, 153, 94, 163, 231, 107, 205, 68, 153, 29, 11, 152, 255, 201, 129, 35, 253, 229, 224, 172, 26, 142, 45, 190, 30, 151, 244, 199, 93, 46, 180, 222, 87, 122, 153, 189, 189, 170}}, //EpHFJPTxVGrafvPALCGnCRgSbsAjstrB2viZM2A6rnoK
		{PrivateKey: solana.PrivateKey{105, 1, 247, 76, 185, 11, 13, 205, 150, 177, 28, 211, 58, 11, 26, 212, 178, 201, 182, 79, 93, 122, 155, 1, 239, 80, 114, 16, 136, 67, 74, 118, 90, 21, 206, 124, 112, 55, 30, 217, 60, 83, 47, 226, 91, 150, 255, 35, 28, 236, 30, 165, 233, 113, 184, 101, 231, 98, 139, 206, 5, 9, 87, 169}},               //74eymC8f2mwXr8dpPDJkLBgWhCqu9g23UHhwD1oR8Dov
		{PrivateKey: solana.PrivateKey{164, 95, 220, 52, 39, 37, 54, 15, 172, 125, 202, 133, 128, 198, 217, 193, 192, 135, 50, 243, 100, 134, 150, 238, 15, 145, 247, 182, 247, 194, 85, 180, 193, 206, 21, 74, 13, 236, 10, 110, 165, 110, 189, 100, 217, 41, 254, 238, 120, 60, 151, 121, 173, 226, 89, 237, 70, 24, 148, 34, 222, 33, 51, 248}},  //E3XwuZow5ygb1M28p3Z44FfA75P573SJ4Q4HuPrTR36F
		{PrivateKey: solana.PrivateKey{185, 152, 213, 100, 129, 36, 9, 148, 254, 52, 106, 6, 70, 123, 184, 68, 41, 252, 8, 26, 76, 103, 45, 233, 50, 182, 145, 217, 231, 162, 206, 234, 218, 78, 205, 75, 194, 144, 69, 148, 172, 67, 133, 150, 192, 100, 42, 106, 57, 4, 169, 12, 200, 224, 152, 1, 189, 164, 138, 144, 162, 11, 50, 109}},         //FhBa7fArATzJZnvs2mHUFvEfVW2j1Q3KPvMmE44H9i6g
		{PrivateKey: solana.PrivateKey{95, 134, 200, 176, 194, 50, 198, 196, 244, 68, 83, 216, 120, 124, 131, 61, 140, 235, 140, 25, 14, 45, 81, 187, 130, 241, 209, 50, 32, 63, 185, 54, 77, 253, 252, 254, 137, 8, 63, 206, 205, 160, 153, 246, 120, 188, 89, 227, 247, 44, 17, 183, 62, 115, 217, 20, 67, 245, 249, 217, 221, 6, 147, 216}},      //6FT2ERVVXhELckMhJotrPoBJnjXkHjGGQwpJqzobsBC3
		{PrivateKey: solana.PrivateKey{222, 251, 55, 38, 143, 17, 132, 140, 160, 151, 135, 7, 249, 41, 168, 65, 92, 52, 171, 129, 251, 210, 228, 119, 235, 147, 229, 111, 158, 44, 159, 21, 233, 26, 39, 55, 65, 24, 9, 120, 206, 87, 49, 128, 29, 129, 101, 69, 151, 78, 82, 93, 230, 12, 191, 244, 210, 244, 159, 253, 51, 255, 31, 170}},         //Ggw8Akn6nwh3Qf6LwaVtGMZnwr9eueNa5tT5m35w6tdo
		{PrivateKey: solana.PrivateKey{216, 223, 115, 152, 63, 65, 135, 117, 245, 144, 7, 80, 184, 168, 232, 133, 164, 120, 1, 20, 16, 133, 193, 163, 39, 254, 192, 8, 159, 104, 119, 70, 126, 207, 58, 3, 232, 125, 235, 204, 91, 222, 94, 60, 222, 78, 96, 182, 145, 225, 73, 219, 2, 192, 64, 117, 18, 32, 95, 156, 6, 192, 223, 204}},           //9Y1eFLa19iRJG5rrGuCtF6XZdgko5Q52XJKevDNczS3V
		{PrivateKey: solana.PrivateKey{99, 33, 233, 239, 202, 26, 58, 220, 243, 88, 249, 210, 163, 139, 44, 161, 172, 238, 172, 177, 75, 87, 233, 24, 35, 176, 183, 141, 87, 234, 226, 200, 223, 139, 198, 201, 65, 15, 4, 123, 212, 14, 124, 114, 6, 78, 89, 36, 82, 75, 158, 21, 84, 33, 214, 84, 237, 252, 15, 70, 235, 98, 11, 131}},            //G3dY4wQ8qJaAM18cxr3q1e57J88aRq9DjNt1mwNEKH5U
	}
*/

func main() {
	pump()
	//washTrade(0)

	/*for _, w := range WALLETS_VOLUME_PHASE {
		balance := GetVaultBalance(w.PublicKey().String(), "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG", "devnet")

		raydium.FullRaydiumDevSwap(uint64(balance*1e9), w.PrivateKey.String(), false, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")

	}*/
}

func Wallets() {
	var wallets [][]byte
	for i := 0; i < 100; i++ {
		wallet := solana.NewWallet()
		jsonBytes, _ := json.Marshal(wallet.PrivateKey[:])
		wallets = append(wallets, jsonBytes)
	}

	for _, w := range wallets {
		fmt.Printf("{PrivateKey: solana.PrivateKey{%v}},\n", w)
	}
}

func GetAmountToStabilize(current_q_price, previous_q_price float64, amount uint64, QtoB string) uint64 {
	percent_change := utils.CalculatePricePercentageChange(current_q_price, previous_q_price)
	if percent_change > 5 || percent_change < -5 {

		current_q_supply, current_b_supply := common.GetSupply(ARMCHE_LP)
		target_q_supply, _ := CalculateTargetSupply(previous_q_price, current_q_supply, current_b_supply)

		amount = uint64(target_q_supply - current_q_supply)

		if ShouldBuy(percent_change) {
			QtoB = "true"
		} else {
			QtoB = "false"
		}
	} else {
		amount = 0
	}

	return amount
}

func MarketMake() {
	target_volume_Phase_1 := float64(1000000)
	var trading_volume float64 = 0

	washTrade(trading_volume, target_volume_Phase_1)
	pump()

}

func sellOrBuy() bool {
	var AtoB bool

	sell_or_buy := rand.Intn(2)
	if sell_or_buy == 0 { //buy
		fmt.Println("Swapping Q (selling SOL for ARM)")
		AtoB = true
	} else if sell_or_buy == 1 { //sell
		fmt.Println("Swapping Q (selling SOL for ARM)")
		AtoB = true
	} else {
		fmt.Println("Swapping B (selling ARM for SOL)")
		AtoB = false
	}
	return AtoB
}

type Bbalance struct {
	Q float64
	B float64
}

func initBalances(wallets []solana.Wallet) map[string]Bbalance {
	var balances = make(map[string]Bbalance)

	for _, w := range wallets {
		pubkey := w.PublicKey()
		q_bal := float64(solana_methods.GetBalance(pubkey.String(), "devnet").Result.Value) / 1e9
		b_bal := GetVaultBalance(w.PublicKey().String(), "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG", "devnet")
		balances[w.PrivateKey.String()] = Bbalance{Q: q_bal, B: b_bal}
	}

	return balances
}

func pump() {
	wallets := WALLETS_VOLUME_PHASE
	balances := initBalances(wallets)

	var trading_volume float64 = 0

	/*
		for w, bal := range balances {
			fmt.Printf("%s, Q: %f, B: %f\n", solana.MustPrivateKeyFromBase58(w).PublicKey().String(), bal.Q, bal.B)
		}*/

	for {
		wallet := wallets[rand.Intn(len(wallets))]

		var SOLtoARM bool
		if balances[wallet.PrivateKey.String()].B == 0 {

			SOLtoARM = true
		} else {
			SOLtoARM = sellOrBuy()
		}

		wallet_q_bal := balances[wallet.PrivateKey.String()].Q
		wallet_b_bal := balances[wallet.PrivateKey.String()].B

		var amount uint64
		trading_volume, amount = pumpDetermineTradeAmount(SOLtoARM, wallet_q_bal, wallet_b_bal, trading_volume)
		fmt.Printf("TRADING VOL: %f\n", trading_volume)

		fmt.Printf("%s, Q: %f, B: %f, SOLtoARM: %t, Amount: %f\n", wallet.PublicKey().String(), wallet_q_bal, wallet_b_bal, SOLtoARM, float64(amount)/1e9)

		raydium.FullRaydiumDevSwap(amount, wallet.PrivateKey.String(), SOLtoARM, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")

		time.Sleep(time.Duration(rand.Intn(61)) * time.Second)
	}

}

func pumpDetermineTradeAmount(SOLtoARM bool, maxSOL, maxTOKEN, trading_vol float64) (float64, uint64) {
	var amount float64
	last_q_supply, last_b_supply, err := common.GetSupp("5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q", "devnet", "v4")
	if err != nil {
		log.Fatalf("error getting supply: %v\n", err)

	}

	var max float64
	if SOLtoARM {
		max = maxSOL
	} else {
		max = maxTOKEN
	}

	solprice := solana_methods.GetSOLTokenPrice()
	fmt.Printf("sol price: %f\n", solprice)

	if SOLtoARM {

		amount = rand.Float64() * max / 2

		trading_vol += amount * solprice
	} else {

		tokenPerSOL := last_q_supply / last_b_supply
		//maxTokens := maxSOL * tokenPerSOL
		//amount = rand.Float64() * maxTokens

		amount = rand.Float64() * max / 2
		fmt.Printf("tokenPerSOL: %f\n", tokenPerSOL)

		trading_vol += amount * tokenPerSOL * solprice
	}
	fmt.Printf("amount: %f\n", amount)

	return trading_vol, uint64(amount * 1e9)
}

func buy(pkey string, amount float64) {
	raydium.FullRaydiumDevSwap(uint64(amount*1e9), pkey, true, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")
}

func sell(pkey string, amount float64) {
	raydium.FullRaydiumDevSwap(uint64(amount*1e9), pkey, false, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")
}

func washTrade(trading_volume, target_volume float64) {
	wallets := WALLETS_VOLUME_PHASE

	var q_price float64
	var b_price float64

	//private_key := MASTER_PRIVKEY
	var failed_wallets [][]string

	q_supply, b_supply := common.GetSupply(ARMCHE_LP)
	//k := q_supply * b_supply
	for i := 0; i < len(wallets); i++ {
		fmt.Println("Trading volume:", trading_volume)
		if trading_volume >= target_volume {
			return
		}
		if i == len(wallets)-1 {
			i = 0
		}
		w := wallets[i]
		private_key := w.PrivateKey.String()
		q_price = q_supply / b_supply
		b_price = b_supply / q_supply
		amount_tokenQ_to_swap := determineTradeAmount(q_supply) //Token A
		q_balance := float64(solana_methods.GetBalance(w.PublicKey().String(), "devnet").Result.Value) / 1e9
		b_balance := GetVaultBalance(w.PublicKey().String(), "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG", "devnet")

		fmt.Printf("q_balance: %f\n", q_balance)
		fmt.Printf("b_balance: %f\n", b_balance)
		fmt.Printf("q_price: %f\n", q_price)
		fmt.Printf("b_price: %f\n", b_price)
		fmt.Printf("q supply: %f\n", q_supply)
		fmt.Printf("b supply: %f\n", b_supply)
		fmt.Printf("first amount: %f (%d)\n", amount_tokenQ_to_swap, uint64(amount_tokenQ_to_swap*1e9))

		trading_volume += (amount_tokenQ_to_swap * solana_methods.GetSOLTokenPrice()) * 2

		new_q_supply := q_supply + amount_tokenQ_to_swap
		new_b_supply := (q_supply * b_supply) / new_q_supply
		amount_tokenB_to_receive := b_supply - new_b_supply

		//amount_tokenB_to_receive := amount_tokenQ_to_swap * b_price //Token B
		fmt.Printf("Amount Q to swap: %f\n", amount_tokenQ_to_swap)
		err := raydium.FullRaydiumDevSwap(uint64(amount_tokenQ_to_swap*1e9), private_key, true, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")
		if err != nil {
			failed_wallets = append(failed_wallets, []string{w.PublicKey().String(), w.PrivateKey.String()})
			continue
		}

		time.Sleep(7000 * time.Millisecond)
		q_supply = new_q_supply
		b_supply = new_b_supply
		q_price = q_supply / b_supply
		b_price = b_supply / q_supply
		q_balance = float64(solana_methods.GetBalance(w.PublicKey().String(), "devnet").Result.Value) / 1e9
		b_balance = GetVaultBalance(w.PublicKey().String(), "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG", "devnet")

		fmt.Printf("q_balance: %f\n", q_balance)
		fmt.Printf("b_balance: %f\n", b_balance)
		fmt.Printf("q_price_after: %f\n", q_price)
		fmt.Printf("b_price_after: %f\n", b_price)
		fmt.Printf("q_supply_after: %f\n", q_supply)
		fmt.Printf("b_supply_after: %f\n", b_supply)

		amount_tokenB_to_swap_back := amount_tokenB_to_receive //Token A
		trading_volume += amount_tokenQ_to_swap

		new_b_supply = b_supply + amount_tokenB_to_swap_back
		new_q_supply = (q_supply * b_supply) / new_b_supply
		amount_tokenQ_to_receive_back := q_supply - new_q_supply

		fmt.Printf("second amount: %f (%d)\n", amount_tokenB_to_swap_back, uint64(amount_tokenB_to_swap_back*10e8))

		fmt.Sprintf("amount token q to receive back%f\n", amount_tokenQ_to_receive_back)
		fmt.Printf("Amount B to swap: %f\n", amount_tokenB_to_swap_back)
		err = raydium.FullRaydiumDevSwap(uint64(amount_tokenB_to_swap_back*1e9), private_key, false, "5zseWHNxWp1b3aChp96jsLGgHz3ZEVDynDzniSm2ki3q")
		if err != nil {
			failed_wallets = append(failed_wallets, []string{w.PublicKey().String(), w.PrivateKey.String()})
			continue
		}
		time.Sleep(7000 * time.Millisecond)

		q_supply = new_q_supply
		b_supply = new_b_supply

		fmt.Println()

	}

	json_fail, _ := json.MarshalIndent(failed_wallets, " ", "  ")

	fmt.Println(string(json_fail))
}
func AmountOut(amountIn, reserveIn, reserveOut, fee float64) float64 {
	amountInAfterFee := amountIn * (1 - fee)
	return (reserveOut) / (reserveIn + amountInAfterFee)
}
func AmountIn(amountOut, reserveIn, reserveOut, fee float64) float64 {
	if amountOut >= reserveOut {
		panic("amountOut too large")
	}
	return (reserveIn * amountOut) / (reserveOut - amountOut) // round UP
}
func waitForBalanceChange(token_account string, expected_amount float64) {

	balance := solana_methods.GetTokenAccountBalance(token_account, "devnet").Result.Value.UiAmount
	initial_balance := balance

	fmt.Printf("initial balance: %f\n", initial_balance)
	expected_balance := balance + expected_amount
	for {
		balance = solana_methods.GetTokenAccountBalance(token_account, "devnet").Result.Value.UiAmount
		fmt.Printf("balance: %f\n", balance)
		fmt.Printf("expected_balance: %f\n", expected_balance)

		if balance != initial_balance && utils.CalculatePricePercentageChange(balance, initial_balance) < 5 {
			break
		}

		time.Sleep(250 * time.Millisecond)
	}
}

func determineTradeAmount(q float64) float64 {
	var percent_of_supply float64
	for {
		percent_of_supply = float64(rand.Intn(6))

		if percent_of_supply == 0 {
			continue
		} else {
			break
		}
	}

	return utils.GetPercentOf(percent_of_supply, q)
}
