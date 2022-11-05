package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	anyninesv1 "github.com/mmertdogann/dummy-operator/api/v1"
)

var _ = Describe("Dummy controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		DummyName      = "test-dummy"
		DummyNamespace = "default"
		timeout        = time.Second * 10
		interval       = time.Millisecond * 250
	)

	Context("When creating Dummy", func() {
		It("Should check CR specs", func() {
			By("By creating a new Dummy")
			ctx := context.Background()
			dummy := &anyninesv1.Dummy{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "anynines.interview.com/v1",
					Kind:       "Dummy",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      DummyName,
					Namespace: DummyNamespace,
				},
				Spec: anyninesv1.DummySpec{
					Message: "I'm just a dummy",
				},
			}
			Expect(k8sClient.Create(ctx, dummy)).Should(Succeed())

			dummyLookupKey := types.NamespacedName{Name: DummyName, Namespace: DummyNamespace}
			createdDummy := &anyninesv1.Dummy{}

			// We'll need to retry getting this newly created Dummy, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, dummyLookupKey, createdDummy)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdDummy.Spec.Message).Should(Equal("I'm just a dummy"))
		})
	})
})
