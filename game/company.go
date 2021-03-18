package game

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	zeroMessage "github.com/wdvxdr1123/ZeroBot/message"
	"gorm.io/gorm"
)

func createCompany(c *zero.Ctx) (err error) {
	var p Player
	if err = db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return
	}
	var b Company
	if err = db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return
	}
	if b.OwnerID != 0 {
		c.Send("您已经创建过公司了")
		return
	}
	message := fmt.Sprintf("请选择您要创建的公司\n1：%s（%dD币）\n2：%s（%dD币）\n3：%s（%dD币）\n0：结束", companyScales[0].Name, companyScales[0].Cost, companyScales[1].Name, companyScales[1].Cost, companyScales[2].Name, companyScales[2].Cost)
	c.Send(zeroText(message))
	prompting := true
	var answer int32
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
			c.Send("请发送要选择的公司序号，或发送“0”结束")
		}
	}
	if p.DCoin < companyScales[answer].Cost {
		c.Send("您的资金不够创建这个公司哦")
		return
	}
	var playerEfficiency int32 = 0
	if p.PlayerTag1 == playerTag商业头脑 {
		playerEfficiency += 2_000
	}
	if p.EmploymentCompanyOwnerID < 0 {
		// NPC 公司均为加班公司
		playerEfficiency -= 10_000
	} else if p.EmploymentCompanyOwnerID > 0 {
		var employmentCompany Company
		if err = db.First(&employmentCompany, p.EmploymentCompanyOwnerID).Error; err != nil {
			return
		}
		if employmentCompany.Overtime {
			playerEfficiency -= 10_000
		}
	}
	err = db.Transaction(func(tx *gorm.DB) (e error) {
		p.DCoin -= companyScales[answer].Cost
		b = Company{
			OwnerID:             p.ID,
			Name:                p.Nickname + "的公司",
			ScaleID:             answer,
			SalarySetting:       0,
			Overtime:            false,
			RecruitmentProgress: 180,
			EmployeeCounts:      0,
			EmployeeProlificacy: 0,
			EmployeeEfficiency:  0,
			OwnerEfficiency:     playerEfficiency,
		}
		if e = tx.Save(&p).Error; e != nil {
			return
		}
		if e = tx.Save(&b).Error; e != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	c.Send("创建成功")
	return
}

func companyRecruit(c *zero.Ctx) (err error) {
	var p Player
	if err = db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return
	}
	var b Company
	if err = db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return
	}
	if b.OwnerID == 0 {
		c.Send("您还没有创建公司呢")
		return
	}
	// 玩家
	var apply CompanyApplication
	if err = db.Where(&CompanyApplication{CompanyOwnerID: c.Event.UserID}).Limit(1).Find(&apply, c.Event.UserID).Error; err != nil {
		return
	}
	if apply.Applicant != 0 {
		var applicant Player
		if err = db.First(&applicant, apply.Applicant).Error; err != nil {
			return
		}
		message := fmt.Sprintf(`您的公司有新的应聘者：
%s
生产力：%d
人格特点：%s，%s

录用后预计公司收益：%d

请选择是否录用
1：录用
2：拒绝
0：结束`, applicant.Nickname, applicant.Prolificacy, applicant.PlayerTag1, applicant.PlayerTag2, b.EstimateIncomeWithPlayer(&applicant))
		c.Send(append(zeroUserAvatar(apply.Applicant), zeroText(message)...))
		for {
			reply := c.Get("")
			switch reply {
			case "1":
				if err = applicant.ChangeCompany(c.Event.UserID); err != nil {
					return
				}
				c.Send(zeroText(fmt.Sprintf("您已成功录用%s", applicant.Nickname)))
				if p.PlayerTag1 == playerTag人格魅力 {
					b.RecruitmentProgress += 100
					if e := db.Save(&b).Error; e != nil {
						log.Warning(e)
					}
				}
				return
			case "2":
				if err = db.Delete(&apply).Error; err != nil {
					return
				}
				c.Send("已拒绝")
				return
			case "0":
				c.Send("已结束")
				return
			default:
				c.Send("请选择是否录用，或发送“0”结束")
			}
		}
	}

	// NPC
	if b.SalarySetting <= 0 {
		c.Send("你的公司还没有设置薪水呢")
		return
	}
	if b.RecruitmentProgress < 100 {
		// 招募进度不足
		c.Send("还没有人应聘您的岗位哦")
		return
	}
	b.RecruitmentProgress -= 100
	if err = db.Save(&b).Error; err != nil {
		return
	}
	npc, err := FindNpc(b.SalarySetting, &p)
	if err != nil {
		return
	}
	message := fmt.Sprintf(`您的公司有新的应聘者：
%s
生产力：%d
工资期望：%d
团队效率：%s%%
流动率：%s%%

录用后预计公司收益：%d

请选择是否录用：
1：录用
2：拒绝
0：拒绝`, npc.Name, npc.Prolificacy, npc.SalaryExpectation, ThousandthStr(npc.Teamwork), TenthStr(npc.Turnover), b.EstimateIncomeWithNpc(npc))
	c.Send(append(zeroMessage.Message{zeroMessage.Image(npc.GetAvatarPath())}, zeroText(message)...))
	for {
		reply := c.Get("")
		switch reply {
		case "1":
			return db.Transaction(func(tx *gorm.DB) error {
				lock := roleLock.Get(npc.ID)
				lock.Lock()
				defer lock.Unlock()
				npc, e := FindNpc(b.SalarySetting, &p)
				if e != nil {
					return e
				}
				if npc.EmploymentCompanyOwnerID != 0 {
					// 并发冲突了
					c.Send("来晚了一步哦，应聘者已经走了")
					return nil
				}
				npc.EmploymentCompanyOwnerID = c.Event.UserID
				b.EmployeeCounts += 1
				b.EmployeeProlificacy += npc.Prolificacy
				b.EmployeeEfficiency += npc.Teamwork
				if e := tx.Save(&npc).Error; e != nil {
					return e
				}
				if e := tx.Save(&b).Error; e != nil {
					return e
				}
				if p.PlayerTag1 == playerTag人格魅力 {
					b.RecruitmentProgress += 100
					if e := tx.Save(&b).Error; e != nil {
						return e
					}
				}
				availableNpcList.Remove(npc.ID)
				c.Send(zeroText(fmt.Sprintf("您已成功录用%s", npc.Name)))
				return nil
			})
		case "2", "0":
			c.Send("已拒绝")
			return
		default:
			c.Send("请选择是否录用，或发送“0”结束")
		}
	}
}

