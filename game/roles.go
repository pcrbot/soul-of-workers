package game

import (
	"fmt"
	"math/rand"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func createRole(c *zero.Ctx) error {
	var player Player
	if err := db.Limit(1).Find(&player, c.Event.UserID).Error; err != nil {
		return err
	}
	if player.ID != 0 {
		c.Send("您已经创建过角色了")
		return nil
	}
	player = Player{
		ID:                       c.Event.UserID,
		Nickname:                 c.Event.Sender.NickName,
		DCoin:                    100,
		Prolificacy:              50,
		Education:                0,
		Relationship:             10,
		EmploymentCompanyOwnerID: 0,
	}
	tagList := randomSixTags()
	prompt := fmt.Sprintf("请在以下角色中选择一个初始角色：\n\n1：%s，%s\n2：%s，%s\n3：%s，%s\n\n9：重选\n0：结束", tagList[0], tagList[1], tagList[2], tagList[3], tagList[4], tagList[5])
	c.Send(zeroText(prompt))
	prompting := true
	for prompting {
		reply := c.Get("")
		switch reply {
		case "1":
			player.PlayerTag1 = tagList[0]
			player.PlayerTag2 = tagList[1]
			prompting = false
		case "2":
			player.PlayerTag1 = tagList[2]
			player.PlayerTag2 = tagList[3]
			prompting = false
		case "3":
			player.PlayerTag1 = tagList[4]
			player.PlayerTag2 = tagList[5]
			prompting = false
		case "9":
			tagList = randomSixTags()
			prompt = fmt.Sprintf("1：%s，%s\n2：%s，%s\n3：%s，%s\n\n9：重选\n0：结束", tagList[0], tagList[1], tagList[2], tagList[3], tagList[4], tagList[5])
			c.Send(zeroText(prompt))
		case "0":
			c.Send("已结束创建")
			return nil
		default:
			c.Send("请发送要选择的角色序号，或发送“0”结束")
		}
	}
	if err := db.Create(&player).Error; err != nil {
		return err
	}
	c.Send("创建成功")
	return nil
}

func roleStatus(c *zero.Ctx) error {
	var p Player
	if err := db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return err
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return nil
	}
	var occupied string
	if p.EmploymentCompanyOwnerID == 0 {
		occupied = "无业"
	} else if p.EmploymentCompanyOwnerID > 0 {
		var company Company
		var boss Player
		if err := db.First(&company, p.EmploymentCompanyOwnerID).Error; err != nil {
			return err
		}
		if err := db.First(&boss, p.EmploymentCompanyOwnerID).Error; err != nil {
			return err
		}
		occupied = fmt.Sprintf("在%s的%s工作", boss.Nickname, company.GetScaleName())
	} else if p.EmploymentCompanyOwnerID < 0 {
		company := GetNpcCompany(p.EmploymentCompanyOwnerID)
		occupied = fmt.Sprintf("在%s（NPC）工作", company.Name)
	}
	message := fmt.Sprintf("%s的状态：\nＤ币：　　%d\n生产力：　%d\n教育：　　%d级\n人缘：　　%d\n所在公司：%s\n人格特点：%s，%s", p.Nickname, p.DCoin, p.GetProlificacy(), p.Education, p.Relationship, occupied, p.PlayerTag1, p.PlayerTag2)
	c.Send(zeroText(message))
	return nil
}

func roleRename(c *zero.Ctx) error {
	var p Player
	if err := db.Limit(1).Find(&p, c.Event.UserID).Error; err != nil {
		return err
	}
	if p.ID == 0 {
		c.Send("您还没有创建角色呢")
		return nil
	}
	newName := c.State["args"].(string)
	if newName == "" {
		newName = c.Get("请发送新的角色昵称")
	}
	p.Nickname = newName
	if err := db.Save(&p).Error; err != nil {
		return err
	}
	c.Send("改名完成，您的新昵称为：\n" + newName)
	return nil
}

func randomSixTags() (list [6]playerTag) {
	for i := 0; i < 3; i++ {
		list[2*i] = playerTagsPositive[rand.Intn(len(playerTagsPositive))]
		list[2*i+1] = playerTagsNegative[rand.Intn(len(playerTagsNegative))]
	}
	return
}
