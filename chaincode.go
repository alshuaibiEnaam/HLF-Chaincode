package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type eFogServiceProvider struct {
	eFSP_ID         string  `json:"efsp_id"`
	ResponseTime    float32 `json:"responseTime"`
	Availability    int     `json:"availability"`
	Throughput      float32 `json:"throughput"`
	Reliability     int     `json:"reliability"`
	Latency         float32 `json:"latency"`
	ReputationScore string  `json:"reputationScore"`
	Owner           string  `json:"owner"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []eFogServiceProvider{
		{eFSP_ID: "SP1161", ResponseTime: 41, Availability: 97, Throughput: 43.1, Reliability: 73, Latency: 1, Owner: "SP1161", ReputationScore: "4"},
		{eFSP_ID: "SP1606", ResponseTime: 67, Availability: 86, Throughput: 41, Reliability: 73, Latency: 5, Owner: "SP1606", ReputationScore: "3"},
		{eFSP_ID: "SP1495", ResponseTime: 57, Availability: 86, Throughput: 40.1, Reliability: 73, Latency: 1, Owner: "SP1495", ReputationScore: "3"},
		{eFSP_ID: "SP1586", ResponseTime: 82, Availability: 83, Throughput: 41.2, Reliability: 67, Latency: 2, Owner: "SP1586", ReputationScore: "5"},
		{eFSP_ID: "SP614", ResponseTime: 89, Availability: 83, Throughput: 40.4, Reliability: 78, Latency: 3, Owner: "SP614", ReputationScore: "5"},
		{eFSP_ID: "SP951", ResponseTime: 82, Availability: 83, Throughput: 83.7, Reliability: 80, Latency: 2, Owner: "SP951", ReputationScore: "6"},
	}

	for _, eFogServiceProvider := range assets {
		assetJSON, err := json.Marshal(eFogServiceProvider)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(eFogServiceProvider.eFSP_ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) RegisterNew_eFSP(ctx contractapi.TransactionContextInterface, efsp_id string, responseTime float32, availability int, throughput float32, reliability int, latency float32, owner string, reputationScore string) error {
	exists, err := s.eFSPExists(ctx, efsp_id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the provider %s already exists", efsp_id)
	}

	asset := eFogServiceProvider{
		eFSP_ID:         efsp_id,
		ResponseTime:    responseTime,
		Availability:    availability,
		Throughput:      throughput,
		Reliability:     reliability,
		Latency:         latency,
		ReputationScore: reputationScore,
		Owner:           owner,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(efsp_id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) Read_eFSP(ctx contractapi.TransactionContextInterface, efsp_id string) (*eFogServiceProvider, error) {
	assetJSON, err := ctx.GetStub().GetState(efsp_id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the provider %s does not exist", efsp_id)
	}

	var asset eFogServiceProvider
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) Update_eFSP(ctx contractapi.TransactionContextInterface, efsp_id string, responseTime float32, availability int, throughput float32, reliability int, latency float32, owner string, reputationScore string) error {
	exists, err := s.eFSPExists(ctx, efsp_id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", efsp_id)
	}

	// overwriting original asset with new asset
	asset := eFogServiceProvider{
		eFSP_ID:         efsp_id,
		ResponseTime:    responseTime,
		Availability:    availability,
		Throughput:      throughput,
		Reliability:     reliability,
		Latency:         latency,
		Owner:           owner,
		ReputationScore: reputationScore,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(efsp_id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) Delete_eFSP(ctx contractapi.TransactionContextInterface, efsp_id string) error {
	exists, err := s.eFSPExists(ctx, efsp_id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", efsp_id)
	}

	return ctx.GetStub().DelState(efsp_id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) eFSPExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) GetAll_eFSP(ctx contractapi.TransactionContextInterface) ([]*eFogServiceProvider, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*eFogServiceProvider
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset eFogServiceProvider
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) QueryAssetByRepuationScore(ctx contractapi.TransactionContextInterface, minReputationScore string) ([]*eFogServiceProvider, error) {

	queryString := fmt.Sprintf(`{"selector":{"reputationScore":{"$gte":"%s"}}}`, minReputationScore)

	return getQueryResultForQueryString(ctx, queryString)

}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*eFogServiceProvider, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*eFogServiceProvider, error) {
	var assets []*eFogServiceProvider
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset eFogServiceProvider
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
