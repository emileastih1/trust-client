package io.trust.client.api;

import io.trust.client.invoker.ApiClient;

import io.trust.client.model.GetVASPResponse;
import io.trust.client.model.GetVASPResponse1;
import io.trust.client.model.HealthGet403Response;

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

@jakarta.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2026-03-25T10:19:00.496410100+01:00[Europe/Paris]", comments = "Generator version: 7.9.0")
public class VaspsApi {
    private ApiClient apiClient;

    public VaspsApi() {
        this(new ApiClient());
    }

    @Autowired
    public VaspsApi(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return apiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * Retrieve the information of all VASPs
     * Retrieve the information of all VASPs.
     * <p><b>200</b> - Successfully get the information of all VASPs
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return GetVASPResponse1
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1VaspsGet_0RequestCreation() throws RestClientResponseException {
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

        ParameterizedTypeReference<GetVASPResponse1> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/vasps", HttpMethod.GET, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Retrieve the information of all VASPs
     * Retrieve the information of all VASPs.
     * <p><b>200</b> - Successfully get the information of all VASPs
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return GetVASPResponse1
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public GetVASPResponse1 v1VaspsGet_0() throws RestClientResponseException {
        ParameterizedTypeReference<GetVASPResponse1> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1VaspsGet_0RequestCreation().body(localVarReturnType);
    }

    /**
     * Retrieve the information of all VASPs
     * Retrieve the information of all VASPs.
     * <p><b>200</b> - Successfully get the information of all VASPs
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return ResponseEntity&lt;GetVASPResponse1&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<GetVASPResponse1> v1VaspsGet_0WithHttpInfo() throws RestClientResponseException {
        ParameterizedTypeReference<GetVASPResponse1> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1VaspsGet_0RequestCreation().toEntity(localVarReturnType);
    }

    /**
     * Retrieve the information of all VASPs
     * Retrieve the information of all VASPs.
     * <p><b>200</b> - Successfully get the information of all VASPs
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1VaspsGet_0WithResponseSpec() throws RestClientResponseException {
        return v1VaspsGet_0RequestCreation();
    }
    /**
     * Retrieve the information of a specific VASP
     * This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   
     * <p><b>200</b> - Successfully retrieve the information of the VASP
     * <p><b>400</b> - Invalid request, parameter validation failure. For example, the provided vasp id is not a valid UUID
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - VASP not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param vaspId The UUID of the VASP that your are querying for
     * @return GetVASPResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1VaspsVaspIdGet_0RequestCreation(String vaspId) throws RestClientResponseException {
        Object postBody = null;
        // verify the required parameter 'vaspId' is set
        if (vaspId == null) {
            throw new RestClientResponseException("Missing the required parameter 'vaspId' when calling v1VaspsVaspIdGet_0", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
        }
        // create path and map variables
        final Map<String, Object> pathParams = new HashMap<>();

        pathParams.put("vaspId", vaspId);

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

        ParameterizedTypeReference<GetVASPResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/vasps/{vaspId}", HttpMethod.GET, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Retrieve the information of a specific VASP
     * This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   
     * <p><b>200</b> - Successfully retrieve the information of the VASP
     * <p><b>400</b> - Invalid request, parameter validation failure. For example, the provided vasp id is not a valid UUID
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - VASP not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param vaspId The UUID of the VASP that your are querying for
     * @return GetVASPResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public GetVASPResponse v1VaspsVaspIdGet_0(String vaspId) throws RestClientResponseException {
        ParameterizedTypeReference<GetVASPResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1VaspsVaspIdGet_0RequestCreation(vaspId).body(localVarReturnType);
    }

    /**
     * Retrieve the information of a specific VASP
     * This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   
     * <p><b>200</b> - Successfully retrieve the information of the VASP
     * <p><b>400</b> - Invalid request, parameter validation failure. For example, the provided vasp id is not a valid UUID
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - VASP not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param vaspId The UUID of the VASP that your are querying for
     * @return ResponseEntity&lt;GetVASPResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<GetVASPResponse> v1VaspsVaspIdGet_0WithHttpInfo(String vaspId) throws RestClientResponseException {
        ParameterizedTypeReference<GetVASPResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1VaspsVaspIdGet_0RequestCreation(vaspId).toEntity(localVarReturnType);
    }

    /**
     * Retrieve the information of a specific VASP
     * This endpoint retrieves the information of a specific VASP.  VASP(S) can query the GET /addresses endpoint to know the VASP(R) who claims the custody of an address.  Once VASP(S) knows VASP(R)&#39;s identity and verifies the Proof of Address Ownership, it can query this endpoint to know VASP(R)&#39;s PII endpoint.   The PII information should be encrypted with the VASP(R)&#39;s publicKey using the JOSE spec before being transmitted to VASP(R)&#39;s PII endpoint.   
     * <p><b>200</b> - Successfully retrieve the information of the VASP
     * <p><b>400</b> - Invalid request, parameter validation failure. For example, the provided vasp id is not a valid UUID
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - VASP not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param vaspId The UUID of the VASP that your are querying for
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1VaspsVaspIdGet_0WithResponseSpec(String vaspId) throws RestClientResponseException {
        return v1VaspsVaspIdGet_0RequestCreation(vaspId);
    }
}
