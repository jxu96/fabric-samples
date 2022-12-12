package contract

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (block *DataBlock) to_bytes() ([]byte, error) {
	bs, err := json.Marshal(*block)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode data block to bytes: %s\n%v", block.ID, err)
	}

	return bs, nil
}

func (block *DataBlock) from_bytes(key string, bs []byte) error {
	err := json.Unmarshal(bs, block)
	if err != nil {
		return fmt.Errorf("Failed to parse data block: %s\n%v", key, err)
	}

	return nil
}

func (block *DataBlockPrivate) to_bytes() ([]byte, error) {
	bs, err := json.Marshal(*block)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode data block to bytes: %s\n%v", block.ID, err)
	}

	return bs, nil
}

func (block *DataBlockPrivate) from_bytes(key string, bs []byte) error {
	err := json.Unmarshal(bs, block)
	if err != nil {
		return fmt.Errorf("Failed to parse data block: %s\n%v", key, err)
	}

	return nil
}

func (l *DataBlockLedger) require_certification_read(ctx contractapi.TransactionContextInterface) error {
	/**
		TODO
		verify certification to grant read access to enable:
		- reveal private section of data blocks
	**/
	return nil
}

func (l *DataBlockLedger) require_certification_write(ctx contractapi.TransactionContextInterface) error {
	/**
		TODO
		verify certification to grant write access to enable:
		- register, update data blocks
		- reveal private section of data blocks
	**/
	return nil
}

func (l *DataBlockLedger) get_tx_transient(ctx contractapi.TransactionContextInterface) (map[string][]byte, error) {
	m, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient.\n%v", err)
	}
	return m, nil
}

func (l *DataBlockLedger) get_tx_sender(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetStub().GetCreator()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve sender identity from current transaction context.\n%v", err)
	}

	return string(id), nil
}

func (l *DataBlockLedger) get_tx_timestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("Failed to retrive timestamp from current transaction context.\n%v", err)
	}

	return ts.String(), nil
}

func get_client_id(ctx contractapi.TransactionContextInterface) (string, error) {
	bs, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrive ID from the client.\n%v", err)
	}
	id, err := base64.StdEncoding.DecodeString(bs)
	if err != nil {
		return "", fmt.Errorf("Failed to decode ID of the client.\n%v", err)
	}
	return string(id), nil
}

func get_client_msp_id(ctx contractapi.TransactionContextInterface) (string, error) {
	client, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve MSP ID from the client.\n%v", err)
	}
	return client, nil
}

func get_peer_msp_id() (string, error) {
	peer, err := shim.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve MSP ID from the peer.\n%v", err)
	}
	return peer, nil
}

func assert_client_matches_peer(client string, peer string) error {
	if client != peer {
		return fmt.Errorf("Client from org %s is not authorized to access private data from an org %s peer", client, peer)
	}
	return nil
}
