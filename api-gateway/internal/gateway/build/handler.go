package build

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/types"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuildHandler struct {
	Client BuildClient
}

func NewHandler(client BuildClient) *BuildHandler {
	return &BuildHandler{
		Client: client,
	}
}

/**
* Get all builds for community viewing.
**/

func (h *BuildHandler) GetCommunityBuildsHandler(c *gin.Context) {
	// defaults
	pageNo := 1
	pageSize := 20

	// parse query pagination querystrings to ints
	if pageNoQuery := c.Query("page_no"); pageNoQuery != "" {
		pageNo, _ = strconv.Atoi(pageNoQuery)
	}

	if pageSizeQuery := c.Query("page_size"); pageSizeQuery != "" {
		pageSize, _ = strconv.Atoi(pageSizeQuery)
	}

	// query strings
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	search := c.Query("search")
	skillQuery := c.Query("skill")
	minRatingQuery := c.Query("min_rating")
	ratingCategory := c.Query("rating_category")

	// validate querystrings
	skillId, err := uuid.Parse(skillQuery)
	if skillQuery != "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Skill in querystring was not a valid uuid, error: %s", err.Error())})
		return
	}

	var minRating int
	if minRatingQuery != "" {
		minRating, err = strconv.Atoi(minRatingQuery)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("minRating in querystring was not a valid integer, error: %s", err.Error())})
			return
		}

		fmt.Println("minRating:", minRating)

		if minRating < 1 || minRating > 10 {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"statusCode": http.StatusBadRequest,
					"message":    "min_rating needs to be in the range 1-10."})
			return
		}
	}

	// Get builds from gRPC client
	grpcReq := &pb.GetCommunityBuildsRequest{
		PageNo:         int32(pageNo),
		PageSize:       int32(pageSize),
		SortOrder:      sortOrder,
		SortBy:         sortBy,
		Search:         search,
		SkillId:        skillId.String(),
		MinRating:      int32(minRating),
		RatingCategory: string(types.RatingCategory(ratingCategory)),
	}
	builds, err := h.Client.GetCommunityBuilds(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all community builds: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all builds for member.", "result": builds})
}

/**
* Create build for a signed-in member.
**/
func (h *BuildHandler) CreateBuildHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	// memberId, _ := c.Get("userId")
	var createBuildReq CreateBuildRequest

	if err := c.ShouldBindJSON(&createBuildReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	tagIdsStr := make([]string, len(createBuildReq.TagIds))
	for i, tagId := range createBuildReq.TagIds {
		tagIdsStr[i] = tagId.String()
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.CreateBuildRequest{
		MemberId:     userIdStr.(string),
		SkillId:      createBuildReq.SkillId.String(),
		TagIds:       tagIdsStr,
		Title:        createBuildReq.Title,
		Description:  createBuildReq.Description,
		ClassId:      createBuildReq.ClassId.String(),
		AscendancyId: createBuildReq.AscendancyId.String(),
	}

	build, err := h.Client.CreateBuild(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create a build: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created the build.", "result": build})

}

// /**
// * Updates an existing build for a signed-in member.
// **/
func (h *BuildHandler) UpdateBuildHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	buildIdQuery := c.Param("id")

	buildId, err := uuid.Parse(buildIdQuery)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", buildId)})
		return
	}

	var request UpdateBuildRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	tags := make([]string, len(request.TagIds))
	for _, tag := range request.TagIds {
		tags = append(tags, tag.String())
	}
	grpcReq := &pb.UpdateBuildRequest{
		MemberId:     userIdStr.(string),
		Id:           buildId.String(),
		Tags:         tags,
		SkillId:      request.SkillId.String(),
		Title:        *request.Title,
		Description:  *request.Description,
		ClassId:      request.ClassId.String(),
		AscendancyId: request.AscendancyId.String(),
	}

	_, err = h.Client.UpdateBuild(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update a build: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully updated the build."})
}

// /**
// * Get list of builds by a signed-in member's ID.
// **/
func (h *BuildHandler) GetBuildsForMemberHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	grpcReq := &pb.GetBuildsForMemberRequest{
		MemberId: userIdStr.(string),
	}

	builds, err := h.Client.GetBuildsForMember(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all builds for memberId %s: %s", userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all builds for member.", "result": builds})
}

// /**
// * Get all information for a single build by ID for a particular member.
// **/
func (h *BuildHandler) GetBuildInfoForMemberHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	// fmt.Printf("memberId: %s, id: %s\n", memberId, id)

	grpcReq := &pb.GetBuildInfoForMemberRequest{
		MemberId: userIdStr.(string),
		Id:       id.String(),
	}

	build, err := h.Client.GetBuildInfoForMember(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all build information for memberId %s: %s", userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved build for member.", "result": build})
}

