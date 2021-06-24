package route

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/util"
	"github.com/gokins-main/gokins/yml"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type ApiController struct{}

func (ApiController) GetPath() string {
	return "/api"
}
func (c *ApiController) Routes(g gin.IRoutes) {
	g.POST("/builds", util.GinReqParseJson(c.test))
}
func (ApiController) test(c *gin.Context) {
	all, err := ioutil.ReadAll(c.Request.Body)
	y := &yml.YML{}
	err = yaml.Unmarshal(all, y)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	marshal, err := yaml.Marshal(y)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	//TODO insert db
	b := &runtime.Build{}
	err = yaml.Unmarshal(marshal, b)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	err = prebuild(b)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	engine.Mgr.BuildEgn().Put(b)
	c.JSON(200, gin.H{
		"msg": b,
	})
}

func prebuild(b *runtime.Build) error {
	if b == nil {
		return errors.New("build is empty")
	}
	if b.Stages == nil || len(b.Stages) <= 0 {
		return errors.New("stages is empty")
	}
	pipelineId := utils.NewXid()
	buildId := utils.NewXid()
	b.Id = buildId
	b.Repo = &runtime.Repository{
		Name:     "SuperHeroJim",
		Token:    "",
		Sha:      "",
		CloneURL: "https://gitee.com/SuperHeroJim/gokins-test.git",
	}
	for _, stage := range b.Stages {
		stage.Id = utils.NewXid()
		stage.PipelineId = pipelineId
		stage.BuildId = buildId
		for _, step := range stage.Steps {
			step.Id = utils.NewXid()
			step.StageId = stage.Id
			step.BuildId = buildId
		}
	}
	return nil
}
