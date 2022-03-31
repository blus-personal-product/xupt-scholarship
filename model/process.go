package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"xupt-scholarship/db"
	"xupt-scholarship/mvc_struct"
	"xupt-scholarship/utils"
)

type ProcessModel struct {
}

func (p *ProcessModel) CreateProcess(data mvc_struct.ProcessFormData, userId string) BaseModelFmtData {
	info, _ := json.Marshal(data)
	process := db.Procedure{
		CurrentStep: []byte("[]"),
		UserId:      userId,
		Info:        info,
		History:     []byte("[]"),
	}
	result := db.Mysql.Create(&process)
	return HandleDBData(result, process.ID)
}

type ProcedureModelFormData struct {
	Id       int                        `json:"id"`
	Form     mvc_struct.ProcessFormData `json:"form"`
	UserId   string                     `json:"user_id"`
	CreateAt string                     `json:"create_at"`
	EditAt   string                     `json:"edit_at"`
}

func (p *ProcessModel) GetProcessFormData(id int) BaseModelFmtData {
	var processData mvc_struct.ProcessFormData
	var processInfo db.Procedure
	var result *gorm.DB
	if id == -1 {
		result = db.Mysql.Last(&processInfo)
	} else {
		result = db.Mysql.First(&processInfo, id)
	}
	json.Unmarshal(processInfo.Info, &processData)
	return HandleDBData(result, ProcedureModelFormData{
		Id:       id,
		Form:     processData,
		UserId:   processInfo.UserId,
		CreateAt: utils.FmtTimeByUnix(processInfo.CreateAt),
		EditAt:   utils.FmtTimeByUnix(processInfo.UpdateAt),
	})
}

type CurrentYearProcessData struct {
	History []mvc_struct.ProcessHistoryItem `json:"history"`
	UserId  string                          `json:"user_id"`
	Id      int                             `json:"id"`
}

type ProcessStatusRes struct {
	Status     string `json:"status"`
	ProcessId  int    `json:"process_id"`
	Editable   bool   `json:"editable"`
	Createable bool   `json:"createable"`
}

// GetCurrentYearProcess 获取当前学年的评定流程
func (p *ProcessModel) GetCurrentYearProcess(userId string, identity string) BaseModelFmtData {
	var procedure db.Procedure
	var stepHistory []mvc_struct.ProcessHistoryItem
	var processInfo mvc_struct.ProcessFormData
	yearTime := GetCurrentYear("")
	result := db.Mysql.Where("create_at > ?", yearTime).First(&procedure)
	status := "not_create"
	isLate := false
	procedureId := -1
	creatorId := userId
	if result.Error == nil {
		procedureId = procedure.ID
		creatorId = procedure.UserId
		json.Unmarshal(procedure.History, &stepHistory)
		json.Unmarshal(procedure.Info, &processInfo)
		isLate = GetIsLate(processInfo.Form.IndividualApplicationStage.Date[0])
		if isLate {
			status = "pre_start"
		}
		if len(stepHistory) > 0 {
			status = "opened"
		}
	}
	res := ProcessStatusRes{
		Status:     status,
		ProcessId:  procedureId,
		Editable:   userId == creatorId && (!isLate),
		Createable: identity == "manager" && creatorId == userId && procedureId == -1,
	}
	return HandleDBData(result, res)
}

func (p *ProcessModel) UpdateProcessFormData(id int, info mvc_struct.ProcessFormData) BaseModelFmtData {
	jsonInfo, _ := json.Marshal(info)
	var procedure db.Procedure
	result := db.Mysql.Model(&procedure).Where("id = ?", id).Update("info", jsonInfo)
	return HandleDBData(result, id)
}

type stepData struct {
	History []mvc_struct.ProcessHistoryItem `json:"history"`
	Current mvc_struct.ProcessHistoryItem   `json:"current"`
}

func (p *ProcessModel) GetProcessStep(id int) BaseModelFmtData {
	var procedure db.Procedure
	var stepHistory []mvc_struct.ProcessHistoryItem
	var currentStep mvc_struct.ProcessHistoryItem
	result := db.Mysql.First(&procedure, id)
	json.Unmarshal(procedure.CurrentStep, &currentStep)
	json.Unmarshal(procedure.History, &stepHistory)
	return HandleDBData(result, stepData{
		History: stepHistory,
		Current: currentStep,
	})
}