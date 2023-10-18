package models

type Renown1Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown1Struct) populate() {
	r.Level1 = 100
	r.Level2 = 220
	r.Level3 = 360
	r.Level4 = 520
	r.Level5 = 700
	r.Level6 = 900
	r.Level7 = 1120
	r.Level8 = 1360
	r.Level9 = 1620
	r.Level10 = 1900
}

func GetRenown1() *Renown1Struct {
	r := new(Renown1Struct)
	r.populate()
	return r
}

func (r Renown1Struct) determineRenownLevel(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown2Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown2Struct) populate() {
	r.Level1 = 2200
	r.Level2 = 2520
	r.Level3 = 2860
	r.Level4 = 3220
	r.Level5 = 3600
	r.Level6 = 4000
	r.Level7 = 4420
	r.Level8 = 4860
	r.Level9 = 5320
	r.Level10 = 5800
}

func GetRenown2() *Renown2Struct {
	r := new(Renown2Struct)
	r.populate()
	return r
}

func (r Renown2Struct) determineRenown2Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown3Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown3Struct) populate() {
	r.Level1 = 6300
	r.Level2 = 6820
	r.Level3 = 7360
	r.Level4 = 7920
	r.Level5 = 8500
	r.Level6 = 9100
	r.Level7 = 9720
	r.Level8 = 10360
	r.Level9 = 11020
	r.Level10 = 11700
}

func GetRenown3() *Renown3Struct {
	r := new(Renown3Struct)
	r.populate()
	return r
}

func (r Renown3Struct) determineRenown3Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64

	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown4Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown4Struct) populate() {
	r.Level1 = 12400
	r.Level2 = 13120
	r.Level3 = 13860
	r.Level4 = 14620
	r.Level5 = 15400
	r.Level6 = 16200
	r.Level7 = 17020
	r.Level8 = 17860
	r.Level9 = 18720
	r.Level10 = 19600
}

func GetRenown4() *Renown4Struct {
	r := new(Renown4Struct)
	r.populate()
	return r
}

func (r Renown4Struct) determineRenown4Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown5Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown5Struct) populate() {
	r.Level1 = 20500
	r.Level2 = 21420
	r.Level3 = 22360
	r.Level4 = 23320
	r.Level5 = 24300
	r.Level6 = 25300
	r.Level7 = 26320
	r.Level8 = 27360
	r.Level9 = 28420
	r.Level10 = 29500
}

func GetRenown5() *Renown5Struct {
	r := new(Renown5Struct)
	r.populate()
	return r
}

func (r Renown5Struct) determineRenown5Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown6Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown6Struct) populate() {
	r.Level1 = 30600
	r.Level2 = 31720
	r.Level3 = 32860
	r.Level4 = 34020
	r.Level5 = 35200
	r.Level6 = 36400
	r.Level7 = 37620
	r.Level8 = 38860
	r.Level9 = 40120
	r.Level10 = 41400
}

func GetRenown6() *Renown6Struct {
	r := new(Renown6Struct)
	r.populate()
	return r
}

func (r Renown6Struct) determineRenown6Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown7Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown7Struct) populate() {
	r.Level1 = 42700
	r.Level2 = 44020
	r.Level3 = 45360
	r.Level4 = 46720
	r.Level5 = 48100
	r.Level6 = 49500
	r.Level7 = 50920
	r.Level8 = 52360
	r.Level9 = 53820
	r.Level10 = 55300
}

func GetRenown7() *Renown7Struct {
	r := new(Renown7Struct)
	r.populate()
	return r
}

func (r Renown7Struct) determineRenown7Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown8Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown8Struct) populate() {
	r.Level1 = 56800
	r.Level2 = 58320
	r.Level3 = 59860
	r.Level4 = 61420
	r.Level5 = 63000
	r.Level6 = 64600
	r.Level7 = 66220
	r.Level8 = 67860
	r.Level9 = 69520
	r.Level10 = 71200
}

func GetRenown8() *Renown8Struct {
	r := new(Renown8Struct)
	r.populate()
	return r
}

func (r Renown8Struct) determineRenown8Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown9Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown9Struct) populate() {
	r.Level1 = 72900
	r.Level2 = 74620
	r.Level3 = 76360
	r.Level4 = 78120
	r.Level5 = 79900
	r.Level6 = 81700
	r.Level7 = 83520
	r.Level8 = 85360
	r.Level9 = 87220
	r.Level10 = 89100
}

