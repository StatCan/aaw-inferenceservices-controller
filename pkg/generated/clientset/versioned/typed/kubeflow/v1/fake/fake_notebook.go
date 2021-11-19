/*
The MIT License (MIT)

Copyright © 2020 Her Majesty the Queen in Right of Canada, as represented by the Minister of Statistics Canada

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	kubeflowv1 "github.com/StatCan/inferenceservices-controller/pkg/apis/kubeflow/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNotebooks implements NotebookInterface
type FakeNotebooks struct {
	Fake *FakeKubeflowV1
	ns   string
}

var notebooksResource = schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1", Resource: "notebooks"}

var notebooksKind = schema.GroupVersionKind{Group: "kubeflow.org", Version: "v1", Kind: "Notebook"}

// Get takes name of the notebook, and returns the corresponding notebook object, and an error if there is any.
func (c *FakeNotebooks) Get(ctx context.Context, name string, options v1.GetOptions) (result *kubeflowv1.Notebook, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(notebooksResource, c.ns, name), &kubeflowv1.Notebook{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubeflowv1.Notebook), err
}

// List takes label and field selectors, and returns the list of Notebooks that match those selectors.
func (c *FakeNotebooks) List(ctx context.Context, opts v1.ListOptions) (result *kubeflowv1.NotebookList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(notebooksResource, notebooksKind, c.ns, opts), &kubeflowv1.NotebookList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kubeflowv1.NotebookList{ListMeta: obj.(*kubeflowv1.NotebookList).ListMeta}
	for _, item := range obj.(*kubeflowv1.NotebookList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested notebooks.
func (c *FakeNotebooks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(notebooksResource, c.ns, opts))

}

// Create takes the representation of a notebook and creates it.  Returns the server's representation of the notebook, and an error, if there is any.
func (c *FakeNotebooks) Create(ctx context.Context, notebook *kubeflowv1.Notebook, opts v1.CreateOptions) (result *kubeflowv1.Notebook, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(notebooksResource, c.ns, notebook), &kubeflowv1.Notebook{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubeflowv1.Notebook), err
}

// Update takes the representation of a notebook and updates it. Returns the server's representation of the notebook, and an error, if there is any.
func (c *FakeNotebooks) Update(ctx context.Context, notebook *kubeflowv1.Notebook, opts v1.UpdateOptions) (result *kubeflowv1.Notebook, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(notebooksResource, c.ns, notebook), &kubeflowv1.Notebook{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubeflowv1.Notebook), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNotebooks) UpdateStatus(ctx context.Context, notebook *kubeflowv1.Notebook, opts v1.UpdateOptions) (*kubeflowv1.Notebook, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(notebooksResource, "status", c.ns, notebook), &kubeflowv1.Notebook{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubeflowv1.Notebook), err
}

// Delete takes name of the notebook and deletes it. Returns an error if one occurs.
func (c *FakeNotebooks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(notebooksResource, c.ns, name), &kubeflowv1.Notebook{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNotebooks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(notebooksResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &kubeflowv1.NotebookList{})
	return err
}

// Patch applies the patch and returns the patched notebook.
func (c *FakeNotebooks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *kubeflowv1.Notebook, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(notebooksResource, c.ns, name, pt, data, subresources...), &kubeflowv1.Notebook{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubeflowv1.Notebook), err
}
