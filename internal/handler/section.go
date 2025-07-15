package handler

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/httputil"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"net/http"
)

type SectionDefault struct {
	sv service.SectionService
}

func NewSectionHandler(sv service.SectionService) *SectionDefault {
	return &SectionDefault{sv: sv}
}

func (h *SectionDefault) FindAllSections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sections, err := h.sv.FindAllSections(ctx)

	if err != nil {
		response.Error(w, err)
		return
	}

	sectionDoc := make([]section.ResponseSection, 0, len(sections))
	for _, s := range sections {
		secDoc := mappers.SectionToResponseSection(s)
		sectionDoc = append(sectionDoc, secDoc)
	}
	response.JSON(w, http.StatusOK, sectionDoc)

}

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

func (h *SectionDefault) CreateSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var sectionReq section.PostSection
	err := request.JSON(r, &sectionReq)

	if err != nil {
		response.Error(w, err)
		return
	}

	err1 := validators.ValidateSectionRequest(sectionReq)
	if err1 != nil {
		response.Error(w, err1)
		return
	}

	sec := mappers.RequestSectionToSection(sectionReq)
	newSection, err2 := h.sv.CreateSection(ctx, sec)
	if handleApiError(w, err2) {
		return
	}

	response.JSON(w, http.StatusCreated, mappers.SectionToResponseSection(*newSection))
}

func (h *SectionDefault) UpdateSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputil.ParseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var sec section.PatchSection
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
