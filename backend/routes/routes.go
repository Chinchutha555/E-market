package routes

import (
	"net/http"

	"emarket-api/databases"
	"emarket-api/middlewares"
	"emarket-api/responses"

	authHandlerPkg "emarket-api/src/auth/handlers"
	authRepoPkg "emarket-api/src/auth/repositories"
	authUsecasePkg "emarket-api/src/auth/usecases"

	stallHandlerPkg "emarket-api/src/stall/handlers"
	stallRepoPkg "emarket-api/src/stall/repositories"
	stallUsecasePkg "emarket-api/src/stall/usecases"

	bookingHandlerPkg "emarket-api/src/booking/handlers"
	bookingRepoPkg "emarket-api/src/booking/repositories"
	bookingUsecasePkg "emarket-api/src/booking/usecases"
	debtHandlerPkg "emarket-api/src/debt/handlers"
	notificationHandlerPkg "emarket-api/src/notification/handlers"

	appHandlerPkg "emarket-api/src/application/handlers"
	appRepoPkg "emarket-api/src/application/repositories"
	appUsecasePkg "emarket-api/src/application/usecases"

	appTypeHandlerPkg "emarket-api/src/application_type/handlers"
	appTypeRepoPkg "emarket-api/src/application_type/repositories"
	appTypeUsecasePkg "emarket-api/src/application_type/usecases"

	contractHandlerPkg "emarket-api/src/contract/handlers"
	contractRepoPkg "emarket-api/src/contract/repositories"
	contractUsecasePkg "emarket-api/src/contract/usecases"

	marketHandlerPkg "emarket-api/src/market/handlers"
	marketRepoPkg "emarket-api/src/market/repositories"
	marketUsecasePkg "emarket-api/src/market/usecases"

	itemHandlerPkg "emarket-api/src/item/handlers"
	itemRepoPkg "emarket-api/src/item/repositories"
	itemUsecasePkg "emarket-api/src/item/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middlewares.Cors())
	r.Use(middlewares.Recovery())
	r.Use(middlewares.LogRequest())
	r.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "private, no-store, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Next()
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{
			"status":  "healthy",
			"service": "emarket-api",
		}, http.StatusOK))
	})

	// Dependency Injections
	db := databases.DB

	authRepo := authRepoPkg.NewRepository(db)
	authUC := authUsecasePkg.NewUseCase(authRepo)
	authHandler := authHandlerPkg.NewHandler(authUC)

	stallRepo := stallRepoPkg.NewRepository(db)
	stallUC := stallUsecasePkg.NewUseCase(stallRepo)
	stallHandler := stallHandlerPkg.NewHandler(stallUC)

	bookingRepo := bookingRepoPkg.NewRepository(db)
	bookingUC := bookingUsecasePkg.NewUseCase(bookingRepo)
	bookingHandler := bookingHandlerPkg.NewHandler(bookingUC)
	notificationHandler := notificationHandlerPkg.NewHandler(db)
	debtHandler := debtHandlerPkg.NewHandler(db)

	appRepo := appRepoPkg.NewRepository(db)
	appUC := appUsecasePkg.NewUseCase(appRepo)
	appHandler := appHandlerPkg.NewHandler(appUC)

	appTypeRepo := appTypeRepoPkg.NewRepository(db)
	appTypeUC := appTypeUsecasePkg.NewUseCase(appTypeRepo)
	appTypeHandler := appTypeHandlerPkg.NewHandler(appTypeUC)

	contractRepo := contractRepoPkg.NewRepository(db)
	contractUC := contractUsecasePkg.NewUseCase(contractRepo)
	contractHandler := contractHandlerPkg.NewHandler(contractUC)

	marketRepo := marketRepoPkg.NewRepository(db)
	marketUC := marketUsecasePkg.NewUseCase(marketRepo)
	marketHandler := marketHandlerPkg.NewHandler(marketUC)

	itemRepo := itemRepoPkg.NewRepository(db)
	itemUC := itemUsecasePkg.NewUseCase(itemRepo)
	itemHandler := itemHandlerPkg.NewHandler(itemUC)

	// API Groups
	api := r.Group("/api/v1")
	{
		// Auth Routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middlewares.AuthorizeJWT(), authHandler.GetProfile)
			authGroup.POST("/refresh_token", authHandler.RefreshToken)
		}

		// Market Routes
		marketGroup := api.Group("/markets")
		{
			marketGroup.GET("", marketHandler.GetMarkets)
			marketGroup.GET("/mine", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), marketHandler.GetMyMarkets)
			marketGroup.GET("/:id", marketHandler.GetMarketByID)
			marketGroup.GET("/:id/zones", marketHandler.GetMarketZones)
			marketGroup.GET("/:id/gis-config", marketHandler.GetGISConfig)
			marketGroup.GET("/:id/gis-config/draft", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireMarketAccess(db, "id"), marketHandler.GetGISDraft)
			marketGroup.PUT("/:id/gis-config", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireMarketAccess(db, "id"), marketHandler.SaveGISConfig)
		}
		api.GET("/zones", marketHandler.GetMarketZones)

		// Stall Routes
		stallGroup := api.Group("/stalls")
		{
			stallGroup.GET("", stallHandler.GetStalls)
			stallGroup.GET("/:id", stallHandler.GetStallByID)
			stallGroup.POST("", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireBodyMarketAccess(db), stallHandler.CreateStall)
			stallGroup.PUT("/:id", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireStallAccess(db, "id"), stallHandler.UpdateStall)
			stallGroup.GET("/:id/history", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireStallAccess(db, "id"), stallHandler.GetStallHistory)
		}

		// Public Endpoints (PDPA Masked & Capacity Aggregates)
		api.GET("/bookings/waiting-queue/:ident", stallHandler.GetWaitingQueue)
		api.GET("/bookings/stall-queue/:ident", stallHandler.GetWaitingQueue)
		api.GET("/bookings/visit-capacity", bookingHandler.GetVisitCapacity)
		api.GET("/bookings/policy-settings", bookingHandler.GetPolicySettings)

		// Booking Routes
		bookingGroup := api.Group("/bookings", middlewares.AuthorizeJWT(), middlewares.ScopeMarketQuery(db))
		{
			bookingGroup.GET("", bookingHandler.GetBookings)
			bookingGroup.PUT("/policy-settings", middlewares.RequireRoles("admin", "staff", "market_director"), bookingHandler.UpdatePolicySettings)
			bookingGroup.POST("/policy-settings", middlewares.RequireRoles("admin", "staff", "market_director"), bookingHandler.UpdatePolicySettings)
			bookingGroup.GET("/:id", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.GetBookingByID)
			bookingGroup.POST("", bookingHandler.CreateBooking)
			bookingGroup.POST("/waitlist", bookingHandler.CreateWaitlist)
			bookingGroup.GET("/waitlist/mine", bookingHandler.GetMyWaitlist)
			bookingGroup.GET("/waitlist/stall/:ident", bookingHandler.GetWaitlistByStall)
			bookingGroup.POST("/waitlist/stall/:ident/advance", middlewares.RequireRoles("admin", "staff", "market_director"), bookingHandler.AdvanceWaitlistByStall)
			bookingGroup.DELETE("/waitlist/:id", bookingHandler.WithdrawWaitlist)
			bookingGroup.POST("/waitlist/:id/decline", bookingHandler.DeclineWaitlist)
			bookingGroup.POST("/clear", middlewares.RequireRoles("admin", "staff"), bookingHandler.ClearAll)
			bookingGroup.POST("/:id/stalls/:ident/cancel", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.CancelBookingStall)
			bookingGroup.DELETE("/:id/stalls/:ident", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.CancelBookingStall)
			bookingGroup.POST("/:id/approve", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.Approve)
			bookingGroup.POST("/:id/reject", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.Reject)
			bookingGroup.POST("/:id/mark-visited", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.MarkVisited)
			bookingGroup.POST("/:id/visited", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.MarkVisited)
			bookingGroup.POST("/:id/no-show", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.NoShow)
			bookingGroup.POST("/:id/noshow", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.NoShow)
			bookingGroup.POST("/:id/decide", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.Decide)
			bookingGroup.POST("/:id/officer-record-proxy", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.OfficerRecordProxy)
			bookingGroup.POST("/:id/cancel", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.Cancel)
			bookingGroup.POST("/:id/select-slot", middlewares.RequireResourceAccess(db, "booking"), bookingHandler.SelectSlot)
			bookingGroup.PATCH("/:id/status", middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireResourceAccess(db, "booking"), bookingHandler.UpdateBookingStatus)
		}

		notificationGroup := api.Group("/notifications", middlewares.AuthorizeJWT())
		{
			notificationGroup.GET("", notificationHandler.GetNotifications)
			notificationGroup.POST("/read-all", notificationHandler.MarkAllRead)
			notificationGroup.POST("/:id/read", notificationHandler.MarkRead)
		}

		debtGroup := api.Group("/debts", middlewares.AuthorizeJWT())
		{
			debtGroup.GET("", debtHandler.List)
			debtGroup.POST("", middlewares.RequireRoles("admin", "staff", "market_director"), debtHandler.Create)
		}

		// Application Routes
		appGroup := api.Group("/applications", middlewares.AuthorizeJWT(), middlewares.ScopeMarketQuery(db))
		{
			appGroup.GET("", appHandler.GetApplications)
			appGroup.GET("/:id", middlewares.RequireResourceAccess(db, "application"), appHandler.GetApplicationByID)
			appGroup.GET("/:id/officer-signature", middlewares.RequireResourceAccess(db, "application"), appHandlerPkg.OfficerSignature(db))
			appGroup.GET("/:id/id-card-copy", middlewares.RequireResourceAccess(db, "application"), appHandlerPkg.IDCardCopy(db))
			appGroup.POST("", middlewares.RequireRoles("tenant", "admin", "staff", "market_director"), appHandler.CreateApplication)
			appGroup.POST("/:id/sign", middlewares.RequireRoles("tenant"), middlewares.RequireResourceAccess(db, "application"), appHandlerPkg.SignApplication(db))
			appGroup.POST("/:id/resubmit", middlewares.RequireRoles("tenant"), middlewares.RequireResourceAccess(db, "application"), appHandlerPkg.Resubmit(db))
			appGroup.POST("/:id/action", middlewares.RequireApplicationReviewer(db), appHandler.Action)
			appGroup.POST("/clear", middlewares.RequireRoles("admin", "staff"), func(c *gin.Context) {
				c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"cleared": true}, http.StatusOK))
			})
			appGroup.PATCH("/:id/status", func(c *gin.Context) { c.AbortWithStatus(http.StatusMethodNotAllowed) })
		}

		// Application Type Routes (จัดการประเภทคำร้อง)
		appTypeGroup := api.Group("/application-types")
		{
			appTypeGroup.GET("", appTypeHandler.GetApplicationTypes)
			appTypeGroup.GET("/:id", appTypeHandler.GetApplicationTypeByID)
			appTypeGroup.POST("", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff"), appTypeHandler.CreateApplicationType)
			appTypeGroup.PUT("/:id", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff"), appTypeHandler.UpdateApplicationType)
			appTypeGroup.DELETE("/:id", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff"), appTypeHandler.DeleteApplicationType)
			appTypeGroup.PATCH("/:id/toggle", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff"), appTypeHandler.ToggleActive)
			appTypeGroup.POST("/:id/toggle", middlewares.AuthorizeJWT(), middlewares.RequireRoles("admin", "staff"), appTypeHandler.ToggleActive)
		}

		// Contract Routes
		contractGroup := api.Group("/contracts", middlewares.AuthorizeJWT(), middlewares.ScopeMarketQuery(db))
		{
			contractGroup.GET("", contractHandler.GetContracts)
			contractGroup.GET("/:id", middlewares.RequireResourceAccess(db, "contract"), contractHandler.GetContractByID)
			contractGroup.POST("", middlewares.RequireRoles("admin", "staff", "market_director", "tenant"), middlewares.RequireStallPayloadAccess(db), contractHandler.CreateContract)
			contractGroup.POST("/:id/approve", middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireResourceAccess(db, "contract"), contractHandler.ApproveContract)
			contractGroup.POST("/:id/reject", middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireResourceAccess(db, "contract"), contractHandler.RejectContract)
			contractGroup.POST("/:id/terminate", middlewares.RequireRoles("admin", "staff", "market_director"), middlewares.RequireResourceAccess(db, "contract"), contractHandler.TerminateContract)
		}

		// Item Routes (Blueprint demonstration module)
		itemGroup := api.Group("/items")
		{
			itemGroup.GET("/datatable", itemHandler.GetItems)
			itemGroup.GET("/:id", itemHandler.GetItemByID)
			itemGroup.POST("", middlewares.AuthorizeJWT(), itemHandler.CreateItem)
			itemGroup.PUT("/:id", middlewares.AuthorizeJWT(), itemHandler.UpdateItem)
			itemGroup.DELETE("/:id", middlewares.AuthorizeJWT(), itemHandler.DeleteItem)
		}
	}

	return r
}
