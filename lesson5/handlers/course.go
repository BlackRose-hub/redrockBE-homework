package handlers

import (
	"course-system/config"
	"course-system/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCourses 获取课程列表（公开或需认证）
func GetCourses(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, teacher, capacity, created_at 
		FROM courses 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取课程列表失败",
			"data": nil,
		})
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name, &course.Teacher, &course.Capacity, &course.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "数据解析失败",
				"data": nil,
			})
			return
		}
		courses = append(courses, course)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": courses,
	})
}

// CreateCourse 创建课程（仅管理员）
func CreateCourse(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Teacher  string `json:"teacher" binding:"required"`
		Capacity int    `json:"capacity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	result, err := config.DB.Exec(
		"INSERT INTO courses(name, teacher, capacity) VALUES (?, ?, ?)",
		req.Name, req.Teacher, req.Capacity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "课程创建失败",
			"data": nil,
		})
		return
	}

	courseID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "课程创建成功",
		"data": gin.H{
			"course_id": courseID,
			"name":      req.Name,
			"teacher":   req.Teacher,
			"capacity":  req.Capacity,
		},
	})
}
