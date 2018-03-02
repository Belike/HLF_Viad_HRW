package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}
type Good struct {
	Barcode      string `json:"barcode"`
	Beschreibung string `json:"beschreibung"`
	Menge        string `json:"menge"`
	Produzent    string `json:"produzent"`
	Status       string `json:"status"`
	Owner        string `json:"owner"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryGoodsByOwner" {
		return s.queryGoodsByOwner(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	}

	return shim.Error("Smart Contract Funktion nicht gefunden.")
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	goods := []Good{
		Good{Barcode: "AAAA1", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "AAAA2", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "AAAA3", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "BBBB1", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "BBBB2", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "BBBB3", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "CCCC1", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
		Good{Barcode: "CCCC2", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
		Good{Barcode: "CCCC3", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
	}
	i := 0
	for i < len(goods) {
		goodsAsBytes, _ := json.Marshal(goods[i])
		APIstub.PutState(goods[i].Barcode, goodsAsBytes)
		fmt.Println("Added ", goods[i])
	}
	return shim.Success(nil)
}
func (s *SmartContract) queryGoodsByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Parameteranzahl falsch. Benötige mindestens 1")
	}
	startKey := ""
	endKey := ""

	if args[0] == "Produzent" {
		startKey = "AAAA1"
		endKey = "AAAA99"
	} else if args[0] == "Lieferant" {
		startKey = "BBBB1"
		endKey = "BBBB99"
	} else {
		//Verkäufer
		startKey = "CCCC1"
		endKey = "CCCC99"
	}
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
