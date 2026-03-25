package io.trust.client.starter;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.trust.client.invoker.ApiClient;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.beans.factory.config.AbstractFactoryBean;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestClient;

@RequiredArgsConstructor
@ConditionalOnMissingBean(ApiClient.class)
public class TrustClientFactoryBean extends AbstractFactoryBean<ApiClient> {

    private final TrustClientProperties properties;
    private final ObjectProvider<RestClient.Builder> restClientBuilderProvider;
    private final ObjectMapper objectMapper;

    @Override
    public Class<?> getObjectType() {
        return ApiClient.class;
    }

    @Override
    protected ApiClient createInstance() throws Exception {
        RestClient.Builder builder = restClientBuilderProvider.getIfAvailable(RestClient::builder);
        builder.messageConverters(converters -> {
            converters.add(0, new MappingJackson2HttpMessageConverter(objectMapper));
        });

        ApiClient apiClient = new ApiClient(builder.build(), objectMapper, ApiClient.createDefaultDateFormat());
        apiClient.setBasePath(properties.getBasePath());
        if (properties.getDnsName() != null) {
            apiClient.addDefaultHeader("x-vasp-dns-name", properties.getDnsName());
            apiClient.addDefaultHeader("grpc-metadata-x-vasp-dns-name", properties.getDnsName());
        }
        properties.getAdditionalHeaders().forEach(apiClient::addDefaultHeader);
        return apiClient;
    }
}
