# Go Mirror Proxy

This project implements a gateway that proxies requests to a main backend service while mirroring the same traffic to multiple other backend services. This is often referred to as "shadowing."

## Architecture

The system is orchestrated using `docker-compose` and consists of the following services:

-   **`gateway`**: The main entry point for all traffic. It listens on port `8080`.
-   **`tcp` (Main Proxy)**: The primary backend service that receives the original request from the gateway. Its response is the one sent back to the original caller.
-   **`http` (Mirror)**: A secondary service that receives a copy (a mirror) of the incoming traffic from the gateway. Its response is ignored by the gateway.
-   **`reverse` (Mirror)**: Another secondary service that also receives a mirrored copy of the traffic for shadowing purposes. Its response is also ignored.

## Traffic Flow

1.  A client sends a request to the `gateway` on port `8080`.
2.  The `gateway` forwards the request to the main proxy service (`tcp`).
3.  Simultaneously, the `gateway` sends copies of the request to all configured mirror services (`http` and `reverse`).
4.  The `gateway` receives a response from the `tcp` service and forwards it back to the client.
5.  Responses from the mirror services are discarded.

This setup is useful for testing new backend services with production traffic without affecting the user-facing responses.
