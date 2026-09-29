package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miriyaladhanwinn/on-duty/internal/auth"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
	"github.com/miriyaladhanwinn/on-duty/internal/middleware"
	"github.com/miriyaladhanwinn/on-duty/internal/service"
)

type Server struct {
	store   db.Store
	service *service.ODService
	router  *gin.Engine
}

func NewServer(store db.Store) *Server {
	s := &Server{
		store:   store,
		service: service.NewODService(store),
		router:  gin.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiV1 := s.router.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "Amrita CSE On-Duty System", "timestamp": time.Now().Unix()})
		})

		// Auth endpoints
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", s.handleLogin)
			authGroup.POST("/register", s.handleRegister)
			authGroup.GET("/me", middleware.AuthMiddleware(), s.handleMe)
		}

		// Categories
		apiV1.GET("/categories", s.handleListCategories)

		// Protected OD endpoints
		od := apiV1.Group("/od", middleware.AuthMiddleware())
		{
			od.POST("/requests", middleware.RequireRoles(db.RoleStudent), s.handleCreateODRequest)
			od.GET("/requests", s.handleListODRequests)
			od.GET("/requests/:id", s.handleGetODRequest)
			od.POST("/requests/:id/action", middleware.RequireRoles(db.RoleFacultyAdvisor, db.RoleClubCoordinator, db.RoleHOD), s.handleProcessAction)
			od.GET("/stats", s.handleGetStats)
		}
	}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) handleLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.store.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or credentials"})
		return
	}

	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or credentials"})
		return
	}

	token, err := auth.GenerateToken(user, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"roll_no":    user.RollNo,
			"department": user.Department,
			"role":       user.Role,
		},
	})
}

type RegisterInput struct {
	Email      string      `json:"email" binding:"required,email"`
	Password   string      `json:"password" binding:"required,min=6"`
	Name       string      `json:"name" binding:"required"`
	RollNo     string      `json:"roll_no"`
	Department string      `json:"department"`
	Role       db.UserRole `json:"role" binding:"required"`
}

func (s *Server) handleRegister(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	dept := input.Department
	if dept == "" {
		dept = "CSE"
	}

	newUser := &db.User{
		Email:        input.Email,
		PasswordHash: hash,
		Name:         input.Name,
		RollNo:       input.RollNo,
		Department:   dept,
		Role:         input.Role,
	}

	created, err := s.store.CreateUser(c.Request.Context(), newUser)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": created})
}

func (s *Server) handleMe(c *gin.Context) {
	val, _ := c.Get(middleware.ClaimsKey)
	claims := val.(*auth.Claims)
	c.JSON(http.StatusOK, gin.H{"user": claims})
}

func (s *Server) handleListCategories(c *gin.Context) {
	cats, err := s.store.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": cats})
}

func (s *Server) handleCreateODRequest(c *gin.Context) {
	val, _ := c.Get(middleware.ClaimsKey)
	claims := val.(*auth.Claims)

	var input service.CreateODInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req, err := s.service.SubmitRequest(c.Request.Context(), claims, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"request": req})
}

func (s *Server) handleListODRequests(c *gin.Context) {
	val, _ := c.Get(middleware.ClaimsKey)
	claims := val.(*auth.Claims)

	ctx := c.Request.Context()

	if claims.Role == db.RoleStudent {
		requests, err := s.store.ListODRequestsByStudent(ctx, claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"requests": requests})
		return
	}

	// For faculty/coordinators/HOD, return pending requests relevant to their role
	var statusFilter db.ODStatus
	switch claims.Role {
	case db.RoleFacultyAdvisor:
		statusFilter = db.StatusPendingAdvisor
	case db.RoleClubCoordinator:
		statusFilter = db.StatusPendingCoordinator
	case db.RoleHOD:
		statusFilter = db.StatusPendingHOD
	default:
		statusFilter = db.StatusPendingAdvisor
	}

	requests, err := s.store.ListPendingODRequestsByStatus(ctx, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func (s *Server) handleGetODRequest(c *gin.Context) {
	requestID := c.Param("id")
	req, timeline, err := s.service.GetRequestWithTimeline(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OD request not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"request": req, "timeline": timeline})
}

func (s *Server) handleProcessAction(c *gin.Context) {
	val, _ := c.Get(middleware.ClaimsKey)
	claims := val.(*auth.Claims)
	requestID := c.Param("id")

	var input service.ActionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req, err := s.service.ProcessApproval(c.Request.Context(), claims, requestID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": req, "message": "Approval action processed successfully"})
}

func (s *Server) handleGetStats(c *gin.Context) {
	stats, err := s.store.GetDepartmentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
