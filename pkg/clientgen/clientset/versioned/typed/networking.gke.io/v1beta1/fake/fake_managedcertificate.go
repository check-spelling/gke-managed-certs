/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/GoogleCloudPlatform/gke-managed-certs/pkg/apis/networking.gke.io/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeManagedCertificates implements ManagedCertificateInterface
type FakeManagedCertificates struct {
	Fake *FakeNetworkingV1beta1
	ns   string
}

var managedcertificatesResource = schema.GroupVersionResource{Group: "networking.gke.io", Version: "v1beta1", Resource: "managedcertificates"}

var managedcertificatesKind = schema.GroupVersionKind{Group: "networking.gke.io", Version: "v1beta1", Kind: "ManagedCertificate"}

// Get takes name of the managedCertificate, and returns the corresponding managedCertificate object, and an error if there is any.
func (c *FakeManagedCertificates) Get(name string, options v1.GetOptions) (result *v1beta1.ManagedCertificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(managedcertificatesResource, c.ns, name), &v1beta1.ManagedCertificate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ManagedCertificate), err
}

// List takes label and field selectors, and returns the list of ManagedCertificates that match those selectors.
func (c *FakeManagedCertificates) List(opts v1.ListOptions) (result *v1beta1.ManagedCertificateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(managedcertificatesResource, managedcertificatesKind, c.ns, opts), &v1beta1.ManagedCertificateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ManagedCertificateList{ListMeta: obj.(*v1beta1.ManagedCertificateList).ListMeta}
	for _, item := range obj.(*v1beta1.ManagedCertificateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested managedCertificates.
func (c *FakeManagedCertificates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(managedcertificatesResource, c.ns, opts))

}

// Create takes the representation of a managedCertificate and creates it.  Returns the server's representation of the managedCertificate, and an error, if there is any.
func (c *FakeManagedCertificates) Create(managedCertificate *v1beta1.ManagedCertificate) (result *v1beta1.ManagedCertificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(managedcertificatesResource, c.ns, managedCertificate), &v1beta1.ManagedCertificate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ManagedCertificate), err
}

// Update takes the representation of a managedCertificate and updates it. Returns the server's representation of the managedCertificate, and an error, if there is any.
func (c *FakeManagedCertificates) Update(managedCertificate *v1beta1.ManagedCertificate) (result *v1beta1.ManagedCertificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(managedcertificatesResource, c.ns, managedCertificate), &v1beta1.ManagedCertificate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ManagedCertificate), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeManagedCertificates) UpdateStatus(managedCertificate *v1beta1.ManagedCertificate) (*v1beta1.ManagedCertificate, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(managedcertificatesResource, "status", c.ns, managedCertificate), &v1beta1.ManagedCertificate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ManagedCertificate), err
}

// Delete takes name of the managedCertificate and deletes it. Returns an error if one occurs.
func (c *FakeManagedCertificates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(managedcertificatesResource, c.ns, name), &v1beta1.ManagedCertificate{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeManagedCertificates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(managedcertificatesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.ManagedCertificateList{})
	return err
}

// Patch applies the patch and returns the patched managedCertificate.
func (c *FakeManagedCertificates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ManagedCertificate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(managedcertificatesResource, c.ns, name, pt, data, subresources...), &v1beta1.ManagedCertificate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.ManagedCertificate), err
}
