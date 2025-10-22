package raydium

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"

	sol_methods "amm/solana_methods"

	"amm/types"
)

type Maester struct {
	*rpc.Client
}

func Raydium_V4_Swap(pkey, cluster string) error {

	var amm_program string
	switch cluster {
	case "mainnet":
		amm_program = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
	case "devnet":
		amm_program = "DRaya7Kj3aMWQSy19kSjvmuwq9docCHofyP9kanQGaav"
	}

	LPmint := "97p3acpjRH9kZJ189fRHzVzmWiTSsSoGZeCwFNmHDQhf"
	//LPmintA := "ED5nyyWEzpPPiWimP8vYm7sD7TD3LAt3Q3gRTWHzPJBY"
	//LPmintB := "So11111111111111111111111111111111111111112"

	privateKey := solana.MustPrivateKeyFromBase58(pkey)

	var res types.DeserializedRayV4
	for {
		ress, err := sol_methods.DecodeRayV4TokenAccount(solana.MustPublicKeyFromBase58(LPmint), cluster)
		if err != nil {
			continue
		} else {
			res = ress
			break
		}

	}

	swap := types.SwapInfo{
		PayerPkey:   privateKey,             //solana.MustPrivateKeyFromBase58(os.Getenv("MAIN_SOL_PKEY")),
		PayerPubkey: privateKey.PublicKey(), //solana.MustPrivateKeyFromBase58(os.Getenv("MAIN_SOL_PKEY")).PublicKey(),
		Pool: types.PoolInfo{
			Address: solana.MustPublicKeyFromBase58(LPmint), //pool address
			LPMintA: solana.MustPublicKeyFromBase58(res.EncodedInfo.BaseMint),
			LPMintB: solana.MustPublicKeyFromBase58(res.EncodedInfo.QuoteMint),

			AmmProgram:    solana.MustPublicKeyFromBase58(amm_program),
			AmmPool:       solana.MustPublicKeyFromBase58(LPmint),
			AmmAuthority:  solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"), //getAMMAuthority(solana.MustPublicKeyFromBase58(LPmint)),
			AmmOpenOrders: solana.MustPublicKeyFromBase58(res.EncodedInfo.OpenOrders),
			AmmCoinVault:  solana.MustPublicKeyFromBase58(res.EncodedInfo.BaseVault),
			AmmPcVault:    solana.MustPublicKeyFromBase58(res.EncodedInfo.QuoteVault),

			MarketProgram:     solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketProgramId),
			Market:            solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId),
			MarketBids:        sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).Bids,
			MarketAsks:        sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).Asks,
			MarketEventQueue:  sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).EventQ,
			MarketCoinVault:   sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).CoinVault,
			MarketPcVault:     sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).PcVault,
			MarketVaultSigner: sol_methods.GetVaultSigner(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketProgramId), sol_methods.GetMarketData(solana.MustPublicKeyFromBase58(res.EncodedInfo.MarketId), cluster).VaultSignerNonce),
		},

		AmountIn:     uint64(50000),
		MinAmountOut: uint64(1),
		RPC:          os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET"),
	}

	fmt.Println("Starting to build instructions...")

	totalInstructions := Raydium_Swap(swap, cluster)

	for i, ix := range totalInstructions {
		fmt.Printf("Instruction #%d:\n", i+1)
		fmt.Printf("  ProgramID: %s\n", ix.ProgramID())
		fmt.Printf("  Accounts:\n")
		for j, acct := range ix.Accounts() {
			fmt.Printf("    %d: %s (Signer: %v, Writable: %v)\n", j, acct.PublicKey, acct.IsSigner, acct.IsWritable)
		}
		data, err := ix.Data()
		if err != nil {
			log.Fatalf("Could not get account data: %v\n", err)
		}
		fmt.Printf("  Data: %x\n", data)
		fmt.Println("-----")
	}

	//sol_methods.SimTransaction(totalInstructions, swap.PayerPubkey, swap.PayerPkey, "mainnet")
	siga, err := sol_methods.SignAndSendTransaction(totalInstructions, privateKey.PublicKey(), privateKey, cluster)

	fmt.Println("swap hash:", siga)
	return err

}

