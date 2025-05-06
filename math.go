package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/yangge2333/go-docx"
)

// ================== 配置参数 ==================
const (
	questionCount = 100 // 生成题目数量
	animalTypes   = 100 // 动物种类数
	lifeScenarios = 20  // 生活场景模板数
	shapeTypes    = 20  // 几何图形种类
)

var legRange = [2]int{2, 4}

// ================== 数据结构 ==================
type Animal struct {
	Name string
	Legs int
}

// 动物名称列表
var animalNames = []string{
	"鸡", "兔", "鸭", "企鹅", "鸵鸟", "袋鼠", "马", "牛", "羊", "猪",
	"狗", "猫", "猴子", "猩猩", "长颈鹿", "大象", "狮子", "老虎", "熊", "狐狸",
	"蛇", "蜥蜴", "乌龟", "青蛙", "蟾蜍", "蝙蝠", "老鹰", "猫头鹰", "鸽子", "鹦鹉",
	"金鱼", "鲤鱼", "草鱼", "鲨鱼", "鲸鱼", "海豚", "海豹", "海狮", "海象", "企鹅",
	"孔雀", "火鸡", "鹌鹑", "鹧鸪", "画眉", "八哥", "喜鹊", "乌鸦", "啄木鸟", "猫头鹰",
	"袋鼠", "树袋熊", "鸭嘴兽", "针鼹", "食蚁兽", "穿山甲", "刺猬", "豪猪", "老鼠", "松鼠",
	"兔子", "龙猫", "荷兰猪", "仓鼠", "龙猫", "豚鼠", "水豚", "骆驼", "牦牛", "梅花鹿",
	"驯鹿", "麋鹿", "犀牛", "河马", "野猪", "野狼", "豺狼", "狐狸", "浣熊", "臭鼬",
	"貂", "水獭", "海狸", "熊猫", "棕熊", "北极熊", "黑熊", "眼镜蛇", "响尾蛇", "蝮蛇",
	"蟒蛇", "蚺蛇", "鳄鱼", "蜥蜴", "变色龙", "壁虎", "青蛙", "蟾蜍", "娃娃鱼", "蝾螈",
}

// 按腿数分组的动物列表
var twoLeggedAnimals []string
var fourLeggedAnimals []string

// ================== 初始化数据 ==================
var (
	// 100种动物库（网页8动物类型扩展）
	animals = generateAnimals(animalTypes)

	// 20种生活场景模板（网页5数学建模思路）
	lifeTemplates = make([]func() (string, float64), lifeScenarios)

	// 20种几何图形生成器（网页12几何库参考）
	shapeGenerators = make([]func() (string, float64), shapeTypes)
)

func init() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 初始化按腿数分组的动物列表
	for _, name := range animalNames {
		legs := rand.Intn(legRange[1]-legRange[0]+1) + legRange[0]
		if legs == 2 {
			twoLeggedAnimals = append(twoLeggedAnimals, name)
		} else {
			fourLeggedAnimals = append(fourLeggedAnimals, name)
		}
	}

	// 初始化生活场景模板
	for i := 0; i < lifeScenarios; i++ {
		lifeTemplates[i] = createLifeScenario(i)
	}

	// 初始化几何图形生成器
	for i := 0; i < shapeTypes; i++ {
		shapeGenerators[i] = geometryProblem
	}
}

// ================== 核心功能模块 ==================

// 1. 四则运算及混合运算（网页6基础运算扩展）
func generateBasic() (string, float64) {
	// 基础四则运算占比 50%，混合运算占比 20%
	if rand.Float64() < 0.5 {
		return basicArithmetic()
	} else if rand.Float64() < 0.7 {
		return mixedOperation()
	}
	return basicArithmetic()
}

// 记录已经生成的题目
var generatedQuestions = make(map[string]bool)

