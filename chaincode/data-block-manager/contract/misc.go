package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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
