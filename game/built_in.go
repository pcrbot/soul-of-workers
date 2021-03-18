package game

type playerTag int32

const (
	playerTag工作狂人 playerTag = iota
	playerTag精力充沛
	playerTag人格魅力
	playerTag慧眼识珠
	playerTag商业头脑
	playerTag职场精英
	playerTag随机应变
	playerTag理财达人
	playerTag粗心大意
	playerTag形单影只
	playerTag榆木脑袋
	playerTag笨嘴拙舌
)

var (
	playerTagsPositive = []playerTag{playerTag工作狂人, playerTag精力充沛, playerTag人格魅力, playerTag慧眼识珠, playerTag商业头脑, playerTag职场精英, playerTag随机应变, playerTag理财达人}
	playerTagsNegative = []playerTag{playerTag粗心大意, playerTag形单影只, playerTag榆木脑袋, playerTag笨嘴拙舌}
	playerTagsNameMap  = []string{
		playerTag工作狂人: "工作狂人",
		playerTag精力充沛: "精力充沛",
		playerTag人格魅力: "人格魅力",
		playerTag慧眼识珠: "慧眼识珠",
		playerTag商业头脑: "商业头脑",
		playerTag职场精英: "职场精英",
		playerTag随机应变: "随机应变",
		playerTag理财达人: "理财达人",
		playerTag粗心大意: "粗心大意",
		playerTag形单影只: "形单影只",
		playerTag榆木脑袋: "榆木脑袋",
		playerTag笨嘴拙舌: "笨嘴拙舌",
	}
)

var npcCompanies = []NpcCompany{
	{
		// ID: -1
		Name:                   "水泥工",
		Salary:                 5,
		ProlificacyRequirement: 50,
	},
	{
		// ID: -2
		Name:                   "会所嫩模",
		Salary:                 5,
		ProlificacyRequirement: 50,
	},
	{
		// ID: -3
		Name:                   "清洁工",
		Salary:                 6,
		ProlificacyRequirement: 60,
	},
	{
		// ID: -4
		Name:                   "快递员",
		Salary:                 7,
		ProlificacyRequirement: 75,
	},
	{
		// ID: -5
		Name:                   "保安",
		Salary:                 8,
		ProlificacyRequirement: 95,
	},
	{
		// ID: -6
		Name:                   "司机",
		Salary:                 9,
		ProlificacyRequirement: 120,
	},
	{
		// ID: -7
		Name:                   "程序员",
		Salary:                 10,
		ProlificacyRequirement: 160,
	},
}

var companyScales = []CompanyScale{
	{
		// ID: 0
		Name: "擦鞋小队",
		Cost: 100,
		Efficiency: func(i int32) float32 {
			return 1.0 - (float32(i) * 0.2)
		},
	},
	{
		// ID: 1
		Name: "施工队",
		Cost: 300,
		Efficiency: func(i int32) float32 {
			return 0.95 - (float32(i) * 0.1)
		},
	},
	{
		// ID: 2
		Name: "甜品店",
		Cost: 1000,
		Efficiency: func(i int32) float32 {
			return 0.93 - (float32(i) * 0.05)
		},
	},
	{
		// ID: 3
		Name: "杂货铺",
		Cost: 4000,
		Efficiency: func(i int32) float32 {
			return 0.90 - (float32(i) * 0.03)
		},
	},
	{
		// ID: 4
		Name: "超市",
		Cost: 20_000,
		Efficiency: func(i int32) float32 {
			return 0.85 - (float32(i) * 0.015)
		},
	},
	{
		// ID: 5
		Name: "经销商",
		Cost: 100_000,
		Efficiency: func(i int32) float32 {
			return 0.80 - (float32(i) * 0.008)
		},
	},
	{
		// ID: 6
		Name: "连锁餐饮店",
		Cost: 1_000_000,
		Efficiency: func(i int32) float32 {
			return 0.70 - (float32(i) * 0.005)
		},
	},
	{
		// ID: 7
		Name: "房地产集团",
		Cost: 20_000_000,
		Efficiency: func(i int32) float32 {
			return 0.60 - (float32(i) * 0.004)
		},
	},
	{
		// ID: 8
		Name: "国际商业巨头",
		Cost: 500_000_000,
		Efficiency: func(i int32) float32 {
			return 0.50 - (float32(i) * 0.0035)
		},
	},
}

func (t playerTag) String() string {
	return playerTagsNameMap[t]
}

// GetNpcCompany 通过ID（负数）获得NPC公司
func GetNpcCompany(id int64) NpcCompany {
	return npcCompanies[^id]
}
