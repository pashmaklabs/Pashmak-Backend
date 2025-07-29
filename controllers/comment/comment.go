package controllers_comment

import (
	"log"
	"net/http"
	"strconv"

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
	isLoggedIn := false
	var userpayload services_auth.UserInfo
	userinfo, _ := c.Get("user")
	if userinfo != nil{
		isLoggedIn = true
		userpayload = userinfo.(services_auth.UserInfo)
	}
	
	token := c.Param("placeToken")
	paginator, pagedComments, err := cc.CommentService.GetCommentsByPlaceToken(c, token, userpayload, isLoggedIn)
	if err != nil {
		if err.Error() == "no comments found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "دیدگاهی برای این مکان ثبت نشده است",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"comments": pagedComments,
		"Pagination": paginator.PageInfo,
	})
}

func (cc *CommentController) SetNewReaction(c *gin.Context){
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_comment.AddReactionRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected AddReactionRequest, got %T", validatedData)
		return
	}

	userinfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}

	userpayload := userinfo.(services_auth.UserInfo)
	commentID, _ := strconv.Atoi(c.Param("id"))

	// id, _ := strconv.Atoi(c.Param("id"))

	err := cc.CommentService.AddReaction(userpayload, uint(commentID), body.ReactionType)
	if err != nil{
		if err.Error() == "comment not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"message": "کامنت یافت نشد",
			})
			return
		}else{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			return
		}
	}
	c.Status(http.StatusOK)
}

func (cc *CommentController) AddNewComment(c *gin.Context){
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_comment.AddCommentRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected AddCommentRequest, got %T", validatedData)
		return
	}

	placeToken := c.Param("id")
	userinfo, exists := c.Get("user")
	
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"message": "ابتدا باید وارد شوید",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	err := cc.CommentService.AddNewComment(placeToken, userpayload, body)
	if err != nil{
		if err.Error() == "place not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"message": "مکان یافت نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیر منتظره ای رخ داد",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"message": "دیدگاه با موفقیت ثبت شد",
	})
}

func (cc *CommentController) RemoveReaction(c *gin.Context){
	userinfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}

	userpayload := userinfo.(services_auth.UserInfo)
	commentID, _ := strconv.Atoi(c.Param("id"))



	err := cc.CommentService.RemoveRection(userpayload, uint(commentID))
	if err != nil{
		if err.Error() == "comment not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "دیدگاه یافت نشد",
			})
			return
		}else{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			return
		}
	}
	c.Status(http.StatusOK)
}

func (cc *CommentController) ReportComment(c *gin.Context){
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_comment.SendReportRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected AddCommentRequest, got %T", validatedData)
		return
	}

	commentID, _ := strconv.Atoi(c.Param("id"))
	userinfo, exists := c.Get("user")
	
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"message": "ابتدا باید وارد شوید",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	err := cc.CommentService.ReportComment(userpayload, commentID, body.Reason)
	if err != nil{
		if err.Error() == "comment not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "دیدگاه یافت نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیر منتظره ای رخ داد",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"message": "گزارش با موفقیت ثبت شد",
	})
}

func (cc *CommentController) GetReportedComments(c *gin.Context){
	status := c.Query("status")
	paginator, pagedReports, err := cc.CommentService.GetReportedComments(c, status)
	if err != nil {
		if err.Error() == "no reports found"{
			if pagedReports == nil {
				pagedReports = []serializers_comment.ReportedCommentsResponse{}
			}
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"reports": pagedReports,
				"paginator": paginator,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"reports": pagedReports,
		"paginator": paginator,
	})
}


func (cc *CommentController) ChangeReportStatus(c *gin.Context){
	reportId := c.Param("id")
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_comment.ChangeReportStatus)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected AddCommentRequest, got %T", validatedData)
		return
	}

	err := cc.CommentService.ChangeReportStatus(body.Status, reportId)
	if err != nil {
		if err.Error() == "report not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"message": "گزارش یافت نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *CommentController) GetCommentsByUser(c *gin.Context){
	userinfo, exists := c.Get("user")
	
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"message": "ابتدا باید وارد شوید",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	comments, err := cc.CommentService.GetCommentsByUser(userpayload)
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