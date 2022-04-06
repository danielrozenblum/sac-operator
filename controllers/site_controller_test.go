package controllers

import (
	"context"
	"fmt"
	"time"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/rand"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"

	"k8s.io/apimachinery/pkg/types"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getSiteSpec(numberOfConnectors int, connectorNamespace string) accessv1.SiteSpec {
	return accessv1.SiteSpec{
		NumberOfConnectors: numberOfConnectors,
	}
}

var _ = Describe("SiteName controller", func() {

	const (
		SiteNamePrefix            = "test-operator-site"
		SiteNamespace             = "default"
		ConnectorsNamespace       = "default"
		InitialNumberOfConnectors = 0

		timeout  = time.Second * 20
		duration = time.Second * 10
		interval = time.Second * 1
	)

	PContext("When creating new site", func() {
		site := &accessv1.Site{}
		siteDto := &dto.SiteDTO{}
		BeforeEach(func() {
			siteName := fmt.Sprintf("%s-%s", SiteNamePrefix, rand.String(4))
			site = &accessv1.Site{
				TypeMeta: metav1.TypeMeta{
					Kind:       "access.secure-access-cloud.symantec.com/v1",
					APIVersion: "Site",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      siteName,
					Namespace: SiteNamespace,
				},
				Spec:   getSiteSpec(InitialNumberOfConnectors, ConnectorsNamespace),
				Status: accessv1.SiteStatus{},
			}
		})
		AfterEach(func() {
			if siteDto.ID != "" {
				fmt.Fprintf(GinkgoWriter, "deleting site %s in k8s\n", site.Name)
				err := sacClient.DeleteSite(siteDto.ID)
				Expect(err).NotTo(HaveOccurred())
			}
		})
		It("Should create site in sac and update status with SAC site ID", func() {
			Skip("skipping for now")
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, site)).Should(Succeed())
			siteLookupKey := types.NamespacedName{Name: site.Name, Namespace: site.Namespace}
			createdSite := &accessv1.Site{}
			Eventually(func() bool {
				fmt.Fprintf(GinkgoWriter, "looking for site %s in k8s\n", site.Name)
				err := k8sClient.Get(ctx, siteLookupKey, createdSite)
				if err != nil {
					fmt.Fprintf(GinkgoWriter, "got error %s\n", err)
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(createdSite.Spec.NumberOfConnectors).Should(Equal(site.Spec.NumberOfConnectors))
			fmt.Fprintf(GinkgoWriter, "site spec seems to be ok %+v\n", createdSite.Spec)
			Eventually(func() bool {
				var err error
				fmt.Fprintf(GinkgoWriter, "looking for site in sac %s\n", createdSite.Name)
				siteDto, err = sacClient.FindSiteByName(createdSite.Name)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(siteDto.Name).Should(Equal(site.Name))

			numberOfConnectors := 2
			By("increasing number of connectors", func() {
				err := k8sClient.Get(ctx, siteLookupKey, createdSite)
				Expect(err).NotTo(HaveOccurred())
				createdSite.Spec.NumberOfConnectors = numberOfConnectors
				err = k8sClient.Update(ctx, createdSite)
				Expect(err).NotTo(HaveOccurred())
			})
			var connectorsInSac []string
			fmt.Fprintf(GinkgoWriter, "looking for connectors site in sac %s\n", createdSite.Name)
			Eventually(func() bool {
				var err error
				connectorsInSac, err = sacClient.ListConnectorsBySite(createdSite.Name)
				if err != nil {
					return false
				}
				fmt.Fprintf(GinkgoWriter, "connectors in sac %+v\n", connectorsInSac)
				if numberOfConnectors != len(connectorsInSac) {
					return false
				}
				return true
			}, 3*time.Minute, 5*time.Second).Should(BeTrue())
			Expect(createdSite.Spec.NumberOfConnectors).Should(Equal(numberOfConnectors))

			Eventually(func() bool {
				podList := &corev1.PodList{}
				err := k8sClient.List(ctx, podList)
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(GinkgoWriter, "looking for connector pods in cluster %+v\n", podList)
				healthyConnectors := 0
				for _, pod := range podList.Items {
					if pod.Status.Phase != corev1.PodRunning {
						continue
					}
					annotations := pod.GetAnnotations()
					if val, ok := annotations[fmt.Sprintf("%s/%s", connector_deployer.AnnotationPrefix, "site")]; ok {
						if val == createdSite.Name {
							healthyConnectors++
						}
					}
				}
				if healthyConnectors == createdSite.Spec.NumberOfConnectors {
					return true
				}
				return false
			}, 3*time.Minute, 5*time.Second).Should(BeTrue())
		})
	})

})
