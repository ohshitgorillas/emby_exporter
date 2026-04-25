# Emby exporter

A Prometheus exporter for [Emby Media Server](https://emby.media/), written
in Go. Forked from
[williamclot/emby_exporter](https://github.com/williamclot/emby_exporter)
with substantially expanded metric coverage (per-session, per-transcode,
per-library, scheduled tasks, plugin updates, exporter health).

## Metrics

### Server / scrape health

| Metric | Type | Labels |
|---|---|---|
| `emby_info` | gauge=1 | `server`, `hostname`, `version` |
| `emby_up` | gauge | — |
| `emby_scrape_duration_seconds` | gauge | — |

### Library

| Metric | Type | Labels |
|---|---|---|
| `emby_movie_count` | gauge | `server` |
| `emby_series_count` | gauge | `server` |
| `emby_episode_count` | gauge | `server` |
| `emby_user_count` | gauge | `server` |
| `emby_device_count` | gauge | `server` |
| `emby_library_item_count` | gauge | `server`, `library`, `collection_type`, `type` |
| `emby_failed_login` | gauge | `server`, `user` |

### Sessions

| Metric | Type | Labels |
|---|---|---|
| `emby_stream_count` | gauge | `server` |
| `emby_session_info` | gauge=1 | `server`, `session_id`, `user`, `client`, `device`, `device_type`, `app_version`, `media_type`, `play_method`, `is_paused`, `item_type` |
| `emby_session_position_seconds` | gauge | `server`, `session_id`, `user` |
| `emby_session_runtime_seconds` | gauge | `server`, `session_id`, `user` |
| `emby_video_bitrate` | gauge | `server`, `user` |
| `emby_audio_bitrate` | gauge | `server`, `user` |
| `emby_audio_channels` | gauge | `server`, `user` |

### Transcoding

| Metric | Type | Labels |
|---|---|---|
| `emby_transcoding_count` | gauge | `server` |
| `emby_transcoding_bitrate` | gauge | `server`, `user` |
| `emby_transcode_info` | gauge=1 | `server`, `session_id`, `user`, `hw_decoder`, `hw_encoder`, `video_codec`, `audio_codec`, `container`, `reasons` |
| `emby_transcode_bitrate` | gauge | `server`, `session_id`, `user` |
| `emby_transcode_completion_percent` | gauge | `server`, `session_id`, `user` |
| `emby_transcode_cpu_usage` | gauge | `server`, `session_id`, `user` |
| `emby_transcode_framerate` | gauge | `server`, `session_id`, `user` |
| `emby_transcode_width_pixels` | gauge | `server`, `session_id`, `user` |
| `emby_transcode_height_pixels` | gauge | `server`, `session_id`, `user` |

### Scheduled tasks and plugins

| Metric | Type | Labels |
|---|---|---|
| `emby_scheduled_task_info` | gauge=1 | `server`, `task`, `category`, `state`, `last_status` |
| `emby_scheduled_task_last_duration_seconds` | gauge | `server`, `task` |
| `emby_scheduled_task_last_run_timestamp_seconds` | gauge | `server`, `task` |
| `emby_plugin_update_available` | gauge=1 | `server`, `name`, `version` |

## Usage

### Options

```
emby_exporter [options]

  --telemetry.addr   listen address              (default ":9162")
  --telemetry.path   metrics URL path            (default "/metrics")
  --emby.addr        URL of the Emby HTTP API    (required)
  --emby.token       Emby API token              (required)
  --emby.verifyTLS   verify Emby TLS certificate (default true)
  --health           run a healthcheck and exit
  -h, --help         show help
```

### Docker

Pre-built multi-arch images (`linux/amd64`, `linux/arm64`) are published to
GitHub Container Registry on every push to `main` and on tags:

```
ghcr.io/ohshitgorillas/emby_exporter:latest
ghcr.io/ohshitgorillas/emby_exporter:0.1.0
```

### docker-compose

```yaml
services:
  emby_exporter:
    image: ghcr.io/ohshitgorillas/emby_exporter:latest
    container_name: emby_exporter
    command:
      - --emby.addr=http://emby:8096
      - --emby.token=${EMBY_TOKEN}
    ports:
      - "127.0.0.1:9162:9162"
    healthcheck:
      test: ["CMD", "/bin/emby_exporter", "-health"]
      interval: 30s
      timeout: 5s
      retries: 3
    restart: unless-stopped
```

### Prometheus scrape config

```yaml
scrape_configs:
  - job_name: emby
    scrape_interval: 30s
    static_configs:
      - targets: ['emby_exporter:9162']
```

## Building from source

```
go build -mod=vendor ./...
go test -mod=vendor ./...
```

## License

Same as upstream (see `LICENSE`).
