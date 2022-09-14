package service

import (
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
)

type TemplateService struct{}

func init() {
	lib.QueryTemplate = func(name lib.EventName) []model.CoreTemplate {
		return TemplateService{}.ListByEvent(string(name))
	}
}

func (TemplateService) Create(info *model.CoreTemplate) bool {
	return model.DB().Create(info).RowsAffected > 0
}

func (TemplateService) Modify(info *model.CoreTemplate) bool {
	return 0 != info.ID && model.DB().Model(info).Where(&model.CoreTemplate{ID: info.ID}).Update(info).RowsAffected > 0
}

func (TemplateService) Status(id uint, status int) bool {
	return 0 != id && model.DB().Where(&model.CoreTemplate{ID: id}).Update("status", status).RowsAffected > 0
}

func (TemplateService) Page(start, end int) (count int, infos []model.CoreTemplate) {
	model.DB().Model(&model.CoreTemplate{}).Order("id desc").Count(&count).Offset(start).Limit(end).Find(&infos)
	return
}

func (TemplateService) InfoById(id uint) *model.CoreTemplate {
	if 0 == id {
		return nil
	}

	info := &model.CoreTemplate{}
	if model.DB().Model(info).Where(&model.CoreTemplate{ID: id}).Take(info).RowsAffected > 0 {
		return info
	}

	return nil
}

func (TemplateService) InfoByAlias(alias string) *model.CoreTemplate {
	if "" == alias {
		return nil
	}

	info := &model.CoreTemplate{}
	if model.DB().Where(&model.CoreTemplate{Alias: alias}).Take(info).RowsAffected > 0 {
		return info
	}

	return nil
}

func (TemplateService) ListByEvent(event string) []model.CoreTemplate {
	if "" == event {
		return nil
	}

	var result []model.CoreTemplate
	model.DB().Model(&model.CoreTemplate{}).Where("`event` like ? and `status` = ?", "%"+event+"%", 1).Find(&result)
	return result
}
