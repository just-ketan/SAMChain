package main
import(
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


// smartcontract structure
type SmartContract struct{
	contractapi.Contract
}

// asset structure (public)
type Asset struct{
	AssetID string `json:"assetID"`
}

// creataAsset -> stores public + prvate data
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
// public data
	asset := Asset{
		AssetID: assetID,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil{
		return err
	}

	// store in public ledger
	err = ctx.GetStub().PutState(assetID, assetJSON)
	if err != nil{
		return fmt.Errorf("failed to store asset: %v", err)
	}

// private data
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil{
		return err
	}

	// get amt from transient map
	amountJSON, ok := transientMap["amount"]
	if !ok{
		return fmt.Errorf("amount not found in tansient map")
	}
	
	// store in private collection
	err = ctx.GetStub().PutPrivateData("collectionTransaction", assetID, amountJSON)
	if err != nil{
		return fmt.Errorf("failed to store private data: %v", err)
	}
	return nil
}

// read public asset
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string)(*Asset, error){

	data, err := ctx.GetStub().GetState(assetID)
	if err != nil{
		return nil, err
	}
	if data == nil{
		return nil, fmt.Errorf("asset not found")
	}
	
	var asset Asset
	_ = json.Unmarshal(data, &asset)

	return &asset, nil
}

// read private amount
func (s *SmartContract) ReadPrivateAmount(ctx contractapi.TransactionContextInterface, assetID string) ( string, error ) {

	data, err := ctx.GetStub().GetPrivateData("collectionTransaction", assetID)
	if err != nil{
		return "", err
	}
	if data == nil{
		return "", fmt.Errorf("private data not accessible")
	}

	return string(data), nil
}

func main(){
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil{
		panic(err.Error())
	}
	if err := chaincode.Start(); err != nil {
 		panic(err.Error())
	}
}
