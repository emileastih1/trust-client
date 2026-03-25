# AddressesApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**v1AddressesAddressDelete_0**](AddressesApi.md#v1AddressesAddressDelete_0) | **DELETE** /v1/addresses/{address} | Delete a Record of Address Custody |
| [**v1AddressesAddressGet_0**](AddressesApi.md#v1AddressesAddressGet_0) | **GET** /v1/addresses/{address} | Retrieve the custody of an address and the proof of address ownership |
| [**v1AddressesAddressPut_0**](AddressesApi.md#v1AddressesAddressPut_0) | **PUT** /v1/addresses/{address} | Claim custody of an (address/chain) |



## v1AddressesAddressDelete_0

> DeleteAddressOwnershipResponse v1AddressesAddressDelete_0(address, chain)

Delete a Record of Address Custody

This endpoint deletes a record of address custody from your VASP&#39;s filter

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressesApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressesApi apiInstance = new AddressesApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            DeleteAddressOwnershipResponse result = apiInstance.v1AddressesAddressDelete_0(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressesApi#v1AddressesAddressDelete_0");
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


## v1AddressesAddressGet_0

> GetAddressOwnershipResponse v1AddressesAddressGet_0(address, chain)

Retrieve the custody of an address and the proof of address ownership

This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressesApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressesApi apiInstance = new AddressesApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        String chain = "string"; // String | The blockchain of the address, e.g., BITCOIN, ETHEREUM.
        try {
            GetAddressOwnershipResponse result = apiInstance.v1AddressesAddressGet_0(address, chain);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressesApi#v1AddressesAddressGet_0");
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


## v1AddressesAddressPut_0

> CreateAddressOwnershipResponse v1AddressesAddressPut_0(address, createAddressOwnershipRequest)

Claim custody of an (address/chain)

Claim custody of an (address/chain)

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressesApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressesApi apiInstance = new AddressesApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address.
        CreateAddressOwnershipRequest createAddressOwnershipRequest = new CreateAddressOwnershipRequest(); // CreateAddressOwnershipRequest | 
        try {
            CreateAddressOwnershipResponse result = apiInstance.v1AddressesAddressPut_0(address, createAddressOwnershipRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressesApi#v1AddressesAddressPut_0");
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

