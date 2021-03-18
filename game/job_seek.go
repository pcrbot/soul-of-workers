package game

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"gorm.io/gorm"
)

func seekJob(c *zero.Ctx) error {
	var p Player
	if err := db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return err
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return nil
	}
	var offset int64
	prolificacy := p.GetProlificacy()
	switch {
	case prolificacy < 60:
		offset = 0
	case prolificacy < 75:
		offset = 1
	case prolificacy < 95:
		offset = 2
	case prolificacy < 120:
		offset = 3
	default:
		offset = 4
	}
	job := [3]NpcCompany{
		npcCompanies[offset],
		npcCompanies[offset+1],
		npcCompanies[offset+2],
	}
	message := "您可以选择的工作有："
	for i, company := range job {
		message += fmt.Sprintf("\n%d：%s（要求生产力：%d，薪水：%d）", i+1, company.Name, company.ProlificacyRequirement, company.Salary)
	}
	message += "\n0：结束"
	c.Send(c.Send(message))
	prompting := true
	var answer int64
	for prompting {
		reply := c.Get("")
		switch reply {
		case "1":
			answer = 0
			prompting = false
		case "2":
			answer = 1
			prompting = false
		case "3":
			answer = 2
			prompting = false
		case "0":
			c.Send("已结束")
			return nil
		default:
			c.Send("请发送要选择的工作序号，或发送“0”结束")
		}
	}
	application := job[answer]
	if prolificacy < application.ProlificacyRequirement {
		// 生产力不够
		c.Send("你的申请被拒绝了，还是找找其他的岗位吧")
		return nil
	}
	lock := roleLock.Get(p.ID)
	lock.Lock()
	defer lock.Unlock()
	p.EmploymentCompanyOwnerID = ^(offset + answer)
	if err := db.Save(&p).Error; err != nil {
		return err
	}
	c.Send(zeroText(fmt.Sprintf("申请成功，您已成为一名%s，现在的薪水是：每周期%dD币", application.Name, application.Salary)))
	return nil
}

func joinCompany(c *zero.Ctx) error {
	companyOwnerID, ok := findAt(c.Event.Message)
	if !ok {
		c.Send("需要at一名玩家才能加入别人的公司")
		return nil
	}
	var p, owner Player
	var b Company
	var app CompanyApplication
	if err := db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return err
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return nil
	}
	if err := db.Limit(1).Find(&b, companyOwnerID).Error; err != nil {
		return err
	}
	if b.OwnerID == 0 {
		c.Send("这名玩家还没有创建公司呢")
		return nil
	}
	lock := roleLock.Get(p.ID)
	lock.Lock()
	defer lock.Unlock()
	if err := db.First(&owner, companyOwnerID).Error; err != nil {
		return err
	}
	if err := db.Limit(1).Find(&app, c.Event.UserID).Error; err != nil {
		return err
	}
	var plusMessage string
	if app.Applicant != 0 {
		// 有未完成的申请
		var ex Company
		if err := db.First(&ex, app.CompanyOwnerID).Error; err != nil {
			return err
		}
		plusMessage = fmt.Sprintf("\n同时，您向%s的申请已取消", ex.Name)
	}
	app = CompanyApplication{
		Applicant:      c.Event.UserID,
		CompanyOwnerID: companyOwnerID,
	}
	if err := db.Save(&app).Error; err != nil {
		return err
	}
	message := fmt.Sprintf("您已申请加入%s，请等待%s的回复", b.Name, owner.Nickname)
	c.Send(zeroText(message + plusMessage))
	return nil
}

func resign(c *zero.Ctx) error {
	var p Player
	if err := db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return err
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return nil
	}
	lock := roleLock.Get(p.ID)
	lock.Lock()
	defer lock.Unlock()
	if p.EmploymentCompanyOwnerID == 0 {
		c.Send("您还没有加入公司呢")
		return nil
	} else if p.EmploymentCompanyOwnerID < 0 {
		// NPC 的公司
		p.EmploymentCompanyOwnerID = 0
		company := GetNpcCompany(p.EmploymentCompanyOwnerID)
		if err := db.Save(&p).Error; err != nil {
			return err
		}
		c.Send(zeroText(fmt.Sprintf("你放弃了%s的工作，现在你是自由的了", company.Name)))
	} else if p.EmploymentCompanyOwnerID > 0 {
		// 玩家的公司
		err := db.Transaction(func(tx *gorm.DB) (e error) {
			var b Company
			if e = tx.First(&b, p.EmploymentCompanyOwnerID).Error; e != nil {
				return
			}
			p.EmploymentCompanyOwnerID = 0
			b.EmployeeCounts -= 1
			if e = tx.Save(&p).Error; e != nil {
				return
			}
			if e = tx.Save(&b).Error; e != nil {
				return
			}
			return
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Player) ChangeCompany(companyOwnerID int64) error {
	lock := roleLock.Get(p.ID)
	lock.Lock()
	defer lock.Unlock()
	return db.Transaction(func(tx *gorm.DB) error {
		if p.EmploymentCompanyOwnerID > 0 {
			var ex Company
			if err := tx.First(&ex, p.EmploymentCompanyOwnerID).Error; err != nil {
				return err
			}
			ex.EmployeeCounts -= 1
			ex.EmployeeProlificacy -= p.Prolificacy
			if p.PlayerTag1 == playerTag职场精英 {
				ex.EmployeeEfficiency -= 5000
			}
			if err := tx.Save(&ex).Error; err != nil {
				return err
			}
		}
		var b Company
		if err := tx.First(&b, companyOwnerID).Error; err != nil {
			return err
		}
		b.EmployeeCounts += 1
		b.EmployeeProlificacy += p.Prolificacy
		if p.PlayerTag1 == playerTag职场精英 {
			b.EmployeeEfficiency += 5000
		}
		if err := tx.Save(&b).Error; err != nil {
			return err
		}
		p.EmploymentCompanyOwnerID = companyOwnerID
		if err := tx.Save(&p).Error; err != nil {
			return err
		}
		return nil
	})
}
