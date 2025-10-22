package main

import (
	"fmt"
	"log"
	r "math/rand"

	"amm/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	"context"

	"os"

	"github.com/joho/godotenv"
)

type Maester struct {
	*rpc.Client
}

func maisn() {
	fmt.Println("Runnin...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	maestro()
}

//3enQJKSDb3doDF5xHXiLvWniAcKGJK36SNLQywkADCdQ
//Ed4tBvkhuvAmx4mHd27RAHVM3kicTUDGwPXf7GgcJ9v2
//4qKv4xDBLsRHTS2R1X2o3e46WU3gxCG93uCutwDSRmuo
///7Hr2mQZzsDGV6DHRMUbiX4aivuMjPfw2jGz7RGxXJ9ZV
//FErwirmBnBj1aKrdetkeku9xWYYzQhbmztmor7tFx7MS

type Brain struct {
	RPC *rpc.Client
	WS  *ws.Client
	Dev solana.PrivateKey
}

func maestro() {

	rpcClient := rpc.New(rpc.DevNet_RPC)
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		log.Panicf("Failed to connect to ws: %v", err)
	}
	dev_pk, err := solana.PrivateKeyFromBase58(os.Getenv("DEV_PRIVATE_KEY"))
	if err != nil {
		log.Panicf("error decoding string to privatekey: %v", err)
	}

	brain := Brain{
		RPC: rpcClient,
		WS:  wsClient,
		Dev: dev_pk,
	}

	erro := distribute_funds(brain, 100, 10)
	if erro != nil {
		log.Fatalf("There was an error while distributing the funds: %v", erro)
	}
}

func distribute_funds(brain Brain, amount_wallets_1 int, percentage_of_balance_to_transfer uint64) error {

	storage := make(map[string][]map[string]interface{})
	sigs := make(map[string][]solana.Signature)

	var initial []utils.Storage
	var secondary []utils.Storage
	var tertiary []utils.Storage

	utils.UnloadStorage("../global.json", storage)
	utils.UnloadSigs("../sigs.json", sigs)

	balance := func(pkey solana.PublicKey) uint64 {
		balance, err := brain.RPC.GetBalance(
			context.TODO(),
			pkey,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			log.Fatalf("error getting balance: %v", err)
		}
		return balance.Value
	}

	fmt.Println("Initial phase start...")

	init_balance := balance(brain.Dev.PublicKey())
	initial_tokenAmounts := random_amounts((init_balance/100)*percentage_of_balance_to_transfer, amount_wallets_1, 3000000+2000000) //0.003 rent + 0.002 fee

	if len(initial_tokenAmounts) != amount_wallets_1 {
		return fmt.Errorf("the number of 'Amounts' does not match the 'Amount of wallets'\n")
	}

	for i := 1; i < amount_wallets_1; i++ {
		account := solana.NewWallet()

		fee := getFeeTransfer(brain, Transaction{
			Sender:    brain.Dev,
			Recipient: account.PublicKey(),
			Amount:    initial_tokenAmounts[i],
		})

		sig, err := transaction(brain, Transaction{
			Sender:    brain.Dev,
			Recipient: account.PublicKey(),
			Amount:    initial_tokenAmounts[i] - fee,
		})
		if err != nil {
			return fmt.Errorf("an error appeared while sending a transaction: %v\n", err)
		}
		fmt.Printf("Signature: %s\n", sig)

		new := map[string]interface{}{
			"publickey":  account.PublicKey(),
			"privatekey": account.PrivateKey,
		}
		storage["initial"] = append(storage["initial"], new)
		sigs["sigs"] = append(sigs["sigs"], sig)

		initial = append(initial, utils.Storage{account.PublicKey(), account.PrivateKey})
	}
	fmt.Println("\nInitial phase end")

	//secondary phase
	fmt.Println("Secondary phase start...")
	for _, pair := range initial {
		account := solana.NewWallet()
		balance := balance(pair.Private_key.PublicKey())

		fee := getFeeTransfer(brain, Transaction{
			Sender:    pair.Private_key,
			Recipient: account.PublicKey(),
			Amount:    balance,
		})

		sig, err := transaction(brain, Transaction{
			Sender:    pair.Private_key,
			Recipient: account.PublicKey(),
			Amount:    balance - fee,
		})
		if err != nil {
			return fmt.Errorf("An error appeared while sending a transaction: %v", err)
		}
		fmt.Printf("Signature: %s\n", sig)

		new := map[string]interface{}{
			"publickey":  account.PublicKey(),
			"privatekey": account.PrivateKey,
		}
		storage["secondary"] = append(storage["secondary"], new)
		secondary = append(secondary, utils.Storage{account.PublicKey(), account.PrivateKey})
		sigs["sigs"] = append(sigs["sigs"], sig)
	}
	fmt.Println("\nSecondary phase end")

	//tertiary phase
	fmt.Println("Tertiary phase start...")
	for _, pair := range secondary {
		account := solana.NewWallet()
		balance := balance(pair.Private_key.PublicKey())

		fee := getFeeTransfer(brain, Transaction{
			Sender:    pair.Private_key,
			Recipient: account.PublicKey(),
			Amount:    balance,
		})

		sig, err := transaction(brain, Transaction{
			Sender:    pair.Private_key,
			Recipient: account.PublicKey(),
			Amount:    balance - fee,
		})
		if err != nil {
			return fmt.Errorf("An error appeared while sending a transaction: %v", err)
		}
		fmt.Printf("Signature: %s\n", sig)

		new := map[string]interface{}{
			"publickey":  account.PublicKey(),
			"privatekey": account.PrivateKey,
		}
		storage["tertiary"] = append(storage["tertiary"], new)
		tertiary = append(tertiary, utils.Storage{account.PublicKey(), account.PrivateKey})
		sigs["sigs"] = append(sigs["sigs"], sig)
	}
	fmt.Println("\nTertiary phase end")

	defer func() {
		err0 := utils.LoadStorage("global.json", storage)
		err1 := utils.LoadSigs("sigs.json", sigs)

		if err0 != nil && err1 != nil {
			log.Panicf("There was an error loading new vars:\n%v\n%v", err0, err1)
		}

		fmt.Println("Function end...")

	}()
	return nil
}

