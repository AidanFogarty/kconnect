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

package app

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/fidelity/kconnect/api/v1alpha1"
	"github.com/fidelity/kconnect/pkg/history"
	"github.com/fidelity/kconnect/pkg/k8s/kubeconfig"
	"github.com/fidelity/kconnect/pkg/printer"
)

type HistoryQueryInput struct {
	HistoryConfig
	KubernetesConfig
	HistoryQueryConfig

	ClusterProvider  *string `json:"cluster-provider,omitempty"`
	IdentityProvider *string `json:"identity-provider,omitempty"`

	ProviderID *string `json:"provider-id,omitempty"`
	HistoryID  *string `json:"id,omitempty"`
	Alias      *string `json:"alias,omitempty"`

	Flags map[string]string `json:"flags,omitempty"`

	Output *printer.OutputPrinter `json:"output,omitempty"`
}

func (a *App) QueryHistory(ctx context.Context, input *HistoryQueryInput) error {
	zap.S().Debug("querying history")

	list, err := a.historyStore.GetAllSortedByLastUsed()
	if err != nil {
		return fmt.Errorf("getting history entries: %w", err)
	}

	filterSpec := createFilter(input.Filter)

	if err := history.FilterHistory(list, filterSpec); err != nil {
		return fmt.Errorf("filtering history list: %w", err)
	}

	objPrinter, err := printer.New(*input.Output)
	if err != nil {
		return fmt.Errorf("getting printer for output %s: %w", *input.Output, err)
	}

	if *input.Output == printer.OutputPrinterTable {

		currentContexID, err := a.getCurrentContextID(input.Kubeconfig)
		if err != nil {
			zap.S().Warnf("Error getting current context ID: %s", err)
		}
		return objPrinter.Print(list.ToTable(currentContexID), os.Stdout)
	}

	return objPrinter.Print(list, os.Stdout)
}

func (a *App) getCurrentContextID(kubecfg string) (string, error) {
	currentContext, err := kubeconfig.GetCurrentContext(kubecfg)
	if err != nil {
		return "", err
	}
	if currentContext == nil || currentContext.Extensions == nil {
		return "", nil
	}
	currentContextHistoryReference, err := v1alpha1.GetHistoryReferenceFromContext(currentContext)
	if err != nil {
		return "", err
	}
	currentContexID := currentContextHistoryReference.EntryID
	return currentContexID, nil
}
