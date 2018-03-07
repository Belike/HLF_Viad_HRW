package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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

var statBarcode = 10

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryGoodsByOwner" {
		return s.queryAllGoods(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "changeOwner" {
		return s.changeOwner(APIstub, args)
	} else if function == "createGood" {
		return s.createGood(APIstub, args)
	}

	return shim.Error("Smart Contract Funktion nicht gefunden.")
}
func (s *SmartContract) createGood(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Parameteranzahl falsch. Erwarte 5")
	}
	localBarcode := ""
	if statBarcode >= 10 && statBarcode < 100 {
		localBarcode = "100" + strconv.Itoa(statBarcode)
	} else if statBarcode >= 100 && statBarcode < 1000 {
		localBarcode = "10" + strconv.Itoa(statBarcode)
	}
	var good = Good{Barcode: localBarcode, Beschreibung: args[0], Menge: args[1], Produzent: args[2], Status: args[3], Owner: args[4]}
	goodsAsBytes, _ := json.Marshal(good)
	APIstub.PutState("GOOD"+strconv.Itoa(statBarcode), goodsAsBytes)
	statBarcode = statBarcode + 1
	return shim.Success(nil)
}
func (s *SmartContract) changeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	goods := []Good{
		Good{Barcode: "10000", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "10001", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "10002", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Prod", Owner: "ProduzentA"},
		Good{Barcode: "10003", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "10004", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "10005", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Lief", Owner: "LieferantA"},
		Good{Barcode: "10006", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
		Good{Barcode: "10007", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
		Good{Barcode: "10008", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
		Good{Barcode: "10009", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Done", Owner: "VerkaufA"},
	}
	i := 0
	for i < len(goods) {
		goodsAsBytes, _ := json.Marshal(goods[i])
		APIstub.PutState("GOOD"+strconv.Itoa(i), goodsAsBytes)
		fmt.Println("Added ", goods[i])
		i = i + 1
	}
	return shim.Success(nil)
}
func (s *SmartContract) queryAllGoods(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Parameteranzahl falsch. Benötige mindestens 1")
	}
	startKey := "GOOD0"
	endKey := "GOOD999"

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
