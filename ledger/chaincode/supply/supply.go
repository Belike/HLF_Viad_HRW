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
type MAX struct {
	ID string `json:"id"`
}

func getStatus(input string) (status string) {
	if input == "ProduzentA" {
		status = "Produziert"
	} else if input == "LieferantA" {
		status = "Lieferung"
	} else if input == "VerkaufA" {
		status = "Geliefert"
	}
	return
}

var id = 0

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	max := MAX{}
	max.ID = "0"
	maxAsBytes, _ := json.Marshal(max)
	APIstub.PutState("ID", maxAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryAllGoods" {
		return s.queryAllGoods(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "changeOwnerByID" {
		return s.changeOwnerByID(APIstub, args)
	} else if function == "createGood" {
		return s.createGood(APIstub, args)
	} else if function == "changeOwnerByBarcode" {
		return s.changeOwnerByBarcode(APIstub, args)
	} else if function == "getMaxID" {
		return s.getMaxID(APIstub)
	}

	return shim.Error("Smart Contract Funktion nicht gefunden.")
}
func (s *SmartContract) getMaxID(APIstub shim.ChaincodeStubInterface) sc.Response {
	maxAsBytes, _ := APIstub.GetState("ID")
	max := MAX{}
	json.Unmarshal(maxAsBytes, &max)

	return shim.Success([]byte(max.ID))
}
func (s *SmartContract) changeOwnerByBarcode(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Parameteranzahl falsch. Erwarte 2")
	}
	maxAsBytes, _ := APIstub.GetState("ID")
	max := MAX{}
	json.Unmarshal(maxAsBytes, &max)
	id, _ = strconv.Atoi(max.ID)
	i := 0
	for i < id {
		goodsAsBytes, _ := APIstub.GetState(strconv.Itoa(i))
		good := Good{}
		json.Unmarshal(goodsAsBytes, &good)
		if good.Barcode == args[0] {
			good.Owner = args[1]
			good.Status = getStatus(args[1])
			goodsAsBytes, _ = json.Marshal(good)
			APIstub.PutState(strconv.Itoa(i), goodsAsBytes)
			break
		}
		i = i + 1
	}
	return shim.Success(nil)
}
func (s *SmartContract) createGood(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 6 {
		return shim.Error("Parameteranzahl falsch. Erwarte 6")
	}
	maxAsBytes, _ := APIstub.GetState("ID")
	max := MAX{}
	json.Unmarshal(maxAsBytes, &max)
	id, _ = strconv.Atoi(max.ID)
	var good = Good{Barcode: args[0], Beschreibung: args[1], Menge: args[2], Produzent: args[3], Status: args[4], Owner: args[5]}
	goodsAsBytes, _ := json.Marshal(good)
	APIstub.PutState(max.ID, goodsAsBytes)
	id = id + 1
	max.ID = strconv.Itoa(id)
	maxAsBytes, _ = json.Marshal(max)
	APIstub.PutState("ID", maxAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) changeOwnerByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	goodsAsBytes, _ := APIstub.GetState(args[0])
	good := Good{}
	json.Unmarshal(goodsAsBytes, &good)
	good.Owner = args[1]
	good.Status = getStatus(good.Owner)
	goodsAsBytes, _ = json.Marshal(good)

	APIstub.PutState(args[0], goodsAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	maxToCheckAsBytes, _ := APIstub.GetState("ID")
	maxToCheck := MAX{}
	json.Unmarshal(maxToCheckAsBytes, &maxToCheck)
	if maxToCheck.ID != "0" {
		return shim.Error("Ledger ist bereits initialisiert")
	}
	goods := []Good{
		Good{Barcode: "10000", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Produziert", Owner: "ProduzentA"},
		Good{Barcode: "10001", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Produziert", Owner: "ProduzentA"},
		Good{Barcode: "10002", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Produziert", Owner: "ProduzentA"},
		Good{Barcode: "10003", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Lieferung", Owner: "LieferantA"},
		Good{Barcode: "10004", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Lieferung", Owner: "LieferantA"},
		Good{Barcode: "10005", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Lieferung", Owner: "LieferantA"},
		Good{Barcode: "10006", Beschreibung: "Holzfigur", Menge: "100", Produzent: "ProduzentA", Status: "Geliefert", Owner: "VerkaufA"},
		Good{Barcode: "10007", Beschreibung: "Steinskulptur", Menge: "80", Produzent: "ProduzentA", Status: "Geliefert", Owner: "VerkaufA"},
		Good{Barcode: "10008", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Geliefert", Owner: "VerkaufA"},
		Good{Barcode: "10009", Beschreibung: "Plastikhaus", Menge: "10", Produzent: "ProduzentA", Status: "Geliefert", Owner: "VerkaufA"},
		Good{Barcode: "10010", Beschreibung: "Plastikhaus2", Menge: "5", Produzent: "ProduzentA", Status: "Geliefert", Owner: "VerkaufA"},
	}
	i := 0
	for i < len(goods) {
		goodsAsBytes, _ := json.Marshal(goods[i])
		APIstub.PutState(strconv.Itoa(id), goodsAsBytes)
		id = id + 1
		fmt.Println("Added ", goods[i])
		i = i + 1
	}
	var max = MAX{ID: strconv.Itoa(id)}
	maxAsBytes, _ := json.Marshal(max)
	APIstub.PutState("ID", maxAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) queryAllGoods(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Parameteranzahl falsch. Benötige mindestens 1")
	}
	startKey := "0"
	endKey := "999"

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
		buffer.WriteString("\n")
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
