# HealthApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**healthGet**](HealthApi.md#healthGet) | **GET** /health | Health check endpoint |



## healthGet

> HealthResponse healthGet()

Health check endpoint

This is the health check endpoint

### Example

```java
// Import classes:
import io.trust.client.invoker.ApiClient;
import io.trust.client.invoker.ApiException;
import io.trust.client.invoker.Configuration;
import io.trust.client.invoker.models.*;
import io.trust.client.api.HealthApi;

public class Example {
    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://localhost");

        HealthApi apiInstance = new HealthApi(defaultClient);
        try {
            HealthResponse result = apiInstance.healthGet();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling HealthApi#healthGet");
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

[**HealthResponse**](HealthResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **403** | Not authorized |  -  |
| **500** | Internal service error |  -  |
| **0** | An unexpected error response. |  -  |