// /**
// * Get all information for a single build by id community version.
// **/
func (h *BuildHandler) GetBuildInfoByIdHandler(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	fmt.Printf("id: %s\n", id)

	grpcReq := &pb.GetBuildInfoRequest{
		Id: id.String(),
	}

	build, err := h.Client.GetBuildInfo(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all build information: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved build for member.", "result": build})
}

// /**
// * Adds primary, secondary, and other skills and links to an existing build.
// **/
func (h *BuildHandler) AddSkillLinksToBuildHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	var request AddSkillsToBuildRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	links := make([]string, len(request.MainSkillLinks.Links))
	for _, link := range request.MainSkillLinks.Links {
		links = append(links, link.String())
	}
	MainSkillLinks := &pb.SkillLinks{
		SkillLinkName: request.MainSkillLinks.SkillLinkName,
		Skill:         request.MainSkillLinks.Skill.String(),
		Links:         links,
	}

	additionsLinks := make([]*pb.SkillLinks, len(request.AdditionalSkills))
	for _, additionsLink := range request.AdditionalSkills {

		subLinks := make([]string, len(additionsLink.Links))
		for _, subLink := range additionsLink.Links {
			subLinks = append(subLinks, subLink.String())
		}
		additionsLinks = append(additionsLinks, &pb.SkillLinks{
			SkillLinkName: additionsLink.SkillLinkName,
			Skill:         additionsLink.Skill.String(),
			Links:         subLinks,
		})
	}
	AdditionalSkills := []*pb.SkillLinks{}

	grpcReq := &pb.AddSkillLinksToBuildRequest{
		MemberId:         userIdStr.(string),
		Id:               id.String(),
		MainSkillLinks:   MainSkillLinks,
		AdditionalSkills: AdditionalSkills,
	}

	_, err = h.Client.AddSkillLinksToBuild(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting add skills to builds, buildId %s: memberId: %s, error: %s", id, userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully added skills to build for member."})
}

// /**
// * Update set.
// **/
func (h *BuildHandler) UpdateItemSetsToBuildHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	var request AddItemsToBuildRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	grpcReq := &pb.UpdateItemSetsToBuildRequest{
		MemberId:   userIdStr.(string),
		Id:         id.String(),
		Weapon:     request.Weapon,
		Shield:     request.Shield,
		Helmet:     request.Helmet,
		BodyArmour: request.BodyArmour,
		Boots:      request.Boots,
		Gloves:     request.Gloves,
		Belt:       request.Belt,
		Amulet:     request.Amulet,
		LeftRing:   request.LeftRing,
		RightRing:  request.RightRing,
	}

	_, err = h.Client.UpdateItemSetsToBuild(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting add items to builds, buildId %s: memberId: %s, error: %s", id, userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully added items to build for member."})
}

// /**
// * Updates a specific build's skill links.
// **/
// func (h *BuildHandler) UpdateBuildSkillLinksHandler(c *gin.Context) {
// 	memberId, _ := c.Get("userId")

// 	idParam := c.Param("id")

// 	id, err := uuid.Parse(idParam)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
// 		return
// 	}

// 	var request UpdateSkillsToBuildRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
// 		return
// 	}

// 	err = h.Service.UpdateBuildSkillLinksService(memberId.(uuid.UUID), id, request)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all builds for memberId %s: %s", memberId, err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all builds for member."})
// }

// /**
// * Deletes build by member Id.
// **/

func (h *BuildHandler) DeleteBuildForMemberHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")

	idParam := c.Param("id")

	buildId, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", buildId)})
		return
	}
	grpcReq := &pb.DeleteBuildByMemberRequest{
		MemberId: userIdStr.(string),
		Id:       buildId.String(),
	}

	_, err = h.Client.DeleteBuildByMember(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all builds for memberId %s: %s", userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": fmt.Sprintf("Successfully deleted build with build id: %s.", buildId)})
}

// /**
// * Quick example setup for quick creation of extra handlers.
// **/
// func (h *BuildHandler) GetBuildsTemplate(c *gin.Context) {
// 	memberId, _ := c.Get("userId")

// 	builds, err := h.Service.GetBuildsForMemberService(memberId.(uuid.UUID))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all builds for memberId %s: %s", memberId, err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all builds for member.", "result": builds})
// }

// /**
// * Publish a build for a member by Id.
// **/
func (h *BuildHandler) PublishBuildHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	idParams := c.Param("id")

	id, err := uuid.Parse(idParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.\n", id)})
		return
	}

	grpcReq := &pb.PublishBuildRequest{
		MemberId: userIdStr.(string),
		Id:       id.String(),
	}
	_, err = h.Client.PublishBuild(c.Request.Context(), grpcReq)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Could not publish build due to error: %s\n", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully published build.", "result": "success"})
}
