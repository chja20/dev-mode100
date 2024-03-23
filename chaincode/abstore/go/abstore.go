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
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ABstore Chaincode implementation
type ABstore struct {
	contractapi.Contract
}

func (t *ABstore) Init(ctx contractapi.TransactionContextInterface, aCert string, Aval string, bCert string, Bval string, cCert string, Cval string) error {
    fmt.Println("ABstore Init")
    var err error
    // Initialize the chaincode
    fmt.Printf("Aval = %s, Bval = %s, Cval = %s\n", Aval, Bval, Cval)
    // Write the state to the ledger
    err = ctx.GetStub().PutState(aCert, []byte(Aval))
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(bCert, []byte(Bval))
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(cCert, []byte(Cval))
    if err != nil {
        return err
    }

    return nil
}

// Transaction makes payment of X units from A to B
func (t *ABstore) Invoke(ctx contractapi.TransactionContextInterface, aCert, bCert, cCert string) error {
    var err error
    var Aval string
    var Bval string
    var Cval string
    // Get the state from the ledger
    // TODO: will be nice to have a GetAllState call to ledger
    Avalbytes, err := ctx.GetStub().GetState(aCert)
    if err != nil {
        return fmt.Errorf("Failed to get state: %v", err)
    }
    if Avalbytes == nil {
        return fmt.Errorf("Entity not found")
    }
    Aval = string(Avalbytes)

    Bvalbytes, err := ctx.GetStub().GetState(bCert)
    if err != nil {
        return fmt.Errorf("Failed to get state: %v", err)
    }
    if Bvalbytes == nil {
        return fmt.Errorf("Entity not found")
    }
    Bval = string(Bvalbytes)

    Cvalbytes, err := ctx.GetStub().GetState(cCert)
    if err != nil {
        return fmt.Errorf("Failed to get state: %v", err)
    }
    if Cvalbytes == nil {
        return fmt.Errorf("Entity not found")
    }
    Cval = string(Cvalbytes)
    
    // Perform the execution
    // Aval = Aval - X
    // Bval = Bval + X - ( X / 10 )
    // Cval = Cval + ( X / 10 )
    // fmt.Printf("Aval = %d, Bval = %d, Cval = %d\n", Aval, Bval, Cval)

    // Write the state back to the ledger
    err = ctx.GetStub().PutState(aCert, []byte(Aval))
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(bCert, []byte(Bval))
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(cCert, []byte(Cval))
    if err != nil {
        return err
    }

    return nil
}

// Delete  an entity from state
// func (t *ABstore) Delete(ctx contractapi.TransactionContextInterface, A string) error {

// 	// Delete the key from the state in ledger
// 	err := ctx.GetStub().DelState(A)
// 	if err != nil {
// 		return fmt.Errorf("Failed to delete state")
// 	}

// 	return nil
// }

// Query callback representing the query of a chaincode
func (t *ABstore) Query(ctx contractapi.TransactionContextInterface, aCert string) (string, error) {
	var err error
	// Get the state from the ledger
	Avalbytes, err := ctx.GetStub().GetState(aCert)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + aCert + "\"}"
		return "", errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + aCert + "\"}"
		return "", errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + aCert + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return string(Avalbytes), nil
}

// func (t *ABstore) GetAllQuery(ctx contractapi.TransactionContextInterface) ([]string, error) {
//     resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
//     if err != nil {
//         return nil, err
//     }
//     defer resultsIterator.Close()
//     var wallet []string
//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, err
//         }
//         jsonResp := "{\"Name\":\"" + string(queryResponse.Key) + "\",\"Amount\":\"" + string(queryResponse.Value) + "\"}"
//         wallet = append(wallet, jsonResp)
//     }
//     return wallet, nil
// }

func main() {
	cc, err := contractapi.NewChaincode(new(ABstore))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting ABstore chaincode: %s", err)
	}
}