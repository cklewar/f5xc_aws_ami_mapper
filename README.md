# F5XC AWS AMI Mapper Tool

Mapper tool generates mapping between F5XC gateway types and latest publicly available AWS AMIs for later consumption in
for example Terraform.

## Usage

### Positional command line arguments

1. Argument points to `settings` file to use
2. Argument `WriteMappingFile`. If set to `true` writes json data to file. Filename take from `settings` file
3. Argument `PrettyPrintMapping`. If set to `true` prints json data to stdout

### Settings file

- `certifiedHardwareUrl` most current certified hardware list
- `awsRegionsFile` region file to obtain list of AWS regions from
- `mappingsFile` name to write json data to

```json
{
  "certifiedHardwareUrl": "https://vesio.blob.core.windows.net/releases/certified-hardware/aws.yml",
  "awsRegionsFile": "regions.json",
  "mappingsFile": "mapping.json"
}
```

### Run example

Example points to settings file `settings.json` in current directory and saves json data to given file set in settings
and prints json data to stdout.

```bash
go run main.go settings.json true true
```

## Example output

```bash
Successfully opened settings.json
Client: MachineImages status code: 200
Successfully Opened regions.json
Create mapping for region: af-south-1
Create mapping for region: af-south-1 --> Done
Create mapping for region: ap-south-1
Create mapping for region: ap-south-1 --> Done
Create mapping for region: eu-north-1
Create mapping for region: eu-north-1 --> Done
Create mapping for region: eu-west-3
Create mapping for region: eu-west-3 --> Done
Create mapping for region: eu-south-1
Create mapping for region: eu-south-1 --> Done
Create mapping for region: eu-west-2
Create mapping for region: eu-west-2 --> Done
Create mapping for region: eu-west-1
Create mapping for region: eu-west-1 --> Done
Create mapping for region: ap-northeast-3
Create mapping for region: ap-northeast-3 --> Done
Create mapping for region: ap-northeast-2
Create mapping for region: ap-northeast-2 --> Done
Create mapping for region: me-south-1
Create mapping for region: me-south-1 --> Done
Create mapping for region: ap-northeast-1
Create mapping for region: ap-northeast-1 --> Done
Create mapping for region: ca-central-1
Create mapping for region: ca-central-1 --> Done
Create mapping for region: sa-east-1
Create mapping for region: sa-east-1 --> Done
Create mapping for region: ap-east-1
Create mapping for region: ap-east-1 --> Done
Create mapping for region: ap-southeast-1
Create mapping for region: ap-southeast-1 --> Done
Create mapping for region: ap-southeast-2
Create mapping for region: ap-southeast-2 --> Done
Create mapping for region: eu-central-1
Create mapping for region: eu-central-1 --> Done
Create mapping for region: ap-southeast-3
Create mapping for region: ap-southeast-3 --> Done
Create mapping for region: us-east-1
Create mapping for region: us-east-1 --> Done
Create mapping for region: us-east-2
Create mapping for region: us-east-2 --> Done
Create mapping for region: us-west-1
Create mapping for region: us-west-1 --> Done
Create mapping for region: us-west-2
Create mapping for region: us-west-2 --> Done
```