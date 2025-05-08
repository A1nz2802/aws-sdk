# AWS SDK Go Examples

This project demonstrates the usage of AWS SDK for Go (v2). It provides examples of common DynamoDB operations and AWS service integrations.

## Prerequisites

- Go 1.24.2 or later
- AWS account and credentials configured
- Environment variables set up (see ``.env.example``)

## Installation

1. Clone the repository.
2. Copy ``.env.example`` to .env and fill in your AWS credentials.
3. Install dependencies:
    ```bash
    go mod download
    ```

## Project Structure
- ``dynamodb/`` - DynamoDB operations and queries
- ``rds/`` - RDS related operations
- ``main.go`` - Example usage and entry point

## Environment Variables

Create a .env file with the following variables:

```css
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=your_region
```

## Usage
Run the example:
```sh
go run main.go
```