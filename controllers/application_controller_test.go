package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"bitbucket.org/accezz-io/sac-operator/model"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/util/rand"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("application controller", func() {

	const (
		applicationNamePrefix = "test-operator-application"
		ApplicationNamespace  = "default"

		timeout  = time.Second * 30
		interval = time.Second * 1
	)

	PContext("for invalid application", func() {
		application := &accessv1.Application{}
		It("Should return an error and do nothing", func() {
			fmt.Fprintf(GinkgoWriter, "looking for application in sac %s\n", application.Name)
		})
	})

	Context("for any valid application kind in k8s", Ordered, func() {
		applicationDTO := &dto.ApplicationDTO{}
		serviceName := "test"
		serviceNamespace := "test"
		servicePort := "443"
		serviceSchema := "https"
		accessPolicies := []string{"rnd-policy", "product-policy"}
		spec := accessv1.NewApplicationSpecBuilder().
			Service(accessv1.Service{Name: serviceName, Namespace: serviceNamespace, Port: servicePort, Schema: serviceSchema}).
			SiteName("operator-site-1").
			AccessPolicies(accessPolicies).
			ActivityPolicies([]string{}).
			IsVisible(true).
			IsNotificationEnabled(true).
			Enabled(true).
			ApplicationType(model.HTTP).
			Build()
		application := &accessv1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       "access.secure-access-cloud.symantec.com/v1",
				APIVersion: "Application",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", applicationNamePrefix, rand.String(4)),
				Namespace: ApplicationNamespace,
			},
			Spec:   spec,
			Status: accessv1.ApplicationStatus{},
		}
		applicationLookupKey := types.NamespacedName{Name: application.Name, Namespace: application.Namespace}
		AfterEach(func() {
			if applicationDTO.ID != "" {
				fmt.Fprintf(GinkgoWriter, "deleting application %s in sac\n", applicationDTO.Name)
				err := sacClient.DeleteApplication(applicationDTO.ID)
				Expect(err).NotTo(HaveOccurred())
			}

			if useCluster {
				fmt.Fprintf(GinkgoWriter, "deleting application %s in k8s\n", applicationDTO.Name)
				err := k8sClient.Delete(context.Background(), application)
				Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			}
		})
		It("Should CREATE UPDATE and DELETE the application", func() {
			By("Creating the application", func() {
				Expect(k8sClient.Create(ctx, application)).Should(Succeed())
				Eventually(func(g Gomega) {
					applicationFromSAC, err := sacClient.FindApplicationByName(application.Name)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(applicationFromSAC.ID).NotTo(BeEmpty())
					g.Expect(validateApplicationInSac(applicationFromSAC, application)).Should(BeNil())
					g.Expect(validateApplicationInK8S(applicationFromSAC, applicationLookupKey)).Should(BeNil())
				}).WithPolling(1 * time.Second).WithTimeout(30 * time.Second).
					Should(Succeed())
			})
			By("Updating the application", func() {
				Expect(k8sClient.Get(ctx, applicationLookupKey, application)).Should(Succeed())
				application.Spec.SiteName = "operator-site-2"
				application.Spec.Service.Name = "test-2"
				Expect(k8sClient.Update(ctx, application)).Should(Succeed())
				Eventually(func(g Gomega) {
					applicationFromSAC, err := sacClient.FindApplicationByName(application.Name)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(applicationFromSAC.ID).NotTo(BeEmpty())
					g.Expect(validateApplicationInSac(applicationFromSAC, application)).Should(BeNil())
					g.Expect(validateApplicationInK8S(applicationFromSAC, applicationLookupKey)).Should(BeNil())
				}).WithPolling(1 * time.Second).WithTimeout(30 * time.Second).
					Should(Succeed())
			})
			By("Deleting the application", func() {
				Expect(k8sClient.Get(ctx, applicationLookupKey, application)).Should(Succeed())
				Expect(k8sClient.Delete(ctx, application)).Should(Succeed())
				Eventually(func(g Gomega) {
					_, err := sacClient.FindApplicationByName(application.Name)
					g.Expect(errors.Is(err, sac.ErrorNotFound)).Should(BeTrue())
				}).WithPolling(1 * time.Second).WithTimeout(30 * time.Second).
					Should(Succeed())
			})

		})
	})

})

