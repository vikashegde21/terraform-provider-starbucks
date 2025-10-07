# Terraform Provider for Starbucks

A Terraform provider for managing Starbucks infrastructure including stores, employees, menu items, inventory, and promotions.

## Features

- **Store Management**: Create and manage store locations with full configuration
- **Employee Management**: Manage employee (partner) lifecycle and assignments
- **Menu Items**: Configure menu offerings across locations
- **Inventory**: Track and manage inventory across stores
- **Promotions**: Create and manage promotional campaigns

## Installation

### Build from Source

```bash
git clone https://github.com/yourusername/terraform-provider-starbucks
cd terraform-provider-starbucks
make install
```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    starbucks = {
      source  = "registry.terraform.io/starbucks/starbucks"
      version = "~> 0.1"
    }
  }
}

provider "starbucks" {
  api_key  = var.starbucks_api_key  # or set STARBUCKS_API_KEY env var
  endpoint = "https://api.starbucks.com/v1"
  region   = "us-west-2"
  timeout  = 30
}
```

## Development

### Prerequisites

- Go 1.21+
- Terraform 1.5+

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

## License

Apache 2.0
