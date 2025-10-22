package raydium

import (
	"amm/solana_methods"
	sol_methods "amm/solana_methods"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
)

const ProgramATA = "ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"
const ProgramSystem = "11111111111111111111111111111111"
const ProgramToken = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
const TokenWSOL = "So11111111111111111111111111111111111111112"
const Program2022Token = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"

func FullRaydiumDevSwap(amount uint64, signer_pkey string, SOLtoARM bool, lp string) error {
	var total_Instructions []solana.Instruction

	signer := solana.MustPrivateKeyFromBase58(signer_pkey)

	var a string
	var b string

	if SOLtoARM {
		a = "So11111111111111111111111111111111111111112"
		b = "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG"
	} else {
		a = "4NqyhaZwvs54822ML3n814XzBreHi4xVG5hQ8jPV63nG"
		b = "So11111111111111111111111111111111111111112"
	}

	userSourceTokenAcc, _, _ := solana.FindAssociatedTokenAddress(signer.PublicKey(), solana.MustPublicKeyFromBase58(a))
	userDestTokenAcc, _, _ := solana.FindAssociatedTokenAddress(signer.PublicKey(), solana.MustPublicKeyFromBase58(b))

	fmt.Println("addy:", solana.MustPrivateKeyFromBase58(signer_pkey).PublicKey())
	fmt.Println("pkey:", signer_pkey)
	fmt.Println("userSourceTokenAcc:", userSourceTokenAcc)
	fmt.Println("userDestTokenAcc:", userDestTokenAcc)

	/*if len(solana_methods.GetAccountInfo(userSourceTokenAcc, "devnet").Result.Value.Data[0]) == 0 {
		log.Fatal("Error ohoh:", userSourceTokenAcc)
	}
	if len(solana_methods.GetAccountInfo(userDestTokenAcc, "devnet").Result.Value.Data[0]) == 0 {
		log.Fatal("Error ohoh:", userDestTokenAcc)
	}*/

	var ata_Creation_Instruction []solana.Instruction
	var ata_account solana.PublicKey

	if SOLtoARM {
		ata_account = userSourceTokenAcc

		ata_Creation_Instruction = CreateATADev(signer.PublicKey().String(), ata_account.String(), amount) //5FY1Nvfzhkn8YajwHZPFzgrpQRdfxGRyj8x82NU6SWGC
		total_Instructions = append(total_Instructions, ata_Creation_Instruction...)

	} else {
		ata_account = userDestTokenAcc

		//ata_Creation_Instruction = CreateATADev(signer.PublicKey().String(), ata_account.String(), amount) //5FY1Nvfzhkn8YajwHZPFzgrpQRdfxGRyj8x82NU6SWGC
		//total_Instructions = append(total_Instructions, ata_Creation_Instruction...)
	}

	swaap := InnerDevSwapInstructions(amount, signer, userSourceTokenAcc, userDestTokenAcc, a, b, lp, SOLtoARM)
	total_Instructions = append(total_Instructions, swaap...)

	saaw := ThirdDevSwapInstruction(signer.PublicKey().String(), ata_account.String()) //5FY1Nvfzhkn8YajwHZPFzgrpQRdfxGRyj8x82NU6SWGC
	total_Instructions = append(total_Instructions, saaw...)

	/*for i, ix := range total_Instructions {
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
	}*/

	siga, err := sol_methods.SignAndSendTransaction(total_Instructions, signer.PublicKey(), signer, "devnet")
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	fmt.Println("swap res:", siga)

	return nil
}

/*
00002d31010000000 0d8cf2deb24000000

008096980000000000 047b207912000000
  8096980000000000 8096980000000000
008096980000000000 570d7d7b12000000
00002d310100000000 d8cf2deb24000000

78855437444
10000000
*/

