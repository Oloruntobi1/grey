## 1a. GitHub Link

https://github.com/Oloruntobi1/grey

## 1b. Thought Process Behind Code Application Structure

This is a sample Golang application designed following Go community best practices. I have used and practiced these principles extensively, continually updating myself with new insights and adding my own logical and objective preferences developed over many years.

Below is a skeletal diagram of the directory structure I used:

### Top-Level Folders

**appconstants**: This folder typically contains a file that holds all the constants used throughout the application.

**cmd**: This folder contains the entry points to all the runnable applications attached to the project. In my submitted example, it includes the wallet-app server and a seeder application used to populate the database with dummy data for local testing.

**internal**: Widely used in the Go community, this folder stores business logic and other modules intended for internal use only. Modules in this folder cannot be used by other applications, which is beneficial for applications running in a microservices environment.

- **internal/config**: Stores configurations used throughout the project.
- **internal/db**: Contains migrations, SQL files, and Go files related to database work. Subfolders include:
  - `internal/db/migrations`
  - `internal/db/query`
  - `internal/db/sqlc`
- **internal/models**: Contains data models based on use cases for the application.
- **internal/repositories**: Abstracts interaction with a database or datastore. Contains files:
  - `internal/repositories/user.go`
  - `internal/repositories/wallet.go`
- **internal/transport**: Includes the transport layers used in the application, such as HTTP and gRPC. Subfolders include:
  - `internal/transport/grpc`
  - `internal/transport/http`
    - `internal/transport/http/domains`
      - `internal/transport/http/domains/users`
      - `internal/transport/http/domains/wallets`
    - `internal/transport/http/handlers`

Each of these subfolders contains its own domains and handlers. This structure supports flexibility and loose coupling. The domain folders include users, wallets, transfers, and transactions. The grpc folder is intentionally left empty.

**pkg**: Contains packages like logger, metrics, and other observability tools, tailored to the application's needs. Subfolders include:
- `pkg/logger`
- `pkg/metrics`
- `pkg/otel`
- `pkg/tracer`

**providers**: Contains client modules used to communicate with external providers, such as banks. This folder is left empty by design.

**tests**: Contains integration tests. This folder is left empty by design.


### Extra Files

- **pkg/otel/confs**: Contains necessary configuration files for OtelCollector, Prometheus, and Grafana. These tools are used for observability and monitoring.
  - **OtelCollector**: Collects and exports telemetry data such as traces and metrics.
  - **Prometheus**: A monitoring system and time-series database.
  - **Grafana**: A data visualization tool that integrates with Prometheus.
- **.env**: Stores environment variables needed for the application. This file should be stored securely and not left in Git, cloud, or any public space to prevent sensitive data leaks.
- **.gitignore**: Specifies intentionally untracked files to ignore in Git. This prevents certain files from being tracked and uploaded to the repository.
- **go.mod**: Defines the module's dependencies and the versions of the modules required for the project.
- **go.sum**: Contains the expected cryptographic checksums of the content of specific module versions.
- **Makefile**: Contains a set of directives used with the `make` build automation tool to automate the compilation and management tasks.
- **README.md**: Provides an overview of the project, including installation instructions, usage, and other relevant information.
- **sqlc.yaml**: Configuration file for `sqlc`, a tool used to generate Go code from SQL queries.

### Things Undone

- No links to Postman Collection or Open API spec.
- Not checking for HTTP methods in requests.
- Not returning trace IDs to clients when the server returns errors.
- No example metrics set up in the code.
- Left `dev.env` on purpose for testing.
- etc

## How To Run
Please read `run.md`.

## 2

## Architecture Diagram

Below is a link to a high-level architecture diagram of the wallet system:
https://miro.com/app/board/uXjVK-338So=/

## Detailed Explanation of Components

### Load Balancer
- **Purpose**: Distributes incoming traffic across multiple instances of microservices, preventing any single service from becoming a bottleneck and ensuring high availability.
- **Examples**: NGINX, HAProxy, AWS ELB.

### API Gateway

- **Purpose**: Manages and routes external requests to the appropriate microservices, enforces security policies, and provides rate limiting.
- **Examples**: Kong, NGINX, AWS API Gateway.

