package solana_methods

import (
	"amm/types"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gagliardetto/solana-go"
)

type TokenData struct {
	USDPrice       float64 `json:"usdPrice"`
	BlockID        int     `json:"blockId"`
	Decimals       int     `json:"decimals"`
	PriceChange24h float64 `json:"priceChange24h"`
}
type TokenPrices map[string]TokenData

func GetSOLTokenPrice() float64 {
	res, err := http.Get("https://lite-api.jup.ag/price/v3?ids=So11111111111111111111111111111111111111112")
	if err != nil {
		log.Fatalf("error gerring price: %v", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	var prices TokenPrices
	err = json.Unmarshal(body, &prices)
	if err != nil {
		log.Fatal(err)
	}

	return prices["So11111111111111111111111111111111111111112"].USDPrice
}

func GetTokenPrice(token_mint string) float64 {
	res, err := http.Get(fmt.Sprintf("https://lite-api.jup.ag/price/v3?ids=%s", token_mint))
	if err != nil {
		log.Fatalf("error gerring price: %v", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var prices TokenPrices
	err = json.Unmarshal(body, &prices)
	if err != nil {
		log.Fatal(err)
	}

	return prices[token_mint].USDPrice
}

func GetMarketData(marketId solana.PublicKey, cluster string) types.MarketState {
	response := GetAccountInfo(marketId.String(), cluster)

	res := DecodeOpenBookData(response.Result.Value.Data[0])

	return res
}
func DecodeOpenBookData(data string) types.MarketState {
	bin, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalf("error decoding openbook data")
	}

	var m types.MarketState
	buf := bytes.NewReader(bin[5:])

	// Adjust endianness if needed; Solana usually uses little endian
	if err := binary.Read(buf, binary.LittleEndian, &m); err != nil {
		panic(err)
	}
	return m
}

func Uint64ToLittleEndianBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

func GetVaultSigner(marketID, programID solana.PublicKey, nonce uint64) solana.PublicKey {
	seed := [][]byte{
		marketID.Bytes(),
		Uint64ToLittleEndianBytes(nonce),
	}
	pda, err := solana.CreateProgramAddress(seed, programID)
	if err != nil {
		log.Fatalf("error creating vault signer PDA: %v", err)
	}
	return pda
}

func GetAMMAuthority(ammID solana.PublicKey) solana.PublicKey {
	seed := [][]byte{
		[]byte("amm authority"),
		ammID.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seed, ammID)
	if err != nil {
		log.Fatalf("error getting amm auth: %v\n", err)
	}

	return pda
}

func GetTokenSupply(token_account_address, cluster string) (float64, uint8) {
	balance := GetTokenAccountBalance(
		token_account_address,
		cluster,
	)

	supply, _ := strconv.ParseFloat(balance.Result.Value.Amount, 64)
	formatted_a_supply := supply / math.Pow(10, float64(balance.Result.Value.Decimals))

	return formatted_a_supply, balance.Result.Value.Decimals
}

func GetQuoteTokenPriceInLP(lp_mint, cluster string, pool any) float64 {
	var base_token_mint string
	var base_token_account string
	var quote_token_account string

	var A_amount float64
	var B_amount float64

	fmt.Println(reflect.TypeOf(pool).String())
	switch reflect.TypeOf(pool).String() {
	case "types.DecodedCPMM":
		switch pool.(types.DecodedCPMM).TokenMint0 {
		case "So11111111111111111111111111111111111111112":

			base_token_account = pool.(types.DecodedCPMM).TokenVault0
			quote_token_account = pool.(types.DecodedCPMM).TokenVault1
			base_token_mint = pool.(types.DecodedCPMM).TokenMint0
		case "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v":
			base_token_account = pool.(types.DecodedCPMM).TokenVault0
			quote_token_account = pool.(types.DecodedCPMM).TokenVault1
			base_token_mint = pool.(types.DecodedCPMM).TokenMint0
		default:
			base_token_account = pool.(types.DecodedCPMM).TokenVault1
			quote_token_account = pool.(types.DecodedCPMM).TokenVault0
			base_token_mint = pool.(types.DecodedCPMM).TokenMint1
		}

		A_amount, _ = GetTokenSupply(base_token_account, cluster)
		B_amount, _ = GetTokenSupply(quote_token_account, cluster)
	case "types.DeserializedRayV4":
		switch pool.(types.DeserializedRayV4).EncodedInfo.BaseMint {
		case "So11111111111111111111111111111111111111112":
			fmt.Println("SALAM")
			fmt.Println(pool.(types.DeserializedRayV4).EncodedInfo.BaseMint)
			base_token_account = pool.(types.DeserializedRayV4).EncodedInfo.QuoteVault
			quote_token_account = pool.(types.DeserializedRayV4).EncodedInfo.BaseVault
			base_token_mint = pool.(types.DeserializedRayV4).EncodedInfo.QuoteMint
		case "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v":
			base_token_account = pool.(types.DeserializedRayV4).EncodedInfo.QuoteVault
			quote_token_account = pool.(types.DeserializedRayV4).EncodedInfo.BaseVault
			base_token_mint = pool.(types.DeserializedRayV4).EncodedInfo.QuoteMint
		default:
			fmt.Println("BYE")
			fmt.Println(pool.(types.DeserializedRayV4).EncodedInfo.BaseMint)
			base_token_account = pool.(types.DeserializedRayV4).EncodedInfo.BaseVault
			quote_token_account = pool.(types.DeserializedRayV4).EncodedInfo.QuoteVault
			base_token_mint = pool.(types.DeserializedRayV4).EncodedInfo.BaseMint
		}

		A_amount, _ = GetTokenSupply(base_token_account, cluster)
		B_amount, _ = GetTokenSupply(quote_token_account, cluster)

	}
	fmt.Println("base:", base_token_mint)
	fmt.Println("A_amount:", A_amount)
	fmt.Println("B_amount:", B_amount)

	tokenA_price_in_lp := (A_amount / B_amount) * GetTokenPrice(base_token_mint) // quoteTokenSupply / baseTokenSupply * GetTokenPrice(base token)

	return tokenA_price_in_lp
}

func GetAccountBalance() {

}