type Transaction struct {
	Sender    solana.PrivateKey
	Recipient solana.PublicKey
	Amount    uint64 //without accounting for fees so result will be: amount - fees
}

func transaction(brain Brain, info Transaction) (solana.Signature, error) {

	recentBlock, err := brain.RPC.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return [64]byte{}, fmt.Errorf("failed to get recent block: %v", err)
	}

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				info.Amount,
				info.Sender.PublicKey(),
				info.Recipient,
			).Build(),
		},
		recentBlock.Value.Blockhash,
		solana.TransactionPayer(info.Sender.PublicKey()),
	)

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if info.Sender.PublicKey().Equals(key) {
			return &info.Sender
		}
		return nil
	})
	if err != nil {
		return solana.SignatureFromBytes(nil), fmt.Errorf("Failed to sign transaction:", err)
	}
	sig, err := brain.RPC.SendTransaction(
		context.Background(),
		tx,
	)
	if err != nil {
		return solana.SignatureFromBytes(nil), fmt.Errorf("Error sending transaction: %v", err)
	}

	/*
		sig, err := confirm.SendAndConfirmTransaction(
			context.TODO(),
			brain.RPC,
			brain.WS,
			tx,
		)
		if err != nil {
			return [64]byte{}, fmt.Errorf("failed to send transaction: %v", err)
		}
	*/

	return sig, nil
}

func random_amounts(x_to_divide uint64, total_parts int, min_z uint64) []uint64 {
	if x_to_divide < uint64(total_parts)*min_z {
		return []uint64{} // Not enough to distribute at least `z` per part
	}
	remaining := x_to_divide - uint64(total_parts)*min_z // Remaining amount after ensuring min_z per part
	parts := make([]uint64, total_parts)

	rawParts := make([]uint64, total_parts) // Store initial random values
	var sumRaw uint64

	for i := 0; i < total_parts; i++ {
		rawParts[i] = uint64(r.Intn(100) + 1) // Generate random values (1 to 100)
		sumRaw += rawParts[i]
	}

	var sumParts uint64
	for i := 0; i < total_parts; i++ {
		parts[i] = (rawParts[i] * remaining) / sumRaw
		sumParts += parts[i]
		//fmt.Printf("%d ", parts[i])
	}

	diff := remaining - sumParts
	for i := uint64(0); i < diff; i++ {
		parts[i]++ // Add 1 to some parts to balance rounding error
	}

	for i := 0; i < total_parts; i++ {
		parts[i] += min_z
	}

	return parts

}

func getFeeTransfer(brain Brain, info Transaction) uint64 {
	recentBlock, err := brain.RPC.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("failed to get recent block: %v", err)
	}

	tx, _ := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				info.Amount,
				info.Sender.PublicKey(),
				info.Recipient,
			).Build(),
		},
		recentBlock.Value.Blockhash,
		solana.TransactionPayer(info.Sender.PublicKey()),
	)

	fees, _ := brain.RPC.GetFeeForMessage(context.TODO(), tx.Message.ToBase64(), rpc.CommitmentConfirmed)

	accountInfo, err := brain.RPC.GetAccountInfo(context.TODO(), info.Sender.PublicKey())
	if err != nil {
		log.Fatalf("Failed to fetch account info: %v", err)
	}
	accountSize := uint64(len(accountInfo.Value.Data.GetBinary()))

	rentExemptBalance, err := brain.RPC.GetMinimumBalanceForRentExemption(context.TODO(), accountSize, rpc.CommitmentConfirmed)
	if err != nil {
		log.Fatalf("Failed to calculate rent: %v", err)
	}

	return (*fees.Value) + rentExemptBalance
}

//TODO:
//in the distribute funds function, there is an error that says:
//"Transaction simulation failed: Transaction results in an account (1) with insufficient funds for rent"
//to fix this,
//before the main transaction, I would send a small transaction that contains the minimum rent
//or
//I would need to add the "minimum rent exempt" to the random amounts (random_amount = random + rent) which is not working i think

/*

 */
