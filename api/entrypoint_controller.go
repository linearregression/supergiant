package api

import (
	"net/http"

	"github.com/supergiant/supergiant/core"
)

type EntrypointController struct {
	core *core.Core
}

func (c *EntrypointController) Create(w http.ResponseWriter, r *http.Request) {
	entrypoint := c.core.Entrypoints().New()

	if err := unmarshalBodyInto(w, r, entrypoint); err != nil {
		return
	}

	core.ZeroReadonlyFields(entrypoint)

	err := c.core.Entrypoints().Create(entrypoint)
	if err != nil {
		renderError(w, err, http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, entrypoint)
	if err != nil {
		return
	}
	renderWithStatusCreated(w, body)
}

func (c *EntrypointController) Index(w http.ResponseWriter, r *http.Request) {
	entrypoints, err := c.core.Entrypoints().List()
	if err != nil {
		renderError(w, err, http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, entrypoints)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *EntrypointController) Show(w http.ResponseWriter, r *http.Request) {
	entrypoint, err := loadEntrypoint(c.core, w, r)
	if err != nil {
		return
	}

	body, err := marshalBody(w, entrypoint)
	if err != nil {
		return
	}
	renderWithStatusOK(w, body)
}

func (c *EntrypointController) Update(w http.ResponseWriter, r *http.Request) {
	entrypoint, err := loadEntrypoint(c.core, w, r)
	if err != nil {
		return
	}

	if err := unmarshalBodyInto(w, r, entrypoint); err != nil {
		return
	}

	core.ZeroReadonlyFields(entrypoint)

	if err := entrypoint.Patch(); err != nil {
		renderError(w, err, http.StatusInternalServerError)
		return
	}

	body, err := marshalBody(w, entrypoint)
	if err != nil {
		return
	}
	renderWithStatusAccepted(w, body)
}

func (c *EntrypointController) Delete(w http.ResponseWriter, r *http.Request) {
	entrypoint, err := loadEntrypoint(c.core, w, r)
	if err != nil {
		return
	}
	if err = entrypoint.Delete(); err != nil {
		renderError(w, err, http.StatusInternalServerError)
		return
	}
}
