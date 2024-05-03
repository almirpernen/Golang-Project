package models

type Player struct{
	Name string
	Surname string
	Nickname string 
	Kills int
	Deaths int
	Assists int
	Headshots int
	ADR int
	Rounds int
	ImpRounds int
}

func (p *Player) GetKD() float32{
	return float32(p.Kills) / float32(p.Deaths)
}

func (p *Player) GetImpact() float32{
	var a1 float32 = 2.13 * p.GetKPR()
	var b1 float32 = 0.42 * p.GetAPR()
	return a1 + b1 - 0.41
}

func (p *Player) GetKPR() float32{
	return float32(p.Kills) / float32(p.Rounds)
}

func (p *Player) GetDPR() float32{
	return float32(p.Deaths) / float32(p.Rounds)
}

func (p *Player) GetAPR() float32{
	return float32(p.Assists) / float32(p.Rounds)
}

func (p *Player) GetRating() float32{
	var a1 float32 = 0.3591 * p.GetKPR()
	var b1 float32 = -0.5329 * p.GetDPR()
	var c1 float32 = 0.2372 * p.GetImpact()
	var d1 float32 = 0.0032 * float32(p.ADR)
	var e1 float32 = 0.0073 * p.GetKAST()
	return a1 + b1 + c1 + d1 + e1 + 0.1587
}

func (p *Player) GetKAST() float32{
	return 100 * float32(p.ImpRounds) / float32(p.Rounds)
}