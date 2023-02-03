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

### Run

Example points to settings file `settings.json` in current directory and saves json data to given file set in settings
and prints json data to stdout.

Clone this repository and change into repository directory

```bash
git clone https://github.com/cklewar/f5xc_aws_ami_mapper
```

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
....truncated 
{
  "ingress_egress_gateway": {
    "af-south-1": {
      "ami": "ami-0c22728f79f714ed1",
      "creationDate": "2022-07-01T11:27:36.000Z"
    },
    "ap-east-1": {
      "ami": "ami-0a6cf3665c0612f91",
      "creationDate": "2022-07-01T11:27:36.000Z"
    },
    "ap-northeast-1": {
      "ami": "ami-0384d075a36447e2a",
      "creationDate": "2022-07-01T11:27:35.000Z"
    },
    "ap-northeast-2": {
      "ami": "ami-01472d819351faf92",
      "creationDate": "2022-07-01T11:27:36.000Z"
    },
    "ap-northeast-3": {
      "ami": "",
      "creationDate": ""
    },
....truncated
```