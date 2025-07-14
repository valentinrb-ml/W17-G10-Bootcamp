package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/mappers"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/validators"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/request"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/models/section"
	"net/http"
	"strconv"
)

type SectionDefault struct {
	sv service.SectionService
}

func NewSectionHandler(sv service.SectionService) *SectionDefault {
	return &SectionDefault{sv: sv}
}

func (h *SectionDefault) FindAllSections(w http.ResponseWriter, r *http.Request) {

	sections, err := h.sv.FindAllSections()
	if err != nil {
		response.Error(w, err.ResponseCode, err.Message)
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

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ID param.")
		return
	}

	sec, err1 := h.sv.FindById(id)

	if err1 != nil {
		response.Error(w, err1.ResponseCode, err1.Message)
		return
	}

	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(sec))

}

func (h *SectionDefault) DeleteSection(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ID param.")
		return
	}

	err1 := h.sv.DeleteSection(id)

	if err1 != nil {
		response.Error(w, err1.ResponseCode, err1.Message)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *SectionDefault) CreateSection(w http.ResponseWriter, r *http.Request) {

	var sectionReq section.RequestSection
	err := request.JSON(r, &sectionReq)

	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err1 := validators.ValidateSectionRequest(sectionReq)
	if err1 != nil {
		response.Error(w, err1.ResponseCode, err1.Message)
		return
	}

	sec := mappers.RequestSectionToSection(sectionReq)

	newSection, err2 := h.sv.CreateSection(r.Context(), sec)

	if err2 != nil {
		response.Error(w, err2.ResponseCode, err2.Message)
		return
	}

	response.JSON(w, http.StatusCreated, mappers.SectionToResponseSection(newSection))
}

func (h *SectionDefault) UpdateSection(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ID param.")
		return
	}

	var sec section.RequestSection
	err = request.JSON(r, &sec)

	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err1 := validators.ValidateSectionPatch(sec); err1 != nil {
		response.Error(w, err1.ResponseCode, err1.Message)
		return
	}

	secUpd, err2 := h.sv.UpdateSection(r.Context(), id, sec)

	if err2 != nil {
		response.Error(w, err2.ResponseCode, err2.Message)
		return
	}

	response.JSON(w, http.StatusOK, mappers.SectionToResponseSection(secUpd))

}
