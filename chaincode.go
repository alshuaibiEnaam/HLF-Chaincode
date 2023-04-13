//original chaincode
package chaincode

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"

)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
    contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
//Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type eFogServiceProvider struct {
	eFSP_ID             string      `json:"efsp_id"`
    Name                string      `json:"name"`
    ResponseTime        float32     `json:"responseTime"`
    Availability        int         `json:"availability"`
    Throughput          float32     `json:"throughput"`
    Reliability         int         `json:"reliability"`
    Latency             float32     `json:"latency"`    
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []eFogServiceProvider{
        {eFSP_ID: "SP1161", ResponseTime: 41, Availability: 97, Throughput: 43.1, Reliability: 73, Latency: 1, Name:"SP1161"},
        {eFSP_ID: "SP1606", ResponseTime: 67, Availability: 86, Throughput: 41,   Reliability: 73, Latency: 5, Name:"SP1606"},
        {eFSP_ID: "SP1495", ResponseTime: 57, Availability: 86, Throughput: 40.1, Reliability: 73, Latency: 1, Name:"SP1495"},
        {eFSP_ID: "SP1586", ResponseTime: 82, Availability: 83, Throughput: 41.2, Reliability: 67, Latency: 2, Name:"SP1586"},
        {eFSP_ID: "SP614", ResponseTime: 89, Availability: 83, Throughput: 40.4, Reliability: 78, Latency: 3, Name:"SP614"},
        {eFSP_ID: "SP951", ResponseTime: 82, Availability: 83, Throughput: 83.7, Reliability: 80, Latency: 2, Name:"SP951"},
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
func (s *SmartContract) RegisterNew_eFSP(ctx contractapi.TransactionContextInterface, id string, responseTime float32, availability int, throughput float32, reliability int, latency float32, name string) error {
    exists, err := s.eFSPExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the asset %s already exists", id)
    }

    asset := eFogServiceProvider{
        eFSP_ID : id,
        ResponseTime: responseTime,
        Availability: availability,
        Throughput: throughput,
        Reliability:reliability,
        Latency: latency,
        Name: name,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, rank string) (*eFogServiceProvider, error) {
    assetJSON, err := ctx.GetStub().GetState(rank)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("the asset %s does not exist", rank)
    }

    var asset eFogServiceProvider
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) Update_eFSP(ctx contractapi.TransactionContextInterface, id string, responseTime float32, availability int, throughput float32, reliability int, latency float32, name string) error {
    exists, err := s.eFSPExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", id)
    }

    // overwriting original asset with new asset
    asset := eFogServiceProvider{
		eFSP_ID : id,
        ResponseTime: responseTime,
        Availability: availability,
        Throughput: throughput,
        Reliability:reliability,
        Latency: latency,
        Name: name,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) Delete_eFSP(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.eFSPExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) eFSPExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
/*func (s *SmartContract) updateRank(ctx contractapi.TransactionContextInterface, id string, newRank string) (string, error) {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
        return "", err
    }

    odlerRank := asset.Rank
    asset.Rank = newRank

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return "", err
    }

    err = ctx.GetStub().PutState(id, assetJSON)
    if err != nil {
        return "", err
    }

    return odlerRank, nil
}
*/
// GetAllAssets returns all assets found in world state
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

func (s *SmartContract) QueryAssetByRank (ctx contractapi.TransactionContextInterface, rank string) ([]*eFogServiceProvider, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"Rank\":\"%v\"}}", rank)

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

