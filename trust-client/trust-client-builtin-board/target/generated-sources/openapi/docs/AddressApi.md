# AddressApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**v1AddressOwnershipProofsAddressPut_1**](AddressApi.md#v1AddressOwnershipProofsAddressPut_1) | **PUT** /v1/address_ownership_proofs/{address} | Submit a proof of address ownership for a previously claimed (address/chain) |
| [**v1AddressesAddressDelete_1**](AddressApi.md#v1AddressesAddressDelete_1) | **DELETE** /v1/addresses/{address} | Delete a Record of Address Custody |
| [**v1AddressesAddressGet_1**](AddressApi.md#v1AddressesAddressGet_1) | **GET** /v1/addresses/{address} | Retrieve the custody of an address and the proof of address ownership |
| [**v1AddressesAddressPut_1**](AddressApi.md#v1AddressesAddressPut_1) | **PUT** /v1/addresses/{address} | Claim custody of an (address/chain) |



## v1AddressOwnershipProofsAddressPut_1

> CreateAddressOwnershipProofResponse v1AddressOwnershipProofsAddressPut_1(address, createAddressOwnershipProofRequest)

Submit a proof of address ownership for a previously claimed (address/chain)

This is the singular endpoint for creating an address ownership proof

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressApi apiInstance = new AddressApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address
        CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest = new CreateAddressOwnershipProofRequest(); // CreateAddressOwnershipProofRequest | 
        try {
            CreateAddressOwnershipProofResponse result = apiInstance.v1AddressOwnershipProofsAddressPut_1(address, createAddressOwnershipProofRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressApi#v1AddressOwnershipProofsAddressPut_1");
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


## v1AddressesAddressDelete_1

> DeleteAddressOwnershipResponse v1AddressesAddressDelete_1(address, chain)

Delete a Record of Address Custody

This endpoint deletes a record of address custody from your VASP&#39;s filter

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressApi apiInstance = new AddressApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            DeleteAddressOwnershipResponse result = apiInstance.v1AddressesAddressDelete_1(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressApi#v1AddressesAddressDelete_1");
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


## v1AddressesAddressGet_1

> GetAddressOwnershipResponse v1AddressesAddressGet_1(address, chain)

Retrieve the custody of an address and the proof of address ownership

This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressApi apiInstance = new AddressApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            GetAddressOwnershipResponse result = apiInstance.v1AddressesAddressGet_1(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressApi#v1AddressesAddressGet_1");
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


## v1AddressesAddressPut_1

> CreateAddressOwnershipResponse v1AddressesAddressPut_1(address, createAddressOwnershipRequest)

Claim custody of an (address/chain)

Claim custody of an (address/chain)

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressApi apiInstance = new AddressApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        CreateAddressOwnershipRequest createAddressOwnershipRequest = new CreateAddressOwnershipRequest(); // CreateAddressOwnershipRequest | 
        try {
            CreateAddressOwnershipResponse result = apiInstance.v1AddressesAddressPut_1(address, createAddressOwnershipRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressApi#v1AddressesAddressPut_1");
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

