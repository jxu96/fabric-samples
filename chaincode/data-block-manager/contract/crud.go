package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func data_block_to_bytes(block *DataBlock) ([]byte, error) {
	bs, err := json.Marshal(*block)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode data block to bytes: %s\n%v", block.ID, err)
	}

	return bs, nil
}

func bytes_to_data_block(key string, bs []byte) (*DataBlock, error) {
	var block DataBlock
	err := json.Unmarshal(bs, &block)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse to json format: %s\n%v", key, err)
	}

	return &block, nil
}

func (l *DataBlockLedger) get(ctx contractapi.TransactionContextInterface, key string) ([]byte, error) {
	bs, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from ledger: %s\n%v", key, err)
	}

	return bs, nil
}

func (l *DataBlockLedger) set(ctx contractapi.TransactionContextInterface, key string, bs []byte) error {
	err := ctx.GetStub().PutState(key, bs)
	if err != nil {
		return fmt.Errorf("Failed to set value on ledger: %s\n%v", key, err)
	}

	return nil
}

func (l *DataBlockLedger) del(ctx contractapi.TransactionContextInterface, key string) error {
	err := ctx.GetStub().DelState(key)
	if err != nil {
		return fmt.Errorf("Failed to delete value from ledger: %s\n%v", key, err)
	}

	return nil
}

func (l *DataBlockLedger) check_data_block_exists(ctx contractapi.TransactionContextInterface, key string) ([]byte, error) {
	bs, err := l.get(ctx, key)
	if err != nil {
		return nil, err
	}

	if bs == nil {
		return nil, fmt.Errorf("Data block ID does not exist on ledger: %s", key)
	}

	return bs, nil
}

func (l *DataBlockLedger) check_data_block_not_exists(ctx contractapi.TransactionContextInterface, key string) error {
	bs, err := l.get(ctx, key)
	if err != nil {
		return err
	}

	if bs != nil {
		return fmt.Errorf("Data block ID already exists on ledger: %s", key)
	}

	return nil
}

func (l *DataBlockLedger) read_data_block(ctx contractapi.TransactionContextInterface, key string) (*DataBlock, error) {
	bs, err := l.check_data_block_exists(ctx, key)
	if err != nil {
		return nil, err
	}

	return bytes_to_data_block(key, bs)
}

// get data blocks on given keys
//
// will not raise error if not existed
func (l *DataBlockLedger) read_data_block_batch(ctx contractapi.TransactionContextInterface, keys []string) ([]*DataBlock, error) {
	var e error = nil
	var blocks []*DataBlock = make([]*DataBlock, len(keys))

	for i, key := range keys {
		block, err := l.read_data_block(ctx, key)
		if err != nil {
			if e == nil {
				e = err
			} else {
				e = fmt.Errorf("%v\n\n%v", e, err)
			}
		}
		blocks[i] = block
	}

	return blocks, e
}

func (l *DataBlockLedger) read_data_block_all(ctx contractapi.TransactionContextInterface) ([]*DataBlock, error) {
	blocks := []*DataBlock{}

	resp, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("Failed to read from ledger: %v", err)
	}
	defer resp.Close()

	for resp.HasNext() {
		query, err := resp.Next()
		if err != nil {
			return nil, fmt.Errorf("Failed to query from ledger: %v", err)
		}

		block, err := bytes_to_data_block(query.Key, query.Value)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (l *DataBlockLedger) create_data_block(ctx contractapi.TransactionContextInterface, block *DataBlock) error {
	err := l.check_data_block_not_exists(ctx, block.ID)
	if err != nil {
		return err
	}

	bs, err := data_block_to_bytes(block)
	if err != nil {
		return err
	}

	return l.set(ctx, block.ID, bs)
}

func (l *DataBlockLedger) update_data_block(ctx contractapi.TransactionContextInterface, block *DataBlock) error {
	_, err := l.check_data_block_exists(ctx, block.ID)
	if err != nil {
		return err
	}

	bs, err := data_block_to_bytes(block)
	if err != nil {
		return err
	}

	return l.set(ctx, block.ID, bs)
}

func (l *DataBlockLedger) delete_data_block(ctx contractapi.TransactionContextInterface, key string) error {
	_, err := l.check_data_block_exists(ctx, key)
	if err != nil {
		return err
	}

	return l.del(ctx, key)
}

// func (l *DataBlockLedger) exists_t(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
// 	bs, err := l.get(ctx, key)
// 	return bs != nil, err
// }
