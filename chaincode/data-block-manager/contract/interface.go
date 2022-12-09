package contract

import (
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

var (
	CollectionDataBlock       string = "collectionDataBlock"
	CollectionDataBlockPublic string = "collectionDataBlockPublic"
)

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

func (l *DataBlockLedger) Init(ctx contractapi.TransactionContextInterface) error {
	return l.Register(ctx, InitialDataBlocks)
}

func (l *DataBlockLedger) Register(ctx contractapi.TransactionContextInterface, inputs []DataBlockInput) error {
	err := l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	sender, err := l.get_tx_sender(ctx)
	if err != nil {
		return err
	}

	timestamp, err := l.get_tx_timestamp(ctx)
	if err != nil {
		return err
	}

	preserved := DataBlockPreserved{
		OwnerID:   sender,
		Timestamp: timestamp,
	}

	for _, input := range inputs {
		block := DataBlock{
			input,
			preserved,
		}

		err := l.create_data_block(ctx, &block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *DataBlockLedger) Update(ctx contractapi.TransactionContextInterface, inputs []DataBlockInput) error {
	err := l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	timestamp, err := l.get_tx_timestamp(ctx)
	if err != nil {
		return err
	}

	for _, input := range inputs {
		block := DataBlock{
			input,
			DataBlockPreserved{
				Timestamp: timestamp,
			},
		}

		err := l.update_data_block(ctx, &block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *DataBlockLedger) Remove(ctx contractapi.TransactionContextInterface, keys []string) error {
	err := l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := l.delete_data_block(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *DataBlockLedger) Query(ctx contractapi.TransactionContextInterface, keys []string) ([]*DataBlock, error) {
	err := l.require_certification_read(ctx)
	visitor_mode := err != nil

	blocks := []*DataBlock{}
	for _, key := range keys {
		block, err := l.read_data_block(ctx, key)
		if err != nil {
			return nil, err
		}
		// if read access if not granted, hide the private section
		if visitor_mode {
			block.DataBlockPrivate = DataBlockPrivate{}
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (l *DataBlockLedger) QueryAll(ctx contractapi.TransactionContextInterface) ([]*DataBlock, error) {
	err := l.require_certification_read(ctx)
	visitor_mode := err != nil

	blocks, err := l.read_data_block_all(ctx)
	if err != nil {
		return nil, err
	}
	if visitor_mode {
		for _, block := range blocks {
			block.DataBlockPrivate = DataBlockPrivate{}
		}
	}

	return blocks, nil
}
