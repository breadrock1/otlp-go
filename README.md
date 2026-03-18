# Otlp-go library

## Overview

There is library to enable observability on rust service with following abilities:
- local and remote tracing (supporting jaeger);
- local and remote logging (supporting loki).

## Quick Start

1. Need to include this repository into project Cargo.toml manifest file:

   Install library using go get:
   ```shell
   go get <address-to-library>
    ```
   
    example:
    
    ```shell
    go get 192.168.0.84:3000/breadrock1/otlp-go
    ```

2. Past example code into main function (see [examples](examples/simple)):

3. Use it!
