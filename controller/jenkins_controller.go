package controller

import (
	"cicd/db"
	"cicd/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type PipeStep struct {
	Type    string `json:"type" bson:"type"`
	Content string `json:"content" bson:"content"`
	Url     string `json:"url" bson:"url"`
}

type PipeNode struct {
	Width  int64       `json:"width" bson:"width"`
	Height int64       `json:"height" bson:"height"`
	Name   string      `json:"name" bson:"name"`
	Type   string      `json:"type" bson:"type"`
	Parent string      `json:"parent" bson:"parent"`
	IsLast bool        `json:"isLast" bson:"isLast"`
	Weight int64       `json:"weight" bson:"weight"`
	Steps  []*PipeStep `json:"steps" bson:"steps"`
}

type PipeEdgs struct {
	Source string `json:"source" bson:"source"`
	Target string `json:"target" bson:"target"`
}

type Pipeline struct {
	Line    string      `json:"line" bson:"line"`
	Nodes   []*PipeNode `json:"nodes" bson:"nodes"`
	Edges   []*PipeEdgs `json:"edges" bson:"edges"`
	JobName string      `json:"jobName" bson:"jobName"`
	PipeId  uint64      `json:"pipeId" bson:"pipeId"`
}

func formatPipelineStep(step *PipeStep) (script string) {

	if step.Type == "git" {
		script = "steps {git '" + step.Url + "' }"
		return script
	}

	if step.Type == "sh" {
		script = "steps {sh ' " + step.Content + "' }"
		return script
	}
	return ""
}

func formatPipelineNode(node *PipeNode) (script string) {

	steps := ""
	for _, step := range node.Steps {
		steps += formatPipelineStep(step)
	}

	script = "stage( '" + node.Name + "' ){" + steps + "}"

	return script
}

func formatPipelineScript(pipeline *Pipeline) (script string) {

	lines := strings.Split(pipeline.Line, ">")

	stages := ""
	var nodes = make(map[string]*PipeNode)
	for _, node := range pipeline.Nodes {
		nodes[node.Name] = node
	}
	for _, line := range lines {
		solts := strings.Split(line, ",")

		for _, solt := range solts {
			node, _ := nodes[solt]
			if node.Type == "header" || node.Type == "footer" || node.Type == "solt" {
				continue
			} else {
				stages += formatPipelineNode(node)
			}
		}

	}
	script = "pipeline {   docker {   image 'node'  args '-p 3000:3000'  }" + stages + "}"
	return
}

func CreatePipeline(c *gin.Context) {

	res := &Pipeline{}
	err := c.BindJSON(res)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, "参数解析失败")
		return
	}

	res.PipeId = utils.GetIntId()

	res.JobName = "pipeline-test"

	err = db.DbClient.C("pipeline").Insert(res)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, "保存失败")
		return
	}

	instance, err := utils.GetInstances(nil, "http://localhost:8080", "root", "root")

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.Error())
		return
	}

	pFolder, err := instance.Instance.GetFolder("cicd")
	if err != nil {
		pFolder, err = instance.Instance.CreateFolder("cicd")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, "文件夹创建失败")
			return
		}
	}

	config := `<?xml version="1.0" encoding="utf-8"?>
<flow-definition plugin="workflow-job"> 
  <actions> 
    <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobAction plugin="pipeline-model-definition"/> 
  </actions>  
  <description>Pipeline Job Sample</description>  
  <keepDependencies>false</keepDependencies>  
  <properties> 
    <hudson.plugins.jira.JiraProjectProperty plugin="jira"/>  
    <com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty plugin="gitlab-plugin"> 
      <gitLabConnection/> 
    </com.dabsquared.gitlabjenkins.connection.GitLabConnectionProperty> 
  </properties>  
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps"> 
    <script>`
	config += formatPipelineScript(res)
	config += `</script>  
    <sandbox>false</sandbox> 
  </definition>  
  <triggers/>  
  <disabled>false</disabled> 
</flow-definition>`

	job, err := instance.Instance.CreateJobInFolder(config, res.JobName, pFolder.GetName())

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, "job 创建失败")
		return
	}

	c.JSON(http.StatusOK, utils.Ok(job.GetName()))

}
