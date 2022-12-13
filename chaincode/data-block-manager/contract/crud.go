package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func readFromPublic(ctx contractapi.TransactionContextInterface, key string) ([]byte, error) {
	bs, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf(`Failed to read from the public ledger with key "%s" : %v`, key, err)
	}

	return bs, nil
}

func readFromCollection(ctx contractapi.TransactionContextInterface, collection string, key string) ([]byte, error) {
	bs, err := ctx.GetStub().GetPrivateData(collection, key)
	if err != nil {
		return nil, fmt.Errorf(`Failed to read from collection "%s" with key "%s" : %v`, collection, key, err)
	}

	return bs, nil
}

func readFromPublicByRange(ctx contractapi.TransactionContextInterface, start string, end string, max int) ([][]byte, error) {
	result := [][]byte{}
	counter := 0

	it, err := ctx.GetStub().GetStateByRange(start, end)
	if err != nil {
		return nil, fmt.Errorf(`Failed to read from the public ledger by range "%s" -> "%s" : %v`, start, end, err)
	}
	defer it.Close()

	for it.HasNext() && counter < max {
		item, err := it.Next()
		if err != nil {
			return nil, err
		}
		result = append(result, item.Value)
		counter++
	}

	return result, nil
}

func readFromCollectionByRange(ctx contractapi.TransactionContextInterface, collection string, start string, end string, max int) ([][]byte, error) {
	result := [][]byte{}
	counter := 0

	it, err := ctx.GetStub().GetPrivateDataByRange(collection, start, end)
	if err != nil {
		return nil, fmt.Errorf(`Failed to read from collection "%s" by range "%s" -> "%s" : %v`, collection, start, end, err)
	}
	defer it.Close()

	for it.HasNext() && counter < max {
		item, err := it.Next()
		if err != nil {
			return nil, err
		}
		result = append(result, item.Value)
		counter++
	}

	return result, nil
}

func createFromPublic(ctx contractapi.TransactionContextInterface, key string, value []byte) error {
	bs, err := readFromPublic(ctx, key)
	if err != nil {
		return err
	}
	if bs != nil {
		return fmt.Errorf(`Failed to create from the public ledger : key already exists "%s"`, key)
	}

	if err := ctx.GetStub().PutState(key, value); err != nil {
		return fmt.Errorf(`Failed to create from the public ledger with key "%s" : %v`, key, err)
	}

	return nil
}

func createFromCollection(ctx contractapi.TransactionContextInterface, collection string, key string, value []byte) error {
	bs, err := readFromCollection(ctx, collection, key)
	if err != nil {
		return err
	}
	if bs != nil {
		return fmt.Errorf(`Failed to create from collection "%s" : key already exists "%s"`, collection, key)
	}

	if err := ctx.GetStub().PutPrivateData(collection, key, value); err != nil {
		return fmt.Errorf(`Failed to write to collection "%s" with key "%s" : %v`, collection, key, err)
	}

	return nil
}

func deleteFromPublic(ctx contractapi.TransactionContextInterface, key string) error {
	err := ctx.GetStub().DelState(key)
	if err != nil {
		return fmt.Errorf(`Failed to remove from the public ledger with key "%s" : %v`, key, err)
	}

	return nil
}

func deleteFromCollection(ctx contractapi.TransactionContextInterface, collection string, key string) error {
	err := ctx.GetStub().DelPrivateData(collection, key)
	if err != nil {
		return fmt.Errorf(`Failed to remove from collection "%s" with key "%s" : %v`, collection, key, err)
	}

	return nil
}

// func (l *DataBlockLedger) read_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
// 	bs, err := l.check_data_block_exists(ctx, collection, key)
// 	if err != nil {
// 		return err
// 	}

// 	return block.from_bytes(key, bs)
// }

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

// func (l *DataBlockLedger) read_data_block_by_range(ctx contractapi.TransactionContextInterface, collection string, start string, end string, max int) error {
// 	map_, err := l.get_by_range(ctx, collection, start, end, max)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l *DataBlockLedger) create_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
// 	err := l.check_data_block_not_exists(ctx, collection, key)
// 	if err != nil {
// 		return err
// 	}

// 	bs, err := block.to_bytes()
// 	if err != nil {
// 		return err
// 	}

// 	return l.set(ctx, collection, key, bs)
// }

// func (l *DataBlockLedger) update_data_block(ctx contractapi.TransactionContextInterface, collection string, key string, block DataBlockInterface) error {
// 	_, err := l.check_data_block_exists(ctx, collection, key)
// 	if err != nil {
// 		return err
// 	}

// 	bs, err := block.to_bytes()
// 	if err != nil {
// 		return err
// 	}

// 	return l.set(ctx, collection, key, bs)
// }

// func (l *DataBlockLedger) delete_data_block(ctx contractapi.TransactionContextInterface, collection, key string) error {
// 	_, err := l.check_data_block_exists(ctx, collection, key)
// 	if err != nil {
// 		return err
// 	}

// 	return l.del(ctx, collection, key)
// }
