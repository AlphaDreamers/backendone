package handler

import (
	"github.com/SwanHtetAungPhyo/gis/internal/model/req"
	"github.com/SwanHtetAungPhyo/gis/internal/model/resp"
	"github.com/SwanHtetAungPhyo/gis/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GigHandler struct {
	log *logrus.Logger
	srv *service.GigService
}

func NewGigHandler(
	log *logrus.Logger,
	srv *service.GigService,
) *GigHandler {
	return &GigHandler{log: log, srv: srv}
}

func (gh *GigHandler) CreateGig(c *fiber.Ctx) error {
	var req *req.CreateGigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resp.Response{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	gigToCreate, err := gh.srv.CreateGig(req)
	if err != nil {
		gh.log.Debug("git create gig err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(resp.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp.Response{
		Status:  fiber.StatusCreated,
		Message: "gig created successfully",
		Data:    gigToCreate,
	})
}

func (gh *GigHandler) GetAllGigs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	paganition, models, err := gh.srv.GetAllGigByPaganition(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resp.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	gh.log.Debug("paganition:", paganition)
	return c.Status(fiber.StatusOK).JSON(resp.Response{
		Status:  fiber.StatusOK,
		Message: "Gig retrieval success with paganition",
		Data:    models,
	})
}

func (gh *GigHandler) AddPackage(c *fiber.Ctx) error {
	gigRawId := c.Params("gig_id")
	var req *req.GigPackageRequest
	if gigRawId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(resp.Response{
			Status:  fiber.StatusBadRequest,
			Message: "gig id required in the param",
		})
	}

	gigUuid := uuid.MustParse(gigRawId)
	gig, err := gh.srv.AddPackageToGig(gigUuid, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resp.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp.Response{
		Status:  fiber.StatusOK,
		Message: "gig id successfully added",
		Data:    gig,
	})
}
