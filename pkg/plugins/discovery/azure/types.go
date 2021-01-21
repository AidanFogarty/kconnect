/*
Copyright 2021 The kconnect Authors.

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

// AzureEnvironment is a type that represents an Azure environment
type AzureEnvironment string

var (
	// AzureEnvironmentPublicCloud is the general Azure public cloud
	AzureEnvironmentPublicCloud = AzureEnvironment("public")
	// AzureEnvironmentChinaCloud is the public Azure cloud in China
	AzureEnvironmentChinaCloud = AzureEnvironment("china")
	// AzureEnvironmentUSGovCloud is the US Government specific Azure cloud
	AzureEnvironmentUSGovCloud = AzureEnvironment("usgov")
	// AzureEnvironmentStackCloud is the Azure stack cloud
	AzureEnvironmentStackCloud = AzureEnvironment("stack")
)

// LoginType is a type that denotes the type of user login
type LoginType string

var (
	// LoginTypeDeviceCode is for using a device code to login
	LoginTypeDeviceCode = LoginType("devicecode")
	// LoginTypeServicePrincipal is for using a service principal to login
	LoginTypeServicePrincipal = LoginType("spn")
	// LoginTypeResourceOwnerPassword is for using the resource owner password type
	LoginTypeResourceOwnerPassword = LoginType("ropc")
	// LoginTypeManagedServiceIdentity is for using the managed service identity to login
	LoginTypeManagedServiceIdentity = LoginType("msi")
)
