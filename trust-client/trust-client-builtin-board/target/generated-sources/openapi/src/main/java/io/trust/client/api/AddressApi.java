package io.trust.client.api;

import io.trust.client.invoker.ApiClient;

import io.trust.client.model.CreateAddressOwnershipProofRequest;
import io.trust.client.model.CreateAddressOwnershipProofResponse;
import io.trust.client.model.CreateAddressOwnershipRequest;
import io.trust.client.model.CreateAddressOwnershipResponse;
import io.trust.client.model.DeleteAddressOwnershipResponse;
import io.trust.client.model.GetAddressOwnershipResponse;
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
public class AddressApi {
    private ApiClient apiClient;

    public AddressApi() {
        this(new ApiClient());
    }

    @Autowired
    public AddressApi(ApiClient apiClient) {
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
    private ResponseSpec v1AddressOwnershipProofsAddressPut_1RequestCreation(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        Object postBody = createAddressOwnershipProofRequest;
        // verify the required parameter 'address' is set
        if (address == null) {
            throw new RestClientResponseException("Missing the required parameter 'address' when calling v1AddressOwnershipProofsAddressPut_1", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
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
    public CreateAddressOwnershipProofResponse v1AddressOwnershipProofsAddressPut_1(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipProofResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressOwnershipProofsAddressPut_1RequestCreation(address, createAddressOwnershipProofRequest).body(localVarReturnType);
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
    public ResponseEntity<CreateAddressOwnershipProofResponse> v1AddressOwnershipProofsAddressPut_1WithHttpInfo(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipProofResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressOwnershipProofsAddressPut_1RequestCreation(address, createAddressOwnershipProofRequest).toEntity(localVarReturnType);
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
    public ResponseSpec v1AddressOwnershipProofsAddressPut_1WithResponseSpec(String address, CreateAddressOwnershipProofRequest createAddressOwnershipProofRequest) throws RestClientResponseException {
        return v1AddressOwnershipProofsAddressPut_1RequestCreation(address, createAddressOwnershipProofRequest);
    }
    /**
     * Delete a Record of Address Custody
     * This endpoint deletes a record of address custody from your VASP&#39;s filter
     * <p><b>200</b> - Delete succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return DeleteAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1AddressesAddressDelete_1RequestCreation(String address, String chain) throws RestClientResponseException {
        Object postBody = null;
        // verify the required parameter 'address' is set
        if (address == null) {
            throw new RestClientResponseException("Missing the required parameter 'address' when calling v1AddressesAddressDelete_1", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
        }
        // create path and map variables
        final Map<String, Object> pathParams = new HashMap<>();

        pathParams.put("address", address);

        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, String> cookieParams = new LinkedMultiValueMap<>();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<>();

        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "chain", chain));
        
        final String[] localVarAccepts = { 
            "application/json"
        };
        final List<MediaType> localVarAccept = apiClient.selectHeaderAccept(localVarAccepts);
        final String[] localVarContentTypes = { };
        final MediaType localVarContentType = apiClient.selectHeaderContentType(localVarContentTypes);

        String[] localVarAuthNames = new String[] {  };

        ParameterizedTypeReference<DeleteAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/addresses/{address}", HttpMethod.DELETE, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Delete a Record of Address Custody
     * This endpoint deletes a record of address custody from your VASP&#39;s filter
     * <p><b>200</b> - Delete succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return DeleteAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public DeleteAddressOwnershipResponse v1AddressesAddressDelete_1(String address, String chain) throws RestClientResponseException {
        ParameterizedTypeReference<DeleteAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressDelete_1RequestCreation(address, chain).body(localVarReturnType);
    }

    /**
     * Delete a Record of Address Custody
     * This endpoint deletes a record of address custody from your VASP&#39;s filter
     * <p><b>200</b> - Delete succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return ResponseEntity&lt;DeleteAddressOwnershipResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<DeleteAddressOwnershipResponse> v1AddressesAddressDelete_1WithHttpInfo(String address, String chain) throws RestClientResponseException {
        ParameterizedTypeReference<DeleteAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressDelete_1RequestCreation(address, chain).toEntity(localVarReturnType);
    }

    /**
     * Delete a Record of Address Custody
     * This endpoint deletes a record of address custody from your VASP&#39;s filter
     * <p><b>200</b> - Delete succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1AddressesAddressDelete_1WithResponseSpec(String address, String chain) throws RestClientResponseException {
        return v1AddressesAddressDelete_1RequestCreation(address, chain);
    }
    /**
     * Retrieve the custody of an address and the proof of address ownership
     * This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.
     * <p><b>200</b> - Get address ownership succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - Address not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return GetAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1AddressesAddressGet_1RequestCreation(String address, String chain) throws RestClientResponseException {
        Object postBody = null;
        // verify the required parameter 'address' is set
        if (address == null) {
            throw new RestClientResponseException("Missing the required parameter 'address' when calling v1AddressesAddressGet_1", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
        }
        // create path and map variables
        final Map<String, Object> pathParams = new HashMap<>();

        pathParams.put("address", address);

        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, String> cookieParams = new LinkedMultiValueMap<>();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<>();

        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "chain", chain));
        
        final String[] localVarAccepts = { 
            "application/json"
        };
        final List<MediaType> localVarAccept = apiClient.selectHeaderAccept(localVarAccepts);
        final String[] localVarContentTypes = { };
        final MediaType localVarContentType = apiClient.selectHeaderContentType(localVarContentTypes);

        String[] localVarAuthNames = new String[] {  };

        ParameterizedTypeReference<GetAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/addresses/{address}", HttpMethod.GET, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Retrieve the custody of an address and the proof of address ownership
     * This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.
     * <p><b>200</b> - Get address ownership succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - Address not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return GetAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public GetAddressOwnershipResponse v1AddressesAddressGet_1(String address, String chain) throws RestClientResponseException {
        ParameterizedTypeReference<GetAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressGet_1RequestCreation(address, chain).body(localVarReturnType);
    }

    /**
     * Retrieve the custody of an address and the proof of address ownership
     * This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.
     * <p><b>200</b> - Get address ownership succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - Address not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return ResponseEntity&lt;GetAddressOwnershipResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<GetAddressOwnershipResponse> v1AddressesAddressGet_1WithHttpInfo(String address, String chain) throws RestClientResponseException {
        ParameterizedTypeReference<GetAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressGet_1RequestCreation(address, chain).toEntity(localVarReturnType);
    }

    /**
     * Retrieve the custody of an address and the proof of address ownership
     * This endpoint retrieves the VASP who has custody of an address. You must query with the sha512 hash of the address. Additionally, a proof of ownership (or an IOU) is returned so VASP(S) can self-verify.
     * <p><b>200</b> - Get address ownership succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>404</b> - Address not found
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param chain The blockchain of the address, e.g., BITCOIN, ETHEREUM.
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1AddressesAddressGet_1WithResponseSpec(String address, String chain) throws RestClientResponseException {
        return v1AddressesAddressGet_1RequestCreation(address, chain);
    }
    /**
     * Claim custody of an (address/chain)
     * Claim custody of an (address/chain)
     * <p><b>200</b> - Claim address custody succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflict (when the address is already claimed by others)
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param createAddressOwnershipRequest The createAddressOwnershipRequest parameter
     * @return CreateAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    private ResponseSpec v1AddressesAddressPut_1RequestCreation(String address, CreateAddressOwnershipRequest createAddressOwnershipRequest) throws RestClientResponseException {
        Object postBody = createAddressOwnershipRequest;
        // verify the required parameter 'address' is set
        if (address == null) {
            throw new RestClientResponseException("Missing the required parameter 'address' when calling v1AddressesAddressPut_1", HttpStatus.BAD_REQUEST.value(), HttpStatus.BAD_REQUEST.getReasonPhrase(), null, null, null);
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

        ParameterizedTypeReference<CreateAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return apiClient.invokeAPI("/v1/addresses/{address}", HttpMethod.PUT, pathParams, queryParams, postBody, headerParams, cookieParams, formParams, localVarAccept, localVarContentType, localVarAuthNames, localVarReturnType);
    }

    /**
     * Claim custody of an (address/chain)
     * Claim custody of an (address/chain)
     * <p><b>200</b> - Claim address custody succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflict (when the address is already claimed by others)
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param createAddressOwnershipRequest The createAddressOwnershipRequest parameter
     * @return CreateAddressOwnershipResponse
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public CreateAddressOwnershipResponse v1AddressesAddressPut_1(String address, CreateAddressOwnershipRequest createAddressOwnershipRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressPut_1RequestCreation(address, createAddressOwnershipRequest).body(localVarReturnType);
    }

    /**
     * Claim custody of an (address/chain)
     * Claim custody of an (address/chain)
     * <p><b>200</b> - Claim address custody succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflict (when the address is already claimed by others)
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param createAddressOwnershipRequest The createAddressOwnershipRequest parameter
     * @return ResponseEntity&lt;CreateAddressOwnershipResponse&gt;
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseEntity<CreateAddressOwnershipResponse> v1AddressesAddressPut_1WithHttpInfo(String address, CreateAddressOwnershipRequest createAddressOwnershipRequest) throws RestClientResponseException {
        ParameterizedTypeReference<CreateAddressOwnershipResponse> localVarReturnType = new ParameterizedTypeReference<>() {};
        return v1AddressesAddressPut_1RequestCreation(address, createAddressOwnershipRequest).toEntity(localVarReturnType);
    }

    /**
     * Claim custody of an (address/chain)
     * Claim custody of an (address/chain)
     * <p><b>200</b> - Claim address custody succeeded
     * <p><b>400</b> - Invalid request, parameter validation failure.
     * <p><b>403</b> - Not authorized
     * <p><b>409</b> - Conflict (when the address is already claimed by others)
     * <p><b>500</b> - Internal service error
     * <p><b>0</b> - An unexpected error response.
     * @param address The SHA512 of the blockchain address.
     * @param createAddressOwnershipRequest The createAddressOwnershipRequest parameter
     * @return ResponseSpec
     * @throws RestClientResponseException if an error occurs while attempting to invoke the API
     */
    public ResponseSpec v1AddressesAddressPut_1WithResponseSpec(String address, CreateAddressOwnershipRequest createAddressOwnershipRequest) throws RestClientResponseException {
        return v1AddressesAddressPut_1RequestCreation(address, createAddressOwnershipRequest);
    }
}
