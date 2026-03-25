package io.trust.client.api;

import io.trust.client.invoker.ApiClient;

import io.trust.client.model.CreateAddressOwnershipProofRequest;
import io.trust.client.model.CreateAddressOwnershipProofResponse;
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

@jakarta.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2026-03-25T10:05:33.471345300+01:00[Europe/Paris]", comments = "Generator version: 7.9.0")
public class AddressOwnershipProofsApi {
    private ApiClient apiClient;

    public AddressOwnershipProofsApi() {
        this(new ApiClient());
    }

    @Autowired
    public AddressOwnershipProofsApi(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return apiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * Submit a proof of address ownership for a previously claimed (address/chain)
     * This is the singular endpoint for creating an address ownership proof
     * <p><b>200</b> - Create or update proof of address succeeded
     * <p><b>400</b> - Invalid request. For example, the (address, chain, registration_id) combination is invalid or missing fields.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflicted proof status, e.g., change from non-IOU to IOU
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address
     * @param createAddressOwnershipProofRequest The createAddressOwnershipProofRequest parameter
     * @return CreateAddressOwnershipProofResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1AddressOwnershipProofsAddressPut_0RequestCreation(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        Object postBody = createAddressOwnershipProofRequest;
        // verify the required parameter 'address' is set
        if (address == null) {
            throw new RestClientResponseException("Missing the required parameter 'address' when calling v1AddressOwnershipProofsAddressPut_0", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
        }
        // create path and map variables
        final Map<String, Object> pathParams = new HashMap<>();

        pathParams.put("address", address);

        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, String> cookieParams = new LinkedMultiValueMap<>();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<>();

        final String[] localVarAccepts = { 
            "application/json"
        };
        final List<MediaType> localVarAccept = apiClient.selectHeaderAccept(localVarAccepts);
        final String[] localVarContentTypes = { 
            "application/json"
        };
        final MediaType localVarContentType = apiClient.selectHeaderContentType(localVarContentTypes);

        String[] localVarAuthNames = new String[] {  };

        ParameterizedTypeReference<CreateAddressOwnershipProofResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/address_ownership_proofs/{address}", HttpMethod.PUT, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Submit a proof of address ownership for a previously claimed (address/chain)
     * This is the singular endpoint for creating an address ownership proof
     * <p><b>200</b> - Create or update proof of address succeeded
     * <p><b>400</b> - Invalid request. For example, the (address, chain, registration_id) combination is invalid or missing fields.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflicted proof status, e.g., change from non-IOU to IOU
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address
     * @param createAddressOwnershipProofRequest The createAddressOwnershipProofRequest parameter
     * @return CreateAddressOwnershipProofResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public CreateAddressOwnershipProofResponse v1AddressOwnershipProofsAddressPut_0(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipProofResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressOwnershipProofsAddressPut_0RequestCreation(address, createAddressOwnershipProofRequest).body(localVarReturnType);
    }

    /**
     * Submit a proof of address ownership for a previously claimed (address/chain)
     * This is the singular endpoint for creating an address ownership proof
     * <p><b>200</b> - Create or update proof of address succeeded
     * <p><b>400</b> - Invalid request. For example, the (address, chain, registration_id) combination is invalid or missing fields.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflicted proof status, e.g., change from non-IOU to IOU
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address
     * @param createAddressOwnershipProofRequest The createAddressOwnershipProofRequest parameter
     * @return ResponseEntity&lt;CreateAddressOwnershipProofResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<CreateAddressOwnershipProofResponse> v1AddressOwnershipProofsAddressPut_0WithHttpInfo(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipProofResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressOwnershipProofsAddressPut_0RequestCreation(address, createAddressOwnershipProofRequest).toEntity(localVarReturnType);
    }

    /**
     * Submit a proof of address ownership for a previously claimed (address/chain)
     * This is the singular endpoint for creating an address ownership proof
     * <p><b>200</b> - Create or update proof of address succeeded
     * <p><b>400</b> - Invalid request. For example, the (address, chain, registration_id) combination is invalid or missing fields.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflicted proof status, e.g., change from non-IOU to IOU
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address
     * @param createAddressOwnershipProofRequest The createAddressOwnershipProofRequest parameter
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1AddressOwnershipProofsAddressPut_0WithResponseSpec(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        return v1AddressOwnershipProofsAddressPut_0RequestCreation(address, createAddressOwnershipProofRequest);
    }
}
