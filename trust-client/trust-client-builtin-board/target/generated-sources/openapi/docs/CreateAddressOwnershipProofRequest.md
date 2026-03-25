

# CreateAddressOwnershipProofRequest


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**registrationId** | **String** | The registration UUID when you claimed the address by calling PUT /addresses/&lt;address&gt; |  [optional] |
|**chain** | **String** | The blockchain of the address, e.g., BITCOIN, ETHEREUM |  [optional] |
|**signature** | **String** | The signature of the proof of ownership message signed by the private key of the address |  [optional] |
|**prefix** | **String** | The prefix of the claiming VASP&#39;s choice |  [optional] |
|**iou** | **Boolean** | Indicate whether this is an IOU |  [optional] |
|**auxProofData** | [**List&lt;AuxProofData&gt;**](AuxProofData.md) | Additional data supporting the different types of proof of address ownerships, e.g., redeem script. (See Proof of Address Ownership Spec doc for auxiliary data required for all of the supported proof types) |  [optional] |
|**proofType** | **String** | The proof type dictating the logic of proofing and verification according to the spec. Example values: BITCOIN_P2PKH, BITCOIN_P2SH, BITCOIN_P2WPKH, BITCOIN_P2WSH, ETHEREUM_EOA, ETHEREUM_CONTRACT. (See Proof of Address Ownership Spec doc for all the supported proof types) |  [optional] |



