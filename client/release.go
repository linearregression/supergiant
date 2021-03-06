package client

import (
	"path"

	"github.com/supergiant/supergiant/common"
)

type Release common.Release

type ReleaseCollection struct {
	client *Client

	App       *AppResource
	Component *ComponentResource
}

type ReleaseResource struct {
	collection *ReleaseCollection
	*Release
}

type ReleaseList struct {
	Items []*ReleaseResource
}

func (c *ReleaseCollection) path() string {
	return path.Join("apps", common.StringID(c.App.Name), "components", common.StringID(c.Component.Name), "releases")
}

func (r *ReleaseResource) path() string {
	return path.Join(r.collection.path(), common.StringID(r.Timestamp))
}

// Collection-level
//==============================================================================
func (c *ReleaseCollection) New(m *Release) *ReleaseResource {
	return &ReleaseResource{c, m}
}

func (c *ReleaseCollection) List() (*ReleaseList, error) {
	list := new(ReleaseList)
	if err := c.client.Get(c.path(), list); err != nil {
		return nil, err
	}
	// see TODO in instance.go
	for _, release := range list.Items {
		release.collection = c
	}
	return list, nil
}

func (c *ReleaseCollection) Create(m *Release) (*ReleaseResource, error) {
	r := c.New(m)
	if err := c.client.Post(c.path(), m, r.Release); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *ReleaseCollection) MergeCreate(m *Release) (*ReleaseResource, error) {
	r := c.New(m)
	if err := c.client.Patch(c.path(), m, r.Release); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *ReleaseCollection) Get(timestamp common.ID) (*ReleaseResource, error) {
	m := &Release{
		Timestamp: timestamp,
	}
	r := c.New(m)
	if err := c.client.Get(r.path(), r.Release); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *ReleaseCollection) Update(timestamp common.ID, m *Release) (*ReleaseResource, error) {
	r := c.New(&Release{Timestamp: timestamp})
	if err := c.client.Put(r.path(), m, r.Release); err != nil {
		return nil, err
	}
	return r, nil
}

// Resource-level
//==============================================================================
func (r *ReleaseResource) Save() (*ReleaseResource, error) {
	return r.collection.Update(r.Timestamp, r.Release)
}

func (r *ReleaseResource) Delete() error {
	return r.collection.client.Delete(r.path())
}

// Relations
func (r *ReleaseResource) Instances() *InstanceCollection {
	return &InstanceCollection{
		client:    r.collection.client,
		App:       r.collection.App,
		Component: r.collection.Component,
		Release:   r,
	}
}