// 基础四则运算
func basicArithmetic() (string, float64) {
	ops := []string{"+", "-", "×", "÷"}
	var question string
	var result float64

	for {
		op := ops[rand.Intn(len(ops))]
		// 操作数限制为个位数
		a := rand.Intn(9) + 1
		b := rand.Intn(9) + 1

		switch op {
		case "+":
			result = float64(a + b)
		case "-":
			result = float64(a - b)
		case "×":
			result = float64(a * b)
		case "÷":
			// 确保除法没有小数点
			for b == 0 || a%b != 0 {
				a = rand.Intn(9) + 1
				b = rand.Intn(9) + 1
			}
			result = float64(a / b)
		}

		// 确保答案为十位数以内
		if result >= 100 {
			continue
		}

		question = fmt.Sprintf("%d %s %d =", a, op, b)
		// 检查题目是否重复
		if !generatedQuestions[question] {
			generatedQuestions[question] = true
			break
		}
	}

	return question, result
}

// 新增混合运算（网页5算法扩展）
func mixedOperation() (string, float64) {
	patterns := []string{
		"(%d %s %d) %s %d",
		"%d %s (%d %s %d)",
		"%d %s %d %s %d",
	}

	opCombos := [][]string{
		{"+", "×"}, {"-", "÷"}, {"×", "+"}, {"÷", "-"},
	}

	var question string
	var result float64

	for {
		// 随机选择运算模式
		pattern := patterns[rand.Intn(len(patterns))]
		ops := opCombos[rand.Intn(len(opCombos))]

		// 生成运算数，限制为个位数
		a := rand.Intn(9) + 1
		b := rand.Intn(9) + 1
		c := rand.Intn(9) + 1

		// 处理除法没有小数点
		if ops[0] == "÷" {
			for b == 0 || a%b != 0 {
				a = rand.Intn(9) + 1
				b = rand.Intn(9) + 1
			}
		}
		if ops[1] == "÷" {
			switch pattern {
			case "(%d %s %d) %s %d":
				for c == 0 || (a+b)%c != 0 {
					a = rand.Intn(9) + 1
					b = rand.Intn(9) + 1
					c = rand.Intn(9) + 1
				}
			case "%d %s (%d %s %d)":
				for c == 0 || b%c != 0 {
					a = rand.Intn(9) + 1
					b = rand.Intn(9) + 1
					c = rand.Intn(9) + 1
				}
			case "%d %s %d %s %d":
				for c == 0 || b%c != 0 {
					a = rand.Intn(9) + 1
					b = rand.Intn(9) + 1
					c = rand.Intn(9) + 1
				}
			}
		}

		// 构造表达式
		expr := fmt.Sprintf(pattern, a, ops[0], b, ops[1], c)
		question = expr + " ="

		// 计算结果
		switch ops[0] + ops[1] {
		case "+×":
			result = float64(a + b*c)
		case "-÷":
			result = float64((a - b) / c)
		case "×+":
			result = float64(a*b + c)
		case "÷-":
			result = float64(a/b - c)
		}

		// 确保答案为十位数以内
		if result >= 100 {
			continue
		}

		// 检查题目是否重复
		if !generatedQuestions[question] {
			generatedQuestions[question] = true
			break
		}
	}

	return question, result
}

// 2. 动物组合题（网页14动物问题扩展）
func generateAnimalProblem() (string, float64) {
	// 随机选择两种动物，一种 2 条腿，一种 4 条腿
	animalName1 := twoLeggedAnimals[rand.Intn(len(twoLeggedAnimals))]
	animalName2 := fourLeggedAnimals[rand.Intn(len(fourLeggedAnimals))]
	leg1 := 2
	leg2 := 4

	// 生成题目参数（网页1动态规划思路）
	heads := rand.Intn(9) + 1
	totalLegs := heads*2 + rand.Intn(9)*2

	// 计算答案
	animal2Count := (totalLegs - heads*leg1) / (leg2 - leg1)
	animal1Count := heads - animal2Count

	question := fmt.Sprintf("%s(%d条腿)和%s(%d条腿)共%d只，腿共%d条，%s有多少只？",
		animalName1, leg1, animalName2, leg2, heads, totalLegs, animalName1)

	return question, float64(animal1Count)
}