func companyRename(c *zero.Ctx) error {
	var b Company
	if err := db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return err
	}
	if b.OwnerID == 0 {
		c.Send("您还没有创建公司呢")
		return nil
	}
	newName := c.State["args"].(string)
	if newName == "" {
		newName = c.Get("请发送新的公司名称")
	}
	lock := companyLock.Get(c.Event.UserID)
	lock.Lock()
	defer lock.Unlock()
	if err := db.Model(&b).Update("name", newName).Error; err != nil {
		return err
	}
	c.Send("改名完成，您的新公司名为：\n" + newName)
	return nil
}

func companySetSalary(c *zero.Ctx) error {
	var b Company
	if err := db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return err
	}
	if b.OwnerID == 0 {
		c.Send("您还没有创建公司呢")
		return nil
	}
	newSalary := c.State["args"].(string)
	if newSalary == "" {
		newSalary = c.Get("请发送新的工资设置（整数）")
	}
	salary, err := strconv.ParseInt(newSalary, 10, 64)
	if err != nil {
		c.Send("工资必须是整数哦")
		return nil
	}
	lock := companyLock.Get(c.Event.UserID)
	lock.Lock()
	defer lock.Unlock()
	if err = db.Model(&b).Update("salary_setting", salary).Error; err != nil {
		return err
	}
	c.Send("设置成功")
	return nil
}

func companySetOvertime(c *zero.Ctx) error {
	var b Company
	if err := db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return err
	}
	if b.OwnerID == 0 {
		c.Send("您还没有创建公司呢")
		return nil
	}
	if b.Overtime {
		c.Send("您的员工们已经在加班了")
		return nil
	}
	lock := companyLock.Get(c.Event.UserID)
	lock.Lock()
	defer lock.Unlock()
	if err := db.Model(&b).Update("overtime", true).Error; err != nil {
		return err
	}
	efficiency := companyScales[b.ScaleID].Efficiency(b.EmployeeCounts)
	efficiency += 0.00001 * float32(b.EmployeeEfficiency+b.OwnerEfficiency)
	c.Send(fmt.Sprintf("设置成功，您的公司效率上升了 10%%，现在是%.1f%%", efficiency*100))
	return nil
}

func companyUnsetOvertime(c *zero.Ctx) error {
	var b Company
	if err := db.Limit(1).Find(&b, c.Event.UserID).Error; err != nil {
		return err
	}
	if b.OwnerID == 0 {
		c.Send("您还没有创建公司呢")
		return nil
	}
	if !b.Overtime {
		c.Send("您的员工们并没有加班哦")
		return nil
	}
	lock := companyLock.Get(c.Event.UserID)
	lock.Lock()
	defer lock.Unlock()
	if err := db.Model(&b).Update("overtime", false).Error; err != nil {
		return err
	}
	efficiency := companyScales[b.ScaleID].Efficiency(b.EmployeeCounts)
	efficiency += 0.00001 * float32(b.EmployeeEfficiency+b.OwnerEfficiency)
	c.Send(fmt.Sprintf("设置成功，您的公司效率现在是%.1f%%", efficiency*100))
	return nil
}

func (b *Company) EstimateIncome() int64 {
	efficiency := companyScales[b.ScaleID].Efficiency(b.EmployeeCounts)
	efficiency += 0.00001 * float32(b.EmployeeEfficiency+b.OwnerEfficiency)
	volume := int64(float32(b.EmployeeProlificacy) * efficiency)
	salaryExpenditure := int64(b.EmployeeCounts) * b.SalarySetting
	return volume - salaryExpenditure
}

func (b *Company) EstimateIncomeWithPlayer(p *Player) int64 {
	efficiency := companyScales[b.ScaleID].Efficiency(b.EmployeeCounts + 1)
	efficiency += 0.00001 * float32(b.EmployeeEfficiency+b.OwnerEfficiency)
	if p.PlayerTag1 == playerTag职场精英 {
		efficiency += 0.05
	}
	volume := int64(float32(b.EmployeeProlificacy+p.Prolificacy) * efficiency)
	salaryExpenditure := int64(b.EmployeeCounts+1) * b.SalarySetting
	return volume - salaryExpenditure
}

func (b *Company) EstimateIncomeWithNpc(c *Npc) int64 {
	efficiency := companyScales[b.ScaleID].Efficiency(b.EmployeeCounts + 1)
	efficiency += 0.00001 * float32(b.EmployeeEfficiency+b.OwnerEfficiency)
	volume := int64(float32(b.EmployeeProlificacy+c.Prolificacy) * efficiency)
	salaryExpenditure := int64(b.EmployeeCounts+1) * b.SalarySetting
	return volume - salaryExpenditure
}
