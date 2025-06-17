# DevOps CLI Tool

A comprehensive command-line interface for automating common DevOps tasks including Terraform operations, script execution, and CI/CD pipeline management.

## Features

- **Terraform Operations**: Initialize, plan, apply, and destroy Terraform configurations
- **Script Automation**: Execute bash scripts with built-in discovery
- **CI/CD Pipeline Generation**: Create GitHub Actions workflows from templates
- **Git Integration**: Automatically push pipelines to repositories
- **Configuration Management**: Customizable settings for different environments

## Installation

### Prerequisites

- Go 1.19 or later
- Git
- Terraform (if using Terraform features)
- Bash (for script execution)

### Build and Install

```bash
# Clone the repository
git clone <your-repo-url>
cd devops-cli

# Build the binary
make build

# Install system-wide (requires sudo)
make install

# Or install to user directory
make dev-install
```

### Manual Installation

```bash
go mod tidy
go build -o devops-cli .
sudo cp devops-cli /usr/local/bin/
```

## Usage

### Terraform Operations

```bash
# Initialize Terraform in current directory
devops-cli terraform init

# Initialize Terraform in specific directory
devops-cli terraform init /path/to/terraform/dir

# Plan changes
devops-cli terraform plan

# Apply changes (with confirmation prompt)
devops-cli terraform apply

# Destroy resources (with confirmation prompt)
devops-cli terraform destroy
```

### Script Automation

```bash
# Run a specific bash script
devops-cli script bash ./scripts/deploy.sh

# List available scripts in current directory
devops-cli script list

# List scripts in specific directory
devops-cli script list ./automation-scripts/
```

### CI/CD Pipeline Management

```bash
# List available pipeline templates
devops-cli pipeline templates

# Create a pipeline from template
devops-cli pipeline create node    # For Node.js projects
devops-cli pipeline create go      # For Go projects
devops-cli pipeline create python  # For Python projects

# Push pipeline to repository
devops-cli pipeline push
```

### Configuration

```bash
# Show current configuration
devops-cli config show

# Set configuration values
devops-cli config set terraform-path /usr/local/bin/terraform
devops-cli config set workspace-dir ./my-workspace
devops-cli config set git-branch develop
```

## Pipeline Templates

The CLI includes built-in templates for common technology stacks:

### Node.js Template
- npm ci and build
- Unit tests and linting
- Deployment steps

### Go Template
- Module download and build
- Tests and static analysis
- Binary compilation

### Python Template
- Requirements installation
- pytest and flake8
- Package building

## Directory Structure

```
your-project/
├── .github/
│   └── workflows/          # Generated CI/CD pipelines
├── terraform/              # Terraform configurations
├── scripts/               # Automation bash scripts
├── devops-cli             # The CLI binary
└── README.md
```

## Configuration Options

| Key | Description | Default |
|-----|-------------|---------|
| `terraform-path` | Path to Terraform binary | `terraform` |
| `workspace-dir` | Default workspace directory | `./workspace` |
| `git-remote` | Git remote URL | (empty) |
| `git-branch` | Default Git branch | `main` |

## Examples

### Complete Terraform Workflow

```bash
# Initialize and apply Terraform configuration
devops-cli terraform init ./infrastructure
devops-cli terraform plan ./infrastructure
devops-cli terraform apply ./infrastructure
```

### Script Automation Workflow

```bash
# List available automation scripts
devops-cli script list ./automation

# Run deployment script
devops-cli script bash ./automation/deploy-staging.sh
```

### CI/CD Setup Workflow

```bash
# Create a new pipeline for Go project
devops-cli pipeline create go

# Customize the generated pipeline file if needed
# vim .github/workflows/go.yml

# Push to repository
devops-cli pipeline push
```

## Advanced Usage

### Custom Script Discovery

The CLI automatically discovers `.sh` files in directories:

```bash
# This will show all .sh files in the specified directory
devops-cli script list ./my-scripts/
```

### Pipeline Customization

After generating a pipeline, you can customize it by editing the generated YAML file in `.github/workflows/`. The CLI creates a solid foundation that you can extend based on your specific needs.

### Environment-Specific Configuration

You can set different configurations for different environments:

```bash
# Development environment
devops-cli config set workspace-dir ./dev-workspace
devops-cli config set git-branch develop

# Production environment  
devops-cli config set workspace-dir ./prod-workspace
devops-cli config set git-branch main
```

## Troubleshooting

### Common Issues

1. **Terraform not found**: Set the correct path using `devops-cli config set terraform-path /path/to/terraform`

2. **Script permission denied**: Ensure your bash scripts have execute permissions: `chmod +x script.sh`

3. **Git push fails**: Make sure you have proper Git credentials configured and the repository is initialized

### Debug Mode

For verbose output, you can modify the source code to add debug flags or check the execution of individual commands.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [ ] Docker integration for containerized deployments
- [ ] Kubernetes manifest generation and deployment
- [ ] AWS/Azure/GCP cloud provider integrations
- [ ] Configuration file support (YAML/JSON)
- [ ] Plugin system for custom commands
- [ ] Interactive mode for guided operations
- [ ] Logging and audit trail features