func getApplicationFromSAC(applicationName string) (*dto.ApplicationDTO, error) {
	fmt.Fprintf(GinkgoWriter, "looking for application by name %s in sac\n", applicationName)
	count := 5

	for {
		if count == 0 {
			return nil, fmt.Errorf("could not find application %s", applicationName)
		}
		applicationDTO, err := sacClient.FindApplicationByName(applicationName)
		if err != nil {
			time.Sleep(5 * time.Second)
			count--
			continue
		}
		return applicationDTO, nil
	}
}

func validateApplicationInK8S(applicationDTO *dto.ApplicationDTO, applicationLookupKey types.NamespacedName) error {
	sleep := time.Second * 5
	counter := 5

	application := &accessv1.Application{}

	for {
		if counter == 0 {
			return fmt.Errorf("could not validate application in k8s %+v", application)
		}
		err := k8sClient.Get(ctx, applicationLookupKey, application)
		if err != nil {
			return fmt.Errorf("could not find application in k8s %s", applicationLookupKey)
		}
		fmt.Fprintf(GinkgoWriter, "looking for application id in sac %s vs crd status %s \n", applicationDTO.ID, application.Status.Id)
		if application.Status.Id == applicationDTO.ID {
			return nil
		}
		time.Sleep(sleep)
		counter--
	}

}

func validateApplicationInSac(applicationDTO *dto.ApplicationDTO, application *accessv1.Application) error {

	// checking that application was created as expected in sac
	if applicationDTO.Enabled != application.Spec.Enabled {
		return fmt.Errorf("applicationDTO.Enabled application.Spec.Enabled shouldd be equale got %v %v", applicationDTO.Enabled, application.Spec.Enabled)
	}
	if applicationDTO.IsVisible != application.Spec.IsVisible {
		return fmt.Errorf("applicationDTO.IsVisible application.Spec.IsVisible shouldd be equale got %v %v", applicationDTO.IsVisible, application.Spec.IsVisible)
	}
	if applicationDTO.IsNotificationEnabled != application.Spec.IsNotificationEnabled {
		return fmt.Errorf("applicationDTO.IsNotificationEnabled application.Spec.IsNotificationEnabled shouldd be equale got %v %v", applicationDTO.IsNotificationEnabled, application.Spec.IsNotificationEnabled)
	}
	expectedInternalAddress := fmt.Sprintf("%s://%s.%s:%s",
		application.Spec.Service.Schema, application.Spec.Service.Name,
		application.Spec.Service.Namespace,
		application.Spec.Service.Port)
	if applicationDTO.ConnectionSettings.InternalAddress != expectedInternalAddress {
		return fmt.Errorf("applicationDTO.ConnectionSettings.InternalAddress expectedInternalAddress shouldd be equale got %s %s",
			applicationDTO.ConnectionSettings.InternalAddress, expectedInternalAddress,
		)
	}
	// checking binding to site
	site, err := sacClient.FindSiteByName(application.Spec.SiteName)
	if err != nil {
		return fmt.Errorf("could not find site by its name %s", application.Spec.SiteName)
	}
	if !findInSlice(site.ApplicationIDs, applicationDTO.ID) {
		return fmt.Errorf("application is not part of site")

	}

	// checking binding to policies
	accessPolicies, err := sacClient.FindPoliciesByNames(application.Spec.AccessPoliciesNames)
	if err != nil {
		return fmt.Errorf("failed to find access policy in sac %s\n", application.Spec.AccessPoliciesNames)
	}
	activityPolicies, err := sacClient.FindPoliciesByNames(application.Spec.ActivityPoliciesNames)
	if err != nil {
		return fmt.Errorf("failed to find activity policy in sac %s\n", application.Spec.ActivityPoliciesNames)
	}
	policies := append(accessPolicies, activityPolicies...)
	foundCounter := 0
	for i := range policies {
		for i2 := range policies[i].Applications {
			if policies[i].Applications[i2].ID == applicationDTO.ID {
				foundCounter++
				break
			}
		}
	}
	if foundCounter != len(policies) {
		return fmt.Errorf("application is not part of required policies %+v\n", policies)
	}

	return nil

}

func findInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
