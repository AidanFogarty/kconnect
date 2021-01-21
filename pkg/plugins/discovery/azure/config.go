/*
Copyright 2020 The kconnect Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	"context"
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-09-01/containerservice"

	azclient "github.com/fidelity/kconnect/pkg/azure/client"
	"github.com/fidelity/kconnect/pkg/azure/id"
	"github.com/fidelity/kconnect/pkg/provider"
)

const (
	AKSAADServerAppID = "6dae42f8-4368-4678-94ff-3960e28e3630"
)

func (p *aksClusterProvider) GetClusterConfig(ctx *provider.Context, cluster *provider.Cluster, namespace string) (*api.Config, string, error) {
	p.logger.Debug("getting cluster config")

	cfg, err := p.getKubeconfig(ctx.Context, cluster)
	if err != nil {
		return nil, "", fmt.Errorf("getting kubeconfig: %w", err)
	}
	if !p.config.Admin {
		if err := p.addKubelogin(cfg); err != nil {
			return nil, "", fmt.Errorf("adding kubelogin: %w", err)
		}
		p.printLoginDetails()
	}

	if namespace != "" {
		p.logger.Debugw("setting kubernetes namespace", "namespace", namespace)
		cfg.Contexts[cfg.CurrentContext].Namespace = namespace
	}

	return cfg, cfg.CurrentContext, nil
}

func (p *aksClusterProvider) printLoginDetails() {
	switch p.config.LoginType {
	case LoginTypeResourceOwnerPassword:
		fmt.Fprintf(os.Stderr, "\033[33mSet the AAD_USER_PRINCIPAL_NAME and AAD_USER_PRINCIPAL_PASSWORD environment variables before running kubectl\033[0m\n")
	case LoginTypeServicePrincipal:
		fmt.Fprintf(os.Stderr, "\033[33mSet the AAD_SERVICE_PRINCIPAL_CLIENT_ID and AAD_SERVICE_PRINCIPAL_CLIENT_SECRET environment variables before running kubectl\033[0m\n")
	}

	if p.config.AzureEnvironment == AzureEnvironmentStackCloud {
		fmt.Fprintf(os.Stderr, "\033[33mSet the Azure Stack URLs in a config file and set the AZURE_ENVIRONMENT_FILEPATH environment variable to the path of that file\033[0m\n")
	}
}

func (p *aksClusterProvider) addKubelogin(cfg *api.Config) error {
	contextName := cfg.CurrentContext
	context := cfg.Contexts[contextName]
	userName := context.AuthInfo

	execConfig := &api.ExecConfig{
		APIVersion: "client.authentication.k8s.io/v1beta1",
		Command:    "kubelogin",
		Args: []string{
			"get-token",
			"--environment",
			mapAzureEnvironment(p.config.AzureEnvironment),
			"--server-id",
			AKSAADServerAppID,
			"--client-id",
			p.config.ClientID,
			"--tenant-id",
			p.config.TenantID,
			"--login",
			string(p.config.LoginType),
		},
	}

	cfg.AuthInfos = map[string]*api.AuthInfo{
		userName: {
			Exec: execConfig,
		},
	}

	return nil
}

func (p *aksClusterProvider) getKubeconfig(ctx context.Context, cluster *provider.Cluster) (*api.Config, error) {
	resourceID, err := id.FromClusterID(cluster.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing cluster id: %w", err)
	}

	client := azclient.NewContainerClient(resourceID.SubscriptionID, p.authorizer)

	var credentialList containerservice.CredentialResults
	if p.config.Admin {
		credentialList, err = client.ListClusterAdminCredentials(ctx, resourceID.ResourceGroupName, resourceID.ResourceName)
	} else {
		credentialList, err = client.ListClusterUserCredentials(ctx, resourceID.ResourceGroupName, resourceID.ResourceName)
	}
	if err != nil {
		return nil, fmt.Errorf("getting user credentials: %w", err)
	}

	if credentialList.Kubeconfigs == nil || len(*credentialList.Kubeconfigs) < 1 {
		return nil, ErrNoKubeconfigs
	}

	config := *(*credentialList.Kubeconfigs)[0].Value
	kubeCfg, err := clientcmd.Load(config)
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %w", err)
	}

	return kubeCfg, err
}

func mapAzureEnvironment(env AzureEnvironment) string {
	switch env {
	case AzureEnvironmentPublicCloud:
		return "AzurePublicCloud"
	case AzureEnvironmentChinaCloud:
		return "AzureChinaCloud"
	case AzureEnvironmentUSGovCloud:
		return "AzureUSGovernmentCloud"
	case AzureEnvironmentStackCloud:
		return "AzureStackCloud"
	default:
		return ""
	}
}
