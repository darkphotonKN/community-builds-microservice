package config

import (
	"github.com/darkphotonKN/community-builds/internal/auth"
	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/item"
	"github.com/darkphotonKN/community-builds/internal/member"
	"github.com/darkphotonKN/community-builds/internal/rating"
	"github.com/darkphotonKN/community-builds/internal/skill"
	"github.com/darkphotonKN/community-builds/internal/tag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/**
* Sets up API prefix route and all routers.
**/
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// TODO: CORS for development, remove in PROD
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3666"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// base route
	api := router.Group("/api")

	// --- ITEM ---

	// -- Item Setup --
	itemRepo := item.NewItemRepository(DB)
	itemService := item.NewItemService(itemRepo)
	itemHandler := item.NewItemHandler(itemService)

	// -- Item Routes --
	itemRoutes := api.Group("/item")

	// Protected Routes
	itemRoutes.Use(auth.AuthMiddleware())
	itemRoutes.GET("/", itemHandler.GetItemsHandler)
	itemRoutes.POST("/", itemHandler.CreateItemHandler)
	itemRoutes.PATCH("/:id", itemHandler.UpdateItemsHandler)
	itemRoutes.POST("/rare-item", itemHandler.CreateRareItemHandler)

	// --- SKILL ---

	// -- Skill Setup --
	skillRepo := skill.NewSkillRepository(DB)
	skillService := skill.NewSkillService(skillRepo)
	skillHandler := skill.NewSkillHandler(skillService)

	// -- Skill Routes --
	skillRoutes := api.Group("/skill")

	// Protected Routes
	skillRoutes.Use(auth.AuthMiddleware())
	skillRoutes.GET("/", skillHandler.GetSkillsHandler)
	skillRoutes.POST("/", skillHandler.CreateSkillHandler)

	// --- BUILD ---

	// -- Build Setup --
	buildRepo := build.NewBuildRepository(DB)
	buildService := build.NewBuildService(buildRepo, skillService)
	buildHandler := build.NewBuildHandler(buildService)

	// -- Build Routes --
	buildRoutes := api.Group("/build")

	// Public Routes
	buildRoutes.GET("/community", buildHandler.GetCommunityBuildsHandler)
	buildRoutes.GET("/community/:id/info", buildHandler.GetBuildInfoByIdHandler)

	// Protected Routes
	protectedBuildRoutes := buildRoutes.Group("/")
	protectedBuildRoutes.Use(auth.AuthMiddleware())
	protectedBuildRoutes.GET("/", buildHandler.GetBuildsForMemberHandler)
	protectedBuildRoutes.GET("/:id/info", buildHandler.GetBuildInfoForMemberHandler)
	protectedBuildRoutes.POST("/", buildHandler.CreateBuildHandler)
	protectedBuildRoutes.PATCH("/:id", buildHandler.UpdateBuildHandler)
	protectedBuildRoutes.POST("/:id/addSkills", buildHandler.AddSkillLinksToBuildHandler)
	protectedBuildRoutes.PATCH("/:id/updateSet", buildHandler.UpdateItemSetsToBuildHandler)

	// --- TAG ---

	// -- Tag Setup --
	tagRepo := tag.NewTagRepository(DB)
	tagService := tag.NewTagService(tagRepo)
	tagHandler := tag.NewTagHandler(tagService)

	// -- Tag Routes --
	tagRoutes := api.Group("/tag")

	// Protected Routes
	tagRoutes.Use(auth.AuthMiddleware())
	tagRoutes.GET("/", tagHandler.GetTagsHandler)
	tagRoutes.POST("/", tagHandler.CreateTagHandler)
	tagRoutes.PATCH("/:id", tagHandler.UpdateTagsHandler)

	// -- RATING --

	// --- Rating Setup ---
	ratingRepo := rating.NewRatingRepository(DB)
	ratingService := rating.NewRatingService(ratingRepo, buildService)
	ratingHandler := rating.NewRatingHandler(ratingService)

	ratingRoutes := api.Group("/rating")

	ratingRoutes.Use(auth.AuthMiddleware())
	ratingRoutes.POST("/", ratingHandler.CreateRatingByBuildIdHandler)

	// --- MEMBER ---

	// -- Member Setup --
	memberRepo := member.NewMemberRepository(DB)
	memberService := member.NewMemberService(memberRepo)
	memberHandler := member.NewMemberHandler(memberService, ratingService)

	// -- Member Routes --
	memberRoutes := api.Group("/member")

	// Public Routes
	memberRoutes.GET("/:id", memberHandler.GetMemberByIdHandler)
	memberRoutes.POST("/signup", memberHandler.CreateMemberHandler)
	memberRoutes.POST("/signin", memberHandler.LoginMemberHandler)

	// Protected Routes
	protectedMemberRoutes := memberRoutes.Group("/")
	protectedMemberRoutes.Use(auth.AuthMiddleware())
	protectedMemberRoutes.POST("/update-password", memberHandler.UpdatePasswordMemberHandler)
	protectedMemberRoutes.POST("/update-info", memberHandler.UpdateInfoMemberHandler)

	return router
}
