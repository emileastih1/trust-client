

# VASPResponseObject


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**id** | **String** | VASP UUID in the Bulletin Board config |  [optional] |
|**name** | **String** | VASP name |  [optional] |
|**domain** | **String** | The domain is the dNSName in the SAN (Subject Alternative Name) of a VASP&#39;s client certificate. This will be used by all other VASPs&#39; P2P channel and the bulletin Board for authentication during the mTLS handshake |  [optional] |
|**piiEndpoint** | **String** | The endpoint for receiving PII data. This is a full url (host + path) that can be directly used to make the HTTPS API call |  [optional] |
|**piiRequestEndpoint** | **String** | The endpoint for requesting PII data. This is a full url (host + path) that can be directly used to make the HTTPS API call |  [optional] |
|**publicKey** | **String** | The public key for other VASPs to encrypt the PII data using the JWE JOSE spec |  [optional] |
|**lei** | **String** | The Legal Entity Identifier of the VASP |  [optional] |
|**publicKeyVersion** | **Integer** | The version of the public key. Version should be bumped when a VASP rotates and updates its public key |  [optional] |
|**returnAddressEndpoint** | **String** | The endpoint for returning address. This is a full url (host + path) that can be directly used to make the HTTPS API call |  [optional] |
|**returnFundConfirmationEndpoint** | **String** | The endpoint for returning fund confirmation. This is a full url (host + path) that can be directly used to make the HTTPS API call |  [optional] |



