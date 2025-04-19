package controllers_comment

import(
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