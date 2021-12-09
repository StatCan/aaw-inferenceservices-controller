package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/StatCan/inferenceservices-controller/pkg/signals"
	servingv1alpha2 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	servingclientset "github.com/kserve/kserve/pkg/client/clientset/versioned"
	servinginformers "github.com/kserve/kserve/pkg/client/informers/externalversions"
	servingv1alpha2listers "github.com/kserve/kserve/pkg/client/listers/serving/v1beta1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var configMapName string
var configMapKey string
var aliasServiceName string

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Configure DNS resources",
	Long:  `Configure DNS resources for inferences services`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		// Setup signals so we can shutdown cleanly
		stopCh := signals.SetupSignalHandler()

		// Create Kubernetes config
		cfg, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
		if err != nil {
			klog.Fatalf("error building kubeconfig: %v", err)
		}

		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		}

		servingClient, err := servingclientset.NewForConfig(cfg)
		if err != nil {
			klog.Fatalf("error building Serving client: %v", err)
		}

		// Setup informers
		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Minute*5)
		servingInformerFactory := servinginformers.NewSharedInformerFactory(servingClient, time.Minute*5)
		inferenceServicesInformer := servingInformerFactory.Serving().V1beta1().InferenceServices()

		update := func() {
			conf, err := generateDNS(inferenceServicesInformer.Lister())
			if err != nil {
				klog.Errorf("Failed to generate initial DNS config: %v", err)
			}

			if conf == "" {
				klog.Info("No DNS entries...")
				return
			}

			components := strings.Split(configMapName, "/")
			existingConfigMap, err := kubeClient.CoreV1().ConfigMaps(components[0]).Get(ctx, components[1], metav1.GetOptions{})
			if err != nil {
				klog.Errorf("Failed to load existing configmap: %v", err)
			}

			update := false
			updated := existingConfigMap.DeepCopy()
			stubConf := "%s\n"

			if existing, ok := existingConfigMap.Data[configMapKey]; ok {
				if existing != conf {
					update = true

					updated.Data[configMapKey] = fmt.Sprintf(stubConf, conf)
				}
			} else {
				update = true
				updated.Data[configMapKey] = fmt.Sprintf(stubConf, conf)
			}

			if update {
				_, err := kubeClient.CoreV1().ConfigMaps(components[0]).Update(ctx, updated, metav1.UpdateOptions{})
				if err != nil {
					klog.Errorf("Failed to update ConfigMap: %v", err)
				}
			}
		}

		inferenceServicesInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(new interface{}) {
				if inferenceServicesInformer.Informer().HasSynced() {
					update()
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldIS := old.(*servingv1alpha2.InferenceService)
				newIS := new.(*servingv1alpha2.InferenceService)

				if oldIS.ResourceVersion == newIS.ResourceVersion {
					return
				}

				if inferenceServicesInformer.Informer().HasSynced() {
					update()
				}
			},
			DeleteFunc: func(old interface{}) {
				if inferenceServicesInformer.Informer().HasSynced() {
					update()
				}
			},
		})

		// Start informers
		kubeInformerFactory.Start(stopCh)
		servingInformerFactory.Start(stopCh)

		// Wait for caches
		klog.Info("Waiting for informer caches to sync")
		if ok := cache.WaitForCacheSync(stopCh, inferenceServicesInformer.Informer().HasSynced); !ok {
			klog.Fatalf("failed to wait for caches to sync")
		}
		klog.Info("Informer caches synched")

		update()

		// Run the controller
		<-stopCh
	},
}

func generateDNS(inferenceServiceLister servingv1alpha2listers.InferenceServiceLister) (string, error) {
	klog.Info("Generating DNS...")

	inferenceServices, err := inferenceServiceLister.List(labels.Everything())
	if err != nil {
		return "", err
	}

	entries := []string{}

	for _, inferenceService := range inferenceServices {
		if inferenceService.Status.IsReady() {
			entries = append(entries, fmt.Sprintf("rewrite name %s %s", strings.ReplaceAll(inferenceService.Status.URL.Host, "http://", ""), aliasServiceName))
		}
	}

	return strings.Join(entries, "\n"), nil
}

func init() {
	dnsCmd.Flags().StringVar(&configMapName, "config-map", "kube-system/coredns-custom", "Configmap to output DNS config")
	dnsCmd.Flags().StringVar(&configMapKey, "key", "kfserving-ingress.override", "Key in the ConfigMap to output to")
	dnsCmd.Flags().StringVar(&aliasServiceName, "alias-service-name", "istio-ingressgateway.istio-system.svc.cluster.local", "Service to rewrite the hostname to")

	rootCmd.AddCommand(dnsCmd)
}
