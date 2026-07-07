package dev.dragonejt.hakase.clients

import com.palantir.osdk.api.auth.ConfidentialClientAuth
import com.palantir.osdk.internal.api.FoundryConnectionConfig
import dev.dragonejt.hakase.telemetry.LogBase
import dev.dragonejt.hakase_sdk.FoundryClient
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@ConfigurationProperties(prefix = "discord")
data class OntologyProperties(val url: String, val token: String)

@Configuration
class OntologyConfiguration : LogBase() {

    private val clientID = "33976368fd8b43e534fd9e4b0a51c4c0"

    @Bean
    fun ontologyClient(props: OntologyProperties): FoundryClient {
        log.atDebug {
            message =
                "Building Palantir Ontology SDK Client with URL: ${props.url} and client ID: $clientID"
            payload = mapOf("foundry_uri" to props.url, "foundry_client_id" to clientID)
        }
        val auth =
            ConfidentialClientAuth.builder().clientId(clientID).clientSecret(props.token).build()
        val connection = FoundryConnectionConfig.builder().foundryUri(props.url).build()

        return FoundryClient.builder().auth(auth).connectionConfig(connection).build()
    }
}
