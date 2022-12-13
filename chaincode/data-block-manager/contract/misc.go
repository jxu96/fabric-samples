package contract

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// return the name of the implicit organisation-specific private collection
func implicitPrivateDataCollection(mspId string) string {
	return fmt.Sprintf("_implicit_org_%s", mspId)
}

// e.g. "x509::CN=org1admin,OU=admin,O=Hyperledger,ST=North Carolina,C=US::CN=ca.org1.example.com,O=org1.example.com,L=Durham,ST=North Carolina,C=US"
func getClientID(ctx contractapi.TransactionContextInterface) (string, error) {
	bs, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve ID from the client : %v", err)
	}
	id, err := base64.StdEncoding.DecodeString(bs)
	if err != nil {
		return "", fmt.Errorf("Failed to decode ID of the client : %v", err)
	}
	return string(id), nil
}

func getClientMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve MSP ID from the client : %v", err)
	}
	return id, nil
}

func getPeerMSPID() (string, error) {
	id, err := shim.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve MSP ID from the peer : %v", err)
	}

	return id, nil
}

func getTxTransient(ctx contractapi.TransactionContextInterface) (map[string][]byte, error) {
	map_, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve transient : %v", err)
	}

	return map_, nil
}

func getTxTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve timestamp from current transaction context : %v", err)
	}

	return ts.AsTime().String(), nil
}

func requireIdenticalMSPID(ctx contractapi.TransactionContextInterface, mspID *string) error {
	clientMSPID, err := getClientMSPID(ctx)
	if err != nil {
		return err
	}

	if mspID != nil {
		*mspID = clientMSPID
	}

	peerMSPID, err := getPeerMSPID()
	if err != nil {
		return err
	}

	if clientMSPID != peerMSPID {
		return fmt.Errorf(`Client from Org "%s" has no access to service from Org "%s".`, clientMSPID, peerMSPID)
	}

	return nil
}

func requireCertification(ctx contractapi.TransactionContextInterface, cond any) error {
	/**
		TODO
		verify certification to grant access to enable:
		- register, update ledger
		- access private collection
	**/
	return nil
}
