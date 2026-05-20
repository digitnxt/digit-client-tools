# DIGIT CLI

A comprehensive command-line interface for interacting with DIGIT platform services. This CLI tool provides commands for account management, user creation, role assignment, workflow management, template creation, MDMS operations, access control, boundary management, and more.

## Features

- **Account Management**: Create, search, and delete tenant accounts via admin API
- **User Management**: Complete Keycloak user lifecycle (create, update, delete, search, password reset)
- **Role Management**: Create roles and assign them to users in Keycloak
- **Access Control**: Full RBAC/JBAC rule management (create, list, get, delete)
- **Template Management**: Create, search, and delete notification templates (EMAIL, SMS)
- **Workflow Management**: Create complete workflows from YAML, search, and delete process definitions
- **ID Generation**: Create, search, and delete ID generation templates
- **Document Categories**: Create and delete filestore document categories
- **MDMS Operations**: Create schemas, manage master data
- **Boundary Management**: Create and search boundary hierarchies and relationships
- **Registry Management**: Create, search, and delete registry schemas
- **Configuration Management**: Multi-context configuration with authentication
- **Cross-Platform**: Available for Linux, macOS, and Windows

## Installation

Download the latest release for your platform from the [releases page](https://github.com/digitnxt/digit-client-tools/releases):

**Linux (x86_64):**
```bash
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit-client-tools/releases/latest/download/digit-cli_Linux_x86_64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
digit --help
```

**Linux (arm64):**
```bash
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit-client-tools/releases/latest/download/digit-cli_Linux_arm64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
digit --help
```

**macOS (Intel):**
```bash
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit-client-tools/releases/latest/download/digit-cli_Darwin_x86_64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
digit --help
```

**macOS (Apple Silicon):**
```bash
curl -L -o digit-cli.tar.gz "https://github.com/digitnxt/digit-client-tools/releases/latest/download/digit-cli_Darwin_arm64.tar.gz"
tar -xzf digit-cli.tar.gz
sudo mv digit /usr/local/bin/
digit --help
```

**Windows (x86_64):**
```powershell
Invoke-WebRequest -Uri "https://github.com/digitnxt/digit-client-tools/releases/latest/download/digit-cli_Windows_x86_64.zip" -OutFile digit-cli.zip
Expand-Archive digit-cli.zip -DestinationPath digit-cli
# Add the digit-cli folder to your PATH, then verify:
digit --help
```

## Usage

### Quick Start

1. **Create an account:**
```bash
digit create-account --name demoaccount --email test@example.com --server https://digit-lts.digit.org
```

2. **Configure authentication:**
```bash
digit config set --server https://digit-lts.digit.org --account DEMOACCOUNT6 --client-id auth-server --client-secret changeme --username st@example.com --password default
```

3. **Create a notification template:**
```bash
digit create-notification-template --template-id "welcome-email" --version "1.0.0" --type "EMAIL" --subject "Welcome!" --content "Welcome to DIGIT platform"
```

4. **Create an ID generation template:**
```bash
digit create-idgen-template --template-code "user-id" --template "USER-{SEQ}-{RAND}"
```

5. **Create a workflow process:**
```bash
digit create-process --name "Application Workflow" --code "APP-WF-001" --description "Application processing workflow" --version "1.0" --sla 86400
```

6. **Create a complete workflow from YAML:**
```bash
digit create-workflow --file examples/example-workflow.yaml
```

7. **Create MDMS schema and data:**
```bash
# Create MDMS schema from YAML file
digit create-mdms-schema --file examples/example-schema.yaml

# Create MDMS schema using default configuration
digit create-mdms-schema --default

# Create MDMS schema with custom schema code
digit create-mdms-schema --default --code "MY_CUSTOM_SCHEMA"

# Create MDMS data
digit create-mdms-data --file examples/example-mdms-data.yaml
```

8. **Create registry schema:**
```bash
# Create registry schema from YAML file
digit create-registry-schema --file examples/registry-schema.yaml

# Create registry schema using default configuration
digit create-registry-schema --default

# Create registry schema with custom schema code
digit create-registry-schema --default --schema-code "custom-license-registry"
```

9. **Search registry schema:**
```bash
# Search registry schema by code
digit search-registry-schema --schema-code "license-registry"

# Search registry schema by code and version
digit search-registry-schema --schema-code "license-registry" --version "1"
```

10. **Create a user in Keycloak:**
```bash
digit create-user --username johndoe --password mypassword --email john@example.com --account master
```

11. **Create and assign a role:**
```bash
# Create role
digit create-role --role-name "ADMIN" --account master

# Assign role to user
digit assign-role --username johndoe --role-name "ADMIN" --account master
```

## Command Reference

### `digit config`

Manage CLI configuration with multiple subcommands.

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
digit config set --file examples/sample-digit-config.yaml
digit config set --server https://digit-lts.digit.org --account CLI --client-id admin-cli --client-secret mysecret --username user@example.com --password mypassword
```

#### `digit config show`

Show current configuration.

```bash
digit config show
```

#### `digit config get-contexts`

List available contexts from a configuration file.

```bash
digit config get-contexts --file examples/sample-digit-config.yaml
```

#### `digit config use-context`

Switch to a different context.

```bash
digit config use-context <context-name> --file examples/sample-digit-config.yaml
```

---

### `digit create-account`

Create a new tenant account via the admin API.

**Flags:**
- `--name`: Tenant name (required)
- `--email`: Tenant email (required)
- `--password`: Password (optional — server generates one if omitted)
- `--phone`: Phone number (optional)
- `--address`: Address (optional)
- `--city`: City (optional)
- `--state`: State (optional)
- `--pincode`: Pincode (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --server https://digit-lts.digit.org
digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --password Changeme1 --server https://digit-lts.digit.org
digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --phone "+254202229000" --city Nairobi --server https://digit-lts.digit.org
```

---

### `digit search-account`

Search or list tenant accounts with optional filters.

**Flags:**
- `--name`: Filter by tenant name (partial match)
- `--email`: Filter by tenant email (partial match)
- `--page`: Page number, 1-indexed (default: 1)
- `--size`: Results per page (default: 20)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-account --server https://digit-lts.digit.org
digit search-account --name "Nairobi" --server https://digit-lts.digit.org
digit search-account --email admin@nairobi.go.ke --server https://digit-lts.digit.org
digit search-account --page 2 --size 10 --server https://digit-lts.digit.org
```

---

### `digit delete-account`

Permanently delete a tenant account by ID.

**Flags:**
- `--id`: Tenant account ID to delete (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-account --id 3fa85f64-5717-4562-b3fc-2c963f66afa6 --server https://digit-lts.digit.org
```

---

### `digit create-user`

Create a new user in Keycloak.

**Flags:**
- `--username`: Username (required)
- `--password`: Password (required)
- `--email`: Email (required)
- `--account`: Keycloak account/realm name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-user --username johndoe --password mypassword --email john@example.com --account master
digit create-user --username johndoe --password mypassword --email john@example.com --account myrealm --server https://keycloak.example.com
```

---

### `digit search-user`

Search for users in Keycloak.

**Flags:**
- `--account`: Keycloak account name (required)
- `--username`: Username to search (optional — lists all if omitted)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-user --account master --server https://digit-lts.digit.org
digit search-user --username johndoe --account master --server https://digit-lts.digit.org
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
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit update-user --username johndoe --email newemail@example.com --account master --server https://digit-lts.digit.org
digit update-user --username johndoe --first-name John --last-name Doe --enabled=false --account master --server https://digit-lts.digit.org
```

---

### `digit reset-password`

Reset a user's password in Keycloak.

**Flags:**
- `--username`: Username (required)
- `--new-password`: New password (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit reset-password --username johndoe --new-password newpassword123 --account master --server https://digit-lts.digit.org
```

---

### `digit delete-user`

Delete a user from Keycloak.

**Flags:**
- `--username`: Username to delete (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-user --username johndoe --account master --server https://digit-lts.digit.org
```

---

### `digit create-role`

Create a new role in Keycloak.

**Flags:**
- `--role-name`: Role name (required)
- `--account`: Keycloak account name (required)
- `--description`: Description (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-role --role-name "ADMIN" --account master --server https://digit-lts.digit.org
digit create-role --role-name "MANAGER" --description "Manager role" --account master --server https://digit-lts.digit.org
```

---

### `digit assign-role`

Assign a role to a user in Keycloak.

**Flags:**
- `--username`: Username (required)
- `--role-name`: Role name (required)
- `--account`: Keycloak account name (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit assign-role --username johndoe --role-name "ADMIN" --account master --server https://digit-lts.digit.org
```

---

### `digit create-rbac-rule`

Create RBAC/JBAC access control rules via flags or YAML file.

**Flags (flag mode):**
- `--roles`: Comma-separated role names (required)
- `--method`: HTTP method (GET, POST, PUT, DELETE, PATCH)
- `--path`: API path to protect (required)
- `--effect`: ALLOW or DENY (default: ALLOW)
- `--priority`: Rule priority (default: 100)
- `--description`: Rule description
- `--crud`: Create rules for all CRUD methods at once
- `--constraint`: Constraint in `type:key:op:value` format (repeatable)
- `--file`: YAML file with multiple rules
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**YAML file structure:**
```yaml
rules:
  - roleNames: ["SUPERUSER"]
    httpMethod: "GET"
    path: "/workflow/v3/process/definition"
    description: "Allow SUPERUSER to read workflow definitions"
  - roleNames: ["SUPERUSER", "ADMIN"]
    path: "/accounts/v3/tenants"
    crud: true
```

**Examples:**
```bash
# Single rule via flags
digit create-rbac-rule --roles SUPERUSER --method POST --path /workflow/v3/process/definition --server https://digit-lts.digit.org

# All CRUD methods for a path
digit create-rbac-rule --roles SUPERUSER --path /accounts/v3/tenants --crud --server https://digit-lts.digit.org

# With a JBAC constraint
digit create-rbac-rule --roles CITIZEN --method GET --path /filestore/v3/files --constraint "boundary:tenantId:EQ:KA.BLR" --server https://digit-lts.digit.org

# From YAML file
digit create-rbac-rule --file rbac-rules.yaml --server https://digit-lts.digit.org
```

---

### `digit list-rbac-rules`

List RBAC access control rules for your tenant.

**Flags:**
- `--role`: Filter rules by role name (optional)
- `--page`: Page number (default: 0)
- `--size`: Page size (default: 50)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit list-rbac-rules --server https://digit-lts.digit.org
digit list-rbac-rules --role SUPERUSER --server https://digit-lts.digit.org
digit list-rbac-rules --page 0 --size 20 --server https://digit-lts.digit.org
```

---

### `digit get-rbac-rule`

Get the full details of a single RBAC rule by ID.

**Flags:**
- `--id`: Rule UUID (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit get-rbac-rule --id 550e8400-e29b-41d4-a716-446655440000 --server https://digit-lts.digit.org
```

---

### `digit delete-rbac-rule`

Delete a single RBAC rule by ID.

**Flags:**
- `--id`: Rule UUID (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-rbac-rule --id 550e8400-e29b-41d4-a716-446655440000 --server https://digit-lts.digit.org
```

---

### `digit delete-all-rbac-rules`

Delete all RBAC rules for your tenant (destructive — use with caution).

**Flags:**
- `--role`: Only delete rules for this role (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-all-rbac-rules --server https://digit-lts.digit.org
digit delete-all-rbac-rules --role SUPERUSER --server https://digit-lts.digit.org
```

---

### `digit create-idgen-template`

Create a new ID generation template.

**Flags:**
- `--template-code`: Template code (required)
- `--template`: Template pattern e.g. `{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}` (required unless `--default`)
- `--default`: Use built-in default configuration (requires `--template-code`)
- `--scope`: Sequence scope — daily, monthly, yearly, global (default: daily)
- `--start`: Starting number for sequence (default: 1)
- `--padding-length`: Padding length for sequence (default: 4)
- `--padding-char`: Padding character (default: "0")
- `--random-length`: Length of random string (default: 2)
- `--random-charset`: Character set for random string (default: "A-Z0-9")
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-idgen-template --template-code orgId --template "{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}" --server https://digit-lts.digit.org
digit create-idgen-template --default --template-code "my-custom-template" --server https://digit-lts.digit.org
digit create-idgen-template --template-code userId --template "USER-{SEQ}-{RAND}" --scope global --start 100 --server https://digit-lts.digit.org
```

---

### `digit search-idgen-template`

Search for an existing ID generation template by code.

**Flags:**
- `--template-code`: Template code (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-idgen-template --template-code orgId --server https://digit-lts.digit.org
```

---

### `digit delete-idgen-template`

Delete an existing ID generation template.

**Flags:**
- `--template-code`: Template code (required)
- `--version`: Template version to delete (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-idgen-template --template-code orgId --version v2 --server https://digit-lts.digit.org
```

---

### `digit create-filestore-document-category`

Create a new document category in filestore.

**Flags:**
- `--type`: Document type e.g. Identity, Certificate (required)
- `--code`: Document code e.g. AADHAR, BIRTH_CERT (required)
- `--allowed-formats`: Comma-separated allowed file formats e.g. "pdf,jpg,jpeg" (required)
- `--min-size`: Minimum file size e.g. 1KB, 512B (default: 1KB)
- `--max-size`: Maximum file size e.g. 1MB, 10MB (default: 1MB)
- `--sensitive`: Mark as sensitive (default: false)
- `--active`: Mark as active (default: true)
- `--description`: Description (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-filestore-document-category --type Identity --code AADHAR --allowed-formats "pdf,jpg,jpeg,xsm" --min-size 1KB --max-size 1MB --sensitive --active --server https://digit-lts.digit.org
digit create-filestore-document-category --type Certificate --code BIRTH_CERT --allowed-formats "pdf,jpg" --min-size 512KB --max-size 2MB --description "Birth certificate documents" --server https://digit-lts.digit.org
```

---

### `digit delete-filestore-document-category`

Delete a document category by code.

**Flags:**
- `--code`: Category code to delete (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-filestore-document-category --code AADHAR --server https://digit-lts.digit.org
```

---

### `digit create-notification-template`

Create a new notification template.

**Flags:**
- `--template-id`: Template ID (required unless `--file`)
- `--version`: Version e.g. "1.0.0" (required unless `--file` or `--default`)
- `--type`: Template type — EMAIL, SMS (required unless `--file` or `--default`)
- `--subject`: Subject (required unless `--file` or `--default`)
- `--content`: Template content
- `--content-file`: Path to file with template content
- `--html`: Content is HTML (default: false)
- `--default`: Use built-in default EMAIL template (requires `--template-id`)
- `--file`: Path to YAML configuration file
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**YAML file structure:**
```yaml
template-id: "welcome-email"
version: "1.0.0"
type: "EMAIL"
subject: "Welcome to DIGIT Services"
content: |
  <html><body><h1>Welcome!</h1></body></html>
html: true
```

**Examples:**
```bash
digit create-notification-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test" --content "Hello" --server https://digit-lts.digit.org
digit create-notification-template --default --template-id "welcome-email" --server https://digit-lts.digit.org
digit create-notification-template --template-id "sms-otp" --version "1.0.0" --type "SMS" --subject "OTP" --content "Your OTP: {{otp}}" --server https://digit-lts.digit.org
digit create-notification-template --file examples/template-config.yaml --server https://digit-lts.digit.org
```

---

### `digit search-notification-template`

Search for notification templates by template ID.

**Flags:**
- `--template-id`: Template ID to search (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-notification-template --template-id "my-template" --server https://digit-lts.digit.org
```

---

### `digit delete-notification-template`

Delete a notification template by ID and version.

**Flags:**
- `--template-id`: Template ID (required)
- `--version`: Template version (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-notification-template --template-id user-notify --version v1 --server https://digit-lts.digit.org
```

---

### `digit create-workflow`

Create a complete workflow (process + states + actions) from a YAML file or using the built-in default. Uses the single composite `POST /workflow/v3/process/definition` API.

**Flags:**
- `--file`: Path to YAML workflow definition file
- `--default`: Use built-in default workflow (requires `--code`)
- `--code`: Process code for default workflow
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**YAML structure:**
```yaml
workflow:
  process:
    name: "Application Processing Workflow"
    code: "MY_APP_WF"
    description: "Workflow for application processing"
    version: "1.0"
    sla: 86400
  states:
    - code: "INIT"
      name: "Init"
      type: "INITIAL"
      sla: 86400
      actions:
        - code: "APPLY"
          label: "Apply"
          nextState: "PENDINGFORASSIGNMENT"
    - code: "PENDINGFORASSIGNMENT"
      name: "Pending For Assignment"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "ASSIGN"
          label: "Assign"
          nextState: "RESOLVED"
        - code: "REJECT"
          label: "Reject"
          nextState: "REJECTED"
    - code: "RESOLVED"
      name: "Resolved"
      type: "TERMINAL_SUCCESS"
      sla: 0
      actions: []
    - code: "REJECTED"
      name: "Rejected"
      type: "TERMINAL_FAILURE"
      sla: 0
      actions: []
```

State `type` values: `INITIAL`, `INTERMEDIATE`, `TERMINAL_SUCCESS`, `TERMINAL_FAILURE`

**Examples:**
```bash
digit create-workflow --file workflow.yaml --server https://digit-lts.digit.org
digit create-workflow --default --code MY_CUSTOM_WORKFLOW --server https://digit-lts.digit.org
```

---

### `digit search-workflow`

Get a workflow process definition by process code.

**Flags:**
- `--code`: Process code (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-workflow --code MY_WORKFLOW_CODE --server https://digit-lts.digit.org
```

---

### `digit delete-workflow`

Delete a workflow process definition by process code.

**Flags:**
- `--code`: Process code to delete (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-workflow --code TRADE_LICENSE_WF --server https://digit-lts.digit.org
```

---

### `digit create-boundaries`

Create boundaries from a YAML file or default configuration.

**Flags:**
- `--file`: Path to YAML file with boundary data
- `--default`: Use default boundary configuration
- `--code-prefix`: Code prefix for default boundaries (default: DEFAULT)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-boundaries --file boundaries.yaml --server https://digit-lts.digit.org
digit create-boundaries --default --code-prefix "KARNATAKA" --server https://digit-lts.digit.org
```

---

### `digit create-boundary-hierarchy`

Create a boundary hierarchy from a YAML file or default configuration.

**Flags:**
- `--file`: Path to YAML file with boundary hierarchy data
- `--default`: Use default hierarchy configuration
- `--hierarchy-type`: Hierarchy type for default config (default: state-district-hierarchy)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-boundary-hierarchy --file boundary-hierarchy.yaml --server https://digit-lts.digit.org
digit create-boundary-hierarchy --default --server https://digit-lts.digit.org
digit create-boundary-hierarchy --default --hierarchy-type "custom-hierarchy" --server https://digit-lts.digit.org
```

---

### `digit search-boundary-hierarchy`

Search for a boundary hierarchy by hierarchy type.

**Flags:**
- `--hierarchy-type`: Hierarchy type to search (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-boundary-hierarchy --hierarchy-type "state-district-hierarchy" --server https://digit-lts.digit.org
```

---

### `digit create-boundary-relationships`

Create a boundary relationship using command flags.

**Flags:**
- `--code`: Boundary code (required)
- `--hierarchy-type`: Hierarchy type (required)
- `--boundary-type`: Boundary type (required)
- `--parent`: Parent boundary code (optional — null if omitted)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-boundary-relationships --code "STATE1" --hierarchy-type "state-district-hierarchy" --boundary-type "state" --server https://digit-lts.digit.org
digit create-boundary-relationships --code "DISTRICT1" --hierarchy-type "state-district-hierarchy" --boundary-type "district" --parent "STATE1" --server https://digit-lts.digit.org
```

---

### `digit search-boundary-relationships`

Search for boundary relationships.

**Flags:**
- `--hierarchy-type`: Hierarchy type (required)
- `--boundary-type`: Boundary type (required)
- `--codes`: Comma-separated boundary codes (optional)
- `--include-children`: Include children in response (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state" --server https://digit-lts.digit.org
digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state" --codes "STATE1" --include-children --server https://digit-lts.digit.org
```

---

### `digit create-registry-schema`

Create a registry schema from a YAML file or default configuration.

**Flags:**
- `--file`: Path to YAML file with registry schema
- `--default`: Use default registry schema configuration
- `--schema-code`: Schema code for default config (default: license-registry)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-registry-schema --file examples/registry-schema.yaml --server https://digit-lts.digit.org
digit create-registry-schema --default --schema-code "custom-license-registry" --server https://digit-lts.digit.org
```

---

### `digit search-registry-schema`

Search for a registry schema by code and optional version.

**Flags:**
- `--schema-code`: Schema code (required)
- `--version`: Schema version (optional)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-registry-schema --schema-code "license-registry" --server https://digit-lts.digit.org
digit search-registry-schema --schema-code "license-registry" --version "1" --server https://digit-lts.digit.org
```

---

### `digit delete-registry-schema`

Delete a registry schema by code.

**Flags:**
- `--schema-code`: Schema code (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit delete-registry-schema --schema-code "license-registry" --server https://digit-lts.digit.org
```

---

### `digit create-mdms-schema`

Create a new MDMS schema from a YAML file or default configuration.

**Flags:**
- `--file`: Path to YAML file with schema definition
- `--default`: Use default MDMS schema configuration
- `--code`: Schema code for default config (default: RAINMAKER_PGR_ServiceDefs)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit create-mdms-schema --file examples/example-schema.yaml --server https://digit-lts.digit.org
digit create-mdms-schema --default --server https://digit-lts.digit.org
digit create-mdms-schema --default --code "MY_CUSTOM_SCHEMA" --server https://digit-lts.digit.org
```

---

### `digit search-mdms-schema`

Search for an MDMS schema by code.

**Flags:**
- `--code`: Schema code (required)
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-mdms-schema --code "RAINMAKER-PGR.ServiceDefs" --server https://digit-lts.digit.org
```

---

### `digit create-mdms-data`

Create MDMS data entries from a YAML file.

**Flags:**
- `--file`: Path to YAML file with MDMS data (required)
- `--server`: Server URL (overrides config)

**Examples:**
```bash
digit create-mdms-data --file examples/example-mdms-data.yaml --server https://digit-lts.digit.org
```

---

### `digit search-mdms-data`

Search MDMS data by schema code with optional unique identifier filters.

**Flags:**
- `--code`: Schema code (required)
- `--unique-identifiers`: Comma-separated unique identifiers to filter
- `--server`: Server URL (overrides config)
- `--jwt-token`: JWT token (overrides config)

**Examples:**
```bash
digit search-mdms-data --code "common-masters.abcd" --server https://digit-lts.digit.org
digit search-mdms-data --code "common-masters.abcd" --unique-identifiers "Alice1,Alice3" --server https://digit-lts.digit.org
```

---

## API Integration

The CLI integrates with DIGIT v3 services:

| Domain | Endpoint |
|--------|----------|
| Account | `POST/GET/DELETE /accounts/v3/tenants` |
| Workflow | `POST/GET/DELETE /workflow/v3/process/definition` |
| Access Control | `POST/GET/DELETE /access/v3/rbac/rules/` |
| Filestore | `POST/DELETE /filestore/v3/files/document-categories` |
| IDGen | `POST/GET/DELETE /idgen/v3/id/format` |
| Notification | `POST/GET/DELETE /notification/v3/templates` |
| Boundaries | `POST/GET /boundary/v3/boundary`, `/boundary/v3/boundary/hierarchy` |
| MDMS | `POST/GET /mdms-v2/schema/v1`, `/mdms-v2/v2/create` |
| Registry | `POST/GET/DELETE /registry/v3/schema` |
| Keycloak Users | `/admin/realms/{realm}/users` |

All protected endpoints use `Authorization: Bearer <token>`, `X-Tenant-ID`, and `X-User-ID` headers derived from the configured JWT token.

## Project Structure

```
DIGIT-CLI/
├── cmd/
│   ├── root.go                    # Root command
│   ├── config.go                  # config subcommands
│   ├── configSet.go
│   ├── configShow.go
│   ├── configGetContexts.go
│   ├── configUseContext.go
│   ├── createAccount.go           # create/search/delete-account
│   ├── createUser.go              # create/search/update/delete-user, reset-password, create/assign-role
│   ├── createRbacRule.go          # create-rbac-rule
│   ├── manageRbacRules.go         # list/get/delete/delete-all-rbac-rules
│   ├── createWorkflow.go          # create/search-workflow, delete-process
│   ├── createIdGenTemplate.go     # create/search/delete-idgen-template
│   ├── createDocumentCategory.go  # create/delete-filestore-document-category
│   ├── createNotificationTemplate.go # create/search/delete-notification-template
│   ├── createBoundaries.go        # create-boundaries, create/search-boundary-hierarchy, create/search-boundary-relationships
│   ├── createRegistry.go          # create/search/delete-registry-schema
│   ├── createmdms.go              # create/search-mdms-schema, create/search-mdms-data
│   └── build.go                   # build command
├── pkg/
│   ├── config/                    # Configuration management
│   └── jwt/                       # JWT token handling
├── main.go
├── go.mod
├── .goreleaser.yaml
└── README.md
```

## Configuration File

The CLI stores configuration in `~/.digit/config.yaml`. Set it via:

```bash
digit config set --server https://digit-lts.digit.org --account MYACCOUNT --client-id admin-cli --client-secret mysecret --username user@example.com --password mypassword
```

## Available Commands Summary

| Command | Description | Key Flags |
|---------|-------------|-----------|
| **Configuration** | | |
| `config set` | Authenticate and save configuration | `--server`, auth flags or `--file` |
| `config show` | Show current configuration | — |
| `config get-contexts` | List contexts in a config file | `--file` |
| `config use-context` | Switch to a context | `--file`, context name |
| **Account Management** | | |
| `create-account` | Create new tenant account (admin API) | `--name`, `--email` |
| `search-account` | Search/list tenant accounts | `--name`, `--email`, `--page`, `--size` |
| `delete-account` | Delete tenant account by ID | `--id` |
| **User Management** | | |
| `create-user` | Create Keycloak user | `--username`, `--password`, `--email`, `--account` |
| `search-user` | Search Keycloak users | `--account`, `--username` |
| `update-user` | Update Keycloak user | `--username`, `--account`, update fields |
| `reset-password` | Reset user password | `--username`, `--new-password`, `--account` |
| `delete-user` | Delete Keycloak user | `--username`, `--account` |
| **Role Management** | | |
| `create-role` | Create Keycloak role | `--role-name`, `--account` |
| `assign-role` | Assign role to user | `--username`, `--role-name`, `--account` |
| **Access Control** | | |
| `create-rbac-rule` | Create RBAC/JBAC rule(s) | `--roles`, `--path`, `--method` or `--file` |
| `list-rbac-rules` | List RBAC rules for tenant | `--role`, `--page`, `--size` |
| `get-rbac-rule` | Get RBAC rule by ID | `--id` |
| `delete-rbac-rule` | Delete RBAC rule by ID | `--id` |
| `delete-all-rbac-rules` | Delete all RBAC rules (destructive) | `--role` |
| **ID Generation** | | |
| `create-idgen-template` | Create ID generation template | `--template-code`, `--template` |
| `search-idgen-template` | Search ID generation template | `--template-code` |
| `delete-idgen-template` | Delete ID generation template | `--template-code`, `--version` |
| **Document Management** | | |
| `create-filestore-document-category` | Create filestore document category | `--type`, `--code`, `--allowed-formats` |
| `delete-filestore-document-category` | Delete document category | `--code` |
| **Notification Templates** | | |
| `create-notification-template` | Create notification template | `--template-id`, `--version`, `--type`, `--subject`, `--content` |
| `search-notification-template` | Search notification templates | `--template-id` |
| `delete-notification-template` | Delete notification template | `--template-id`, `--version` |
| **Workflow Management** | | |
| `create-workflow` | Create complete workflow from YAML | `--file` or `--default --code` |
| `search-workflow` | Get workflow process definition | `--code` |
| `delete-workflow` | Delete workflow process definition | `--code` |
| **Boundary Management** | | |
| `create-boundaries` | Create boundaries from YAML | `--file` or `--default` |
| `create-boundary-hierarchy` | Create boundary hierarchy | `--file` or `--default` |
| `search-boundary-hierarchy` | Search boundary hierarchy | `--hierarchy-type` |
| `create-boundary-relationships` | Create boundary relationship | `--code`, `--hierarchy-type`, `--boundary-type` |
| `search-boundary-relationships` | Search boundary relationships | `--hierarchy-type`, `--boundary-type` |
| **Registry Management** | | |
| `create-registry-schema` | Create registry schema | `--file` or `--default` |
| `search-registry-schema` | Search registry schema | `--schema-code` |
| `delete-registry-schema` | Delete registry schema | `--schema-code` |
| **MDMS Operations** | | |
| `create-mdms-schema` | Create MDMS schema | `--file` or `--default` |
| `search-mdms-schema` | Search MDMS schema | `--code` |
| `create-mdms-data` | Create MDMS data | `--file` |
| `search-mdms-data` | Search MDMS data | `--code` |
| **Utility** | | |
| `completion` | Generate shell autocompletion | shell type |
| `help` | Help about any command | command name |
