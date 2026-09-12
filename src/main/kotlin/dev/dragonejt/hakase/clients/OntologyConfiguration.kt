package dev.dragonejt.hakase.clients

import com.palantir.osdk.api.UserTokenAuth
import com.palantir.osdk.internal.api.FoundryConnectionConfig
import dev.dragonejt.hakase.telemetry.LogBase
import dev.dragonejt.hakase_sdk.FoundryClient
import dev.dragonejt.hakase_sdk.Ontology5a5029d53a2342fdA36a2cf9e72e73b9 as Ontology
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "ontology")
data class OntologyProperties(val url: String, val clientID: String, val token: String)

@Configuration
@EnableConfigurationProperties(OntologyProperties::class)
class OntologyConfiguration : LogBase() {

    @Bean
    fun ontology(props: OntologyProperties): Ontology {
        log.atDebug {
            message = "Building Palantir Ontology SDK Client with URL: ${props.url} and User Token"
            payload = mapOf("foundry_uri" to props.url)
        }
        val auth = UserTokenAuth.builder().token(props.token).build()
        val connection = FoundryConnectionConfig.builder().foundryUri(props.url).build()

        return FoundryClient.builder().auth(auth).connectionConfig(connection).build().ontology()
    }
}
