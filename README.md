# word-of-wisdom

This application implements a TCP server designed to provide inspirational quotes from the "Word of Wisdom" book. The server is protected from potential DDOS attacks using the Proof of Work (PoW) consensus mechanism. This ensures that clients must perform computational work before gaining access to the server.

## Features

- TCP server with PoW protection against DDOS attacks.
- Upon successful PoW verification, the server sends an inspirational quote.

## Proof of Work (PoW) Algorithm

The chosen PoW algorithm is (SHA-256)[https://en.wikipedia.org/wiki/SHA-2] due to its robustness and widespread adoption as a cryptographic hash function. This algorithm provides a strong foundation for securing the blockchain network by offering computational efficiency and a high level of cryptographic security.

## Server Usage

To run the server, follow these steps:

1. Build the Docker image for the server using the provided Dockerfile.
2. Start the Docker container with the server image.
3. The server will be accessible on port `11001`.

## Client Usage

To interact with the server, you will need to implement a client that can solve the PoW challenge. The client should perform the necessary computations and send the solution to the server for verification.

## Getting Started

To deploy the Word of Wisdom TCP server and client, follow these steps:

1. Clone this repository to your local machine:
2. Build the Docker images for both the server and client using the provided docker-compose file:
    ```bash
    make env
    ```
    This will build the Docker images for the server and client
3. Once the server and client containers are running, the client will automatically execute 10 requests to the server. Each request will involve solving a Proof of Work challenge.
4. After completing the 10 requests, the client will stop automatically.
---
