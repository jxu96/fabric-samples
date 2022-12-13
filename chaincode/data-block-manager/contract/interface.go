package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DatasetMetadataLedger struct {
	contractapi.Contract
}

func (l *DatasetMetadataLedger) Register(ctx contractapi.TransactionContextInterface) error {
	var mspID string
	if err := requireIdenticalMSPID(ctx, &mspID); err != nil {
		return err
	}

	if err := requireCertification(ctx, nil); err != nil {
		return err
	}

	transient, err := getTxTransient(ctx)
	if err != nil {
		return err
	}

	// Read metadata input from transient
	mdInputAsBytes, ok := transient["metadata"]
	if !ok {
		return fmt.Errorf("Dataset metadata not defined in transient.")
	}
	// Read collections to register from transient
	collectionsAsBytes, ok := transient["collections"]
	if !ok {
		return fmt.Errorf("Collections not defined in transient.")
	}

	// Write to public collections defined in transient
	var collections []string
	if err := json.Unmarshal(collectionsAsBytes, &collections); err != nil {
		return fmt.Errorf("Failed to decode collections : %v", err)
	}

	mdPublic := new(DatasetMetadataPublic)
	if err := mdPublic.FromBytes(mdInputAsBytes); err != nil {
		return err
	}
	if err := mdPublic.Validate(); err != nil {
		return err
	}
	mdPublicAsBytes, err := mdPublic.ToBytes()
	if err != nil {
		return err
	}

	for _, collection := range collections {
		if collection != "" {
			if err := createFromCollection(ctx, collection, mdPublic.ID, mdPublicAsBytes); err != nil {
				return err
			}
		} else {
			if err := createFromPublic(ctx, mdPublic.ID, mdPublicAsBytes); err != nil {
				return err
			}
		}
	}

	// Write to implicit collection
	md := new(DatasetMetadata)
	if err := md.FromBytes(mdInputAsBytes); err != nil {
		return err
	}
	if err := md.Validate(); err != nil {
		return err
	}
	mdAsBytes, err := md.ToBytes()
	if err != nil {
		return err
	}

	if err := createFromCollection(ctx, implicitPrivateDataCollection(mspID), md.ID, mdAsBytes); err != nil {
		return err
	}

	return nil
}

func (l *DatasetMetadataLedger) Query(ctx contractapi.TransactionContextInterface, collection string, key string) (*DatasetMetadataPublic, error) {
	if err := requireCertification(ctx, nil); err != nil {
		return nil, err
	}

	md := new(DatasetMetadataPublic)

	if collection != "" {
		mdAsBytes, err := readFromCollection(ctx, collection, key)
		if err != nil {
			return nil, err
		}
		if err := md.FromBytes(mdAsBytes); err != nil {
			return nil, err
		}
	} else {
		mdAsBytes, err := readFromPublic(ctx, key)
		if err != nil {
			return nil, err
		}
		if err := md.FromBytes(mdAsBytes); err != nil {
			return nil, err
		}
	}

	return md, nil
}

func (l *DatasetMetadataLedger) QueryPrivate(ctx contractapi.TransactionContextInterface, key string) (*DatasetMetadata, error) {
	mspID, err := getPeerMSPID()
	if err != nil {
		return nil, err
	}

	if err := requireCertification(ctx, nil); err != nil {
		return nil, err
	}

	md := new(DatasetMetadata)
	mdAsBytes, err := readFromCollection(ctx, implicitPrivateDataCollection(mspID), key)
	if err != nil {
		return nil, err
	}
	if err := md.FromBytes(mdAsBytes); err != nil {
		return nil, err
	}

	return md, nil
}

func (l *DatasetMetadataLedger) QueryByRange(ctx contractapi.TransactionContextInterface, collection string, start string, end string, max int) ([]*DatasetMetadataPublic, error) {
	if err := requireCertification(ctx, nil); err != nil {
		return nil, err
	}

	toMdArray := func(bsArray [][]byte) ([]*DatasetMetadataPublic, error) {
		result := make([]*DatasetMetadataPublic, len(bsArray))
		for i, md := range result {
			if err := md.FromBytes(bsArray[i]); err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	if collection != "" {
		bsArray, err := readFromCollectionByRange(ctx, collection, start, end, max)
		if err != nil {
			return nil, err
		}
		return toMdArray(bsArray)
	} else {
		bsArray, err := readFromPublicByRange(ctx, start, end, max)
		if err != nil {
			return nil, err
		}
		return toMdArray(bsArray)
	}
}
