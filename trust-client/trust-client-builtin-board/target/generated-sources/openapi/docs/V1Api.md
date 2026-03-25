# V1Api

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**v1AddressOwnershipProofsAddressPut**](V1Api.md#v1AddressOwnershipProofsAddressPut) | **PUT** /v1/address_ownership_proofs/{address} | Submit a proof of address ownership for a previously claimed (address/chain) |
| [**v1AddressesAddressDelete**](V1Api.md#v1AddressesAddressDelete) | **DELETE** /v1/addresses/{address} | Delete a Record of Address Custody |
| [**v1AddressesAddressGet**](V1Api.md#v1AddressesAddressGet) | **GET** /v1/addresses/{address} | Retrieve the custody of an address and the proof of address ownership |
| [**v1AddressesAddressPut**](V1Api.md#v1AddressesAddressPut) | **PUT** /v1/addresses/{address} | Claim custody of an (address/chain) |
| [**v1VaspsGet**](V1Api.md#v1VaspsGet) | **GET** /v1/vasps | Retrieve the information of all VASPs |
| [**v1VaspsVaspIdGet**](V1Api.md#v1VaspsVaspIdGet) | **GET** /v1/vasps/{vaspId} | Retrieve the information of a specific VASP |



## v1AddressOwnershipProofsAddressPut

> CreateAddressOwnershipProofResponse v1AddressOwnershipProofsAddressPut(address, createAddressOwnershipProofRequest)

Submit a proof of address ownership for a previously claimed (address/chain)

This is the singular endpoint for creating an address ownership proof

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address
        CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest = new CreateAddressOwnershipProofRequest(); // CreateAddressOwnershipProofRequest | 
        try {
            CreateAddressOwnershipProofResponse result = apiInstance.v1AddressOwnershipProofsAddressPut(address, createAddressOwnershipProofRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1AddressOwnershipProofsAddressPut");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **address** | **String**| The SHA512 of the blockchain address | |
| **createAddressOwnershipProofRequest** | [**CreateAddressOwnershipProofRequest**](CreateAddressOwnershipProofRequest.md)|  | [optional] |

### Return type

[**CreateAddressOwnershipProofResponse**](CreateAddressOwnershipProofResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Create or update proof of address succeeded |  -  |
| **400** | Invalid request. For example, the (address, chain, registration_id) combination is invalid or missing fields. |  -  |
| **403** | Not authorized |  -  |
| **409** | Conflicted proof status, e.g., change from non-IOU to IOU |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |


## v1AddressesAddressDelete

> DeleteAddressOwnershipResponse v1AddressesAddressDelete(address, chain)

Delete a Record of Address Custody

This endpoint deletes a record of address custody from your VASP&#39;s filter

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            DeleteAddressOwnershipResponse result = apiInstance.v1AddressesAddressDelete(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1AddressesAddressDelete");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **address** | **String**| The SHA512 of the blockchain address. | |
| **chain** | **String**| The blockchain of the address, e.g., BITCOIN, ETHEREUM. | [optional] |

### Return type

[**DeleteAddressOwnershipResponse**](DeleteAddressOwnershipResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Delete succeeded |  -  |
| **400** | Invalid request, parameter validation failure. |  -  |
| **403** | Not authorized |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |


## v1AddressesAddressGet

> GetAddressOwnershipResponse v1AddressesAddressGet(address, chain)

Retrieve the custody of an address and the proof of address ownership

This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            GetAddressOwnershipResponse result = apiInstance.v1AddressesAddressGet(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1AddressesAddressGet");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **address** | **String**| The SHA512 of the blockchain address. | |
| **chain** | **String**| The blockchain of the address, e.g., BITCOIN, ETHEREUM. | [optional] |

### Return type

[**GetAddressOwnershipResponse**](GetAddressOwnershipResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Get address ownership succeeded |  -  |
| **400** | Invalid request, parameter validation failure. |  -  |
| **403** | Not authorized |  -  |
| **404** | Address not found |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |


## v1AddressesAddressPut

> CreateAddressOwnershipResponse v1AddressesAddressPut(address, createAddressOwnershipRequest)

Claim custody of an (address/chain)

Claim custody of an (address/chain)

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        CreateAddressOwnershipRequest createAddressOwnershipRequest = new CreateAddressOwnershipRequest(); // CreateAddressOwnershipRequest | 
        try {
            CreateAddressOwnershipResponse result = apiInstance.v1AddressesAddressPut(address, createAddressOwnershipRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1AddressesAddressPut");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **address** | **String**| The SHA512 of the blockchain address. | |
| **createAddressOwnershipRequest** | [**CreateAddressOwnershipRequest**](CreateAddressOwnershipRequest.md)|  | [optional] |

### Return type

[**CreateAddressOwnershipResponse**](CreateAddressOwnershipResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Claim address custody succeeded |  -  |
| **400** | Invalid request, parameter validation failure. |  -  |
| **403** | Not authorized |  -  |
| **409** | Conflict (when the address is already claimed by others) |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |


## v1VaspsGet

> GetVASPResponse1 v1VaspsGet()

Retrieve the information of all VASPs

Retrieve the information of all VASPs.

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        try {
            GetVASPResponse1 result = apiInstance.v1VaspsGet();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1VaspsGet");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**GetVASPResponse1**](GetVASPResponse1.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully get the information of all VASPs |  -  |
| **403** | Not authorized |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |


## v1VaspsVaspIdGet

> GetVASPResponse v1VaspsVaspIdGet(vaspId)

Retrieve the information of a specific VASP

This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.V1Api;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        V1Api apiInstance = new V1Api(defaultClient);
        String vaspId = "string"; // String | The UUID of the VASP that your are querying for
        try {
            GetVASPResponse result = apiInstance.v1VaspsVaspIdGet(vaspId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling V1Api#v1VaspsVaspIdGet");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **vaspId** | **String**| The UUID of the VASP that your are querying for | |

### Return type

[**GetVASPResponse**](GetVASPResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Successfully retrieve the information of the VASP |  -  |
| **400** | Invalid request, parameter validation failure. For example, the provided vasp id is not a valid UUID |  -  |
| **403** | Not authorized |  -  |
| **404** | VASP not found |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |

