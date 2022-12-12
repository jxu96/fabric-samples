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

type DataBlockListInterface interface {
}

var (
	CollectionDataBlock        string = "collectionDataBlock"
	CollectionDataBlockPrivate string = "collectionDataBlockPrivate"
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

func (l *DataBlockLedger) Register(ctx contractapi.TransactionContextInterface) error {
	err := l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	transient, err := l.get_tx_transient(ctx)
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

	for key, data := range transient {
		// register public section in the collection
		block := new(DataBlock)
		block.Owner = sender
		block.Timestamp = timestamp

		err := block.from_bytes(key, data)
		if err != nil {
			return err
		}

		err = l.create_data_block(ctx, CollectionDataBlock, block.ID, block)
		if err != nil {
			return err
		}

		// register private section in permissioned collection
		block_private := new(DataBlockPrivate)
		err = block_private.from_bytes(key, data)
		if err != nil {
			return err
		}

		err = l.create_data_block(ctx, CollectionDataBlockPrivate, block_private.ID, block_private)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *DataBlockLedger) Update(ctx contractapi.TransactionContextInterface) error {
	err := l.require_certification_write(ctx)
	if err != nil {
		return err
	}

	transient, err := l.get_tx_transient(ctx)
	if err != nil {
		return err
	}

	timestamp, err := l.get_tx_timestamp(ctx)
	if err != nil {
		return err
	}

	for key, data := range transient {
		// update public section in the collection
		block := new(DataBlock)
		block.Timestamp = timestamp

		err := block.from_bytes(key, data)
		if err != nil {
			return err
		}

		err = l.update_data_block(ctx, CollectionDataBlock, block.ID, block)
		if err != nil {
			return err
		}

		// update private section in permissioned collection
		block_private := new(DataBlockPrivate)
		err = block_private.from_bytes(key, data)
		if err != nil {
			return err
		}

		err = l.update_data_block(ctx, CollectionDataBlockPrivate, block_private.ID, block_private)
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
		err := l.delete_data_block(ctx, CollectionDataBlock, key)
		if err != nil {
			return err
		}

		err = l.delete_data_block(ctx, CollectionDataBlockPrivate, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *DataBlockLedger) Query(ctx contractapi.TransactionContextInterface, keys []string) ([]*DataBlock, error) {
	err := l.require_certification_read(ctx)
	if err != nil {
		return nil, err
	}

	blocks := make([]*DataBlock, len(keys))
	for i, key := range keys {
		err := l.read_data_block(ctx, CollectionDataBlock, key, blocks[i])
		if err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

func (l *DataBlockLedger) QueryPrivate(ctx contractapi.TransactionContextInterface, keys []string) ([]*DataBlockPrivate, error) {
	err := l.require_certification_read(ctx)
	if err != nil {
		return nil, err
	}

	blocks := make([]*DataBlockPrivate, len(keys))
	for i, key := range keys {
		err := l.read_data_block(ctx, CollectionDataBlockPrivate, key, blocks[i])
		if err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

func (l *DataBlockLedger) QueryByRange(ctx contractapi.TransactionContextInterface, start string, end string, max int) ([]*DataBlock, error) {
	err := l.require_certification_read(ctx)
	if err != nil {
		return nil, err
	}

	map_, err := l.get_by_range(ctx, CollectionDataBlock, start, end, max)
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
