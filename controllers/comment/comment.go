package controllers_comment

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	services_comment "pashmak.com/pashmak/services/comment"
)

type CommentController struct{
	CommentService *services_comment.CommentService
}

func NewCommentController(commentservice *services_comment.CommentService) *CommentController{
	return &CommentController{
		CommentService: commentservice,
	}
}

func (cc *CommentController) GetCommentsByToken(c *gin.Context) {
	token, _ := strconv.Atoi(c.Param("token"))
	comments, err := cc.CommentService.GetCommentsByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": profile,
	})
}