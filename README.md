# F5XC AWS AMI Mapper Tool

Mapper tool generates mapping between F5XC gateway types and latest publicly available AWS AMIs for later consumption in
for example Terraform.

## Command line arguments

```bash
go run main.go --help
  -print
        Print mapping in stdout
  -settings string
        The settings file to use
  -writeHcl
        Write mapping hcl file
  -writeJson
        Write mapping json file
```

## Settings file

- `certifiedHardwareUrl` most current certified hardware list
- `awsRegionsFile` region file to obtain list of AWS regions from
- `mappingsFile` name to write json data to

```json
{
  "certifiedHardwareUrl": "https://vesio.blob.core.windows.net/releases/certified-hardware/aws.yml",
  "awsRegionsFile": "regions.json",
  "mappingsFile": "mapping.json",
  "templateFile": "machine_images.tpl",
  "machineImagesFile": "machine_images.tf"
}
```

## Usage

Clone this repository and change into repository directory

```bash
git clone https://github.com/cklewar/f5xc_aws_ami_mapper
```

### Export - JSON

Example points to settings file `settings.json` in current directory and saves json data to given file set in settings
file.

```bash
go run main.go -settings settings.json -writeJson true
```

#### Example output

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

### Export - HCL

Example points to settings file `settings.json` in current directory and saves hcl data to given file set in settings
file.

```bash
go run main.go -settings settings.json -writeHcl true
```

#### Example output

```hcl
variable "f5xc_ce_machine_image" {
  type = object({
    ingress_gateway = object({
      af-south-1     = string
      ap-east-1      = string
      ap-northeast-1 = string
      ap-northeast-2 = string
      ap-northeast-3 = string
      ap-south-1     = string
      ap-southeast-1 = string
      ap-southeast-2 = string
      ap-southeast-3 = string
      ca-central-1   = string
      eu-central-1   = string
      eu-north-1     = string
      eu-south-1     = string
      eu-west-1      = string
      eu-west-2      = string
      eu-west-3      = string
      me-south-1     = string
      sa-east-1      = string
      us-east-1      = string
      us-east-2      = string
      us-west-1      = string
      us-west-2      = string
    })
    ingress_egress_gateway = object({
      af-south-1     = string
      ap-east-1      = string
      ap-northeast-1 = string
      ap-northeast-2 = string
      ap-northeast-3 = string
      ap-south-1     = string
      ap-southeast-1 = string
      ap-southeast-2 = string
      ap-southeast-3 = string
      ca-central-1   = string
      eu-central-1   = string
      eu-north-1     = string
      eu-south-1     = string
      eu-west-1      = string
      eu-west-2      = string
      eu-west-3      = string
      me-south-1     = string
      sa-east-1      = string
      us-east-1      = string
      us-east-2      = string
      us-west-1      = string
      us-west-2      = string
    })
  })
  default = {
    ingress_gateway = {
      af-south-1     = "ami-0bcfb554a48878b52"
      ap-east-1      = "ami-03cf35954fb9084fc"
      ap-northeast-1 = "ami-07dac882268159d52"
      ap-northeast-2 = "ami-04f6d5781039d2f88"
      ap-northeast-3 = ""
      ap-south-1     = "ami-099c0c7e19e1afd16"
      ap-southeast-1 = "ami-0dba294abe676bd58"
      ap-southeast-2 = "ami-0ae68f561b7d20682"
      ap-southeast-3 = "ami-065fc7b0f6ec02011"
      ca-central-1   = "ami-0ddc009ae69986eb4"
      eu-central-1   = "ami-027625cb269f5d7e9"
      eu-north-1     = "ami-0366c929eb2ac407b"
      eu-south-1     = "ami-00cb6474298a310af"
      eu-west-1      = "ami-01baaca2a3b1b0114"
      eu-west-2      = "ami-05f5a414a42961df6"
      eu-west-3      = "ami-0e1361351f9205511"
      me-south-1     = "ami-0fb5db9d908d231c3"
      sa-east-1      = "ami-09082c4758ef6ec36"
      us-east-1      = "ami-0f94aee77d07b0094"
      us-east-2      = "ami-0660aaf7b6edaa980"
      us-west-1      = "ami-0cf44e35e2aecacb4"
      us-west-2      = "ami-0cba83d31d405a8f5"
    }
    ingress_egress_gateway = {
      af-south-1     = "ami-0c22728f79f714ed1"
      ap-east-1      = "ami-0a6cf3665c0612f91"
      ap-northeast-1 = "ami-0384d075a36447e2a"
      ap-northeast-2 = "ami-01472d819351faf92"
      ap-northeast-3 = ""
      ap-south-1     = "ami-0277ab0b4db359c93"
      ap-southeast-1 = "ami-0d6463ee1e3727e84"
      ap-southeast-2 = "ami-03ff18dfb7f90eb54"
      ap-southeast-3 = "ami-0189e67b4c856e4cd"
      ca-central-1   = "ami-052252c245ff77338"
      eu-central-1   = "ami-06d5e0073d97ecf99"
      eu-north-1     = "ami-006c465449ed98c69"
      eu-south-1     = "ami-0baafa10ffcd081b7"
      eu-west-1      = "ami-090680f491ad6d46a"
      eu-west-2      = "ami-0df8a483722043a41"
      eu-west-3      = "ami-03bd7c41ca1b586a8"
      me-south-1     = "ami-094efc1a78169dd7c"
      sa-east-1      = "ami-07369c4b06cf22299"
      us-east-1      = "ami-089311edbe1137720"
      us-east-2      = "ami-01ba94b5a83adcb35"
      us-west-1      = "ami-092a2a07d2d3a445f"
      us-west-2      = "ami-07252e5ab4023b8cf"
    }
  }
}
```