// 3. 生活应用题（网页5数学建模扩展）
func generateLifeScenario() (string, float64) {
	return lifeTemplates[rand.Intn(len(lifeTemplates))]()
}

// 4. 几何图形题（网页12几何库参考）
func generateGeometry() (string, float64) {
	return shapeGenerators[rand.Intn(len(shapeGenerators))]()
}

// ================== 辅助函数 ==================

// 生成100种动物（网页8动物类型参考）
func generateAnimals(count int) []Animal {
	base := []Animal{}
	for i := 0; i < count; i++ {
		base = append(base, Animal{
			Name: fmt.Sprintf("动物%d", i+1),
			Legs: legRange[0] + rand.Intn(legRange[1]-legRange[0]+1),
		})
	}
	return base
}

// 创建生活场景模板（网页5案例扩展）
func createLifeScenario(index int) func() (string, float64) {
	switch index {
	case 0: // 购物场景
		return func() (string, float64) {
			price1 := rand.Float64()*10 + 5
			quantity1 := rand.Intn(10) + 1
			price2 := rand.Float64()*15 + 8
			quantity2 := rand.Intn(8) + 1
			money := rand.Float64()*(price1*float64(quantity1)+price2*float64(quantity2)) + 20
			question := fmt.Sprintf("买单价为%.2f元的商品%d件，单价为%.2f元的商品%d件，给了%.2f元，应找回多少钱？", price1, quantity1, price2, quantity2, money)
			result := money - (price1*float64(quantity1) + price2*float64(quantity2))
			return question, result
		}
	case 1: // 交通费用
		return func() (string, float64) {
			distance := rand.Float64()*50 + 10
			basePrice := 10.0
			perKmPrice := 2.0
			if distance > 3 {
				distance -= 3
			} else {
				distance = 0
			}
			question := fmt.Sprintf("打车行驶了%.2f公里，起步价%d元（3公里内），超过3公里后每公里%.2f元，需要支付多少钱？", distance+3, int(basePrice), perKmPrice)
			result := basePrice + distance*perKmPrice
			return question, result
		}
	case 2: // 工程问题
		return func() (string, float64) {
			workRate1 := rand.Float64()*0.1 + 0.05
			workRate2 := rand.Float64()*0.1 + 0.05
			days := rand.Intn(20) + 5
			question := fmt.Sprintf("甲的工作效率是%.2f，乙的工作效率是%.2f，两人合作%d天完成了多少工作量？", workRate1, workRate2, days)
			result := (workRate1 + workRate2) * float64(days)
			return question, result
		}
	case 3: // 浓度问题
		return func() (string, float64) {
			solute1 := rand.Float64()*20 + 10
			solvent1 := rand.Float64()*80 + 20
			solute2 := rand.Float64()*15 + 5
			solvent2 := rand.Float64()*60 + 15
			question := fmt.Sprintf("将含溶质%.2f克，溶剂%.2f克的溶液与含溶质%.2f克，溶剂%.2f克的溶液混合，混合后溶液的浓度是多少？", solute1, solvent1, solute2, solvent2)
			result := (solute1 + solute2) / (solute1 + solvent1 + solute2 + solvent2)
			return question, result
		}
	case 4: // 利润问题
		return func() (string, float64) {
			cost := rand.Float64()*50 + 20
			price := rand.Float64()*(cost+20) + cost
			quantity := rand.Intn(10) + 1
			question := fmt.Sprintf("商品成本为%.2f元，售价为%.2f元，卖出%d件，利润是多少？", cost, price, quantity)
			result := (price - cost) * float64(quantity)
			return question, result
		}
	case 5: // 时间问题
		return func() (string, float64) {
			startHour := rand.Intn(24)
			startMinute := rand.Intn(60)
			durationHour := rand.Intn(10)
			durationMinute := rand.Intn(60)
			totalMinute := startMinute + durationMinute
			addHour := totalMinute / 60
			totalMinute %= 60
			endHour := (startHour + durationHour + addHour) % 24
			question := fmt.Sprintf("从%d时%d分开始，经过%d小时%d分后是几点？", startHour, startMinute, durationHour, durationMinute)
			result := float64(endHour*60 + totalMinute)
			return question, result
		}
	case 6: // 速度问题
		return func() (string, float64) {
			distance := rand.Float64()*100 + 50
			time := rand.Float64()*10 + 5
			question := fmt.Sprintf("行驶了%.2f公里，用时%.2f小时，速度是多少？", distance, time)
			result := distance / time
			return question, result
		}
	case 7: // 折扣问题
		return func() (string, float64) {
			originalPrice := rand.Float64()*100 + 50
			discount := rand.Float64()*0.5 + 0.5
			question := fmt.Sprintf("商品原价%.2f元，打%.2f折后价格是多少？", originalPrice, discount*10)
			result := originalPrice * discount
			return question, result
		}
	case 8: // 平均数问题
		return func() (string, float64) {
			num1 := rand.Intn(50) + 10
			num2 := rand.Intn(50) + 10
			num3 := rand.Intn(50) + 10
			question := fmt.Sprintf("%d、%d、%d的平均数是多少？", num1, num2, num3)
			result := float64(num1+num2+num3) / 3
			return question, result
		}
	case 9: // 比例问题
		return func() (string, float64) {
			part := rand.Intn(50) + 10
			total := rand.Intn(100) + 50
			question := fmt.Sprintf("%d占%d的百分之几？", part, total)
			result := float64(part) / float64(total) * 100
			return question, result
		}
	case 10: // 植树问题（两端都种）
		return func() (string, float64) {
			length := rand.Intn(100) + 20
			interval := rand.Intn(10) + 2
			question := fmt.Sprintf("在%d米长的道路两端都种树，每隔%d米种一棵，一共种多少棵树？", length, interval)
			result := float64(length/interval + 1)
			return question, result
		}
	case 11: // 植树问题（两端都不种）
		return func() (string, float64) {
			length := rand.Intn(100) + 20
			interval := rand.Intn(10) + 2
			question := fmt.Sprintf("在%d米长的道路两端都不种树，每隔%d米种一棵，一共种多少棵树？", length, interval)
			result := float64(length/interval - 1)
			return question, result
		}
	case 12: // 年龄问题
		return func() (string, float64) {
			age1 := rand.Intn(30) + 10
			age2 := rand.Intn(30) + 10
			years := rand.Intn(10) + 1
			question := fmt.Sprintf("甲今年%d岁，乙今年%d岁，%d年后甲比乙大几岁？", age1, age2, years)
			result := float64(age1 - age2)
			return question, result
		}
	case 13: // 盈亏问题
		return func() (string, float64) {
			profit := rand.Intn(20) + 5
			loss := rand.Intn(20) + 5
			difference := rand.Intn(5) + 1
			question := fmt.Sprintf("分物品，每人分%d个多%d个，每人分%d个少%d个，有多少人？", difference, profit, difference+1, loss)
			result := float64(profit+loss) / float64(1)
			return question, result
		}
	case 14: // 鸡兔同笼变种（多种动物）
		return func() (string, float64) {
			animal1Legs := 2
			animal2Legs := 4
			animal3Legs := 6
			animal1Count := rand.Intn(10) + 5
			animal2Count := rand.Intn(10) + 5
			animal3Count := rand.Intn(10) + 5
			totalHeads := animal1Count + animal2Count + animal3Count
			totalLegs := animal1Count*animal1Legs + animal2Count*animal2Legs + animal3Count*animal3Legs
			question := fmt.Sprintf("鸡(2条腿)、兔(4条腿)、蜘蛛(6条腿)共%d只，腿共%d条，鸡有多少只？", totalHeads, totalLegs)
			result := float64(animal1Count)
			return question, result
		}
	case 15: // 行程问题（相遇）
		return func() (string, float64) {
			speed1 := rand.Float64()*20 + 10
			speed2 := rand.Float64()*20 + 10
			distance := rand.Float64()*100 + 50
			question := fmt.Sprintf("甲速度为%.2f公里/小时，乙速度为%.2f公里/小时，两人相向而行，相距%.2f公里，多久相遇？", speed1, speed2, distance)
			result := distance / (speed1 + speed2)
			return question, result
		}
	case 16: // 行程问题（追及）
		return func() (string, float64) {
			speed1 := rand.Float64()*20 + 10
			speed2 := rand.Float64()*20 + 10
			if speed2 >= speed1 {
				speed2 = speed1 - 1
			}
			distance := rand.Float64()*100 + 50
			question := fmt.Sprintf("甲速度为%.2f公里/小时，乙速度为%.2f公里/小时，乙在甲前面%.2f公里，甲多久能追上乙？", speed1, speed2, distance)
			result := distance / (speed1 - speed2)
			return question, result
		}
	case 17: // 利息问题
		return func() (string, float64) {
			principal := rand.Float64()*1000 + 500
			rate := rand.Float64()*0.05 + 0.01
			years := rand.Intn(5) + 1
			question := fmt.Sprintf("本金%.2f元，年利率%.2f%%，存%d年，利息是多少？", principal, rate*100, years)
			result := principal * rate * float64(years)
			return question, result
		}
	case 18: // 面积比例问题
		return func() (string, float64) {
			area1 := rand.Float64()*100 + 50
			area2 := rand.Float64()*100 + 50
			question := fmt.Sprintf("甲面积为%.2f平方米，乙面积为%.2f平方米，甲面积是乙面积的几分之几？", area1, area2)
			result := area1 / area2
			return question, result
		}
	case 19: // 人数分配问题
		return func() (string, float64) {
			totalPeople := rand.Intn(100) + 50
			group1Ratio := rand.Float64()*0.5 + 0.1
			question := fmt.Sprintf("总共有%d人，按%.2f的比例分配到一组，这组有多少人？", totalPeople, group1Ratio)
			result := float64(totalPeople) * group1Ratio
			return question, result
		}
	default:
		return basicLifeScenario
	}
}

