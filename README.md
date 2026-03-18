# Otlp-go library

## Overview

There is library to enable observability on rust service with following abilities:
- local and remote tracing (supporting jaeger);
- local and remote logging (supporting loki).

## Quick Start

1. Download library to project dir, for example into `third/`

    ```shell
    git clone <gitlab-url>/breadrock1/otlp-go.git third/otlp-go
    ```

2. Create go.work and include this library into use block:

    ```text
    go 1.25.0

    use (
        third/otlp-go
        .
    )
    ```

3. Update project state:
    ```shell
    go mod tidy
    ```

4. Past example code into main function (see [examples](examples/simple)):

5. Use it!
