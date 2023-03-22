package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-google-go/google/v5/computeinstance"
	"github.com/cdktf/cdktf-provider-google-go/google/v5/computenetwork"
	"github.com/cdktf/cdktf-provider-google-go/google/v5/computesubnetwork"
	"github.com/cdktf/cdktf-provider-google-go/google/v5/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"path/filepath"
)

func NewProvider(stack cdktf.TerraformStack) {
	_ = provider.NewGoogleProvider(stack, jsii.String("google"), &provider.GoogleProviderConfig{
		Project:     jsii.String(varProject),
		Region:      jsii.String(varRegion),
		Zone:        jsii.String(varZone),
		Credentials: jsii.String(filepath.Base(varCredsFile)),
	})
}

func NewComputeNetwork(stack cdktf.TerraformStack) computenetwork.ComputeNetwork {
	return computenetwork.NewComputeNetwork(
		stack,
		jsii.String("terraform_network"),
		&computenetwork.ComputeNetworkConfig{
			Name:                  jsii.String("terraform-network"),
			AutoCreateSubnetworks: jsii.Bool(false),
		})
}

func NewComputeSubNetwork(
	stack cdktf.TerraformStack,
	computeNetwork computenetwork.ComputeNetwork,
) computesubnetwork.ComputeSubnetwork {
	return computesubnetwork.NewComputeSubnetwork(
		stack,
		jsii.String("terraform_subnet"),
		&computesubnetwork.ComputeSubnetworkConfig{
			IpCidrRange: jsii.String("10.20.0.0/16"),
			Name:        jsii.String("terraform-subnetwork"),
			Region:      jsii.String(varRegion),
			Network:     computeNetwork.Id(),
		})
}

func ComputeInstanceBootDisk() computeinstance.ComputeInstanceBootDisk {
	return computeinstance.ComputeInstanceBootDisk{
		InitializeParams: &computeinstance.ComputeInstanceBootDiskInitializeParams{
			Image: jsii.String(varOSImage),
		},
	}
}

func NewComputeInstance(stack cdktf.TerraformStack) computeinstance.ComputeInstance {
	computeNetwork := NewComputeNetwork(stack)
	computeSubNetwork := NewComputeSubNetwork(stack, computeNetwork)
	computeInstanceBootDisk := ComputeInstanceBootDisk()

	return computeinstance.NewComputeInstance(
		stack,
		jsii.String("st_instance"),
		&computeinstance.ComputeInstanceConfig{
			MachineType: jsii.String(varMachineType),
			Name:        jsii.String(varName),
			Zone:        jsii.String(varZone),
			BootDisk:    &computeInstanceBootDisk,
			NetworkInterface: []computeinstance.ComputeInstanceNetworkInterface{
				{
					Network:      computeNetwork.SelfLink(),
					Subnetwork:   computeSubNetwork.SelfLink(),
					AccessConfig: []computeinstance.ComputeInstanceNetworkInterfaceAccessConfig{
						// necessary even empty
					},
				},
			},
			AllowStoppingForUpdate: jsii.Bool(varAllowStoppingForUpdate),
		})
}

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	NewProvider(stack)
	NewComputeInstance(stack)

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "learn-cdktf-gcp")

	app.Synth()
}

const (
	varProject                = "my-project-id"
	varCredsFile              = "credentials.json"
	varRegion                 = "europe-west6"
	varZone                   = "europe-west6-c"
	varMachineType            = "n1-standard-1"
	varName                   = "cdktf-instance"
	varOSImage                = "debian-cloud/debian-11"
	varAllowStoppingForUpdate = false
)
