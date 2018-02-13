/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main


import (
	"fmt"
	"strings"
	"encoding/json"
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var logger = shim.NewLogger("example_cc0")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type mystruct struct {
 ModelNumber string `json:"ModelNumber"`
 SerialNumber string `json:"SerialNumber"`
 CP string `json:"CP"`
 NP string `json:"NP"`
 CFlag string `json:"CFlag"`
 OrderNo string `json:"OrderNo"`
 M_R_D string `json:"M_R_D"`
 Location string `json:"Location"`
 Date_time string `json:"Date_time"`
}

type Products []struct { 
  ModelNumber string `json:"ModelNumber"` 
  SerialNumber string `json:"SerialNumber"` 
  CP string `json:"CP"` 
  NP string `json:"NP"` 
  CFlag string `json:"CFlag"` 
  OrderNo string `json:"OrderNo"` 
  M_R_D string `json:"M_R_D"`
  Location string `json:"Location"`
  Date_time string `json:"Date_time"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### example_cc0 Init ###########")

	_, args := stub.GetFunctionAndParameters()
	var key string


	// Initialize the chaincode
	
	logger.Info("before assiging value")
	
		mystr :=mystruct{}
		mystr.ModelNumber = args[0]
		mystr.SerialNumber = args[1]
		mystr.CP = args[2]
		mystr.NP = args[3]
		mystr.CFlag = args[4]
		mystr.OrderNo = args[5]
		mystr.M_R_D = args[6]
		mystr.Location = args[7]
		mystr.Date_time = args[8]
		
		
		
		key = args[0] + "-" + args[1]
		
		
		jsonAsBytes, _ := json.Marshal(mystr)

		
   
    fmt.Println(mystr)
	fmt.Println(jsonAsBytes)
	fmt.Println(string(jsonAsBytes))

	err := stub.PutState(key, jsonAsBytes)
	if err != nil {
	fmt.Println("Error in putState",err)
	return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if function == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if function == "move" {
		// Deletes an entity from its state
		logger.Infof("MOVE")
		return t.move(stub, args)
	}
		
	if function == "read_everything"{   //read everything, 
		return read_everything(stub)
	}
	
	if function == "getHistory"{   //getHistory , 
	logger.Infof("getHistory")
		return getHistory(stub,args)
	}
	if function == "queryByCP"{   //queryByCP, 
	logger.Infof("queryByCP")
		return queryByCP(stub,args)
	}
	if function == "bulkInsert"{   //bulkInsert, 
	logger.Infof("bulkInsert")
		return bulkInsert(stub,args)
	}	
	
	if function == "update"{   //bulkInsert, 
	logger.Infof("update")
		return update(stub,args)
	}	
	
	
	

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', read_everything, getHistory, queryByCP, update or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', read_everything, getHistory, queryByCP, update or 'move'. But got: %v", args[0]))
}

func bulkInsert(stub shim.ChaincodeStubInterface, args []string) pb.Response {

fmt.Println("************************ bulkInsert  ************************* ")
fmt.Printf("- args value: %s\n", args)

var buffer bytes.Buffer


buffer.WriteString("[")
resp := strings.Join(args,",")
buffer.WriteString(resp)
buffer.WriteString("]") 

resps := buffer.String()
fmt.Printf("- resps value: %s\n", resps)

modelNumber := &Products{}
	_ = json.Unmarshal([]byte(resps), modelNumber)
	for _, value := range *modelNumber {
		fmt.Println(value.ModelNumber)
		fmt.Println(value.SerialNumber)
		fmt.Println(value.CP)
		fmt.Println(value)
		
		jsonAsBytes, _ := json.Marshal(value)
				
				key := value.ModelNumber + "-" +value.SerialNumber

	
				fmt.Println(jsonAsBytes)
				fmt.Println(string(jsonAsBytes))

				err := stub.PutState(key, jsonAsBytes)
				if err != nil {
				fmt.Println("Error in putState",err)
				return shim.Error(err.Error())
				}
}

fmt.Println("************************ bulkInsert ends  ************************* ")


	return shim.Success(nil);

}


// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "KDL-123"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

fmt.Println("************************ READ HISTORY  ************************* ")
fmt.Printf("- args value: %s\n", args)
fmt.Printf("- len(args) value: %s\n", len(args))

	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   mystruct   `json:"value"`
	}
	var history []AuditHistory;
	var mys mystruct

		
	mode1serial := strings.Join(args,"")
	fmt.Printf("- start getHistoryForModel_SerialNumber Combination: %s\n", mode1serial)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(mode1serial)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
fmt.Printf("resultsIterator.HasNext()\n%s", resultsIterator.HasNext())

	for resultsIterator.HasNext() {
	fmt.Printf("Inside loop\n%s")
		historyData, err := resultsIterator.Next()
		fmt.Printf("Inside loop : \n%s", historyData)
		
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
			fmt.Printf("Inside loop historyData.TxId : \n%s", historyData.TxId )
		json.Unmarshal(historyData.Value, &mys)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //product has been deleted
			var emptymys mystruct
			tx.Value = emptymys                 //copy nil product
		} else {
			json.Unmarshal(historyData.Value, &mys) //un stringify it aka JSON.parse()
			tx.Value = mys                      //copy product over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForProduct returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	fmt.Printf("- getHistoryForProduct returning historyAsBytes:\n%s", historyAsBytes)
	return shim.Success(historyAsBytes)
}


func read_everything(stub shim.ChaincodeStubInterface) pb.Response {

fmt.Println("************************ READ EVERYTHING ************************* ")
	type Everything struct {
		Values   []mystruct   `json:"mystructs"`
		
	}
	var everything Everything

	// ---- Get All Values ---- //
	resultIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("resultIterator.HasNext() - ", resultIterator.HasNext())
	
	defer resultIterator.Close()
	
	fmt.Println("resultIterator.HasNext() - ", resultIterator.HasNext())

	for resultIterator.HasNext() {
		aKeyValue, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on id - ", queryKeyAsStr)
		var mys mystruct
		json.Unmarshal(queryValAsBytes, &mys)                   //un stringify it aka JSON.parse()
		everything.Values = append(everything.Values, mys)  //add this to the list
	}
	fmt.Println("order array - ", everything.Values)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	
	fmt.Println("order everything - ", everything)
	fmt.Println("order everythingAsBytes - ", everythingAsBytes)
	return shim.Success(everythingAsBytes)
}

// update only when value is there, else update as counterfeit

func  update(stub shim.ChaincodeStubInterface, args []string) pb.Response {

fmt.Println("************************ UPDATE ************************* ")
	
		var key string

				mystr :=mystruct{}
				mystr.ModelNumber = args[0]
				mystr.SerialNumber = args[1]
				mystr.CP = args[2]
				mystr.NP = args[3]
				mystr.CFlag = args[4]
				mystr.OrderNo = args[5]
				mystr.M_R_D = args[6]
				mystr.Location = args[7]
				mystr.Date_time = args[8]
				
				
		jsonAsBytes, _ := json.Marshal(mystr)
				
		key = args[0] + "-" +args[1]
		
		Exists, err := stub.GetState(key)
		
		
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to process for " + key + "\"}"
		fmt.Println("Error after get state", err)
		return shim.Error(jsonResp)
	}

	if Exists == nil {
	jsonResp1 := "{\"Message\":\"Counter-fiet product.\"}"
	fmt.Println("Counter-fiet product.")
	jsonAsBytes1, _ := json.Marshal(jsonResp1)
		return shim.Success(jsonAsBytes1);
		
	} else {
	
	//Update as the value is found
	fmt.Println("Model-Serial match found")
				fmt.Println(mystr)
				fmt.Println(jsonAsBytes)
				fmt.Println(string(jsonAsBytes))

				err := stub.PutState(key, jsonAsBytes)
				if err != nil {
				fmt.Println("Error in putState",err)
				return shim.Error(err.Error())
		}
		jsonResp2 := "{\"Message\":\"Successfully Updated.\"}"
		jsonAsBytes, _ := json.Marshal(jsonResp2)
		fmt.Println("All Fine")
		return shim.Success(jsonAsBytes);
	}
			
}



func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
		var key string

				mystr :=mystruct{}
				mystr.ModelNumber = args[0]
				mystr.SerialNumber = args[1]
				mystr.CP = args[2]
				mystr.NP = args[3]
				mystr.CFlag = args[4]
				mystr.OrderNo = args[5]
				mystr.M_R_D = args[6]
				mystr.Location = args[7]
				mystr.Date_time = args[8]
				
		jsonAsBytes, _ := json.Marshal(mystr)
				
		key = args[0] + "-" +args[1]

				fmt.Println(mystr)
				fmt.Println(jsonAsBytes)
				fmt.Println(string(jsonAsBytes))

				err := stub.PutState(key, jsonAsBytes)
				if err != nil {
				fmt.Println("Error in putState",err)
				return shim.Error(err.Error())
				}
        return shim.Success(nil);
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}


func  queryByCP(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	

	cp := strings.Join(args,"")

	queryString := fmt.Sprintf("{\"selector\":{\"CP\":\"%s\"}}", cp)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

fmt.Println("************************ QUERY ************************* ")

fmt.Printf("- args value: %s\n", args)
fmt.Printf("- len(args) value: %s\n", len(args))

	var A string // Entities
	var err error

	
  A = strings.Join(args,"")
	

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"No data found for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	//jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	jsonResp :=  string(Avalbytes)
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	logger.Infof("Hey I'm Main Method..")

	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