### Microservices Architecture

- **Purpose**: Allows each service to be developed, deployed, and scaled independently, improving the overall system's flexibility and resilience.
- **Services**: 
  - **User Service**: Manages user data and authentication.
  - **Wallet Service**: Handles wallet operations including balance updates.
  - **Transaction Service**: Processes transactions ensuring atomicity and consistency.
  - **Notification Service**: Sends out notifications regarding transaction status amongst other things.
  - **Analytics Service**: Gathers and analyzes transaction data.

### Database

- **Purpose**: Ensures data persistence and consistency.
- **Components**:
  - **SQL Database**: Ensures strong consistency and ACID transactions for critical data. Example: PostgreSQL.
  - **NoSQL Database**: Handles high write throughput and provides eventual consistency for less critical data. Example: Cassandra.

### Caching Layer

- **Purpose**: Reduces load on the database by caching frequently accessed data, improving response times.
- **Examples**: Redis, Memcached.

### Message Queue

- **Purpose**: Decouples services and enables asynchronous processing, improving fault tolerance and system reliability.
- **Examples**: RabbitMQ, Apache Kafka.

### Distributed Ledger

- **Purpose**: Provides a tamper-proof record of all transactions, ensuring high consistency and auditability.
- **Examples**: AWS QLDB.

### Observability, Telemetry and Logging

- **Purpose**: Provides insights into system performance and health, helping quickly identify and resolve issues.
- **Examples**: Prometheus (monitoring), Jaeger(Traces) ELK Stack (logging).

### Orchestration and Containerization

- **Purpose**: Manages the deployment, scaling, and operations of containerized applications, ensuring consistency across environments.
- **Examples**: Kubernetes, Docker.

### Backup and Recovery

- **Purpose**: Ensures data durability and availability by regularly backing up data and enabling quick recovery in case of failures.

## 3. Diagnosing a Slow Application

To diagnose a slow application, I typically start by introducing tracing using open-source libraries compatible with the application's programming language. Here is a step-by-step approach I typically follow:

1. **Introduce Tracing**:
   - Use libraries such as OpenTelemetry to instrument your code for tracing.
   - Capture traces of your application to understand where time is being spent during execution.

2. **Visualize Traces**:
   - Use a tool like Jaeger, which provides a clean and intuitive UI to display the captured traces.
   - Jaeger helps in visualizing the end-to-end traces, making it easier to identify bottlenecks.

3. **Identify Bottlenecks**:
   - Analyze the traces to see which functions or operations are taking the most time.
   - Look for patterns such as repeated long-running operations or high latency in specific areas.

4. **Run Benchmarks**:
   - Write and execute unit tests with benchmarking to measure the performance of different parts of the application.
   - Use Go’s testing framework to create benchmarks.

5. **Improve Logging**:
   - Enhance the logging in your application to capture more detailed information about the application’s runtime behavior.
   - Ensure logs include context such as request IDs or user IDs to correlate logs with specific traces.

6. **Profile the Application**:
   - Use `pprof`, Go’s profiling tool, to analyze memory profiles, CPU profiles, and other necessary profiles.
   - `pprof` helps in understanding memory usage, CPU usage, and identifying hotspots in the code.

## 4. Understanding Distributed Consensus

Distributed consensus is a complex topic, but it is crucial for ensuring consistency across multiple systems or nodes. Here is a more detailed explanation:

1. **Concept Overview**:
   - Distributed consensus involves multiple computers (or nodes) working together to agree on a single data value or state.
   - It ensures that all nodes in a distributed system have a consistent view of the data, even in the presence of failures.

2. **Real-World Analogy**:
   - Imagine a group of people trying to agree on a decision. They must communicate with each other, propose ideas, and reach a consensus on a single decision.
   - Similarly, in a distributed system, nodes must communicate and agree on a common state.

3. **Challenges**:
   - Network Failures: Nodes might be temporarily unreachable, making it hard to reach a consensus.
   - Latency: Communication delays can slow down the consensus process.
   - Fault Tolerance: The system must handle node failures gracefully, without compromising on the consistency of the data.