package handler

import (
	"net/http"

	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/section"

	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
	models "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
)

// SectionDefault handles HTTP requests related to warehouse sections.
type SectionDefault struct {
	sv     service.SectionService
	logger logger.Logger
}

// SetLogger allows injecting the logger after creation
func (h *SectionDefault) SetLogger(l logger.Logger) {
	h.logger = l
}

// NewSectionHandler creates a new SectionDefault handler with the given service.
func NewSectionHandler(sv service.SectionService) *SectionDefault {
	return &SectionDefault{sv: sv}
}

// FindAllSections handles GET /sections to return all sections.
// - Maps domain sections to response models.
func (h *SectionDefault) FindAllSections(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "section-handler", "Find all sections request received")
	}
	ctx := r.Context()

	sections, err := h.sv.FindAllSections(ctx)

	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to find all sections", err)
		}
		response.ErrorWithRequest(w, r, err)
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
	if h.logger != nil {
		h.logger.Info(r.Context(), "section-handler", "Find by id request received")
	}
	ctx := r.Context()
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to parse id parameter", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	sec, err := h.sv.FindById(ctx, id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to find section by id", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(*sec))

}

// DeleteSection handles DELETE /sections/{id} to remove a section.
// Returns 204 No Content on success.
func (h *SectionDefault) DeleteSection(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "section-handler", "Delete section request received")
	}
	ctx := r.Context()

	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to parse id parameter", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	err1 := h.sv.DeleteSection(ctx, id)

	if err1 != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to delete section", err1)
		}
		response.ErrorWithRequest(w, r, err1)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// CreateSection handles POST /sections to add a new section.
// Validates input, maps request to domain, and returns 201 Created with section.
func (h *SectionDefault) CreateSection(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "section-handler", "Create section request received")
	}

	ctx := r.Context()

	var sectionReq models.PostSection
	if err := httputil.DecodeJSON(r, &sectionReq); err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to decode JSON", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	err := validators.ValidateSectionRequest(sectionReq)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Validation failed for create request", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	sec := mappers.RequestSectionToSection(sectionReq)
	newSection, err := h.sv.CreateSection(ctx, sec)
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to create section", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.SectionToResponseSection(*newSection))
}

// UpdateSection handles PATCH /sections/{id} to update only specified fields.
// Decodes JSON patch, validates and returns updated section.
func (h *SectionDefault) UpdateSection(w http.ResponseWriter, r *http.Request) {
	if h.logger != nil {
		h.logger.Info(r.Context(), "section-handler", "Update section request received")
	}
	ctx := r.Context()
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to parse id parameter", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	var sec models.PatchSection
	if err := httputil.DecodeJSON(r, &sec); err != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to decode JSON", err)
		}
		response.ErrorWithRequest(w, r, err)
		return
	}

	if err1 := validators.ValidateSectionPatch(sec); err1 != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Validation failed for update request", err1)
		}
		response.ErrorWithRequest(w, r, err1)
		return
	}

	secUpd, err2 := h.sv.UpdateSection(ctx, id, sec)

	if err2 != nil {
		if h.logger != nil {
			h.logger.Error(r.Context(), "section-handler", "Failed to update section", err2)
		}
		response.ErrorWithRequest(w, r, err2)
		return
	}
	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(*secUpd))
}
