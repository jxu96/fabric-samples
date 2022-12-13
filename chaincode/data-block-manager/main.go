package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/jxu96/fabric-samples/chaincode/data-block-manager/contract"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&contract.DatasetMetadataLedger{})
	if err != nil {
		log.Panicf("Error creating data-block-manager chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting data-block-manager chaincode: %v", err)
	}
}