func Get_user_source_owner(source_owner, cluster string) solana.PublicKey {
	infa := sol_methods.GetAccountInfo(source_owner, cluster)
	raw, _ := base64.StdEncoding.DecodeString(infa.Result.Value.Data[0])
	jsora, _ := json.MarshalIndent(infa, " ", "  ")
	fmt.Println("RAW:", string(jsora))
	return solana.MustPublicKeyFromBase58(base58.Encode(raw[32:64]))
}

func Raydium_Swap(info types.SwapInfo, cluster string) []solana.Instruction {

	//fmt.Println("nub: ", sol_methods.GetTokenAccounsByOwner("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", "zGh48JtNHVBb5evgoZLXwgPD2Qu4MhkWdJLGDAupump", "PNLCQcVCD26aC7ZWgRyr5ptfaR7bBrWdTFgRWwu2tvF"))
	//fmt.Println("wsol: ", sol_methods.GetTokenAccounsByOwner("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", "So11111111111111111111111111111111111111112", "PNLCQcVCD26aC7ZWgRyr5ptfaR7bBrWdTFgRWwu2tvF"))

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintA) //a is wsol
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintB)

	//info.Pool.UserSourceOwner = Get_user_source_owner(userSourceTokenAcc.String())
	brain := types.Atastruct{
		PayerPrivateKey:    info.PayerPkey,
		LPmints:            [2]solana.PublicKey{info.Pool.LPMintA, info.Pool.LPMintB},
		PoolSourceToken:    info.Pool.UserTokenSource,
		PoolDestToken:      info.Pool.UserTokenDestination,
		UserSourceTokenAcc: userSourceTokenAcc,
		UserDestTokenAcc:   userDestTokenAcc,
	}

	var total_Instructions []solana.Instruction

	compLimitInstr := SetComputeLimit(400_000)
	total_Instructions = append(total_Instructions, compLimitInstr)

	resMap := CheckForATAs(brain, cluster)
	//this for range loop checks if there are ATAs missing, and if it needs to create one, it will create one
	var index int8 = 0
	for _, value := range resMap {

		if !value {
			ATAinstr := CreateNewATA(brain, brain.LPmints[index])
			total_Instructions = append(total_Instructions, ATAinstr)

			info.Pool.UserTokenSource = userSourceTokenAcc
			info.Pool.UserTokenDestination = userDestTokenAcc
			info.Pool.UserSourceOwner = info.PayerPubkey //Get_user_source_owner(userSourceTokenAcc.String())
		} else {
			info.Pool.UserTokenSource = userSourceTokenAcc
			info.Pool.UserTokenDestination = userDestTokenAcc
			info.Pool.UserSourceOwner = info.PayerPubkey //Get_user_source_owner(userSourceTokenAcc.String(), cluster)
		}
		index++
	}

	swapInstr := InnerSwapInstructions(info, cluster)

	jitoTipTransferInstruction := sol_methods.Transfer(1000000, info.PayerPkey, solana.MustPublicKeyFromBase58("Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY"))

	for _, inner := range swapInstr {
		total_Instructions = append(total_Instructions, inner)
	}
	for _, inner := range jitoTipTransferInstruction {
		total_Instructions = append(total_Instructions, inner)
	}
	//these for range loops are necessary because you cannot pass a [][]solana.instruction to create a 'tx'

	return total_Instructions

}

func BuildRaydiumSwapData(amountIn, minAmountOut uint64, discriminator uint8) []byte {
	swap := types.SimpleSwapInstruction{
		Discriminator:    discriminator,
		AmountIn:         amountIn,
		MinimumAmountOut: minAmountOut,
	}

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, swap.Discriminator); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, swap.AmountIn); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, swap.MinimumAmountOut); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}

	return buf.Bytes()
}

