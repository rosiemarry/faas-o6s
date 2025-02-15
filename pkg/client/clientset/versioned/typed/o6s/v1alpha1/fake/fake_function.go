/*
Copyright 2017 OpenFaaS Project

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/
package fake

import (
	v1alpha1 "github.com/openfaas-incubator/faas-o6s/pkg/apis/o6sio/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFunctions implements FunctionInterface
type FakeFunctions struct {
	Fake *FakeO6sV1alpha1
	ns   string
}

var functionsResource = schema.GroupVersionResource{Group: "o6s.io", Version: "v1alpha1", Resource: "functions"}

var functionsKind = schema.GroupVersionKind{Group: "o6s.io", Version: "v1alpha1", Kind: "Function"}

// Get takes name of the function, and returns the corresponding function object, and an error if there is any.
func (c *FakeFunctions) Get(name string, options v1.GetOptions) (result *v1alpha1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(functionsResource, c.ns, name), &v1alpha1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Function), err
}

// List takes label and field selectors, and returns the list of Functions that match those selectors.
func (c *FakeFunctions) List(opts v1.ListOptions) (result *v1alpha1.FunctionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(functionsResource, functionsKind, c.ns, opts), &v1alpha1.FunctionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FunctionList{}
	for _, item := range obj.(*v1alpha1.FunctionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested functions.
func (c *FakeFunctions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(functionsResource, c.ns, opts))

}

// Create takes the representation of a function and creates it.  Returns the server's representation of the function, and an error, if there is any.
func (c *FakeFunctions) Create(function *v1alpha1.Function) (result *v1alpha1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(functionsResource, c.ns, function), &v1alpha1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Function), err
}

// Update takes the representation of a function and updates it. Returns the server's representation of the function, and an error, if there is any.
func (c *FakeFunctions) Update(function *v1alpha1.Function) (result *v1alpha1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(functionsResource, c.ns, function), &v1alpha1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Function), err
}

// Delete takes name of the function and deletes it. Returns an error if one occurs.
func (c *FakeFunctions) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(functionsResource, c.ns, name), &v1alpha1.Function{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFunctions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(functionsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.FunctionList{})
	return err
}

// Patch applies the patch and returns the patched function.
func (c *FakeFunctions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(functionsResource, c.ns, name, data, subresources...), &v1alpha1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Function), err
}
