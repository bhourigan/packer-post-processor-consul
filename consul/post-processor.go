package consul

import (
	"fmt"
	"strings"
	"encoding/json"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
        "github.com/hashicorp/consul/api"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

var builtins = map[string]string{
	"mitchellh.amazonebs": "amazonebs",
	"mitchellh.amazon.instance": "amazoninstance",
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	AwsAccessKey     string `mapstructure:"aws_access_key"`
	AwsSecretKey     string `mapstructure:"aws_secret_key"`
	AwsToken         string `mapstructure:"aws_token"`
	ConsulAddress    string `mapstructure:"consul_address"`
	ConsulScheme     string `mapstructure:"consul_scheme"`
	ConsulToken      string `mapstructure:"consul_token"`

	ProjectName      string `mapstructure:"project_name"`
	ProjectVersion   string `mapstructure:"project_version"`

	tpl *packer.ConfigTemplate
}

type PostProcessor struct {
	config Config
	client *api.Client
	auth aws.Auth
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	_, err := common.DecodeConfig(&p.config, raws...)
	if err != nil {
		return err
	}

	p.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return err
	}
	p.config.tpl.UserVars = p.config.PackerUserVars

	templates := map[string]*string{
		"consul_address":    &p.config.ConsulAddress,
		"consul_scheme":     &p.config.ConsulScheme,
		"consul_token":      &p.config.ConsulToken,
		"aws_access_key":    &p.config.AwsAccessKey,
		"aws_secret_key":    &p.config.AwsSecretKey,
		"aws_token":         &p.config.AwsToken,
		"project_name":      &p.config.ProjectName,
                "project_version":   &p.config.ProjectVersion,
	}

	errs := new(packer.MultiError)
	for key, ptr := range templates {
		*ptr, err = p.config.tpl.Process(*ptr, nil)
		if err != nil {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("Error processing %s: %s", key, err))
		}
	}

	required := map[string]*string{
		"consul_address":      &p.config.ConsulAddress,
		"aws_access_key":      &p.config.AwsAccessKey,
		"aws_secret_key":      &p.config.AwsSecretKey,
		"project_name":        &p.config.ProjectName,
                "project_version":     &p.config.ProjectVersion,
	}

	for key, ptr := range required {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("%s must be set", key))
		}
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	p.auth, err = aws.GetAuth(p.config.AwsAccessKey, p.config.AwsSecretKey)
	if err != nil {
		return err
	}

	if p.config.AwsToken != "" {
		p.config.AwsToken = p.auth.Token
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	_, ok := builtins[artifact.BuilderId()]
	if !ok {
		return nil, false, fmt.Errorf(
			"Unsupported artifact type: %s", artifact.BuilderId())
	}

	ui.Say("Putting build artifacts into consul")

	for _, regions := range strings.Split(artifact.Id(), ",") {
		parts := strings.Split(regions, ":")
		if len(parts) != 2 {
			err := fmt.Errorf("Poorly formatted artifact ID: %s", artifact.Id())
			return nil, false, err
		}

		regionconn := ec2.New(p.auth, aws.Regions[parts[0]])
		ids := []string{parts[1]}
		if images, err := regionconn.Images(ids, nil); err == nil {
			config := api.DefaultConfig()
			config.Address = p.config.ConsulAddress
			config.Datacenter = parts[0]

			if p.config.ConsulScheme != "" {
				config.Scheme = p.config.ConsulScheme
			}

			if p.config.ConsulToken != "" {
				config.Token = p.config.ConsulToken
			}

		        client, err := api.NewClient(config)
		        if err == nil {
				kv := client.KV()
				consul_key_prefix := fmt.Sprintf("aws/%s/%s/%s", images.Images[0].RootDeviceType, p.config.ProjectName, p.config.ProjectVersion)

				ui.Message(fmt.Sprintf("Putting %s image data into consul key prefix %s in datacenter %s",
					parts[1], consul_key_prefix, config.Datacenter))

				consul_data_key := fmt.Sprintf("%s/data", consul_key_prefix)
				ami_data, _ := json.Marshal(images.Images)
				kv_ami_data := &api.KVPair{Key: consul_data_key, Value: ami_data}
				_, err := kv.Put(kv_ami_data, nil)
				if err != nil {
					return artifact, false, err
				}

				consul_ami_key := fmt.Sprintf("%s/ami", consul_key_prefix)
				kv_ami_id := &api.KVPair{Key: consul_ami_key, Value: []byte(parts[1])}
				_, err = kv.Put(kv_ami_id, nil)

				if err != nil {
					return artifact, false, err
				}
			} else {
				return artifact, false, err
		        }
		} else {
			return artifact, false, err
		}
	}

	return artifact, true, nil
}