func InnerSwapInstructions(info types.SwapInfo, cluster string) []solana.Instruction {
	var rayProramId string
	switch cluster {
	case "mainnet":
		rayProramId = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
	case "devnet":
		rayProramId = "DRaya7Kj3aMWQSy19kSjvmuwq9docCHofyP9kanQGaav"
	}
	raydiumProgramID := solana.MustPublicKeyFromBase58(rayProramId)
	instrData := BuildRaydiumSwapData(info.AmountIn, info.MinAmountOut, 9)

	fmt.Println(hex.EncodeToString(instrData))
	fmt.Println("base64 swap data:", base64.StdEncoding.EncodeToString(instrData))

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintA)
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(info.PayerPubkey, info.Pool.LPMintB)

	accounts := []*solana.AccountMeta{

		solana.NewAccountMeta(solana.TokenProgramID, false, false),
		solana.NewAccountMeta(info.Pool.AmmPool, true, false),
		solana.NewAccountMeta(info.Pool.AmmAuthority, false, false),

		solana.NewAccountMeta(info.Pool.AmmOpenOrders, true, false),
		//solana.NewAccountMeta(solana.MustPublicKeyFromBase58("8PdLK3w34qVmQzXrS51mA7Qe3gxT1keBxzqH3iAvMHHX"), true, false),

		solana.NewAccountMeta(info.Pool.AmmCoinVault, true, false),
		solana.NewAccountMeta(info.Pool.AmmPcVault, true, false),

		solana.NewAccountMeta(info.Pool.MarketProgram, false, false),
		solana.NewAccountMeta(info.Pool.Market, true, false),
		solana.NewAccountMeta(info.Pool.MarketBids, true, false),
		solana.NewAccountMeta(info.Pool.MarketAsks, true, false),
		solana.NewAccountMeta(info.Pool.MarketEventQueue, true, false),

		solana.NewAccountMeta(info.Pool.MarketCoinVault, true, false),
		solana.NewAccountMeta(info.Pool.MarketPcVault, true, false),

		solana.NewAccountMeta(info.Pool.MarketVaultSigner, false, false),

		solana.NewAccountMeta(userSourceTokenAcc, true, false),
		solana.NewAccountMeta(userDestTokenAcc, true, false),

		solana.NewAccountMeta(info.PayerPubkey, true, true),
	}
	instructions := []solana.Instruction{}
	instructions = append(instructions, solana.NewInstruction(raydiumProgramID, accounts, instrData))

	return instructions
}

func CreateNewATA(brain types.Atastruct, TokenMintAddress solana.PublicKey) solana.Instruction {
	//this function checks if the ATA that was not found is either "PoolSourcetoken" or "PoolDestToken" and then based on the ATA
	//it appends to 'poolToken' the address from the 'brain' variable, and then a CreateATAinstruction is created
	//var poolToken solana.PublicKey

	/*if PoolToken == "UserSourceTokenAcc" {
		poolToken = brain.UserSourceTokenAcc
		fmt.Println("FOund UserSourceTokenAcc:", brain.UserSourceTokenAcc)
	} else if PoolToken == "UserDestTokenAcc" {
		poolToken = brain.UserDestTokenAcc
		fmt.Println("FOund UserSourceTokenAcc:", brain.UserSourceTokenAcc)

	}*/

	ATA := associatedtokenaccount.NewCreateInstruction(
		brain.PayerPrivateKey.PublicKey(), //brain.PayerPrivateKey.PublicKey(),
		brain.PayerPrivateKey.PublicKey(), //brain.PayerPrivateKey.PublicKey(),
		TokenMintAddress,
	).Build()
	fmt.Printf("ATA Payer (signer & funder): %s\n", brain.PayerPrivateKey.PublicKey())
	fmt.Printf("ATA Owner: %s\n", brain.PayerPrivateKey.PublicKey())
	fmt.Printf("Token Mint: %s\n", TokenMintAddress)

	return ATA
}

func CheckForATAs(brain types.Atastruct, cluster string) map[string]bool {
	//this function checks if there are any existing ATA accounts by calling the getAccountInfo method
	//and checking the length of the data returned; if length=0 --> no existing account, so I need to create an account
	//that means sending an instruction to the tx so it knows that it needs to create an ATA

	state := map[string]bool{
		"UserSourceTokenAcc": true,
		"UserDestTokenAcc":   true,
	}

	resp := sol_methods.GetAccountInfo(brain.UserSourceTokenAcc.String(), cluster)
	if resp.Result.Value.Data[0] == "" {
		state["UserSourceTokenAcc"] = false //ATA not existing
	}
	respa := sol_methods.GetAccountInfo(brain.UserDestTokenAcc.String(), cluster)
	if respa.Result.Value.Data[0] == "" {
		state["UserDestTokenAcc"] = false //ATA not existing
	}

	return state
}

func SetComputeLimit(amount uint32) *computebudget.Instruction {
	computBudget := computebudget.NewSetComputeUnitLimitInstruction(
		amount,
	).Build()
	fmt.Printf("ProgramID: %s\n", computBudget.ProgramID())
	byts, _ := computBudget.Data()
	fmt.Printf("Data: %x\n", byts)

	return computBudget
}
