package game

import (
	"math/rand"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

type availableI64List struct {
	lock sync.RWMutex
	list []int64
}

func (l *availableI64List) RandomChoice() int64 {
	l.lock.RLock()
	defer l.lock.RUnlock()
	k := rand.Intn(len(l.list))
	return l.list[k]
}

func (l *availableI64List) Add(k int64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	index := sort.Search(len(l.list), func(i int) bool {
		return l.list[i] >= k
	})
	l.list = append(append(l.list[0:index], k), l.list[index:]...)
}

func (l *availableI64List) Remove(k int64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	index := sort.Search(len(l.list), func(i int) bool {
		return l.list[i] >= k
	})
	l.list = append(l.list[0:index-1], l.list[index:]...)
}

var availableNpcList availableI64List

func initialAvailableNpcList() error {
	var npcIds []struct {
		ID int64
	}
	if err := db.Where("discovered = 0").Find(&npcIds).Error; err != nil {
		return err
	}
	list := make([]int64, len(npcIds))
	for i, id := range npcIds {
		list[i] = id.ID
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	availableNpcList.list = list
	return nil
}

func RandomAvailableNpc() (*Npc, error) {
	var c Npc
	i:=availableNpcList.RandomChoice()
	err:=db.First(&c,i).Error
	return &c, err
}

func FindNpc(maxSalary int64, player *Player) (*Npc, error) {
	c, err := RandomAvailableNpc()
	if err != nil {
		return nil, err
	}
	err = c.Discover(maxSalary, player)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *Npc) Discover(maxSalary int64, player *Player) error {
	lock := roleLock.Get(c.ID)
	lock.Lock()
	defer lock.Unlock()
	if c.Discovered {
		log.Warning("npc已经被发现，不能重复发现")
		return nil
	}
	var rareThreshold float32 = 0.1
	if player.PlayerTag1 == playerTag慧眼识珠 {
		rareThreshold = 0.3
	}
	rare := rand.Float32() < rareThreshold

	c.SalaryExpectation = randInt64(1, maxSalary)
	fluctuation := rand.Float32()
	if rare {
		fluctuation += 0.5
	} else {
		fluctuation = fluctuation*0.2 + 0.9
	}
	c.Prolificacy = int32(float32(c.SalaryExpectation) * 5 * fluctuation)
	if rare {
		c.Teamwork = randInt32(-1500, 1500)
		c.Turnover = randInt32(-200, 500)
	} else {
		c.Teamwork = randInt32(-500, 500)
		c.Turnover = randInt32(0, 500)
	}

	c.Discovered = true
	return db.Save(&c).Error
}

func (c *Npc) GetAvatarPath() string {
	p, err := filepath.Abs(path.Join("NPC", "avatar", strconv.FormatInt(c.ID, 10)+".jpg"))
	if err != nil {
		p = "https://sqimg.qq.com/qq_product_operations/im/qqlogo/logo.png"
	}
	return p
}
