package solana_methods

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func SignAndSendTransaction(total_Instructions []solana.Instruction, PayerPubkey solana.PublicKey, PayerPkey solana.PrivateKey, cluster string) (solana.Signature, error) {
	var rpc_url string = ""
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}

	//rpc_url = "https://api.devnet.solana.com"
	var recent string
	for {
		recent = GetLatestBlockhash(cluster).Result.Value.Blockhash

		if recent == "" || len(recent) == 0 {
			continue
		} else {
			break
		}
	}

	hash, err := solana.HashFromBase58(recent)
	if err != nil {
		log.Fatalf("error getting must hash azzzzeepe: %v: %s", err, recent)
	}

	tx, err := solana.NewTransaction(
		total_Instructions,
		hash,
		solana.TransactionPayer(PayerPubkey),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(PayerPkey.PublicKey()) {
				return &PayerPkey
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("error signing the tx: %v", err)
	}

	client := rpc.New(rpc_url)
	sig, err := client.SendTransaction(context.Background(), tx)

	if err != nil {
		fmt.Print(fmt.Errorf("err simulating tx: %v", err))
		return solana.Signature{}, err
	}

	return sig, nil
}

func IsTransactionStatus(tx_hash string, commitment rpc.CommitmentType) {

}

func SimTransaction(total_Instructions []solana.Instruction, PayerPubkey solana.PublicKey, PayerPkey solana.PrivateKey, cluster string) *rpc.SimulateTransactionResponse {
	var rpc_url string
	switch cluster {
	case "devnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_DEV")
	case "mainnet":
		rpc_url = os.Getenv("QUICKNODE_HTTP_SOLANA_MAINNET")
	}

	recent := GetLatestBlockhash(cluster)

	tx, err := solana.NewTransaction(
		total_Instructions,
		solana.MustHashFromBase58(recent.Result.Value.Blockhash),
		solana.TransactionPayer(PayerPubkey),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(PayerPkey.PublicKey()) {
				return &PayerPkey
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("error signing the tx: %v", err)
	}

	client := rpc.New(rpc_url)
	resp, err := client.SimulateTransactionWithOpts(context.Background(), tx, &rpc.SimulateTransactionOpts{
		Commitment: rpc.CommitmentConfirmed, // Match Solscan's commitment
		//SigVerify:  false,
	})

	if err != nil {
		log.Fatal("err simulating tx:", err)
	}

	return resp

}

func Transfer(lamports uint64, sender solana.PrivateKey, receiver solana.PublicKey) []solana.Instruction {
	transfInstr := []solana.Instruction{
		system.NewTransferInstruction(
			lamports,
			sender.PublicKey(),
			receiver,
		).Build(),
	}
	return transfInstr
}
