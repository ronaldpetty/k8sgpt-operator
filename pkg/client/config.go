package client

import (
	"context"
	"fmt"
	"strconv"

	rpc "buf.build/gen/go/ronaldpetty/ronk8sgpt/grpc/go/schema/v1/schemav1grpc"
	schemav1 "buf.build/gen/go/ronaldpetty/ronk8sgpt/protocolbuffers/go/schema/v1"
	"github.com/k8sgpt-ai/k8sgpt-operator/api/v1alpha1"
)

func (c *Client) AddConfig(config *v1alpha1.K8sGPT) error {
	client := rpc.NewServerServiceClient(c.conn)
	req := &schemav1.AddConfigRequest{}
	// If multiple caches are configured we pick S3
	// which emulates the behaviour of K8sGPT cli
	if config.Spec.RemoteCache != nil {
		if config.Spec.RemoteCache.S3 != nil {
			req.Cache = &schemav1.Cache{
				CacheType: &schemav1.Cache_S3Cache{
					S3Cache: &schemav1.S3Cache{
						BucketName: config.Spec.RemoteCache.S3.BucketName,
						Region:     config.Spec.RemoteCache.S3.Region,
					},
				},
			}
		} else if config.Spec.RemoteCache.Azure != nil {
			req.Cache = &schemav1.Cache{
				CacheType: &schemav1.Cache_AzureCache{
					AzureCache: &schemav1.AzureCache{
						StorageAccount: config.Spec.RemoteCache.Azure.StorageAccount,
						ContainerName:  config.Spec.RemoteCache.Azure.ContainerName,
					},
				},
			}
		} else if config.Spec.RemoteCache.GCS != nil {
			req.Cache = &schemav1.Cache{
				CacheType: &schemav1.Cache_GcsCache{
					GcsCache: &schemav1.GCSCache{
						BucketName: config.Spec.RemoteCache.GCS.BucketName,
						Region:     config.Spec.RemoteCache.GCS.Region,
						ProjectId:  config.Spec.RemoteCache.GCS.ProjectId,
					},
				},
			}
		}
	}
	if config.Spec.CustomAnalyzers != nil {
		for _, customAnalyzer := range config.Spec.CustomAnalyzers {
			req.CustomAnalyzers = append(req.CustomAnalyzers, &schemav1.CustomAnalyzer{
				Name: customAnalyzer.Name,
				Connection: &schemav1.Connection{
					Url:  customAnalyzer.Connection.Url,
					Port: strconv.Itoa(customAnalyzer.Connection.Port),
				},
			})
		}
	}
	_, err := client.AddConfig(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to call AddConfig RPC: %v", err)
	}

	return nil
}

func (c *Client) RemoveConfig(config *v1alpha1.K8sGPT) error {
	client := rpc.NewServerServiceClient(c.conn)
	req := &schemav1.RemoveConfigRequest{
		Cache:           &schemav1.Cache{},
		Integrations:    nil,
		CustomAnalyzers: make([]*schemav1.CustomAnalyzer, 0),
	}

	_, err := client.RemoveConfig(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to call RemoveConfig RPC: %v", err)
	}

	return nil
}
