package main

import (
	"coursework/internal/selectionAlgorithms"
	"coursework/internal/selectionTypes"
	"fmt"
)

type Utxo struct {
	value  int
	txHash string
	vout   int
}

type Address struct {
	address string
	utxos   []Utxo
}

// i.e. "virtual bytes" required to spend a utxo you own from this address
// see https://bitcoin.stackexchange.com/questions/89385/is-there-a-difference-between-bytes-and-virtual-bytes-vbytes

const (
	DummySpendSize  = 100
	DummyOutputSize = 100
)

func AddressToOutputGroups(addresss []Address) []selectionTypes.OutputGroup {
	ans := make([]selectionTypes.OutputGroup, 0)
	for _, address := range addresss {
		if len(address.utxos) > 0 {
			inputCoins := make([]selectionTypes.InputCoin, 0)
			for _, utxo := range address.utxos {
				inputCoins = append(inputCoins, selectionTypes.InputCoin{
					TxHash:     utxo.txHash,
					Vout:       utxo.vout,
					Value:      utxo.value,
					InputBytes: DummySpendSize,
				})
			}
			temp := selectionTypes.OutputGroup{Address: address.address}
			for i := 0; i < len(inputCoins); i++ {
				temp.InsertCoin(&inputCoins[i])
			}
			ans = append(ans, temp)
		}
	}
	return ans
}

func main() {
	senderUtxos := []Utxo{{
		value:  10000,
		txHash: "340ad7bcf7dfc408fda32e2251fcb7fcbcf022b24daa6cbf11f202170f7748c2",
		vout:   0,
	}, {
		value:  20000,
		txHash: "abdfeabc232a5273545bef856f2edbafd1affbb22d1f8056c150f533426c0b7e",
		vout:   0,
	}, {
		value:  30000,
		txHash: "f25e6fca6236c34ea6eec7ead11e669255bd7a003a1d04463dd81a0eb2cd84df",
		vout:   1,
	},
	}

	senderAddresss := []Address{{"n4VQ5YdHf7hLQ2gWQYYrcxoE5B7nWuDFNF", senderUtxos[:1]},
		{"2N3oefVeg6stiTb5Kh3ozCSkaqmx91FDbsm", senderUtxos[1:]}}

	utxoPool := AddressToOutputGroups(senderAddresss)
	/*Get base transaction size independent of inputs and outputs
	For most simple transactions this will always be 10 bytes
	i.e. Version (4b) + TxOut count (1b) + TxIn count (1b) + Lock time (4b) */
	txBaseBytes := 10
	//recipientAddress := Address{address: "n4VQ5YdHf7hLQ2gWQYYrcxoE5B7nWuDFNF", utxos: make([]Utxo, 0)}
	txOutputBytes := DummyOutputSize

	notInputSizeInBytes := txBaseBytes + txOutputBytes

	//senderChangeAdress := Address{address: "mtXWDB6k5yC5v7TcwKZHB89SUp85yCKshy", utxos: make([]Utxo, 0)}
	changeOutputSizeInBytes := DummyOutputSize
	changeSpendSizeInBytes := DummySpendSize

	var shortTermFee, longTermFee int
	shortTermFee, longTermFee = 10, 10
	for i := 0; i < len(utxoPool); i++ {
		utxoPool[i].SetFee(shortTermFee, longTermFee)
	}
	targetValue := 25000
	result := selectionAlgorithms.SelectCoins(selectionTypes.CoinSelectionParams{
		UtxoPool:                utxoPool,
		TargetValue:             targetValue,
		ShortTermFeeRate:        shortTermFee,
		LongTermFeeRate:         longTermFee,
		ChangeOutputSizeInBytes: changeOutputSizeInBytes,
		ChangeSpendSizeInBytes:  changeSpendSizeInBytes,
		NotInputSizeInBytes:     notInputSizeInBytes,
	})
	fmt.Printf("Target value: %d\n", targetValue)
	for _, i := range result.Outputs {
		fmt.Println(i)
	}
	fmt.Printf("Value: %d\n", result.Value)
	fmt.Printf("Fee: %d\n", result.Fee)
}
