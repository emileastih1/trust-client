package io.trust.client.starter;

import org.springframework.boot.autoconfigure.condition.ConditionMessage;
import org.springframework.boot.autoconfigure.condition.ConditionOutcome;
import org.springframework.boot.autoconfigure.condition.SpringBootCondition;
import org.springframework.boot.context.properties.bind.Binder;
import org.springframework.context.annotation.ConditionContext;
import org.springframework.core.type.AnnotatedTypeMetadata;
import org.springframework.util.StringUtils;

public class TrustClientConfiguredCondition extends SpringBootCondition {

    @Override
    public ConditionOutcome getMatchOutcome(ConditionContext context, AnnotatedTypeMetadata metadata) {
        ConditionMessage.Builder message = ConditionMessage.forCondition("TrustClientConfigured");
        
        String basePath = context.getEnvironment().getProperty("trust.client.base-path");
        if (!StringUtils.hasText(basePath)) {
            // If not found in environment, check if the default value would be used.
            // TrustClientProperties has a default value.
            // However, the requirement says "bind the trust.client properties into the TrustClientConfigurationProperties class and check only if the basePath is not empty to pass the condition."
            // This suggests we SHOULD use binding.
            try {
                TrustClientProperties properties = Binder.get(context.getEnvironment())
                        .bind("trust.client", TrustClientProperties.class)
                        .orElseGet(TrustClientProperties::new);
                basePath = properties.getBasePath();
            } catch (Exception e) {
                // If binding fails (e.g. incomplete yaml), we can fall back to a new instance to get defaults
                basePath = new TrustClientProperties().getBasePath();
            }
        }

        if (StringUtils.hasText(basePath)) {
            return ConditionOutcome.match(message.foundExactly("trust.client.base-path property"));
        }
        return ConditionOutcome.noMatch(message.didNotFind("non-empty trust.client.base-path property").atAll());
    }
}
