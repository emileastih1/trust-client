package io.trust.client.starter;

import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Conditional;
import org.springframework.context.annotation.Import;

@AutoConfiguration
@EnableConfigurationProperties(TrustClientProperties.class)
@Conditional(TrustClientConfiguredCondition.class)
@Import(TrustClientFactoryBean.class)
public class TrustClientAutoConfiguration {
}
