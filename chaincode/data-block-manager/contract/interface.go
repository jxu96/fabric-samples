package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DataBlockLedger struct {
	contractapi.Contract
}

type DataBlockInterface interface {
	to_bytes() ([]byte, error)
	from_bytes(key string, bs []byte) error
}

type DataBlock struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	OwnerOrg    string `json:"owner_org"`
	Timestamp   string `json:"timestamp"`
}

type DataBlockPrivate struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`
}

type DataBlockInput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     []byte `json:"content"`
}

type DataBlockListInterface interface {
}

var public_data_collection = "publicDataBlockCollection"

func implicit_private_data_collection(mspId string) string {
	return fmt.Sprintf("_implicit_org_%s", mspId)
}

var InitialDataBlocks = []DataBlockInput{
	{
		ID:          "data1",
		Title:       "Test data suite 1",
		Description: "Empty dataset for testing",
		Content:     []byte(`{"train": {}, "test": {}}`),
	},
	{
		ID:          "data2",
		Title:       "Test data suite 2",
		Description: "Empty dataset for testing",
		Content:     []byte(`{"train": {}, "test": {}}`),
	},
}

func (l *DataBlockLedger) Register(ctx contractapi.TransactionContextInterface) error {
	client_id, err := get_client_id(ctx)
	if err != nil {
		return err
	}

	client_msp_id, err := get_client_msp_id(ctx)
	if err != nil {
		return err
	}

	peer_msp_id, err := get_peer_msp_id()
	if err != nil {
		return err
	}

	assert_client_matches_peer(client_msp_id, peer_msp_id)

	err = l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	transient, err := l.get_tx_transient(ctx)
	if err != nil {
		return err
	}

	data, ok := transient["data_block"]
	if !ok {
		return fmt.Errorf("Failed to retrive data block information from transient.")
	}

	timestamp, err := l.get_tx_timestamp(ctx)
	if err != nil {
		return err
	}

	block := new(DataBlock)
	block.Owner = client_id
	block.OwnerOrg = client_msp_id
	block.Timestamp = timestamp

	err = block.from_bytes("", data)
	if err != nil {
		return err
	}

	// register in the public collection
	err = l.create_data_block(ctx, public_data_collection, block.ID, block)
	if err != nil {
		return err
	}

	// register in the private collection
	block_private := new(DataBlockPrivate)
	err = block_private.from_bytes("", data)
	if err != nil {
		return err
	}

	err = l.create_data_block(ctx, implicit_private_data_collection(client_msp_id), block_private.ID, block_private)
	if err != nil {
		return err
	}

	return nil
}

func (l *DataBlockLedger) Remove(ctx contractapi.TransactionContextInterface) error {
	client_msp_id, err := get_client_msp_id(ctx)
	if err != nil {
		return err
	}

	peer_msp_id, err := get_peer_msp_id()
	if err != nil {
		return err
	}

	assert_client_matches_peer(client_msp_id, peer_msp_id)

	err = l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	transient, err := l.get_tx_transient(ctx)
	if err != nil {
		return err
	}

	data, ok := transient["data_block"]
	if !ok {
		return fmt.Errorf("Failed to retrieve data block information from transient.")
	}

	block := new(DataBlock)
	err = block.from_bytes("", data)
	if err != nil {
		return err
	}

	err = l.delete_data_block(ctx, public_data_collection, block.ID)
	if err != nil {
		return err
	}

	err = l.delete_data_block(ctx, implicit_private_data_collection(client_msp_id), block.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l *DataBlockLedger) Query(ctx contractapi.TransactionContextInterface, key string) (*DataBlock, error) {
	block := new(DataBlock)
	err := l.read_data_block(ctx, public_data_collection, key, block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (l *DataBlockLedger) QueryPrivate(ctx contractapi.TransactionContextInterface, key string) (*DataBlockPrivate, error) {
	err := l.require_certification_read(ctx)
	if err != nil {
		return nil, err
	}

	block := new(DataBlock)
	err = l.read_data_block(ctx, public_data_collection, key, block)
	if err != nil {
		return nil, err
	}

	peer_msp_id, err := get_peer_msp_id()
	if err != nil {
		return nil, err
	}

	if peer_msp_id != block.OwnerOrg {
		return nil, fmt.Errorf("Org mismatch.")
	}

	block_private := new(DataBlockPrivate)
	err = l.read_data_block(ctx, implicit_private_data_collection(block.OwnerOrg), key, block_private)
	if err != nil {
		return nil, err
	}

	return block_private, nil
}

func (l *DataBlockLedger) QueryByRange(ctx contractapi.TransactionContextInterface, start string, end string, max int) ([]*DataBlock, error) {
	map_, err := l.get_by_range(ctx, public_data_collection, start, end, max)
	if err != nil {
		return nil, err
	}

	blocks := []*DataBlock{}
	for key, val := range map_ {
		block := new(DataBlock)
		err := block.from_bytes(key, val)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}
