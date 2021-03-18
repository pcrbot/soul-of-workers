package game

import (
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func RegisterToZeroBot(command string, h func(*zero.Ctx) error) {
	zero.OnCommand(command).Handle(func(c *zero.Ctx) {
		defer func() {
			if r := recover(); r != nil {
				c.Send("＞﹏＜   发生了内部错误，请查看错误日志")
				log.Error(r)
			}
		}()
		err := h(c)
		if err != nil {
			log.Error(err)
			c.Send("＞﹏＜   发生了内部错误，请查看错误日志")
		}
	})
}

func RegisterCommands() error {
	RegisterToZeroBot("创建角色", createRole)
	RegisterToZeroBot("角色状态", roleStatus)
	RegisterToZeroBot("角色改名", roleRename)
	RegisterToZeroBot("找工作", seekJob)
	RegisterToZeroBot("加入公司", joinCompany)
	RegisterToZeroBot("离职", resign)
	RegisterToZeroBot("创建公司", createCompany)
	RegisterToZeroBot("角色改名", companyRename)
	RegisterToZeroBot("招募员工", companyRecruit)
	RegisterToZeroBot("设置工资", companySetSalary)
	//RegisterToZeroBot("", )
	return nil
}
