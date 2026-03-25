# VaspIdApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**v1VaspsVaspIdGet_1**](VaspIdApi.md#v1VaspsVaspIdGet_1) | **GET** /v1/vasps/{vaspId} | Retrieve the information of a specific VASP |



## v1VaspsVaspIdGet_1

> GetVASPResponse v1VaspsVaspIdGet_1(vaspId)

Retrieve the information of a specific VASP

This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.VaspIdApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        VaspIdApi apiInstance = new VaspIdApi(defaultClient);
        String vaspId = "string"; // String | The UUID of the VASP that your are querying for
        try {
            GetVASPResponse result = apiInstance.v1VaspsVaspIdGet_1(vaspId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling VaspIdApi#v1VaspsVaspIdGet_1");
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

