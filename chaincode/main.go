package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	traceContract := new(TraceContract)

	chaincode, err := contractapi.NewChaincode(traceContract)
	if err != nil {
		log.Panicf("Error creating woodcarving trace chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting woodcarving trace chaincode: %v", err)
	}
}
