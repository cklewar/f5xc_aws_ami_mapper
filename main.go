package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"
)

// const Version = "0.0.1"

const Owner = "434481986642"
const AwsByolVoltstackCombo = "AwsByolVoltstackCombo"
const AwsByolMultiNicVoltmesh = "AwsByolMultiNicVoltmesh"
const AwsByolVoltmesh = "AwsByolVoltmesh"
const IngressGatewayType = "ingress_gateway"
const IngressEgressGatewayType = "ingress_egress_gateway"
const IngressEgressGatewayVolstackType = "ingress_egress_voltstack_gateway"

var AwsType2GwTypeMap = map[string]string{
	AwsByolVoltmesh:         IngressGatewayType,
	AwsByolVoltstackCombo:   IngressEgressGatewayVolstackType,
	AwsByolMultiNicVoltmesh: IngressEgressGatewayType,
}

var (
	settingsFile  = flag.String("settings", "", "The settings file to use")
	writeJsonFile = flag.Bool("writeJson", false, "Write mapping json file")
	writeHclFile  = flag.Bool("writeHcl", false, "Write mapping hcl file")
	printMapping  = flag.Bool("print", false, "Print mapping in stdout")
)

/*type EC2DescribeImagesAPI interface {
	DescribeImages(ctx context.Context,
		params *ec2.DescribeImagesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
}*/

type CertifiedHardware []byte
type AvailableAwsImages []types.Image
type AvailableAwsImage types.Image
type MachineImages map[string]map[string]map[string]string

type CurrentMachineImage struct {
	Ami          string
	CreationDate string
	Region       string
}

type Settings struct {
	CertifiedHardwareUrl string `json:"certifiedHardwareUrl"`
	AwsRegionsFile       string `json:"awsRegionsFile"`
	MappingsFile         string `json:"mappingsFile""`
	TemplateFile         string `json:"templateFile"`
	MachineImagesFile    string `json:"machineImagesFile"`
}

type CertifiedHardwareImages struct {
	CertifiedHardware struct {
		AwsByolVoltmesh struct {
			Vpm struct {
				PrivateNIC string `yaml:"PrivateNIC"`
			} `yaml:"Vpm"`
			Aws struct {
				ImageID []string `yaml:"imageId"`
			} `yaml:"aws"`
		} `yaml:"aws-byol-voltmesh"`
		AwsByolMultiNicVoltmesh struct {
			Vpm struct {
				PrivateNIC string `yaml:"PrivateNIC"`
				InsideNIC  string `yaml:"InsideNIC"`
			} `yaml:"Vpm"`
			Aws struct {
				ImageID []string `yaml:"imageId"`
			} `yaml:"aws"`
		} `yaml:"aws-byol-multi-nic-voltmesh"`
		AwsByolVoltstackCombo struct {
			Vpm struct {
				PrivateNIC string `yaml:"PrivateNIC"`
			} `yaml:"Vpm"`
			Aws struct {
				ImageID []string `yaml:"imageId"`
			} `yaml:"aws"`
		} `yaml:"aws-byol-voltstack-combo"`
	} `yaml:"certifiedHardware"`
}

type AwsRegions struct {
	Regions []struct {
		Endpoint    string `json:"Endpoint"`
		RegionName  string `json:"RegionName"`
		OptInStatus string `json:"OptInStatus"`
	} `json:"Regions"`
}

