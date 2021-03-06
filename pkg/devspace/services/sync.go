package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/devspace/sync"
	"github.com/devspace-cloud/devspace/pkg/util/log"
)

// StartSync starts the syncing functionality
func StartSync(client *kubernetes.Clientset, verboseSync bool, log log.Logger) ([]*sync.SyncConfig, error) {
	config := configutil.GetConfig()
	if config.Dev.Sync == nil {
		return []*sync.SyncConfig{}, nil
	}

	syncConfigs := make([]*sync.SyncConfig, 0, len(*config.Dev.Sync))
	for _, syncPath := range *config.Dev.Sync {
		absLocalPath, err := filepath.Abs(*syncPath.LocalSubPath)
		if err != nil {
			return nil, fmt.Errorf("Unable to resolve localSubPath %s: %v", *syncPath.LocalSubPath, err)
		}

		var labelSelector map[string]*string
		namespace, err := configutil.GetDefaultNamespace(config)
		if err != nil {
			return nil, err
		}

		containerName := ""

		if syncPath.Selector != nil {
			selector, err := configutil.GetSelector(*syncPath.Selector)
			if err != nil {
				log.Fatalf("Error resolving selector name: %v", err)
			}

			labelSelector = *selector.LabelSelector
			if selector.Namespace != nil && *selector.Namespace != "" {
				namespace = *selector.Namespace
			}

			if selector.ContainerName != nil && *selector.ContainerName != "" {
				containerName = *selector.ContainerName
			}
		} else {
			labelSelector = *syncPath.LabelSelector
			if syncPath.Namespace != nil && *syncPath.Namespace != "" {
				namespace = *syncPath.Namespace
			}

			if syncPath.ContainerName != nil && *syncPath.ContainerName != "" {
				containerName = *syncPath.ContainerName
			}
		}

		labels := make([]string, 0, len(labelSelector)-1)
		for key, value := range labelSelector {
			labels = append(labels, key+"="+*value)
		}

		log.StartWait("Sync: Waiting for pods...")
		pod, err := kubectl.GetNewestRunningPod(client, strings.Join(labels, ", "), namespace, time.Second*120)
		log.StopWait()
		if err != nil {
			return nil, fmt.Errorf("Unable to list devspace pods: %v", err)
		} else if pod != nil {
			if len(pod.Spec.Containers) == 0 {
				log.Warnf("Cannot start sync on pod, because selected pod %s/%s has no containers", pod.Namespace, pod.Name)
				continue
			}

			container := &pod.Spec.Containers[0]
			if containerName != "" {
				found := false

				for _, c := range pod.Spec.Containers {
					if c.Name == containerName {
						container = &c
						found = true
						break
					}
				}

				if found == false {
					log.Warnf("Couldn't start sync, because container %s wasn't found in pod %s/%s", containerName, pod.Namespace, pod.Name)
					continue
				}
			}

			syncConfig := &sync.SyncConfig{
				Kubectl:   client,
				Pod:       pod,
				Container: container,
				WatchPath: absLocalPath,
				DestPath:  *syncPath.ContainerPath,
				Verbose:   verboseSync,
			}

			if syncPath.ExcludePaths != nil {
				syncConfig.ExcludePaths = *syncPath.ExcludePaths
			}

			if syncPath.DownloadExcludePaths != nil {
				syncConfig.DownloadExcludePaths = *syncPath.DownloadExcludePaths
			}

			if syncPath.UploadExcludePaths != nil {
				syncConfig.UploadExcludePaths = *syncPath.UploadExcludePaths
			}

			if syncPath.BandwidthLimits != nil {
				if syncPath.BandwidthLimits.Download != nil {
					syncConfig.DownstreamLimit = *syncPath.BandwidthLimits.Download * 1024
				}

				if syncPath.BandwidthLimits.Upload != nil {
					syncConfig.UpstreamLimit = *syncPath.BandwidthLimits.Upload * 1024
				}
			}

			err = syncConfig.Start()
			if err != nil {
				log.Fatalf("Sync error: %s", err.Error())
			}

			log.Donef("Sync started on %s <-> %s (Pod: %s/%s)", absLocalPath, *syncPath.ContainerPath, pod.Namespace, pod.Name)
			syncConfigs = append(syncConfigs, syncConfig)
		}
	}

	return syncConfigs, nil
}
