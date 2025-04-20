package controllers_comment

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	serializers_comment "pashmak.com/pashmak/serializers/comment"
	services_auth "pashmak.com/pashmak/services/auth"
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

func (cc *CommentController) GetCommentsByPlaceToken(c *gin.Context) {
	token := c.Param("token")
	comments, err := cc.CommentService.GetCommentsByPlaceToken(token)
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
		"comments": comments,
	})
}

func (cc *CommentController) SetNewReaction(c *gin.Context){
	reactionType := c.Query("type")

	if reactionType == "like"{

	}else if reactionType == "dislike"{}
	c.JSON(403, reactionType)
}

func (cc *CommentController) AddNewComment(c *gin.Context){
	var body serializers_comment.AddCommentRequest

	if c.Bind(body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	placeToken := c.Param("token")
	userinfo, exists := c.Get("user")
	
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"message": "ابتدا باید وارد شوید",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	err := cc.CommentService.AddNewComment(placeToken, userpayload)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیر منتظره ای رخ داد",
		})
		log.Println(err.Error())
		return
	}
}