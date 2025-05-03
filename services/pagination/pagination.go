package services_paginator

import (
    "github.com/gin-gonic/gin"
    "github.com/rosberry/go-pagination"
    "gorm.io/gorm"
)


func Paginate[T any](objects *gorm.DB, ctx *gin.Context, db *gorm.DB, limit int) ([]T, *pagination.Paginator, error) {
    paginator, err := pagination.New(pagination.Options{
        GinContext:    ctx,
        DB:           db,
        Model:        new(T),
        Limit:        uint(limit),
        DefaultCursor: nil,
    })
    if err != nil {
        return nil, nil, err
    }

    var results []T
    if err := paginator.Find(objects, &results); err != nil {
        return nil, nil, err
    }

    return results, paginator, nil
}