func CreateATADev(signer_public_key, ATA_account string, amount uint64) []solana.Instruction {

	raydiumProgramID := solana.MustPublicKeyFromBase58("DRaybByLpbUL57LJARs3j8BitTxVfzBg351EaMr5UTCd")

	instrData := BuildRaydiumDevSwapData(uint8(5), amount, false)

	//fmt.Println(hex.EncodeToString(instrData))
	//fmt.Println("base64 swap data:", base64.StdEncoding.EncodeToString(instrData))

	accounts := []*solana.AccountMeta{

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(signer_public_key), true, true),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ATA_account), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(TokenWSOL), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramToken), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramATA), false, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramSystem), false, false),
	}
	instructions := []solana.Instruction{}
	instructions = append(instructions, solana.NewInstruction(raydiumProgramID, accounts, instrData))

	return instructions
}

func InnerDevSwapInstructions(amountIn uint64, signer solana.PrivateKey, srcTokenAcc, destTokenAcc solana.PublicKey, a, b, lp string, SOLtoARM bool) []solana.Instruction {

	lp_info, err := solana_methods.DeserializeRayV4(solana_methods.GetAccountInfo(lp, "devnet").Result.Value.Data[0])
	if err != nil {
		log.Fatalf("error getting lp info: %v", err)
	}

	raydiumProgramID := solana.MustPublicKeyFromBase58("DRaybByLpbUL57LJARs3j8BitTxVfzBg351EaMr5UTCd")

	instrData := BuildRaydiumDevSwapData(0, amountIn, true)

	//fmt.Println(hex.EncodeToString(instrData))
	//fmt.Println("base64 swap data:", base64.StdEncoding.EncodeToString(instrData))

	fmt.Println("SOLtoARM?:", SOLtoARM)

	accounts := []*solana.AccountMeta{

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramToken), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(Program2022Token), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramATA), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramSystem), false, false),

		solana.NewAccountMeta(signer.PublicKey(), true, true),

		solana.NewAccountMeta(srcTokenAcc, true, false),
		solana.NewAccountMeta(destTokenAcc, true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58("DRaya7Kj3aMWQSy19kSjvmuwq9docCHofyP9kanQGaav"), false, false),
		solana.NewAccountMeta(destTokenAcc, true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(a), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(b), true, false), //token mint

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58("AhSQwzpJMskCCQ3GkWDZSA69xS5qK5s7jLejhK6sgYf1"), false, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp_info.EncodedInfo.QuoteVault), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp_info.EncodedInfo.BaseVault), true, false),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(lp), true, false),
	}
	instructions := []solana.Instruction{}
	instructions = append(instructions, solana.NewInstruction(raydiumProgramID, accounts, instrData))

	return instructions
}

func BuildRaydiumDevSwapData(discrim uint8, amountIn uint64, swap bool) []byte {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, discrim); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, amountIn); err != nil {
		log.Fatalf("error writing to buf while decoding swap data: %v", err)
	}
	if swap {
		if err := binary.Write(buf, binary.LittleEndian, uint64(0)); err != nil {
			log.Fatalf("error writing to buf while decoding swap data: %v", err)
		}
	}

	return buf.Bytes()
}

func ThirdDevSwapInstruction(signer_public_key, ata string) []solana.Instruction {
	raydiumProgramID := solana.MustPublicKeyFromBase58("DRaybByLpbUL57LJARs3j8BitTxVfzBg351EaMr5UTCd")

	instrData := []byte{6}

	//fmt.Println(hex.EncodeToString(instrData))
	//fmt.Println("base64 swap data:", base64.StdEncoding.EncodeToString(instrData))

	accounts := []*solana.AccountMeta{

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(signer_public_key), true, true),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ata), true, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(signer_public_key), true, true),

		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramToken), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramATA), false, false),
		solana.NewAccountMeta(solana.MustPublicKeyFromBase58(ProgramSystem), false, false),
	}
	instructions := []solana.Instruction{}
	instructions = append(instructions, solana.NewInstruction(raydiumProgramID, accounts, instrData))

	return instructions
}
