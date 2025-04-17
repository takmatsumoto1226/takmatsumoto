package ftn

import "github.com/gin-gonic/gin"

type FTNListRequest struct {
	StartIndex int    `json:"startidx"`
	EndIndex   int    `json:"endindex"`
	StartDate  string `json:"startdate"`
	EndDate    string `json:"enddate"`
}

func FTNListCtx(c *gin.Context) {

}