func basicLifeScenario() (string, float64) {
	return "这是一个基础生活场景题目", 0
}

// 生成几何图形题目
func geometryProblem() (string, float64) {
	shapes := []func() (string, float64){
		// 长方形周长比较
		func() (string, float64) {
			length1 := rand.Intn(20) + 5
			width1 := rand.Intn(20) + 5
			length2 := rand.Intn(20) + 5
			width2 := rand.Intn(20) + 5
			perimeter1 := 2 * (length1 + width1)
			perimeter2 := 2 * (length2 + width2)
			question := fmt.Sprintf("长方形 A 的长是%d，宽是%d，长方形 B 的长是%d，宽是%d，哪个长方形的周长大，大多少？", length1, width1, length2, width2)
			if perimeter1 > perimeter2 {
				return question, float64(perimeter1 - perimeter2)
			}
			return question, float64(perimeter2 - perimeter1)
		},
		// 等边梯形边长计算
		func() (string, float64) {
			upperBase := rand.Intn(10) + 5
			lowerBase := rand.Intn(10) + upperBase + 5
			waist := rand.Intn(10) + 5
			question := fmt.Sprintf("一个等边梯形的上底是%d，下底是%d，腰长是%d，它的周长是多少？", upperBase, lowerBase, waist)
			result := float64(upperBase + lowerBase + 2*waist)
			return question, result
		},
		// 五角星边长总和
		func() (string, float64) {
			sideLength := rand.Intn(10) + 5
			question := fmt.Sprintf("一个正五角星的每条边长度是%d，它的所有边长总和是多少？", sideLength)
			result := float64(5 * sideLength)
			return question, result
		},
		// 正方形周长计算
		func() (string, float64) {
			side := rand.Intn(20) + 5
			question := fmt.Sprintf("一个正方形的边长是%d，它的周长是多少？", side)
			result := float64(4 * side)
			return question, result
		},
		// 三角形周长计算
		func() (string, float64) {
			side1 := rand.Intn(10) + 5
			side2 := rand.Intn(10) + 5
			side3 := rand.Intn(10) + 5
			question := fmt.Sprintf("一个三角形的三条边分别是%d、%d、%d，它的周长是多少？", side1, side2, side3)
			result := float64(side1 + side2 + side3)
			return question, result
		},
		// 圆形周长计算
		func() (string, float64) {
			radius := rand.Intn(10) + 5
			question := fmt.Sprintf("一个圆的半径是%d，它的周长是多少？（π取3.14）", radius)
			result := 2 * 3.14 * float64(radius)
			return question, result
		},
		// 平行四边形周长计算
		func() (string, float64) {
			side1 := rand.Intn(10) + 5
			side2 := rand.Intn(10) + 5
			question := fmt.Sprintf("一个平行四边形相邻两边分别是%d和%d，它的周长是多少？", side1, side2)
			result := float64(2 * (side1 + side2))
			return question, result
		},
		// 菱形周长计算
		func() (string, float64) {
			side := rand.Intn(10) + 5
			question := fmt.Sprintf("一个菱形的边长是%d，它的周长是多少？", side)
			result := float64(4 * side)
			return question, result
		},
		// 正六边形周长计算
		func() (string, float64) {
			side := rand.Intn(10) + 5
			question := fmt.Sprintf("一个正六边形的边长是%d，它的周长是多少？", side)
			result := float64(6 * side)
			return question, result
		},
		// 等腰三角形周长计算
		func() (string, float64) {
			base := rand.Intn(10) + 5
			waist := rand.Intn(10) + 5
			question := fmt.Sprintf("一个等腰三角形的底边长是%d，腰长是%d，它的周长是多少？", base, waist)
			result := float64(base + 2*waist)
			return question, result
		},
		// 长方形面积计算
		func() (string, float64) {
			length := rand.Intn(20) + 5
			width := rand.Intn(20) + 5
			question := fmt.Sprintf("一个长方形的长是%d，宽是%d，它的面积是多少？", length, width)
			result := float64(length * width)
			return question, result
		},
		// 正方形面积计算
		func() (string, float64) {
			side := rand.Intn(20) + 5
			question := fmt.Sprintf("一个正方形的边长是%d，它的面积是多少？", side)
			result := float64(side * side)
			return question, result
		},
	}
	return shapes[rand.Intn(len(shapes))]()
}

