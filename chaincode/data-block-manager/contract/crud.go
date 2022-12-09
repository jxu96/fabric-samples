package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (l *DataBlockLedger) get(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	// bs, err := ctx.GetStub().GetState(key)
	bs, err := ctx.GetStub().GetPrivateData(collection, key)
	if err != nil {
		return nil, fmt.Errorf("Failed to get data block in collection [%s] with key [%s].\n%v", collection, key, err)
	}

	return bs, nil
}

func (l *DataBlockLedger) set(ctx contractapi.TransactionContextInterface, collection string, key string, bs []byte) error {
	// err := ctx.GetStub().PutState(key, bs)
	err := ctx.GetStub().PutPrivateData(collection, key, bs)
	if err != nil {
		return fmt.Errorf("Failed to set data block in collection [%s] with key [%s].\n%v", collection, key, err)
	}

	return nil
}

func (l *DataBlockLedger) del(ctx contractapi.TransactionContextInterface, collection string, key string) error {
	// err := ctx.GetStub().DelState(key)
	err := ctx.GetStub().DelPrivateData(collection, key)
	if err != nil {
		return fmt.Errorf("Failed to delete data block in collection [%s] with key [%s].\n%v", collection, key, err)
	}

	return nil
}

func (l *DataBlockLedger) check_data_block_exists(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	bs, err := l.get(ctx, collection, key)
	if err != nil {
		return nil, err
	}

	if bs == nil {
		return nil, fmt.Errorf("Data block with key [%s] doesn't exist in collection [%s].", key, collection)
	}

	return bs, nil
}

func (l *DataBlockLedger) check_data_block_not_exists(ctx contractapi.TransactionContextInterface, collection string, key string) error {
	bs, err := l.get(ctx, collection, key)
	if err != nil {
		return err
	}

	if bs != nil {
		return fmt.Errorf("Data block with key [%s] already exists in collection [%s].", key, collection)
	}

	return nil
}

func (l *DataBlockLedger) read_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
	bs, err := l.check_data_block_exists(ctx, collection, key)
	if err != nil {
		return err
	}

	return block.from_bytes(key, bs)
}

// func (l *DataBlockLedger) read_data_block_private(ctx contractapi.TransactionContextInterface, collection string, key string) (*DataBlockPrivate, error) {
// 	bs, err := l.check_data_block_exists(ctx, collection, key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	block := new(DataBlockPrivate)
// 	err = block.from_bytes(key, bs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return block, nil
// }

func (l *DataBlockLedger) read_data_block_by_range(ctx contractapi.TransactionContextInterface, collection string, start string, end string) ([]*DataBlock, error) {
	blocks := []*DataBlock{}

	resp, err := ctx.GetStub().GetPrivateDataByRange(collection, start, end)
	if err != nil {
		return nil, fmt.Errorf("Failed to get data block in collection [%s] from key [%s] to key [%s].\n%v", collection, start, end, err)
	}
	defer resp.Close()

	for resp.HasNext() {
		query, err := resp.Next()
		if err != nil {
			return nil, err
		}

		block := new(DataBlock)

		err = block.from_bytes(query.Key, query.Value)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (l *DataBlockLedger) create_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
	err := l.check_data_block_not_exists(ctx, collection, key)
	if err != nil {
		return err
	}

	bs, err := block.to_bytes()
	if err != nil {
		return err
	}

	return l.set(ctx, collection, key, bs)
}

func (l *DataBlockLedger) update_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
	_, err := l.check_data_block_exists(ctx, collection, key)
	if err != nil {
		return err
	}

	bs, err := block.to_bytes()
	if err != nil {
		return err
	}

	return l.set(ctx, collection, key, bs)
}

func (l *DataBlockLedger) delete_data_block(ctx contractapi.TransactionContextInterface, collection, key string) error {
	_, err := l.check_data_block_exists(ctx, collection, key)
	if err != nil {
		return err
	}

	return l.del(ctx, collection, key)
}
