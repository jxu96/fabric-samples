package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type DataBlockPublic struct {
	ID          string `json:"ID"`
	Owner       string `json:"Owner"`
	OwnerID     string `json:"OwnerID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Timestamp   string `json:"Timestamp"`
}

type DataBlock struct {
	DataBlockPublic
	Endpoint string `json:"Endpoint"`
}

var TestBlocks = []DataBlock{
	{
		DataBlockPublic{
			ID:          "data1",
			Owner:       "UiS",
			OwnerID:     "UiSID",
			Title:       "Test data from UiS",
			Description: "This is a test data block",
		},
		"http://localhost:8000",
	},
	{
		DataBlockPublic{
			ID:          "data2",
			Owner:       "UiO",
			OwnerID:     "UiOID",
			Title:       "Test data from UiO",
			Description: "This is a test data block",
		},
		"http://localhost:8001",
	},
}

func (sc *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	for _, data := range TestBlocks {
		if err := sc.Puts(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SmartContract) Register(ctx contractapi.TransactionContextInterface, data_blocks []DataBlock) error {
	for _, data := range data_blocks {
		if existed, err := sc.Exists(ctx, data.ID); err != nil {
			return err
		} else if existed {
			return fmt.Errorf("Failed to register data block due to ID duplication: %s", data.ID)
		}

		if err := sc.Puts(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SmartContract) Update(ctx contractapi.TransactionContextInterface, data_blocks []DataBlock) error {
	for _, data := range data_blocks {
		if existed, err := sc.Exists(ctx, data.ID); err != nil {
			return err
		} else if !existed {
			return fmt.Errorf("Failed to update unexisting data block: %s", data.ID)
		}

		if err := sc.Puts(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SmartContract) Remove(ctx contractapi.TransactionContextInterface, ids []string) error {
	for _, id := range ids {
		if existed, err := sc.Exists(ctx, id); err != nil {
			return err
		} else if !existed {
			return fmt.Errorf("Failed to remove unexisting data block: %s", id)
		}

		if err := ctx.GetStub().DelState(id); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SmartContract) GetView(ctx contractapi.TransactionContextInterface) {

}

func (sc *SmartContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from ledger: %v", err)
	}
	return data != nil, nil
}

func (sc *SmartContract) Puts(ctx contractapi.TransactionContextInterface, data DataBlock) error {
	dataS, err := json.Marshal(data)
	ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(data.ID, dataS)
	if err != nil {
		return fmt.Errorf("Failed to init data block: %s.", data.ID)
	}

	return nil
}
