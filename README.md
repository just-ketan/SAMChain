
# Private Data Collections in Hyperledger Fabric

## 📌 Overview

This project demonstrates the implementation of **Private Data Collections (PDC)** in Hyperledger Fabric. The goal is to store public and private data separately while ensuring controlled access among organizations.

- **Public Data:** Asset ID (visible to all organizations)
- **Private Data:** Transaction amount (visible only to authorized organizations)

This implementation showcases how enterprise blockchain systems maintain **data confidentiality without compromising integrity**.

---

## 🎯 Objectives

- Store asset ID on the public ledger
- Store transaction amount in a private collection
- Configure `collections_config.json`
- Demonstrate controlled data access between organizations
- Prevent unauthorized access to private data

---

## 🧠 Key Concepts

### 🔹 Public Ledger
Stores data visible to all organizations in the channel.

### 🔹 Private Data Collection (PDC)
Stores sensitive data only on authorized peers.

### 🔹 Transient Data
Temporary data used during transaction execution and not stored on ledger.

### 🔹 Endorsement Policy
Defines which organizations must approve a transaction.

### 🔹 Collection Policy
Defines which organizations can access private data.

---

## 🏗️ Project Structure

```yaml
fabric-samples/
│
├── chaincode/privatecc/
│ ├── main.go
│ ├── go.mod
│ ├── collections_config.json
│
├── test-network/
│
├── README.md
```

---

## ⚙️ Setup Instructions

### 1️⃣ Clone Repository

```bash
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples
```

### 2️⃣ Start Fabric Network
```bash 
cd test-network
./network.sh up createChannel -c mychannel -ca
```
### 3️⃣ Install Dependencies
```bash
sudo apt update
sudo apt install -y jq
sudo apt install -y golang-go
````

### 4️⃣ Setup Chaincode
```bash
cd ~/fabric-samples/chaincode/privatecc
go mod init privatecc
go mod tidy
```

## 🔐 Private Data Collection Configuration

collections_config.json
```yaml
[
{
	"name": "collectionTransaction",
	"policy": "OR('Org1MSP.member')",
	"requiredPeerCount": 0,
	"maxPeerCount": 1,
	"blockToLive": 1000000,
	"memberOnlyRead": true
}
]
```

## 🚀 Deploy Chaincode

```bash
cd ~/fabric-samples/test-network

./network.sh down
./network.sh up createChannel -c mychannel -ca

./network.sh deployCC \
-ccn privatecc \
-ccl go \
-ccp ../chaincode/privatecc \
-c mychannel \
-cccg ../chaincode/privatecc/collections_config.json \
-ccep "OR('Org1MSP.peer')"
```

### 🔧 Set Environment Variables
```bash
source scripts/envVar.sh
setGlobals 1
export FABRIC_CFG_PATH=$PWD/../config
```

## 🔐 Encode Private Data
```bash
echo -n '{"value":"5000"}' | base64

Output:
eyJ2YWx1ZSI6IjUwMDAifQ==
```

### 🔄 Invoke Transaction
```bash
peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--tls \
--cafile "$ORDERER_CA" \
-C mychannel \
-n privatecc \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "$PEER0_ORG1_CA" \
--waitForEvent \
--transient '{"amount":"eyJ2YWx1ZSI6IjUwMDAifQ=="}' \
-c '{"function":"CreateAsset","Args":["asset1"]}'
```

### 🔍 Query Public Data
```bash
peer chaincode query \
-C mychannel \
-n privatecc \
-c '{"function":"ReadAsset","Args":["asset1"]}'
Output
{"assetID":"asset1"}
```

### 🔐 Query Private Data (Authorized)
```bash
peer chaincode query \
-C mychannel \
-n privatecc \
-c '{"function":"ReadPrivateAmount","Args":["asset1"]}'
Output
{"value":"5000"}
```

### ❌ Unauthorized Access Test
```bash
setGlobals 2

peer chaincode query \
-C mychannel \
-n privatecc \
-c '{"function":"ReadPrivateAmount","Args":["asset1"]}'
Output
tx creator does not have read access permission
```

## 📊 Results
| Feature             | Result  |
| ------------------- | ------- |
| Asset ID            | Public  |
| Transaction Amount  | Private |
| Authorized Access   | Allowed |
| Unauthorized Access | Denied  |

## 📌 Conclusion
This project successfully demonstrates how Hyperledger Fabric ensures data privacy and controlled access using Private Data Collections. It is a crucial feature for enterprise blockchain systems where confidentiality is required.