package config

import (
	"github.com/darkphotonKN/community-builds/internal/auth"
	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/item"
	"github.com/darkphotonKN/community-builds/internal/member"
	"github.com/darkphotonKN/community-builds/internal/rating"
	"github.com/darkphotonKN/community-builds/internal/skill"
	"github.com/darkphotonKN/community-builds/internal/tag"
	"github.com/gin-gonic/gin"
)

/**
* Sets up API prefix route and all routers.
**/
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// base route
	api := router.Group("/api")

	// -- RATING --

	// --- Rating Setup ---
	ratingRepo := rating.NewRatingRepository(DB)
	ratingService := rating.NewRatingService(ratingRepo)

	// -- MEMBER --

	// --- Member Setup ---
	memberRepo := member.NewMemberRepository(DB)
	memberService := member.NewMemberService(memberRepo)
	memberHandler := member.NewMemberHandler(memberService, ratingService)

	// --- Member Routes ---
	memberRoutes := api.Group("/member")

	memberRoutes.GET("/:id", memberHandler.GetMemberByIdHandler)
	memberRoutes.POST("/signup", memberHandler.CreateMemberHandler)
	memberRoutes.POST("/signin", memberHandler.LoginMemberHandler)
	memberRoutes.Use(auth.AuthMiddleware())
	memberRoutes.POST("/update-password", memberHandler.UpdatePasswordMemberHandler)
	memberRoutes.POST("/update-info", memberHandler.UpdateInfoMemberHandler)

	// -- ITEM --

	// --- Item Setup ---
	itemRepo := item.NewItemRepository(DB)
	itemService := item.NewItemService(itemRepo)
	itemHandler := item.NewItemHandler(itemService)

	// --- Item Routes ---
	itemRoutes := api.Group("/item")
	// Protected Routes
	itemRoutes.Use(auth.AuthMiddleware())
	itemRoutes.GET("/", itemHandler.GetItemsHandler)
	itemRoutes.POST("/", itemHandler.CreateItemHandler)
	itemRoutes.PATCH("/:id", itemHandler.UpdateItemsHandler)

	// -- SKILL --

	// --- Skill Setup ---
	skillRepo := skill.NewSkillRepository(DB)
	skillService := skill.NewSkillService(skillRepo)
	skillHandler := skill.NewSkillHandler(skillService)

	// --- Skill Routes ---
	skillRoutes := api.Group("/skill")
	// Protected Routes
	skillRoutes.Use(auth.AuthMiddleware())
	skillRoutes.GET("/", skillHandler.GetSkillsHandler)
	skillRoutes.POST("/", skillHandler.CreateSkillHandler)

	// -- BUILD --

	// --- Build Setup ---
	buildRepo := build.NewBuildRepository(DB)
	buildService := build.NewBuildService(buildRepo, skillService)
	buildHandler := build.NewBuildHandler(buildService)

	// --- Build Routes ---
	buildRoutes := api.Group("/build")
	// Protected Routes
	buildRoutes.Use(auth.AuthMiddleware())
	buildRoutes.GET("/", buildHandler.GetBuildsForMemberHandler)
	buildRoutes.POST("/", buildHandler.CreateBuildHandler)
	buildRoutes.POST("/:id/addSkills", buildHandler.AddSkillsToBuild)

	// -- TAG --

	// --- Tag Setup ---
	tagRepo := tag.NewTagRepository(DB)
	tagService := tag.NewTagService(tagRepo)
	tagHandler := tag.NewTagHandler(tagService)
	// --- Tag Routes ---
	tagRoutes := api.Group("/tag")
	// Protected Routes
	tagRoutes.Use(auth.AuthMiddleware())
	tagRoutes.GET("/", tagHandler.GetTagsHandler)
	tagRoutes.POST("/", tagHandler.CreateTagHandler)
	tagRoutes.PATCH("/:id", tagHandler.UpdateTagsHandler)

	// -- TAG --

	// --- Tag Setup ---
	tagRepo := tag.NewTagRepository(DB)
	tagService := tag.NewTagService(tagRepo)
	tagHandler := tag.NewTagHandler(tagService)
	// --- Tag Routes ---
	tagRoutes := api.Group("/tag")
	// Protected Routes
	tagRoutes.Use(auth.AuthMiddleware())
	tagRoutes.GET("/", tagHandler.GetTagsHandler)
	tagRoutes.POST("/", tagHandler.CreateTagHandler)
	tagRoutes.PATCH("/:id", tagHandler.UpdateTagsHandler)

	return router
}
