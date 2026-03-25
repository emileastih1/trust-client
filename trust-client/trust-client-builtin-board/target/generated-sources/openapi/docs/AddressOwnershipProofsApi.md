# AddressOwnershipProofsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**v1AddressOwnershipProofsAddressPut_0**](AddressOwnershipProofsApi.md#v1AddressOwnershipProofsAddressPut_0) | **PUT** /v1/address_ownership_proofs/{address} | Submit a proof of address ownership for a previously claimed (address/chain) |



## v1AddressOwnershipProofsAddressPut_0

> CreateAddressOwnershipProofResponse v1AddressOwnershipProofsAddressPut_0(address, createAddressOwnershipProofRequest)

Submit a proof of address ownership for a previously claimed (address/chain)

This is the singular endpoint for creating an address ownership proof

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.AddressOwnershipProofsApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        AddressOwnershipProofsApi apiInstance = new AddressOwnershipProofsApi(defaultClient);
        String address = "string"; // String | The SHA512 of the blockchain address
        CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest = new CreateAddressOwnershipProofRequest(); // CreateAddressOwnershipProofRequest | 
        try {
            CreateAddressOwnershipProofResponse result = apiInstance.v1AddressOwnershipProofsAddressPut_0(address, createAddressOwnershipProofRequest);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AddressOwnershipProofsApi#v1AddressOwnershipProofsAddressPut_0");
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

