package main

import (
	"fmt"
	"strings"
)

type OperationResult string

const (
	成功 OperationResult = "恭喜你抢到泡泡玛特"
	失败 OperationResult = "不好意思没货了"
)

type StockStatus string

const (
	有库存 StockStatus = "有库存"
	无库存 StockStatus = "无库存"
)

type Blindbox struct {
	Name         string
	RealPrice    float64
	SellPrice    float64
	SurfStock    int
	MinSurfStock int
	RealStock    int
	BoxPerCase   int
	Warnsent     bool
}

func (p *Blindbox) TotalPrice() float64 {
	return p.RealPrice * float64(p.RealStock)
}
func (p *Blindbox) ProfitPer() float64 {

	return p.SellPrice - p.RealPrice
}
func (p *Blindbox) TotalProfit() float64 {
	return p.ProfitPer() * float64(p.RealStock)
}
func (p *Blindbox) IsInstock() StockStatus {
	if p.SurfStock > 0 {
		return 有库存
	}
	return 无库存
}
func (p *Blindbox) CostomerInfo() string {
	info := fmt.Sprintf("🎁商品：%s,售价：%.2f元,库存：%d件", p.Name, p.SellPrice, p.SurfStock)
	info += fmt.Sprintf("\n市场估价可达%.2f元！", p.SellPrice*1.5)
	if p.SurfStock <= 5 && !p.Warnsent {
		info += "\n🚨 【库存紧张！欲购从速】"
		p.Warnsent = true
	}
	if p.SurfStock > p.BoxPerCase {
		info += fmt.Sprintf("\n💫 推荐整端购买（%d个一盒）！隐藏款概率大幅提升！", p.BoxPerCase)
	}
	return info
}
func (p *Blindbox) BossInfo() string {
	return fmt.Sprintf("💼 商品: %s, 实际价值: %.2f元, 售卖价格: %.2f元, 表面库存: %d件, 真实库存: %d件\n  , 单件利润: %.2f元, 总利润: %.2f元, 总价值: %.2f元", p.Name, p.RealPrice, p.SellPrice, p.SurfStock, p.RealStock, p.ProfitPer(), p.TotalProfit(), p.TotalPrice())
}
func (p *Blindbox) CheckStockWarn() string {
	if p.SurfStock < 10 && p.RealStock >= 50 {
		return "⚠️  【奸商提示】表面库存紧张，但真实库存充足，可以继续饥饿营销"
	}
	if p.SurfStock <= 5 {
		return "🚨 【紧急】表面库存即将售罄！考虑补充表面库存"
	}
	if p.SurfStock > p.BoxPerCase {
		return "💡 【销售机会】库存充足，可以推广整端购买！利润率可观"
	}
	return ""
}
func (p *Blindbox) Restock(amount int) {
	p.RealStock += amount
	SurfaceIncreace := amount / 10
	if SurfaceIncreace < 1 {
		SurfaceIncreace = 1
	}
	p.SurfStock += SurfaceIncreace
	p.Warnsent = false
}
func (p *Blindbox) Sell(amount int) (result OperationResult, s string) {
	isWholecase := amount == p.BoxPerCase
	isMultipleCase := amount > p.BoxPerCase && amount%p.BoxPerCase == 0
	if amount > p.SurfStock {
		if isWholecase {
			return 失败, s
		}
		return 失败, s
	}
	p.SurfStock -= amount
	p.RealStock -= amount
	if p.SurfStock < p.MinSurfStock && p.RealStock > 0 {
		borrowAmount := p.MinSurfStock - p.SurfStock
		if borrowAmount > p.RealStock/20 {
			borrowAmount = p.RealStock
		}
		if borrowAmount > 0 {
			p.SurfStock += borrowAmount
		}
	}
	if isMultipleCase {
		caseCount := amount / p.BoxPerCase
		return 成功, fmt.Sprintf("🎊 太幸运了！抢到%d整端盲盒！隐藏款在向你招手！", caseCount)
	} else if isWholecase {
		return 成功, s
	} else {
		if p.SurfStock <= 3 {
			return 成功, " 恭喜抢到盲盒！最后几件了，整端购买机会更大哦！"
		} else if p.SurfStock <= 10 {
			return 成功, "⭐ 恭喜抢到盲盒！库存紧张，整端购买隐藏款概率更高！"
		} else {
			return 成功, "✨ 恭喜抢到盲盒！祝你好运！"
		}
	}
}
func (p *Blindbox) SellWholeCase() (result OperationResult, message string) {
	return p.Sell(p.BoxPerCase)
}
func (p *Blindbox) SellMultipleCase(caseCount int) (result OperationResult, message string) {
	return p.Sell(p.BoxPerCase * caseCount)
}
func (p *Blindbox) AdjustSurfStock(NewSurf int) {
	if NewSurf <= p.RealStock {
		p.SurfStock = NewSurf
		p.Warnsent = false
	}
}
func CreatBlindBox(name string, actualprice float64, sellprice float64, realstock int) *Blindbox {
	surfstock := realstock / 15
	if surfstock < 5 {
		surfstock = 5
	}
	return &Blindbox{
		Name:         name,
		RealPrice:    actualprice,
		SellPrice:    sellprice,
		SurfStock:    surfstock,
		RealStock:    realstock,
		MinSurfStock: 3,
		BoxPerCase:   9,
		Warnsent:     false,
	}
}
func main() {
	product := CreatBlindBox("新生日记盲盒内部渠道", 99.00, 120.00, 500)
	fmt.Println("===奸商盲盒销售系统（暴利版）===")
	fmt.Println("\n顾客看到的信息：")
	fmt.Println(product.CostomerInfo())
	fmt.Printf("库存状态：%s\n", product.IsInstock())
	fmt.Println("\n奸商后台消息")
	fmt.Println(product.BossInfo())
	if warn := product.CheckStockWarn(); warn != "" {
		fmt.Printf("\n🔔 %s\n", warn)
	}
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("\n===销售测试===")
	result, message := product.Sell(2)
	fmt.Printf("单个购买两个：%s-%s\n", result, message)
	fmt.Printf("当前表面库存：%d件\n", product.SurfStock)
	fmt.Println(product.CostomerInfo())
	result, message = product.SellWholeCase()
	fmt.Printf("\n整端购买%d个：%s-%s\n", product.BoxPerCase, result, message)
	fmt.Printf("当前表面库存：%d件\n", product.SurfStock)
	if warn := product.CheckStockWarn(); warn != "" {
		fmt.Printf("【系统提示】%s\n", warn)
	}
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("\n===多端购买测试===")
	fmt.Println("顾客看到的消息：")
	fmt.Println(product.CostomerInfo())
	result, message = product.SellMultipleCase(2)
	fmt.Printf("\n购买两整端（%d个）；%s-%s\n", product.BoxPerCase*2, result, message)
	fmt.Printf("当前表面库存：%d件\n", product.SurfStock)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("\n===奸商后台操作===")
	fmt.Println(product.BossInfo())
	fmt.Println("奸商调整库存为8件(刚好不够整端）")
	product.AdjustSurfStock(8)
	fmt.Println("顾客看到的信息：")
	fmt.Println(product.CostomerInfo())
	result, message = product.SellWholeCase()
	fmt.Printf("\n尝试整端购买:%s-%s\n", result, message)
	fmt.Println("奸商悄悄补货200件")
	product.Restock(200)
	fmt.Println("补货后奸商信息：")
	fmt.Println(product.BossInfo())
	fmt.Println("补货后顾客看到的信息：")
	fmt.Println(product.CostomerInfo())
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("===完整销售模拟：===")
	product = CreatBlindBox("新生日记内部货源", 99.00, 120.00, 200)
	salse := []struct {
		amount int
		recall string
	}{
		{2, "单买试试手气"},
		{9, "整端购买"},
		{1, "单买"},
		{18, "两端"},
		{3, "单买"},
	}
	totalprofit := 0.0
	for i, sale := range salse {
		fmt.Printf("\n第%d次交易：%s\n", i+1, sale.recall)
		result, msg := product.Sell(sale.amount)
		fmt.Printf("结果：%s-%s\n", result, msg)
		fmt.Printf("表面库存%d件\n", product.SurfStock)
		profit := product.ProfitPer() * float64(sale.amount)
		totalprofit += profit
		fmt.Printf("本次利润：%.2f元\n", profit)
		fmt.Println(product.CostomerInfo())
		if warn := product.CheckStockWarn(); warn != "" {
			fmt.Printf("💡 %s\n", warn)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
	fmt.Println("\n🎊 最终统计:")
	fmt.Println(product.BossInfo())
	fmt.Printf("模拟销售总利润: %.2f元\n赚死奸商了", totalprofit)
}
