package io.trust.client.api;

import io.trust.client.invoker.ApiClient;

import io.trust.client.model.HealthGet403Response;
import io.trust.client.model.HealthResponse;

import java.util.HashMap;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.web.client.RestClient.ResponseSpec;
import org.springframework.web.client.RestClientResponseException;
import org.springframework.core.io.FileSystemResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;

@jakarta.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2026-03-25T10:05:33.471345300+01:00[Europe/Paris]", comments = "Generator version: 7.9.0")
public class HealthApi {
    private ApiClient apiClient;

    public HealthApi() {
        this(new ApiClient());
    }

    @Autowired
    public HealthApi(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return apiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * Health check endpoint
     * This is the health check endpoint
     * <p><b>200</b> - OK
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return HealthResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec healthGetRequestCreation() throws RestClientResponseException {
        Object postBody = null;
        // create path and map variables
        final Map<String, Object> pathParams = new HashMap<>();

        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, String> cookieParams = new LinkedMultiValueMap<>();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<>();

        final String[] localVarAccepts = { 
            "application/json"
        };
        final List<MediaType> localVarAccept = apiClient.selectHeaderAccept(localVarAccepts);
        final String[] localVarContentTypes = { };
        final MediaType localVarContentType = apiClient.selectHeaderContentType(localVarContentTypes);

        String[] localVarAuthNames = new String[] {  };

        ParameterizedTypeReference<HealthResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/health", HttpMethod.GET, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Health check endpoint
     * This is the health check endpoint
     * <p><b>200</b> - OK
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return HealthResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public HealthResponse healthGet() throws RestClientResponseException {
        ParameterizedTypeReference<HealthResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return healthGetRequestCreation().body(localVarReturnType);
    }

    /**
     * Health check endpoint
     * This is the health check endpoint
     * <p><b>200</b> - OK
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return ResponseEntity&lt;HealthResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<HealthResponse> healthGetWithHttpInfo() throws RestClientResponseException {
        ParameterizedTypeReference<HealthResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return healthGetRequestCreation().toEntity(localVarReturnType);
    }

    /**
     * Health check endpoint
     * This is the health check endpoint
     * <p><b>200</b> - OK
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec healthGetWithResponseSpec() throws RestClientResponseException {
        return healthGetRequestCreation();
    }
}