// ================== 主程序 ==================
func main() {
	// 创建一个新的 A4 格式的 DOCX 文件
	w := docx.NewA4()

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")

	// 添加日期到标题栏左上角
	datePara := w.AddParagraph()
	datePara.Justification("left")
	datePara.AddText("日期：").Size("30").Bold().Italic()
	datePara.AddText(currentDate).Underline("____________").Size("25").Bold()

	// 添加小学数学题相关信息
	headerPara := w.AddParagraph()
	headerPara.Justification("center")
	headerPara.AddText("小学数学题").Size("28").Bold()

	infoPara := w.AddParagraph()
	infoPara.Justification("right")
	infoPara.AddText("姓名____________  班级____________   成绩____________")

	// 正则表达式，用于匹配包含两个或以上运算符的表达式
	mixedOpRegex := regexp.MustCompile(`[+\-×÷].*[+\-×÷]`)

	// 按类型分类存储题目和答案
	questionMap := make(map[string][]string)
	answerMap := make(map[string][]string)

	for i := 1; i <= questionCount; i++ {
		var qType, question string
		var ans float64

		// 基础四则运算和混合运算占比 70%，其他占比 30%
		if rand.Float64() < 0.7 {
			question, ans = generateBasic()
			if mixedOpRegex.MatchString(question) {
				qType = "混合运算"
			} else {
				qType = "四则运算"
			}
		} else {
			randType := rand.Intn(3)
			switch randType {
			case 0:
				question, ans = generateAnimalProblem()
				qType = "动物问题"
			case 1:
				question, ans = generateLifeScenario()
				qType = "生活应用"
			case 2:
				question, ans = generateGeometry()
				qType = "几何图形"
			}
		}

		// 存储题目和答案
		questionMap[qType] = append(questionMap[qType], question)
		answerMap[qType] = append(answerMap[qType], fmt.Sprintf("%.1f", ans))
	}

	categoryOrder := []string{"四则运算", "混合运算", "动物问题", "生活应用", "几何图形"}
	categoryChineseNumbers := []string{"一", "二", "三", "四", "五"}

	// 按类型输出题目
	for index, qType := range categoryOrder {
		questions := questionMap[qType]
		if len(questions) > 0 {
			// 添加类型标题
			titlePara := w.AddParagraph()
			titlePara.AddText(fmt.Sprintf("%s 、%s（%d题）", categoryChineseNumbers[index], qType, len(questions))).Size("24").Bold()

			if qType == "四则运算" || qType == "混合运算" {
				// 四则运算和混合运算三列显示
				for i := 0; i < len(questions); i += 3 {
					para := w.AddParagraph()
					// 第一题
					if i < len(questions) {
						para.AddText(fmt.Sprintf("第%d题：%s _____", i+1, questions[i]))
					}
					// 填充空格使其对齐
					para.AddText("    ")
					// 第二题
					if i+1 < len(questions) {
						para.AddText(fmt.Sprintf("第%d题：%s _____", i+2, questions[i+1]))
					}
					para.AddText("    ")
					// 第三题
					if i+2 < len(questions) {
						para.AddText(fmt.Sprintf("第%d题：%s _____", i+3, questions[i+2]))
					}
				}
			} else {
				// 其他类型题目正常输出
				for i, q := range questions {
					para := w.AddParagraph()
					para.AddText(fmt.Sprintf("第%d题：%s _____", i+1, q))
				}
			}
			// 添加空段落分隔
			w.AddParagraph()
		}
	}

	// 添加空段落分隔
	w.AddParagraph()
	w.AddParagraph()
	w.AddParagraph()
	w.AddParagraph()
	w.AddParagraph()
	w.AddParagraph()

	// 在最后添加答案标题
	answerTitlePara := w.AddParagraph()
	answerTitlePara.AddText("答案").Size("28").Bold()

	// 按类型输出答案，横着写，每个答案中间两个空格
	for index, qType := range categoryOrder {
		answers := answerMap[qType]
		if len(answers) > 0 {
			// 添加类型标题
			answerTypePara := w.AddParagraph()
			answerTypePara.AddText(fmt.Sprintf("%s 、%s（%d题）", categoryChineseNumbers[index], qType, len(answers))).Size("24").Bold()

			// 横着输出答案
			answerPara := w.AddParagraph()
			for i, a := range answers {
				if i > 0 {
					answerPara.AddText("  ")
				}
				answerPara.AddText(fmt.Sprintf("第%d题答案：%s", i+1, a))
			}
			// 添加空段落分隔
			w.AddParagraph()
		}
	}

	// 保存 DOCX 文件
	f, err := os.Create("math_questions.docx")
	if err != nil {
		fmt.Println("创建文件时出错:", err)
		return
	}
	defer f.Close()

	_, err = w.WriteTo(f)
	if err != nil {
		fmt.Println("保存文件时出错:", err)
		return
	}

	fmt.Println("已生成", questionCount, "道题目到 math_questions.docx")
}
