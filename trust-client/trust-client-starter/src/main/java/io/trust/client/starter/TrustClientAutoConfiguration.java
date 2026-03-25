package io.trust.client.starter;

import io.trust.client.invoker.ApiClient;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.web.client.RestClient;

@AutoConfiguration
@EnableConfigurationProperties(TrustClientProperties.class)
public class TrustClientAutoConfiguration {

    @Bean
    @ConditionalOnMissingBean(ApiClient.class)
    public TrustClientFactoryBean trustApiClient(TrustClientProperties properties,
                                                 ObjectProvider<RestClient.Builder> restClientBuilderProvider,
                                                 ObjectMapper objectMapper) {
        return new TrustClientFactoryBean(properties, restClientBuilderProvider, objectMapper);
    }
}
