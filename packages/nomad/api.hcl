job "api" {
  datacenters = ["${gcp_zone}"]
  node_pool = "api"
  priority = 90

  group "api-service" {
    network {
      port "api" {
        static = "${port_number}"
      }

      %{ if prevent_colocation }
      port "scheduling-block" {
        // This port is used to block scheduling of jobs with the same block on the same node.
        // We use this to block API and Loki from being scheduled on the same node.
        static = 40234
      }
      %{ endif }
    }

    constraint {
      operator  = "distinct_hosts"
      value     = "true"
    }

    service {
      name = "api"
      port = "${port_number}"
      task = "start"

      check {
        type     = "http"
        name     = "health"
        path     = "/health"
        interval = "3s"
        timeout  = "3s"
        port     = "${port_number}"
      }
    }

%{ if update_stanza }
    # An update stanza to enable rolling updates of the service
    update {
      # The number of extra instances to run during the update
      max_parallel     = 1
      # Allows to spawn new version of the service before killing the old one
      canary           = 1
      # Time to wait for the canary to be healthy
      min_healthy_time = "10s"
      # Time to wait for the canary to be healthy, if not it will be marked as failed
      healthy_deadline = "30s"
      # Whether to promote the canary if the rest of the group is not healthy
      auto_promote     = true
    }
%{ endif }

    task "start" {
      driver       = "docker"
      # If we need more than 30s we will need to update the max_kill_timeout in nomad
      # https://developer.hashicorp.com/nomad/docs/configuration/client#max_kill_timeout
      kill_timeout = "30s"
      kill_signal  = "SIGTERM"

      resources {
        memory_max = 4096
        memory     = 2048
        cpu        = 2000
      }

      env {
        ORCHESTRATOR_PORT              = "${orchestrator_port}"
        TEMPLATE_MANAGER_HOST          = "${template_manager_host}"
        POSTGRES_CONNECTION_STRING     = "${postgres_connection_string}"
        SUPABASE_JWT_SECRETS           = "${supabase_jwt_secrets}"
        CLICKHOUSE_CONNECTION_STRING   = "${clickhouse_connection_string}"
        ENVIRONMENT                    = "${environment}"
        POSTHOG_API_KEY                = "${posthog_api_key}"
        ANALYTICS_COLLECTOR_HOST       = "${analytics_collector_host}"
        ANALYTICS_COLLECTOR_API_TOKEN  = "${analytics_collector_api_token}"
        LOKI_ADDRESS                   = "${loki_address}"
        OTEL_TRACING_PRINT             = "${otel_tracing_print}"
        LOGS_COLLECTOR_ADDRESS         = "${logs_collector_address}"
        NOMAD_TOKEN                    = "${nomad_acl_token}"
        OTEL_COLLECTOR_GRPC_ENDPOINT   = "${otel_collector_grpc_endpoint}"
        ADMIN_TOKEN                    = "${admin_token}"
        REDIS_URL                      = "${redis_url}"
        REDIS_CLUSTER_URL              = "${redis_cluster_url}"
        DNS_PORT                       = "${dns_port_number}"
        SANDBOX_ACCESS_TOKEN_HASH_SEED = "${sandbox_access_token_hash_seed}"


%{ if launch_darkly_api_key != "" }
        LAUNCH_DARKLY_API_KEY         = "${launch_darkly_api_key}"
%{ endif }

        # This is here just because it is required in some part of our code which is transitively imported
        TEMPLATE_BUCKET_NAME          = "skip"
      }

      config {
        network_mode = "host"
        image        = "${api_docker_image}"
        ports        = ["${port_name}"]
        args         = [
          "--port", "${port_number}",
        ]
      }
    }

    task "db-migrator" {
      driver = "docker"

      env {
        POSTGRES_CONNECTION_STRING="${postgres_connection_string}"
      }

      config {
        image = "${db_migrator_docker_image}"
      }

      resources {
        cpu    = 250
        memory = 128
      }

      lifecycle {
        hook = "prestart"
        sidecar = false
      }
    }
  }
}
