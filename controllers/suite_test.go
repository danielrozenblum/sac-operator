/*
Copyright 2021.

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

package controllers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/accezz-io/sac-operator/service"

	"bitbucket.org/accezz-io/sac-operator/controllers/access"

	"bitbucket.org/accezz-io/sac-operator/service/sac"

	"bitbucket.org/accezz-io/sac-operator/controllers/access/converter"

	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg        *rest.Config
	k8sClient  client.Client
	sacClient  sac.SecureAccessCloudClient
	testEnv    *envtest.Environment
	ctx        context.Context
	cancel     context.CancelFunc
	useCluster = false
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "application controller suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		BinaryAssetsDirectory: filepath.Join("..", "testbin", "bin"),
		UseExistingCluster:    &useCluster,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = accessv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	sacClientID, sacClientSecret, sacTenantDomain := os.Getenv("SAC_CLIENT_ID"), os.Getenv("SAC_CLIENT_SECRET"), os.Getenv("SAC_TENANT_DOMAIN")
	if sacClientID == "" || sacClientSecret == "" || sacTenantDomain == "" {
		Fail("Missing environment variable required for tests. SAC_CLIENT_ID, SAC_CLIENT_SECRET and SAC_TENANT_DOMAIN must all be set.")
	}

	secureAccessCloudSettings := &sac.SecureAccessCloudSettings{
		ClientID:     sacClientID,
		ClientSecret: sacClientSecret,
		TenantDomain: sacTenantDomain,
	}
	sacClient = sac.NewSecureAccessCloudClientImpl(secureAccessCloudSettings)

	err = (&access.SiteReconcile{
		Client:                  k8sManager.GetClient(),
		Scheme:                  k8sManager.GetScheme(),
		SecureAccessCloudClient: sacClient,
		SiteConverter:           converter.NewSiteConverter(),
		Log:                     ctrl.Log.WithName("test-site-reconcile"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	applicationReconcilerLogger := ctrl.Log.WithName("application-reconcile")
	err = (&access.ApplicationReconciler{
		Client:               k8sManager.GetClient(),
		Scheme:               k8sManager.GetScheme(),
		ApplicationService:   service.NewApplicationServiceImpl(sacClient, applicationReconcilerLogger),
		ApplicationConverter: converter.NewApplicationTypeConverter(),
		Log:                  ctrl.Log.WithName("test-application-reconcile"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
