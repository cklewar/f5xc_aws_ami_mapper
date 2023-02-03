package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"time"
)

// const Version = "0.0.1"
const Owner = "434481986642"
const AwsByolVoltstackCombo = "AwsByolVoltstackCombo"
const AwsByolMultiNicVoltmesh = "AwsByolMultiNicVoltmesh"
const AwsByolVoltmesh = "AwsByolVoltmesh"

// const timeParseLayout = "2006-01-02 15:04:05 +0000 UTC"
const IngressGatewayType = "ingress_gateway"
const IngressEgressGatewayType = "ingress_egress_gateway"
const IngressEgressGatewayVolstackType = "ingress_egress_voltstack_gateway"

var AwsType2GwTypeMap = map[string]string{
	AwsByolVoltmesh:         IngressGatewayType,
	AwsByolVoltstackCombo:   IngressEgressGatewayVolstackType,
	AwsByolMultiNicVoltmesh: IngressEgressGatewayType,
}

type EC2DescribeImagesAPI interface {
	DescribeImages(ctx context.Context,
		params *ec2.DescribeImagesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
}

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

func DescribeAwsImages() (AvailableAwsImages, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
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

func build(currentMachineImage CurrentMachineImage, certifiedHardwareImageId string, availableAwsImage AvailableAwsImage, region string) (CurrentMachineImage, error) {
	if aws.StringValue(availableAwsImage.ImageId) == certifiedHardwareImageId {
		if currentMachineImage == (CurrentMachineImage{}) {
			currentMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
			currentMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
			currentMachineImage.Region = region
		} else if currentMachineImage != (CurrentMachineImage{}) {
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

	return currentMachineImage, nil
}

func (r *CertifiedHardwareImages) Contains(mis MachineImages, dai AvailableAwsImages, region string) (MachineImages, error) {
	var currentIngressEgressMachineImage CurrentMachineImage
	var currentIngressMachineImage CurrentMachineImage

	for _, availableAwsImage := range dai {
		for _, certifiedHardwareImageId := range r.CertifiedHardware.AwsByolMultiNicVoltmesh.Aws.ImageID {
			/*_, err := build(currentIngressEgressMachineImage, certifiedHardwareImageId, AvailableAwsImage(availableAwsImage), region)
			if err != nil {
				return nil, err
			}*/
			if aws.StringValue(availableAwsImage.ImageId) == certifiedHardwareImageId {
				if currentIngressEgressMachineImage == (CurrentMachineImage{}) {
					currentIngressEgressMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
					currentIngressEgressMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
					currentIngressEgressMachineImage.Region = region
				} else if currentIngressEgressMachineImage != (CurrentMachineImage{}) {
					currentDate, err := time.Parse(time.RFC3339, currentIngressEgressMachineImage.CreationDate)
					if err != nil {
						fmt.Println("current: ", err)
					}

					candidateDate, err := time.Parse(time.RFC3339, aws.StringValue(availableAwsImage.CreationDate))
					if err != nil {
						fmt.Println("candidate: ", err)
					}

					if currentDate.Before(candidateDate) {
						currentIngressEgressMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
						currentIngressEgressMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
						currentIngressEgressMachineImage.Region = region
					}
				}
			}
		}

		for _, certifiedHardwareImageId := range r.CertifiedHardware.AwsByolVoltmesh.Aws.ImageID {
			cur, err := build(currentIngressMachineImage, certifiedHardwareImageId, AvailableAwsImage(availableAwsImage), region)
			if err != nil {
				return nil, err
			}
			currentIngressMachineImage = cur
			/*if aws.StringValue(availableAwsImage.ImageId) == certifiedHardwareImageId {
				if currentIngressMachineImage == (CurrentMachineImage{}) {
					currentIngressMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
					currentIngressMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
					currentIngressMachineImage.Region = region
				} else if currentIngressMachineImage != (CurrentMachineImage{}) {
					currentDate, err := time.Parse(time.RFC3339, currentIngressMachineImage.CreationDate)
					if err != nil {
						fmt.Println("current: ", err)
					}

					candidateDate, err := time.Parse(time.RFC3339, aws.StringValue(availableAwsImage.CreationDate))
					if err != nil {
						fmt.Println("candidate: ", err)
					}

					if currentDate.Before(candidateDate) {
						currentIngressMachineImage.Ami = aws.StringValue(availableAwsImage.ImageId)
						currentIngressMachineImage.CreationDate = aws.StringValue(availableAwsImage.CreationDate)
						currentIngressMachineImage.Region = region
					}
				}
			}*/
		}
	}

	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region] = make(map[string]string)
	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region]["creationDate"] = currentIngressEgressMachineImage.CreationDate
	mis[AwsType2GwTypeMap[AwsByolMultiNicVoltmesh]][region]["ami"] = currentIngressEgressMachineImage.Ami
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region] = make(map[string]string)
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region]["creationDate"] = currentIngressMachineImage.CreationDate
	mis[AwsType2GwTypeMap[AwsByolVoltmesh]][region]["ami"] = currentIngressMachineImage.Ami

	return mis, nil
}

func CreateAwsSession(region string) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		fmt.Printf("Got an error in creating session: %s\n", err.Error())
		return err
	}

	_, err = sess.Config.Credentials.Get()
	fmt.Println("SESSION REGION", aws.StringValue(sess.Config.Region))
	if err != nil {
		fmt.Printf("Session credentials not found: %s\n", err.Error())
		return err
	}

	return nil
}

func main() {
	argSettingsFile := os.Args[1]

	var settings Settings
	err := settings.loadSettings(argSettingsFile)
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
	/*var regions = []string{
		"us-west-1",
		//"us-west-2",
		"eu-central-1",
		"ap-southeast-1",
	}*/

	for _, region := range r.Regions { //for _, region := range regions {
		err := CreateAwsSession(region.RegionName)
		// err := CreateAwsSession(region)
		if err != nil {
			os.Exit(1)
		}

		dai, err1 := DescribeAwsImages()
		if err1 != nil {
			os.Exit(1)
		}

		_, err = certifiedHardwareImages.Contains(mis, dai, region.RegionName)
		//_, err = certifiedHardwareImages.Contains(mis, dai, region)
		if err != nil {
			os.Exit(1)
		}
	}
	err = PrettyPrint(mis)
	if err != nil {
		return
	}
	// fmt.Println(mis, len(mis))
}
