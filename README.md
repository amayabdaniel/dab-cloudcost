# dab-cloudcost

Cloud cost analyzer CLI for AWS and GCP.

## Features

- AWS Cost Explorer integration
- GCP Billing API integration
- Resource right-sizing recommendations
- Unused resource detection
- CSV/JSON export

## Installation

```bash
go install github.com/amayabdaniel/dab-cloudcost/cmd/dab-cloudcost@latest
```

## Usage

```bash
# show help
dab-cloudcost --help

# analyze aws costs
dab-cloudcost aws --profile default --days 30

# analyze gcp costs
dab-cloudcost gcp --project my-project --days 30

# export to csv
dab-cloudcost aws --output csv > costs.csv
```

## Development

```bash
# build
make build

# test
make test

# run
make run
```

## License

MIT
