# dab-cloudcost

Cloud cost analyzer CLI for AWS and GCP.

## Features

- AWS Cost Explorer integration
- Cost breakdown by service
- Sorted by cost (highest first)
- Multiple output formats (table, json, csv)
- Filter top N services

## Installation

```bash
go install github.com/amayabdaniel/dab-cloudcost/cmd/dab-cloudcost@latest
```

## Usage

```bash
# show help
dab-cloudcost --help

# analyze aws costs (last 30 days)
dab-cloudcost aws

# analyze last 7 days
dab-cloudcost aws --days 7

# use specific aws profile
dab-cloudcost aws --profile production

# show top 5 services
dab-cloudcost aws --top 5

# output as json
dab-cloudcost aws --output json

# output as csv
dab-cloudcost aws --output csv > costs.csv

# combine flags
dab-cloudcost aws -d 7 -t 10 -o json
```

## Example Output

```
SERVICE                    COST     UNIT
-------                    ----     ----
Amazon EC2                 142.50   USD
Amazon S3                  45.20    USD
AWS Lambda                 12.30    USD
Amazon RDS                 8.50     USD
-------                    ----     ----
TOTAL                      208.50   USD
```

## Development

```bash
# build
make build

# test
make test

# run
make run

# install locally
make install
```

## Requirements

- Go 1.24+
- AWS credentials configured (`aws configure`)

## License

MIT
