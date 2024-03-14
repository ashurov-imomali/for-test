package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/db"
	"main/models"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repos *db.Repository
}

func GetHandler(repository *db.Repository) *Handler {
	return &Handler{repository}
}

func (h *Handler) Test(c *gin.Context) {
	c.String(200, "ok")
}

func (h *Handler) CreateCommissionProfile(c *gin.Context) {
	var request models.ProfileCreatRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse 2 json"})
		return
	}
	var userId int64 = 1
	request.Profile.CreatedBy = userId
	profileId, err := h.Repos.CreateProfile(&request.Profile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't create profile"})
		return
	}
	err = h.Repos.CreateRules(request.Rules, profileId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't add new rules"})
		return
	}
	c.Status(201)
}

func (h *Handler) UpdateCommissionProfile(c *gin.Context) {
	var updProfile models.ProfileCreatRequest
	err := c.ShouldBindJSON(&updProfile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse to json"})
		return
	}
	user := 2
	now := time.Now()
	updProfile.Profile.UpdatedBy = int64(user)
	updProfile.Profile.UpdatedAt = &now

	if updProfile.Profile.Active == false {
		updProfile.Profile.DeletedAt = &now
		err := h.Repos.DeleteProfileAndRules(&updProfile)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "can't delete profile"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
		return
	}
	err = h.Repos.UpdateProfileAndRules(&updProfile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't update profile"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateCommissionRules(c *gin.Context) {
	var rules models.CommissionRules
	err := c.ShouldBindJSON(&rules)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't parse 2 json"})
		return
	}
	var userId int64 = 2
	now := time.Now()
	err = h.Repos.UpdateProfileRules(rules.ProfileId, userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't update profile"})
		return
	}
	if rules.Active != nil && *rules.Active == false {
		if err := h.Repos.DeleteRule(&rules); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "can't delete rule"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted successful"})
		return
	}
	rules.UpdatedAt = &now
	updateRules, err := h.Repos.UpdateRules(&rules)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't update rules"})
		return
	}
	c.JSON(http.StatusOK, updateRules)
}

func (h *Handler) GetAllCommissionProfiles(c *gin.Context) {
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	strPage := c.Query("page")
	strLimit := c.Query("limit")
	profileName := c.Query("profileName")
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		log.Println(err)
	}
	page, err := strconv.Atoi(strPage)
	if err != nil {
		log.Println(err)
	}
	offset, limit := makePagination(page, limit)
	allProfiles, count, err := h.Repos.AllProfiles(dateFrom, dateTo, profileName, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message:": "Can't get all profiles error:" + err.Error()})
		return
	}
	response := struct {
		Profiles []models.ProfileResponse
		Total    int64
	}{allProfiles, *count}

	c.JSON(http.StatusOK, response)
}

func makePagination(page, limit int) (int, int) {
	if limit < 1 {
		limit = 1
	}
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit, limit
}
