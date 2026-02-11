# DIGIT Client Tools

A developer toolkit for the [DIGIT platform](https://github.com/digitnxt) — includes a CLI and client libraries for Java and Go to interact with DIGIT services.

## Components

| Component | Description | Language | Docs |
|-----------|-------------|----------|------|
| [**digit-cli**](./digit-cli) | Command-line interface for DIGIT platform operations | Go | [README](./digit-cli/README.md) |
| [**digit-java-client**](./client-libraries/digit-java-client) | Spring-based client library for DIGIT services | Java 17 | [README](./client-libraries/digit-java-client/README.md) |
| [**digit-go-client**](./client-libraries/digit-go-client) | Go client library for DIGIT services | Go | — |

## DIGIT CLI

A cross-platform CLI for managing DIGIT resources — accounts, users, workflows, templates, schemas, and more.

### Quick Start

```bash
# macOS (Apple Silicon)
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit3/releases/latest/download/digit-cli_Darwin_arm64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/

# Or build from source
cd digit-cli
go build -o digit .
```

### Configure

```bash
# Set up connection to a DIGIT environment
digit config set \
  --server https://digit-lts.digit.org \
  --realm MY_ACCOUNT \
  --client-id auth-server \
  --client-secret changeme \
  --username admin@example.com \
  --password admin
```

### Commands

| Category | Commands |
|----------|----------|
| **Account** | `create-account` |
| **Users** | `create-user`, `search-user`, `update-user`, `delete-user`, `reset-password` |
| **Roles** | `create-role`, `assign-role` |
| **Workflows** | `create-workflow`, `create-process`, `search-process-definition` |
| **Templates** | `create-template`, `search-notification-template` |
| **ID Generation** | `create-idgen-template`, `search-idgen-template` |
| **Documents** | `create-document-category` |
| **MDMS** | `create-schema`, `search-schema`, `create-mdms-data`, `search-mdms-data` |
| **Registry** | `create-registry-schema`, `search-registry-schema`, `delete-registry-schema` |
| **Boundaries** | `create-boundaries` |
| **Config** | `config set`, `config show`, `config get-contexts`, `config use-context` |

See the full [CLI documentation](./digit-cli/README.md) for detailed usage and examples.

## Client Libraries

### Java Client

A Spring Framework 6 client library with 8 service clients, OpenTelemetry tracing, automatic header propagation, and built-in retry/timeout handling.

```xml
<dependency>
    <groupId>com.digit</groupId>
    <artifactId>digit-client</artifactId>
    <version>1.0.0</version>
</dependency>
```

**Services:** Account, Boundary, Workflow, Individual, Filestore, IdGen, MDMS, Notification

See the [Java client documentation](./client-libraries/digit-java-client/README.md) for setup and usage.

### Go Client

A lightweight Go client library providing service clients for DIGIT APIs.

```go
import "github.com/digitnxt/digit3/code/libraries/digit-library/digit"
```

**Services:** Account, Auth, Boundary, Filestore, IdGen, MDMS, Registry, Template, User, Workflow

## Project Structure

```
digit-client-tools/
├── digit-cli/                        # CLI tool (Go)
│   ├── cmd/                          # Command implementations
│   ├── pkg/                          # Reusable packages (api, auth, config, jwt)
│   ├── example-*.yaml                # Example input files
│   └── template-*.yaml               # Notification template examples
└── client-libraries/
    ├── digit-java-client/            # Java Spring client library
    │   └── src/main/java/com/digit/  # Services, config, models
    └── digit-go-client/              # Go client library
        └── digit/                    # Service client packages
```

## Prerequisites

| Component | Requirements |
|-----------|-------------|
| DIGIT CLI | Go 1.21+ |
| Java Client | Java 17+, Maven 3.11+ |
| Go Client | Go 1.21+ |

## Platform Support

The CLI is available for:
- Linux (x86_64, arm64)
- macOS (x86_64, arm64)
- Windows (x86_64, arm64)
