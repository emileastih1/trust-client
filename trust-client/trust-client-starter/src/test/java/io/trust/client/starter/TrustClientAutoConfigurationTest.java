package io.trust.client.starter;

import io.trust.client.invoker.ApiClient;
import org.junit.jupiter.api.Test;
import org.springframework.boot.autoconfigure.AutoConfigurations;
import org.springframework.boot.env.YamlPropertySourceLoader;
import org.springframework.boot.autoconfigure.jackson.JacksonAutoConfiguration;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;

class TrustClientAutoConfigurationTest {

    private final ApplicationContextRunner contextRunner = new ApplicationContextRunner()
            .withConfiguration(AutoConfigurations.of(TrustClientAutoConfiguration.class, JacksonAutoConfiguration.class));

    @Test
    void shouldCreateApiClientWithDefaultConfig() {
        this.contextRunner.run(context -> {
            assertThat(context).hasSingleBean(ApiClient.class);
            ApiClient apiClient = context.getBean(ApiClient.class);
            assertThat(apiClient.getBasePath()).isEqualTo("http://localhost:7000");
        });
    }

    @Test
    void shouldCreateApiClientWithCompleteConfig() {
        this.contextRunner
                .withInitializer(context -> {
                    YamlPropertySourceLoader loader = new YamlPropertySourceLoader();
                    Resource resource = new ClassPathResource("application-complete-config.yaml");
                    try {
                        loader.load("complete-config", resource).forEach(context.getEnvironment().getPropertySources()::addLast);
                    } catch (IOException e) {
                        throw new RuntimeException(e);
                    }
                })
                .run(context -> {
                    assertThat(context).hasSingleBean(ApiClient.class);
                    ApiClient apiClient = context.getBean(ApiClient.class);
                    assertThat(apiClient.getBasePath()).isEqualTo("https://api.trust-complete.io");

                    assertThat(context).hasSingleBean(TrustClientProperties.class);
                    TrustClientProperties properties = context.getBean(TrustClientProperties.class);
                    assertThat(properties.getDnsName()).isEqualTo("testvasp1.com");
                    assertThat(properties.getAdditionalHeaders()).containsEntry("X-Custom-Header", "CustomValue");
                });
    }

    @Test
    void shouldCreateApiClientWithIncompleteConfigUsingDefaults() {
        this.contextRunner
                .withInitializer(context -> {
                    YamlPropertySourceLoader loader = new YamlPropertySourceLoader();
                    Resource resource = new ClassPathResource("application-incomplete-config.yaml");
                    try {
                        loader.load("incomplete-config", resource).forEach(context.getEnvironment().getPropertySources()::addLast);
                    } catch (IOException e) {
                        throw new RuntimeException(e);
                    }
                })
                .run(context -> {
                    assertThat(context).hasSingleBean(ApiClient.class);
                    ApiClient apiClient = context.getBean(ApiClient.class);
                    assertThat(apiClient.getBasePath()).isEqualTo("http://localhost:7000");
                });
    }

    @Test
    void shouldOverrideBasePathFromProperties() {
        this.contextRunner
                .withPropertyValues("trust.client.base-path=https://api.example.com")
                .run(context -> {
                    assertThat(context).hasSingleBean(ApiClient.class);
                    ApiClient apiClient = context.getBean(ApiClient.class);
                    assertThat(apiClient.getBasePath()).isEqualTo("https://api.example.com");
                });
    }

    @Test
    void shouldNotCreateApiClientIfBeanAlreadyExists() {
        this.contextRunner
                .withBean("customApiClient", ApiClient.class, () -> {
                    ApiClient custom = new ApiClient();
                    custom.setBasePath("https://custom.io");
                    return custom;
                })
                .run(context -> {
                    assertThat(context).hasSingleBean(ApiClient.class);
                    assertThat(context).hasBean("customApiClient");
                    ApiClient apiClient = context.getBean(ApiClient.class);
                    assertThat(apiClient.getBasePath()).isEqualTo("https://custom.io");
                });
    }
}