func GetRenown9() *Renown9Struct {
	r := new(Renown9Struct)
	r.populate()
	return r
}

func (r Renown9Struct) determineRenown9Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = false
		i = 0
	}
	return b, i, min, max
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Renown10Struct struct {
	Level1  uint64
	Level2  uint64
	Level3  uint64
	Level4  uint64
	Level5  uint64
	Level6  uint64
	Level7  uint64
	Level8  uint64
	Level9  uint64
	Level10 uint64
}

func (r *Renown10Struct) populate() {
	r.Level1 = 91000
	r.Level2 = 92920
	r.Level3 = 94860
	r.Level4 = 96820
	r.Level5 = 98800
	r.Level6 = 100800
	r.Level7 = 102820
	r.Level8 = 104860
	r.Level9 = 106920
	r.Level10 = 109000
}

func GetRenown10() *Renown10Struct {
	r := new(Renown10Struct)
	r.populate()
	return r
}

func (r Renown10Struct) determineRenown10Level(value uint64) (bool, LevelType, uint64, uint64) {
	var b bool
	var i LevelType
	var min uint64
	var max uint64
	if value < r.Level1 {
		b = true
		i = 0
		min = 0
		max = r.Level1
	}
	if value >= r.Level1 {
		b = true
		i = 1
		min = r.Level1
		max = r.Level2
	}
	if value >= r.Level2 {
		b = true
		i = 2
		min = r.Level2
		max = r.Level3
	}
	if value >= r.Level3 {
		b = true
		i = 3
		min = r.Level3
		max = r.Level4
	}
	if value >= r.Level4 {
		b = true
		i = 4
		min = r.Level4
		max = r.Level5
	}
	if value >= r.Level5 {
		b = true
		i = 5
		min = r.Level5
		max = r.Level6
	}
	if value >= r.Level6 {
		b = true
		i = 6
		min = r.Level6
		max = r.Level7
	}
	if value >= r.Level7 {
		b = true
		i = 7
		min = r.Level7
		max = r.Level8
	}
	if value >= r.Level8 {
		b = true
		i = 8
		min = r.Level8
		max = r.Level9
	}
	if value >= r.Level9 {
		b = true
		i = 9
		min = r.Level9
		max = r.Level10
	}
	if value >= r.Level10 {
		b = true
		i = 9
		min = r.Level10
		max = r.Level10
	}
	return b, i, min, max
}

// DetermineUserRenownLevel renown, level based on amount of xp user has
func DetermineUserRenownLevel(value uint64) (TierType, LevelType, uint64, uint64) {
	// check if user is in renown 1, then return level
	b, i, min, max := GetRenown1().determineRenownLevel(value)
	if b == true {
		return 0, i, min, max
	}

	// check if user is in renown 2, then return level
	b, i, min, max = GetRenown2().determineRenown2Level(value)
	if b == true {
		return 1, i, min, max
	}

	// check if user is in renown 3, then return level
	b, i, min, max = GetRenown3().determineRenown3Level(value)
	if b == true {
		return 2, i, min, max
	}

	// check if user is in renown 4, then return level
	b, i, min, max = GetRenown4().determineRenown4Level(value)
	if b == true {
		return 3, i, min, max
	}

	// check if user is in renown 5, then return level
	b, i, min, max = GetRenown5().determineRenown5Level(value)
	if b == true {
		return 4, i, min, max
	}

	// check if user is in renown 6, then return level
	b, i, min, max = GetRenown6().determineRenown6Level(value)
	if b == true {
		return 5, i, min, max
	}

	// check if user is in renown 7, then return level
	b, i, min, max = GetRenown7().determineRenown7Level(value)
	if b == true {
		return 6, i, min, max
	}

	// check if user is in renown 8, then return level
	b, i, min, max = GetRenown8().determineRenown8Level(value)
	if b == true {
		return 7, i, min, max
	}

	// check if user is in renown 9, then return level
	b, i, min, max = GetRenown9().determineRenown9Level(value)
	if b == true {
		return 8, i, min, max
	}

	// check if user is in renown 10, then return level
	b, i, min, max = GetRenown10().determineRenown10Level(value)
	if b == true {
		return 9, i, min, max
	}
	/// returning 69's as an error identifier
	return 69, 69, 69, 69
}
