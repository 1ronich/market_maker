package common

import (
	"fmt"
	"math"
	"time"

	sol_methods "amm/solana_methods"
)

func CalculatePrices(last_b_supply, last_q_supply float64) (float64, float64) {

	price_in_q := last_q_supply / last_b_supply
	price_in_b := last_b_supply / last_q_supply

	return price_in_q, price_in_b
}

func GetPrices(LP string) (float64, float64) {
	var last_b_supply float64
	var last_q_supply float64
	var err error
	for {

		last_q_supply, last_b_supply, err = GetSupp(LP, "devnet", "v4")
		if err != nil {
			fmt.Println(fmt.Errorf("error getting supply: %v\n", err))
			continue
		}

		if last_b_supply == 0 || last_q_supply == 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		} else {
			break
		}
	}
	current_q_price, current_b_price := CalculatePrices(last_b_supply, last_q_supply)

	return current_q_price, current_b_price
}

func GetSupply(LP string) (float64, float64) {
	var last_b_supply float64
	var last_q_supply float64
	var err error
	for {

		last_q_supply, last_b_supply, err = GetSupp(LP, "devnet", "v4")
		if err != nil {
			fmt.Println(fmt.Errorf("error getting supply: %v\n", err))
			continue
		}

		if last_b_supply == 0 || last_q_supply == 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		} else {
			break
		}
	}

	return last_q_supply, last_b_supply
}

func GetSupp(lp_mint, cluster, lp_type string) (float64, float64, error) { //3 rpc requests

	var b string //base vault
	var q string //quote vault
	tries := 0
	for {
		if tries == 5 {
			return 0, 0, fmt.Errorf("\nmax retries reached \n")
		}
		pool := sol_methods.GetAccountInfo(lp_mint, cluster)

		switch lp_type {
		case "v4":
			deserialized_pool, err := sol_methods.DeserializeRayV4(pool.Result.Value.Data[0])
			if err != nil {
				fmt.Println("error deserializing v4 acc:", err)
				continue
			}

			if deserialized_pool.EncodedInfo.QuoteMint == "So11111111111111111111111111111111111111112" || deserialized_pool.EncodedInfo.QuoteMint == "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" {
				b = deserialized_pool.EncodedInfo.BaseVault
				q = deserialized_pool.EncodedInfo.QuoteVault
			} else {
				b = deserialized_pool.EncodedInfo.QuoteVault
				q = deserialized_pool.EncodedInfo.BaseVault

			}

			baseSupply, _ := sol_methods.GetTokenSupply(b, cluster)
			quoteSupply, _ := sol_methods.GetTokenSupply(q, cluster)

			return quoteSupply, baseSupply, nil

		case "cpmm":
			deserialized_pool, err := sol_methods.DeserializeCPMMaccount(pool.Result.Value.Data[0])
			if err != nil {
				fmt.Println("error deserializing cpmm acc:", err)
				continue
			}

			if deserialized_pool.TokenMint0 == "So11111111111111111111111111111111111111112" || deserialized_pool.TokenMint0 == "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" {
				b = deserialized_pool.TokenVault1
				q = deserialized_pool.TokenVault0
			} else {
				b = deserialized_pool.TokenVault0
				q = deserialized_pool.TokenVault1
			}

			baseSupply, _ := sol_methods.GetTokenSupply(b, cluster)
			quoteSupply, _ := sol_methods.GetTokenSupply(q, cluster)

			return quoteSupply, baseSupply, nil

		}

		tries++

		return 0, 0, nil
	}

}

// Get necessary supply amounts of A & B to reach X price in the base token
func GetAmountsToMovePrice(lp_mint, cluster, lp_type string, percentage_change float64) (float64, float64, float64, float64, error) {
	current_q_supply, current_b_supply, err := GetSupp(lp_mint, cluster, lp_type)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	current_price := current_q_supply / current_b_supply

	new_b_supply, new_q_supply := SupplyNeededToReachTarget(current_b_supply, current_q_supply, (current_price * (1 + (percentage_change / 100))))
	//new_price := new_q_supply / new_b_supply

	fmt.Println("old_x:", current_b_supply)
	fmt.Println("old_y:", current_q_supply)
	fmt.Println("new_x:", new_b_supply)
	fmt.Println("new_y:", new_q_supply)

	//move_amount_b := new_b_supply - current_b_supply
	//move_amount_q := new_q_supply - current_q_supply

	return new_b_supply, new_q_supply, current_b_supply, current_q_supply, nil

}

func SupplyNeededToReachTarget(x, y, target_price float64) (float64, float64) {
	//x := a_supply // Current amount of token A in LP
	//y := b_supply // Current amount of token B in LP
	k := x * y //Constant Product

	newX := math.Sqrt(k / target_price)
	newY := float64(target_price) * newX
	newK := newX * newY

	if math.Floor(k*100)/100 != math.Floor(newK*100)/100 {
		panic(fmt.Sprintf("New constant value (newK) does not match old constant value: %f != %f", k, newK))
	}

	return newX, newY

}
