package game

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type dbData interface {
	Save() error
	SaveTo(*gorm.DB) error
}

type Player struct {
	ID                       int64 // QQ号
	Nickname                 string
	DCoin                    int64
	Prolificacy              int32
	Education                int32
	Relationship             int32
	EmploymentCompanyOwnerID int64 // 正数：玩家，负数：NPC，零：无
	PlayerTag1               playerTag
	PlayerTag2               playerTag
}

type Company struct {
	OwnerID             int64 `gorm:"primaryKey"`
	Name                string
	ScaleID             int32
	SalarySetting       int64
	Overtime            bool
	RecruitmentProgress int32
	EmployeeCounts      int32
	EmployeeProlificacy int32
	EmployeeEfficiency  int32 // 单位： 0.001%
	OwnerEfficiency     int32 // 单位： 0.001%
}

type CompanyScale struct {
	Name       string
	Cost       int64
	Efficiency func(int32) float32
}

type CompanyApplication struct {
	Applicant      int64 `gorm:"primaryKey"`
	CompanyOwnerID int64 `gorm:"index"`
}

type Npc struct {
	ID                       int64 // 负数
	Name                     string
	Discovered               bool // 初始为false，属性确定后变为true
	Prolificacy              int32
	SalaryExpectation        int64
	Education                int32
	Teamwork                 int32 // 单位： 0.001%
	Loyalty                  int32
	Turnover                 int32 // 单位： 0.1%
	EmploymentCompanyOwnerID int64 `gorm:"index"`
}

type NpcNickname struct {
	NickName string `gorm:"primaryKey"`
	ID       int64
}

type NpcCompany struct {
	Name                   string
	Salary                 int32
	ProlificacyRequirement int32
}

func (p *Player) Save() error {
	return db.Save(p).Error
}

func (b *Company) Save() error {
	return db.Save(b).Error
}

func (a *CompanyApplication) Save() error {
	return db.Save(a).Error
}

func (c *Npc) Save() error {
	return db.Save(c).Error
}

func (n *NpcNickname) Save() error {
	return db.Save(n).Error
}

func (p *Player) SaveTo(vdb *gorm.DB) error {
	return vdb.Save(p).Error
}

func (b *Company) SaveTo(vdb *gorm.DB) error {
	return vdb.Save(b).Error
}

func (a *CompanyApplication) SaveTo(vdb *gorm.DB) error {
	return vdb.Save(a).Error
}

func (c *Npc) SaveTo(vdb *gorm.DB) error {
	return vdb.Save(c).Error
}

func (n *NpcNickname) SaveTo(vdb *gorm.DB) error {
	return vdb.Save(n).Error
}

// InitialDB 初始化数据库
func InitialDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("sow_data.sqlite"), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return err
	}
	if err = db.AutoMigrate(&Player{}, &Company{}, &CompanyApplication{}, &Npc{}, &NpcNickname{}); err != nil {
		return err
	}
	if err = initialAvailableNpcList(); err != nil {
		return err
	}
	return nil
}

func SaveAll(data ...dbData) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, d := range data {
			err := d.SaveTo(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Player) GetProlificacy() int32 {
	prolificacy := p.Prolificacy
	if p.PlayerTag1 == playerTag工作狂人 {
		// 提高 50%
		prolificacy = prolificacy + prolificacy/2
	}
	return prolificacy
}

func (b *Company) GetScaleName() string {
	return companyScales[b.ScaleID].Name
}