/*func GetImages(c context.Context, api EC2DescribeImagesAPI, input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return api.DescribeImages(c, input)
}*/

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func (s *Settings) loadSettings(settingsFile string) error {
	jsonFile, err := os.Open(settingsFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Successfully opened %s\n", settingsFile)
	defer func(jsonFile *os.File) error {
		err := jsonFile.Close()
		if err != nil {
			fmt.Printf("Error closing file %s\n", Settings{})
			return err
		}
		return nil
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &s)
	if err != nil {
		fmt.Printf("Error reading settings json data %s\n", err.Error())
		return err
	}

	return nil
}

func (r *AwsRegions) loadAwsRegions(RegionsFile string) error {
	var _regions []string

	jsonFile, err := os.Open(RegionsFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Successfully Opened %s\n", RegionsFile)
	defer func(jsonFile *os.File) error {
		err := jsonFile.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", RegionsFile)
			return err
		}
		return nil
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		fmt.Printf("Error reading regions json data %s\n", err.Error())
		return err
	}

	for _, e := range r.Regions {
		_regions = append(_regions, e.RegionName)
	}

	return nil
}

func (r *CertifiedHardwareImages) Load(list CertifiedHardware) error {
	err := yaml.Unmarshal(list, r)
	if err != nil {
		fmt.Printf("Error reading CertifiedHardware from YAML %s\n", err.Error())
		return err
	}

	return nil
}

func GetCertifiedHardwareList(requestURL string) (CertifiedHardware, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("Client: could not create request: %s\n", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Client: error making http request: %s\n", err)
		return nil, err
	}

	fmt.Printf("Client: MachineImages status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return nil, err
	}
	return resBody, nil
}

func (m MachineImages) Save(filename string) {
	file, _ := json.MarshalIndent(m, "", "  ")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err1 := f.Write(file)
	if err1 != nil {
		panic(err)
	}
	fmt.Printf("Saved mapping data to %s. Wrote %d bytes\n", filename, n)
}

func (m MachineImages) ExportMachineImages2Hcl(machineImageFile string, templateFile string) {
	tpl := template.Must(template.ParseFiles(templateFile))

	f, err := os.Create(machineImageFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//err1 := tpl.Execute(os.Stdout, m)
	err1 := tpl.Execute(f, m)
	if err1 != nil {
		fmt.Println(err1)
	}
}

func DescribeAwsImages(region string) (AvailableAwsImages, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

	if err != nil {
		fmt.Println("configuration error: " + err.Error())
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	input := ec2.DescribeImagesInput{
		ImageIds: []string{},
		Filters: []types.Filter{
			{
				Name:   aws.String("is-public"),
				Values: []string{"true"},
			},
		},
		Owners: []string{Owner},
	}

	result, err := client.DescribeImages(context.Background(), &input)
	if err != nil {
		fmt.Println("Error describing images: " + err.Error())
		return nil, err
	}

	return result.Images, err
}

func build(currentMachineImage *CurrentMachineImage, certifiedHardwareImageId string, availableAwsImage AvailableAwsImage, region string) (CurrentMachineImage, error) {
	if aws.StringValue(availableAwsImage.ImageId) == certifiedHardwareImageId {
		if *currentMachineImage == (CurrentMachineImage{}) {
			currentMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
			currentMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
			currentMachineImage.Region = region
		} else if *currentMachineImage != (CurrentMachineImage{}) {
			currentDate, err := time.Parse(time.RFC3339, currentMachineImage.CreationDate)
			if err != nil {
				fmt.Println("current: ", err)
			}

			candidateDate, err := time.Parse(time.RFC3339, aws.StringValue(availableAwsImage.CreationDate))
			if err != nil {
				fmt.Println("candidate: ", err)
			}

			if currentDate.Before(candidateDate) {
				currentMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
				currentMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
				currentMachineImage.Region = region
			}
		}
	}

	return *currentMachineImage, nil
}

func (r *CertifiedHardwareImages) Contains(mis MachineImages, dai AvailableAwsImages, region string) (MachineImages, error) {
	var currentIngressEgressMachineImage CurrentMachineImage
	var currentIngressMachineImage CurrentMachineImage

	for _, availableAwsImage := range dai {
		for _, certifiedHardwareImageId := range r.CertifiedHardware.AwsByolMultiNicVoltmesh.Aws.ImageID {
			_, err := build(&currentIngressEgressMachineImage, certifiedHardwareImageId, AvailableAwsImage(availableAwsImage), region)
			if err != nil {
				return nil, err
			}

		}

		/*for _, certifiedHardwareImageId := range r.CertifiedHardware.AwsByolVoltmesh.Aws.ImageID {
			_, err := build(&currentIngressMachineImage, certifiedHardwareImageId, AvailableAwsImage(availableAwsImage), region)
			if err != nil {
				return nil, err
			}
		}*/

	}

	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region] = make(map[string]string)
	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region]["creationDate"] = currentIngressEgressMachineImage.CreationDate
	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region]["ami"] = currentIngressEgressMachineImage.Ami
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region] = make(map[string]string)
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region]["creationDate"] = currentIngressMachineImage.CreationDate
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region]["ami"] = currentIngressMachineImage.Ami

	return mis, nil
}

func main() {
	flag.Parse()

	var settings Settings
	err := settings.loadSettings(*settingsFile)
	if err != nil {
		os.Exit(1)
	}

	chl, err2 := GetCertifiedHardwareList(settings.CertifiedHardwareUrl)
	if err2 != nil {
		os.Exit(1)
	}

	var certifiedHardwareImages CertifiedHardwareImages
	err = certifiedHardwareImages.Load(chl)
	if err != nil {
		os.Exit(1)
	}

	var r AwsRegions
	err = r.loadAwsRegions(settings.AwsRegionsFile)
	if err != nil {
		os.Exit(1)
	}

	var mis MachineImages
	mis = make(map[string]map[string]map[string]string)
	mis[IngressGatewayType] = make(map[string]map[string]string)
	mis[IngressEgressGatewayType] = make(map[string]map[string]string)
	mis[IngressEgressGatewayVolstackType] = make(map[string]map[string]string)

	for _, region := range r.Regions {
		fmt.Printf("Create mapping for region: %s\n", region.RegionName)
		dai, err := DescribeAwsImages(region.RegionName)
		if err != nil {
			os.Exit(1)
		}

		_, err1 := certifiedHardwareImages.Contains(mis, dai, region.RegionName)
		if err1 != nil {
			os.Exit(1)
		}
		fmt.Printf("Create mapping for region: %s --> Done\n", region.RegionName)
	}

	fmt.Println(mis["ingress_egress_gateway"])

	if *writeHclFile {
		mis.ExportMachineImages2Hcl(settings.MachineImagesFile, settings.TemplateFile)
	}

	if *writeJsonFile {
		mis.Save(settings.MappingsFile)
	}

	if *printMapping {
		err = PrettyPrint(mis)
	}
}
