package main

import (
	"log"
	"postman-round-2/internal/core"
	"postman-round-2/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

func main() {
	core.ConnectDB()
	err := core.DB.AutoMigrate(&models.User{}, &models.Leave{}, &models.Attendance{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		// write funct for showing error when role mentioned is not one of the given
		user := models.User{Name: req.Name, Email: req.Email, Password: req.Password, Role: req.Role}
		if err := core.DB.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(201, gin.H{"user": user})
	})

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		var user models.User
		if err := core.DB.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "invalid email or password"})
			return
		}
		// verify login details with jwt
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"student_id": user.StudentID,
			"role":       user.Role,
			"exp":        time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "could not login"})
			return
		}
		c.JSON(200, gin.H{"token": tokenString})

	})
	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		if err := core.DB.Find(&users).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch users"})
			return
		}
		c.JSON(200, gin.H{"users": users})
	})

	r.POST("/leaves/apply", func(c *gin.Context) {
		var req struct {
			StudentID uint   `json:"student_id"`
			Reason    string `json:"reason"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		leave := models.Leave{StudentID: req.StudentID, Reason: req.Reason, Status: "pending", StartDate: req.StartDate, EndDate: req.EndDate}
		if err := core.DB.Create(&leave).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to apply leave"})
			return
		}
		c.JSON(201, gin.H{"leave": leave})
	})
	r.GET("/leaves", func(c *gin.Context) {
		var leaves []models.Leave
		if err := core.DB.Find(&leaves).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch leaves"})
			return
		}
		c.JSON(200, gin.H{"leaves": leaves})
	})
	r.POST("/attendance/mark", func(c *gin.Context) {
		var req struct {
			StudentID uint   `json:"student_id"`
			Date      string `json:"date"`
			Present   bool   `json:"present"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		attendance := models.Attendance{StudentID: req.StudentID, Date: req.Date, Present: req.Present}
		if err := core.DB.Create(&attendance).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to mark attendance"})
			return
		}
		c.JSON(201, gin.H{"attendance": attendance})
	})

	r.PATCH("/leaves/:student_id/status", func(c *gin.Context) {
		type StatusUpdate struct {
			Status string `json:"status"`
		}
		var req StatusUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		id := c.Param("id")

		var leave models.Leave
		if err := core.DB.First(&leave, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "leave not found"})
			return
		}

		leave.Status = req.Status
		if err := core.DB.Save(&leave).Error; err != nil {
			c.JSON(500, gin.H{"error": "could not update leave status"})
			return
		}
		c.JSON(200, gin.H{"leave": leave})
	})

	r.GET("attendance/stats/:student_id", func(c *gin.Context) {
		studentID := c.Param("student_id")

		var records []models.Attendance
		if err := core.DB.Where("student_id = ?", studentID).Find(&records).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch attendance"})
			return
		}

		var leaves []models.Leave
		if err := core.DB.Where("student_id = ? AND status = ?", studentID, "approved").Find(&leaves).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch leaves"})
			return
		}
		presentDays := 0
		totalDays := len(records)

		isonLeave := func(date string) bool {
			for _, leave := range leaves {
				if date >= leave.StartDate && date <= leave.EndDate {
					return true
				}
			}
			return false
		}

		for _, att := range records {
			if att.Present && !isonLeave(att.Date) {
				presentDays++
			}
		}

		percent := 0.0
		if totalDays > 0 {
			percent = float64(presentDays) / float64(totalDays) * 100
		}

		c.JSON(200, gin.H{
			"student_id":            studentID,
			"present_days":          presentDays,
			"total_days":            totalDays,
			"attendance_percentage": percent,
		})
	})

	r.Run(":8080")
}
