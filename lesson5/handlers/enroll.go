package handlers

import (
	"course-system/config"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var enrollMutex sync.Mutex

type EnrollRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}

func EnrollCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	enrollMutex.Lock()
	defer enrollMutex.Unlock()

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "事务开始失败",
			"data": nil,
		})
		return
	}

	var capacity, enrolled int
	err = tx.QueryRow("SELECT capacity FROM courses WHERE id = ?", req.CourseID).Scan(&capacity)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "课程不存在",
			"data": nil,
		})
		return
	}

	err = tx.QueryRow("SELECT COUNT(*) FROM enrollments WHERE course_id = ?", req.CourseID).Scan(&enrolled)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}

	if enrolled >= capacity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "课程已满",
			"data": nil,
		})
		return
	}

	var exists int
	err = tx.QueryRow("SELECT COUNT(*) FROM enrollments WHERE student_id = ? AND course_id = ?", userID, req.CourseID).Scan(&exists)
	if err == nil && exists > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "已选过该课程",
			"data": nil,
		})
		return
	}

	_, err = tx.Exec("INSERT INTO enrollments(student_id, course_id) VALUES (?, ?)", userID, req.CourseID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "选课失败",
			"data": nil,
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "选课成功",
		"data": nil,
	})
}

func GetMyCourses(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := config.DB.Query(`
		SELECT c.id, c.name, c.teacher, c.capacity, e.enrolled_at 
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE e.student_id = ?
		ORDER BY e.enrolled_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"data": nil,
		})
		return
	}
	defer rows.Close()

	type CourseEnrollment struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Teacher    string `json:"teacher"`
		Capacity   int    `json:"capacity"`
		EnrolledAt string `json:"enrolled_at"`
	}

	var courses []CourseEnrollment
	for rows.Next() {
		var course CourseEnrollment
		var enrolledAt []uint8
		err := rows.Scan(&course.ID, &course.Name, &course.Teacher, &course.Capacity, &enrolledAt)
		if err != nil {
			continue
		}
		course.EnrolledAt = string(enrolledAt)
		courses = append(courses, course)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": courses,
	})
}

func DropCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID := c.Param("course_id")

	cid, err := strconv.Atoi(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "课程ID错误",
			"data": nil,
		})
		return
	}

	result, err := config.DB.Exec("DELETE FROM enrollments WHERE student_id = ? AND course_id = ?", userID, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "退课失败",
			"data": nil,
		})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "未找到选课记录",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退课成功",
		"data": nil,
	})
}
