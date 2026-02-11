# DIGIT CLI

A comprehensive command-line interface for interacting with DIGIT platform services. This CLI tool provides commands for account management, user creation, role assignment, workflow management, template creation, MDMS operations, and more.

## Features

- **Account Management**: Create and manage DIGIT service accounts
- **User Management**: Complete Keycloak user lifecycle (create, update, delete, search, password reset)
- **Role Management**: Create roles and assign them to users in Keycloak
- **Template Management**: Create and search notification templates (EMAIL, SMS)
- **Workflow Management**: Create processes, states, actions, and complete workflows
- **ID Generation**: Create and manage ID generation templates
- **Document Categories**: Create and manage filestore document categories
- **MDMS Operations**: Create schemas and manage master data
- **Configuration Management**: Multi-context configuration with authentication
- **Cross-Platform**: Available for Linux, macOS, and Windows

## Installation

### Option 1: Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/digitnxt/digit3/releases):

**Linux (x86_64):**
```bash
# Download and install
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit3/releases/latest/download/digit-cli_Linux_x86_64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
chmod +x /usr/local/bin/digit

# Verify installation
digit --help
```

**macOS (Intel):**
```bash
# Download and install
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit3/releases/latest/download/digit-cli_Darwin_x86_64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
chmod +x /usr/local/bin/digit
```

**macOS (Apple Silicon):**
```bash
# Download and install
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit3/releases/latest/download/digit-cli_Darwin_arm64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
chmod +x /usr/local/bin/digit
```

**Windows:**
```powershell
# Download from releases page and extract digit.exe
# Add the directory containing digit.exe to your PATH
```

### Option 2: Build from Source

**Prerequisites:**
- Go 1.21 or later
- Git

**Steps:**
```bash
# Clone the repository
git clone https://github.com/digitnxt/digit3.git
cd digit3/DIGIT-CLI

# Download dependencies
go mod tidy

# Build the binary
go build -o digit .

# Install globally (optional)
sudo mv digit /usr/local/bin/
```

**For development:**
```bash
# Install with go install
go install .
```

## Configuration

### Set Server URL

Before using the CLI, configure the DIGIT server URL:

```bash
digit config set --server http://localhost:8094
```

This will create a configuration file at `~/.digit/config.yaml`.

## Usage

### Quick Start

1. **Configure the server URL:**
```bash
digit config set --server http://localhost:8094
```

2. **Create an account:**
```bash
digit create-account --name kongnew1 --email test@example.com
```

3. **Create a user in Keycloak:**
```bash
digit create-user --username johndoe --password mypassword --email john@example.com --account master
```

4. **Create and assign a role:**
```bash
# Create role
digit create-role --role-name "ADMIN" --account master

# Assign role to user
digit assign-role --username johndoe --role-name "ADMIN" --account master
```

5. **Create a notification template:**
```bash
digit create-template --template-id "welcome-email" --version "1.0.0" --type "EMAIL" --subject "Welcome!" --content "Welcome to DIGIT platform"
```

6. **Create an ID generation template:**
```bash
digit create-idgen-template --template-code "user-id" --template "USER-{SEQ}-{RAND}"
```

7. **Create a workflow process:**
```bash
digit create-process --name "Application Workflow" --code "APP-WF-001" --description "Application processing workflow" --version "1.0" --sla 86400
```

8. **Create a complete workflow from YAML:**
```bash
digit create-workflow --file example-workflow.yaml
```

9. **Create MDMS schema and data:**
```bash
# Create schema
digit create-schema --file example-schema.yaml

# Create data
digit create-mdms-data --file example-mdms-data.yaml
```

10. **Create registry schema:**
```bash
# Create registry schema from YAML file
digit create-registry-schema --file registry-schema.yaml

# Create registry schema using default configuration
digit create-registry-schema --default

# Create registry schema with custom schema code
digit create-registry-schema --default --schema-code "custom-license-registry"
```

