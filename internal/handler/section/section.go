package section

import (
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/section"
	"net/http"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// SectionDefault handles HTTP requests related to warehouse sections.
type SectionDefault struct {
	sv service.SectionService
}

// NewSectionHandler creates a new SectionDefault handler with the given service.
func NewSectionHandler(sv service.SectionService) *SectionDefault {
	return &SectionDefault{sv: sv}
}

// FindAllSections handles GET /sections to return all sections.
// - Maps domain sections to response models.
func (h *SectionDefault) FindAllSections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sections, err := h.sv.FindAllSections(ctx)

	if err != nil {
		response.Error(w, err)
		return
	}

	sectionDoc := make([]models.ResponseSection, 0, len(sections))
	for _, s := range sections {
		secDoc := mappers.SectionToResponseSection(s)
		sectionDoc = append(sectionDoc, secDoc)
	}
	response.JSON(w, http.StatusOK, sectionDoc)

}

// FindById handles GET /sections/{id} to return a section by its ID.
func (h *SectionDefault) FindById(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	sec, err := h.sv.FindById(ctx, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(*sec))

}

// DeleteSection handles DELETE /sections/{id} to remove a section.
// Returns 204 No Content on success.
func (h *SectionDefault) DeleteSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	err1 := h.sv.DeleteSection(ctx, id)

	if err1 != nil {
		response.Error(w, err1)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// CreateSection handles POST /sections to add a new section.
// Validates input, maps request to domain, and returns 201 Created with section.
func (h *SectionDefault) CreateSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sectionReq models.PostSection
	if err := httputil.DecodeJSON(r, &sectionReq); err != nil {
		response.Error(w, err)
		return
	}

	err := validators.ValidateSectionRequest(sectionReq)
	if err != nil {
		response.Error(w, err)
		return
	}

	sec := mappers.RequestSectionToSection(sectionReq)
	newSection, err := h.sv.CreateSection(ctx, sec)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.SectionToResponseSection(*newSection))
}

// UpdateSection handles PATCH /sections/{id} to update only specified fields.
// Decodes JSON patch, validates and returns updated section.
func (h *SectionDefault) UpdateSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var sec models.PatchSection
	if err := httputil.DecodeJSON(r, &sec); err != nil {
		response.Error(w, err)
		return
	}

	if err1 := validators.ValidateSectionPatch(sec); err1 != nil {
		response.Error(w, err1)
		return
	}

	secUpd, err2 := h.sv.UpdateSection(ctx, id, sec)

	if err2 != nil {
		response.Error(w, err2)
		return
	}
	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(*secUpd))
}