11. **Search registry schema:**
```bash
# Search registry schema by code
digit search-registry-schema --schema-code "license-registry"

# Search registry schema by code and version
digit search-registry-schema --schema-code "license-registry" --version "1"
```

12. **Delete registry schema:**
```bash
# Delete registry schema by code
digit delete-registry-schema --schema-code "license-registry"
```

## Command Reference

### `digit config`

Manage CLI configuration with multiple subcommands.

#### `digit config set`

Set configuration values for the CLI.

**Flags:**
- `--server`: Server URL for DIGIT services (required)
- `--jwt-token`: JWT token for authentication

**Examples:**
```bash
digit config set --server http://localhost:8094
digit config set --server https://digit.example.com --jwt-token <your-token>
```

#### `digit config set`

Authenticate with Keycloak and set configuration using YAML file or command-line flags.

**Flags:**
- `--file`: Path to configuration YAML file
- `--server`: Server URL (e.g., https://digit-lts.digit.org)
- `--account`: Keycloak account name
- `--client-id`: Keycloak client ID
- `--client-secret`: Keycloak client secret
- `--username`: Username for authentication
- `--password`: Password for authentication

**Examples:**
```bash
# Using YAML file
digit config set --file sample-digit-config.yaml

# Using command-line flags
digit config set --server https://digit-lts.digit.org --account CLI --client-id admin-cli --client-secret mysecret --username user@example.com --password mypassword
```

#### `digit config show`

Show current configuration.

**Examples:**
```bash
digit config show
```

#### `digit config get-contexts`

List available contexts from configuration file.

**Flags:**
- `--file`: Path to configuration YAML file (required)

**Examples:**
```bash
digit config get-contexts --file sample-digit-config.yaml
```

#### `digit config use-context`

Switch to different context.

**Flags:**
- `--file`: Path to configuration YAML file (required)
- Context name as argument

**Examples:**
```bash
digit config use-context <context-name> --file sample-digit-config.yaml
```

---

### `digit create-account`

Create a new account in DIGIT services.

**Flags:**
- `--name`: Name of the tenant (required)
- `--email`: Email of the tenant (required)
- `--active`: Whether the tenant is active (default: true)
- `--client-id`: Client ID for the request (default: "test-client")
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Basic account creation
digit create-account --name kongnew1 --email test@example.com

# Create inactive account
digit create-account --name kongnew1 --email test@example.com --active=false

# Use custom client ID
digit create-account --name kongnew1 --email test@example.com --client-id custom-client

# Override server URL for single request
digit create-account --name kongnew1 --email test@example.com --server http://different-server:8094
```

---

### `digit create-user`

Create a new user in Keycloak with the specified username, password, email, and realm.

**Flags:**
- `--username`: Username for the new user (required)
- `--password`: Password for the new user (required)
- `--email`: Email for the new user (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Basic user creation
digit create-user --username johndoe --password mypassword --email john@example.com --account master

# With custom server
digit create-user --username johndoe --password mypassword --email john@example.com --account myrealm --server https://keycloak.example.com

# With JWT token
digit create-user --username johndoe --password mypassword --email john@example.com --account master --jwt-token <token>
```

---

### `digit reset-password`

Reset a user's password in Keycloak with the specified username and new password.

**Flags:**
- `--username`: Username of the user (required)
- `--new-password`: New password for the user (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Basic password reset
digit reset-password --username johndoe --new-password newpassword123 --account master

# With custom server
digit reset-password --username johndoe --new-password newpassword123 --account myrealm --server https://keycloak.example.com

# With JWT token
digit reset-password --username johndoe --new-password newpassword123 --account master --jwt-token <token>
```

---

### `digit delete-user`

Delete a user from Keycloak.

**Flags:**
- `--username`: Username of the user to delete (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Basic user deletion
digit delete-user --username johndoe --account master

# With custom server
digit delete-user --username johndoe --account myrealm --server https://keycloak.example.com

# With JWT token
digit delete-user --username johndoe --account master --jwt-token <token>
```

---

### `digit search-user`

Search for users in Keycloak or list all users.

**Flags:**
- `--account`: Keycloak account name (required)
- `--username`: Username to search for (optional, lists all if not provided)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# List all users in account
digit search-user --account master

# Search for specific user
digit search-user --username johndoe --account master

# With custom server
digit search-user --username johndoe --account myrealm --server https://keycloak.example.com

# With JWT token
digit search-user --account master --jwt-token <token>
```

---

### `digit update-user`

Update a user's information in Keycloak.

**Flags:**
- `--username`: Username of the user to update (required)
- `--account`: Keycloak account name (required)
- `--email`: New email address
- `--first-name`: New first name
- `--last-name`: New last name
- `--enabled`: Enable/disable user (true/false)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Update user email
digit update-user --username johndoe --email newemail@example.com --account master

# Update user names
digit update-user --username johndoe --first-name John --last-name Doe --account master

# Enable/disable user
digit update-user --username johndoe --enabled=false --account master

# Update multiple fields
digit update-user --username johndoe --email new@example.com --first-name John --enabled=true --account master

# With custom server and JWT token
digit update-user --username johndoe --email new@example.com --account master --server https://keycloak.example.com --jwt-token <token>
```

---

### `digit create-role`

Create a new role in Keycloak.

**Flags:**
- `--role-name`: Name of the role to create (required)
- `--account`: Keycloak account name (required)
- `--description`: Description of the role (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Create basic role
digit create-role --role-name "ADMIN" --account master

# Create role with description
digit create-role --role-name "MANAGER" --description "Manager role with elevated permissions" --account master

# With custom server and JWT token
digit create-role --role-name "VIEWER" --account master --server https://keycloak.example.com --jwt-token <token>
```

---

### `digit assign-role`

Assign a role to a user in Keycloak.

**Flags:**
- `--username`: Username to assign role to (required)
- `--role-name`: Name of the role to assign (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Assign role to user
digit assign-role --username johndoe --role-name "ADMIN" --account master

# With custom server and JWT token
digit assign-role --username johndoe --role-name "MANAGER" --account master --server https://keycloak.example.com --jwt-token <token>
```

---

### `digit create-idgen-template`

Create a new ID generation template for generating unique IDs.

**Flags:**
- `--template-code`: Template code for the ID generation template (required)
- `--template`: Template pattern for ID generation (required)
- `--scope`: Scope for sequence generation (optional, default: daily)
- `--start`: Starting number for sequence (optional, default: 1)
- `--padding-length`: Padding length for sequence numbers (optional, default: 4)
- `--padding-char`: Padding character for sequence numbers (optional, default: "0")
- `--random-length`: Length of random string (optional, default: 2)
- `--random-charset`: Character set for random string (optional, default: "A-Z0-9")
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Basic ID generation template
digit create-idgen-template --template-code orgId --template "{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}"

# With custom sequence configuration
digit create-idgen-template --template-code userId --template "USER-{SEQ}-{RAND}" --scope global --start 100 --padding-length 6

# With custom random configuration
digit create-idgen-template --template-code docId --template "DOC-{SEQ}-{RAND}" --random-length 4 --random-charset "ABCDEF0123456789"

# With custom server
digit create-idgen-template --template-code orgId --template "{ORG}-{SEQ}" --server http://localhost:8080
```

---

### `digit search-idgen-template`

Search for an existing ID generation template by template code.

**Flags:**
- `--template-code`: Template code to search for (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Search ID generation template
digit search-idgen-template --template-code orgId

# With custom server
digit search-idgen-template --template-code userId --server http://localhost:8080
```

---

### `digit create-document-category`

Create a new document category in filestore.

**Flags:**
- `--type`: Document type (required)
- `--code`: Document code (required)
- `--allowed-formats`: Comma-separated list of allowed file formats (required)
- `--min-size`: Minimum file size in bytes (optional)
- `--max-size`: Maximum file size in bytes (optional)
- `--description`: Description of the document category (optional)
- `--sensitive`: Mark as sensitive document (true/false, default: false)
- `--active`: Mark as active (true/false, default: true)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Basic document category
digit create-document-category --type Identity --code AADHAR --allowed-formats "pdf,jpg,jpeg,xsm" --min-size 1024 --max-size 1024000 --sensitive true --active true

# Certificate document category
digit create-document-category --type Certificate --code BIRTH_CERT --allowed-formats "pdf,jpg" --min-size 512 --max-size 2048000 --description "Birth certificate documents"

# With custom server
digit create-document-category --type Identity --code PAN --allowed-formats "pdf,png" --server http://localhost:8081

# Non-sensitive document category
digit create-document-category --type General --code RECEIPT --allowed-formats "pdf,jpg,png" --min-size 100 --max-size 500000 --sensitive false --active true
```

---

### `digit create-template`

Create a new notification template with the specified parameters.

**Flags:**
- `--template-id`: Template ID for the notification template (required)
- `--version`: Version of the template (required)
- `--type`: Type of template (EMAIL, SMS, etc.) (required)
- `--subject`: Subject of the template (required)
- `--content`: Content of the template (use either --content or --content-file)
- `--content-file`: Path to file containing template content (use either --content or --content-file)
- `--html`: Whether the content is HTML (default: false)
- `--server`: Server URL (overrides config)

**YAML Structure:**
```yaml
template-id: "email-notification-template"
version: "1.0.0"
type: "EMAIL"
subject: "Welcome to DIGIT Services"
content: |
  <html>
  <body>
    <h1>Welcome to DIGIT Services!</h1>
    <p>Dear {{name}},</p>
    <p>Thank you for registering with DIGIT Services. Your account has been successfully created.</p>
    <p>Account Details:</p>
    <ul>
      <li>Username: {{username}}</li>
      <li>Email: {{email}}</li>
      <li>Registration Date: {{date}}</li>
    </ul>
    <p>If you have any questions, please contact our support team.</p>
    <p>Best regards,<br>DIGIT Team</p>
  </body>
  </html>
html: true
# Optional: Override server URL
# server: "http://localhost:8081"
```

**Examples:**
```bash
# Using direct content
digit create-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test Subject" --content "Test Content"

# Using content from file
digit create-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test Subject" --content-file "./template.html" --html=true

# SMS template
digit create-template --template-id "sms-template" --version "1.0.0" --type "SMS" --subject "SMS Alert" --content "Your OTP is: {{otp}}"

# With server override
digit create-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test Subject" --content "Test Content" --server http://localhost:8081

# Using YAML file
digit create-template --file template-config.yaml
```

---

### `digit search-notification-template`

Search for notification templates by template ID.

**Flags:**
- `--template-id`: Template ID to search for (required)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Search for notification template
digit search-notification-template --template-id "my-template"

# With custom server
digit search-notification-template --template-id "email-template" --server http://localhost:8081
```

---

### `digit create-process`

Create a new workflow process with the specified parameters.

**Flags:**
- `--name`: Name of the workflow process (required)
- `--code`: Code of the workflow process (required)
- `--description`: Description of the workflow process (required)
- `--version`: Version of the workflow process (required)
- `--sla`: SLA in seconds for the workflow process (required)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Basic process creation
digit create-process --name "Hello" --code "{{GenProcessId}}" --description "A test process for API validation" --version "1.0" --sla 86400

# Custom process with different parameters
digit create-process --name "MyWorkflow" --code "WF001" --description "Custom workflow" --version "2.0" --sla 3600

# With server override
digit create-process --name "MyWorkflow" --code "WF001" --description "Custom workflow" --version "2.0" --sla 3600 --server http://localhost:9090
```

---

### `digit search-process-definition`

Search for a workflow process definition by process ID.

**Flags:**
- `--id`: Process ID to search definition for (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Search process definition
digit search-process-definition --id "dd2e8cf5-a53e-44b5-82b9-490ac73c50dd"

# With custom server
digit search-process-definition --id "dd2e8cf5-a53e-44b5-82b9-490ac73c50dd" --server http://localhost:9090
```

---

### `digit create-workflow`

Create a complete workflow (process, states, and actions) from a YAML file definition. This command orchestrates multiple API calls to create an entire workflow system.

**Flags:**
- `--file`: Path to YAML file containing workflow definition (required)
- `--server`: Server URL (overrides config)

**YAML Structure:**
```yaml
workflow:
  process:
    name: "Application Processing Workflow"
    code: "{{GenProcessId}}"
    description: "A complete workflow for application processing"
    version: "1.0"
    sla: 86400
  
  states:
    - code: "APPLIED"
      name: "Application Submitted"
      isInitial: true
      isParallel: false
      isJoin: false
      sla: 86400
    
    - code: "VERIFY"
      name: "Under Verification"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
  
  actions:
    - name: "Submit for Verification"
      currentState: "APPLIED"
      nextState: "VERIFY"
      roles: ["APPLICANT"]
      attributeValidation:
        assigneeCheck: false
```

**Examples:**
```bash
# Create complete workflow from YAML file
digit create-workflow --file example-workflow.yaml

# With server override
digit create-workflow --file my-workflow.yaml --server http://localhost:9090
```

---

### `digit create-registry-schema`

Create registry schema from a YAML file definition or using default configuration.

**Flags:**
- `--file`: Path to YAML file containing registry schema data
- `--default`: Use default registry schema configuration
- `--schema-code`: Schema code for default registry schema (default: 'license-registry')
- `--server`: Server URL (overrides config, default: http://localhost:8085)
- `--jwt-token`: JWT token for authentication (overrides config)

**YAML Structure:**
```yaml
schemaCode: "license-registry"
definition:
  $schema: "https://json-schema.org/draft/2020-12/schema"
  type: "object"
  additionalProperties: false
  properties:
    licenseNumber:
      type: "string"
    holderName:
      type: "string"
    issueDate:
      type: "string"
      format: "date"
    expiryDate:
      type: "string"
      format: "date"
    status:
      type: "string"
      enum: ["ACTIVE", "SUSPENDED", "REVOKED"]
  required: ["licenseNumber", "holderName", "issueDate", "status"]
  x-indexes:
    - name: "idx_license_status"
      fieldPath: "status"
      method: "btree"
    - fieldPath: "holderName"
      method: "gin"
```

**Examples:**
```bash
# Create registry schema from YAML file
digit create-registry-schema --file registry-schema.yaml

# Create registry schema using default configuration with custom schema code
digit create-registry-schema --default --schema-code "custom-license-registry"

# Create registry schema using default configuration (uses license-registry)
digit create-registry-schema --default

# With server override
digit create-registry-schema --file registry-schema.yaml --server http://localhost:8085

# With JWT token
digit create-registry-schema --file registry-schema.yaml --jwt-token <your-jwt-token>
```

---

### `digit search-registry-schema`

Search for a registry schema by schema code and optional version.

**Flags:**
- `--schema-code`: Schema code to search for (required)
- `--version`: Version of the schema (optional)
- `--server`: Server URL (overrides config, default: http://localhost:8085)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Search registry schema by code
digit search-registry-schema --schema-code "license-registry"

# Search registry schema by code and version
digit search-registry-schema --schema-code "license-registry" --version "1"

# With server override
digit search-registry-schema --schema-code "license-registry" --server http://localhost:8085

# With JWT token
digit search-registry-schema --schema-code "license-registry" --jwt-token <your-jwt-token>
```

---

### `digit delete-registry-schema`

Delete a registry schema by schema code.

**Flags:**
- `--schema-code`: Schema code to delete (required)
- `--server`: Server URL (overrides config, default: http://localhost:8085)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Delete registry schema by code
digit delete-registry-schema --schema-code "license-registry"

# With server override
digit delete-registry-schema --schema-code "license-registry" --server http://localhost:8085

# With JWT token
digit delete-registry-schema --schema-code "license-registry" --jwt-token <your-jwt-token>
```

---

### `digit create-boundaries`

Create boundaries from a YAML file definition.

**Flags:**
- `--file`: Path to YAML file containing boundary data (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Create boundaries from YAML file
digit create-boundaries --file boundaries.yaml

# With custom server
digit create-boundaries --file boundaries.yaml --server http://localhost:8080
```

---

### `digit create-schema`

Create a new MDMS schema from a YAML file definition.

**Flags:**
- `--file`: Path to YAML file containing schema definition (required)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Create MDMS schema from YAML file
digit create-schema --file example-schema.yaml

# With custom server
digit create-schema --file my-schema.yaml --server http://localhost:8080
```

---

### `digit search-schema`

Search for an MDMS schema by schema code.

**Flags:**
- `--code`: Schema code to search for (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Search MDMS schema
digit search-schema --code "RAINMAKER-PGR.ServiceDefsns"

# With custom server
digit search-schema --code "EMPLOYEE" --server http://localhost:8080
```

---

### `digit create-mdms-data`

Create MDMS data entries from a YAML file definition.

**Flags:**
- `--file`: Path to YAML file containing MDMS data definition (required)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
# Create MDMS data from YAML file
digit create-mdms-data --file example-mdms-data.yaml

# With custom server
digit create-mdms-data --file my-data.yaml --server http://localhost:8080
```

---

### `digit search-mdms-data`

Search MDMS data by schema code and optional unique identifiers.

**Flags:**
- `--code`: Schema code to search MDMS data for (required)
- `--unique-identifiers`: Comma-separated unique identifiers to filter data (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token for authentication (overrides config)

**Examples:**
```bash
# Search MDMS data by schema code
digit search-mdms-data --code "common-masters.abcd"

# Search MDMS data with unique identifiers
digit search-mdms-data --code "common-masters.abcd" --unique-identifiers "Alice1,Alice3"

# With custom server
digit search-mdms-data --code "EMPLOYEE" --server http://localhost:8080
```

## Project Structure

```
DIGIT-CLI/
├── cmd/                           # Command implementations
│   ├── root.go                   # Root command and CLI setup
│   ├── config.go                 # Configuration management commands
│   ├── configSet.go              # Authentication-based config setting
│   ├── configShow.go             # Show current configuration
│   ├── configGetContexts.go      # List available contexts
│   ├── configUseContext.go       # Switch contexts
│   ├── createAccount.go          # Account creation command
│   ├── createUser.go             # User management commands (CRUD + roles)
│   ├── createTemplate.go         # Template management commands
│   ├── createWorkflow.go         # Workflow management commands
│   ├── createIdGenTemplate.go    # ID generation template commands
│   ├── createDocumentCategory.go # Document category management
│   └── createmdms.go             # MDMS schema and data commands
├── pkg/                          # Shared packages
│   ├── api/                      # API client utilities
│   ├── auth/                     # Authentication handling
│   ├── config/                   # Configuration management
│   └── jwt/                      # JWT token handling
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
├── .goreleaser.yaml              # Release configuration
├── example-*.yaml                # Example configuration files
└── README.md                     # This documentation
```

## API Integration

The CLI integrates with DIGIT services using the following API endpoints:

- **Account Creation**: `POST /account`
  - Headers: `Content-Type: application/json`, `X-Client-Id: <client-id>`
  - Payload: JSON with tenant information

- **User Creation**: Keycloak Admin API
  - Endpoint: `/admin/realms/{realm}/users`
  - Headers: `Authorization: Bearer <jwt-token>`, `Content-Type: application/json`
  - Payload: JSON with user information

- **Template Creation**: `POST /template`
  - Headers: `Content-Type: application/json`
  - Payload: JSON with template information

- **Workflow Creation**: `POST /workflow/v3/process`
  - Headers: `Content-Type: application/json`, `X-Tenant-ID: <tenant-id>`
  - Payload: JSON with workflow process information

## Development

### Adding New Commands

The CLI is designed to be extensible. To add new commands:

1. Create a new file in the `cmd/` directory (e.g., `cmd/newCommand.go`)
2. Define your command using the Cobra library
3. Add the command to the root command in the `init()` function
4. Follow the existing patterns for configuration and API integration

### Dependencies

- [Cobra](https://github.com/spf13/cobra): CLI framework
- [Resty](https://github.com/go-resty/resty): HTTP client
- [YAML v3](https://github.com/go-yaml/yaml): YAML parsing for configuration

## Configuration File

The CLI stores configuration in `~/.digit/config.yaml`:

```yaml
server: http://localhost:8094
```

## Error Handling

The CLI provides clear error messages for common issues:

- Missing required flags
- Network connectivity problems
- Invalid server URLs
- Configuration file issues

## Available Commands Summary

| Command | Description | Key Flags |
|---------|-------------|-----------|
| **Configuration** |
| `config set` | Authenticate and set configuration | `--file` or auth flags |
| `config show` | Show current configuration | None |
| `config get-contexts` | List available contexts | `--file` |
| `config use-context` | Switch to different context | `--file`, context name |
| **Account Management** |
| `create-account` | Create new DIGIT account | `--name`, `--email`, `--active` |
| **User Management** |
| `create-user` | Create new Keycloak user | `--username`, `--password`, `--email`, `--account` |
| `reset-password` | Reset user password in Keycloak | `--username`, `--new-password`, `--account` |
| `delete-user` | Delete user from Keycloak | `--username`, `--account` |
| `search-user` | Search users in Keycloak | `--account`, `--username` (optional) |
| `update-user` | Update user information in Keycloak | `--username`, `--account`, update fields |
| **Role Management** |
| `create-role` | Create new role in Keycloak | `--role-name`, `--account`, `--description` |
| `assign-role` | Assign role to user in Keycloak | `--username`, `--role-name`, `--account` |
| **ID Generation** |
| `create-idgen-template` | Create ID generation template | `--template-code`, `--template` |
| `search-idgen-template` | Search ID generation template | `--template-code` |
| **Document Management** |
| `create-document-category` | Create filestore document category | `--type`, `--code`, `--allowed-formats` |
| **Template Management** |
| `create-template` | Create notification template | `--template-id`, `--version`, `--type`, `--subject`, `--content` |
| `search-notification-template` | Search notification templates | `--template-id` |
| **Workflow Management** |
| `create-process` | Create workflow process | `--name`, `--code`, `--description`, `--version`, `--sla` |
| `search-process-definition` | Search workflow process definition | `--id` |
| `create-workflow` | Create complete workflow from YAML | `--file` |
| **Boundary Management** |
| `create-boundaries` | Create boundaries from YAML | `--file` |
| **Registry Management** |
| `create-registry-schema` | Create registry schema from YAML | `--file` or `--default`, `--schema-code` |
| `search-registry-schema` | Search registry schema by code | `--schema-code`, `--version` |
| `delete-registry-schema` | Delete registry schema by code | `--schema-code` |
| **MDMS Operations** |
| `create-schema` | Create MDMS schema from YAML | `--file` |
| `search-schema` | Search MDMS schema by code | `--code` |
| `create-mdms-data` | Create MDMS data from YAML | `--file` |
| `search-mdms-data` | Search MDMS data by schema code | `--code`, `--unique-identifiers` |
| **Utility** |
| `completion` | Generate shell autocompletion | Shell type |
| `help` | Help about any command | Command name |

## Future Enhancements

This CLI is designed to be extended with additional commands for:

- **Bulk Operations**: Batch processing for multiple users, roles, or templates
- **Interactive Mode**: Guided CLI experience with prompts and validation
- **Template Management**: Update and delete notification templates
- **Workflow Management**: Update, delete, and list workflow processes
- **MDMS Management**: Update and delete MDMS schemas and data
- **Account Management**: Account verification, status checks, and updates
- **Reporting**: Generate reports on users, roles, workflows, and system usage
- **Import/Export**: Bulk import/export of configurations and data
- **Validation**: Pre-flight validation of YAML files and configurations
- **Monitoring**: Health checks and status monitoring of DIGIT services
