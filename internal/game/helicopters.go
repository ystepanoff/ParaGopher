package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/utils"
)

type Helicopter struct {
	x, y        float32
	leftToRight bool
	lastDrop    time.Time
}

func (g *Game) drawHelicopter(screen *ebiten.Image, h *Helicopter) {
	op := &ebiten.DrawImageOptions{}
	dx := float64(g.helicopterImage.Bounds().Dx())
	dy := float64(g.helicopterImage.Bounds().Dy())
	if !h.leftToRight {
		op.GeoM.Scale(-1.0, 1.0)
		op.GeoM.Translate(dx, 0.0)
	}
	op.GeoM.Translate(float64(h.x)-dx/2.0, float64(h.y)-dy/2.0)
	screen.DrawImage(g.helicopterImage, op)

	imgW := float32(dx)
	bodyOffset := float32(config.HelicopterTailWidth) + float32(config.HelicopterBodyWidth)/2.0 - imgW/2.0
	if !h.leftToRight {
		bodyOffset = -bodyOffset
	}
	rotorCX := h.x + bodyOffset
	rotorY := h.y - float32(dy)/2.0 + 1

	phase := float64(time.Now().UnixMilli()) * 0.012
	rotorHalfLen := float32(math.Abs(math.Cos(phase))) * float32(config.HelicopterRotorLen) / 2.0
	if rotorHalfLen < 1.0 {
		rotorHalfLen = 1.0
	}

	vector.StrokeLine(
		screen,
		rotorCX-rotorHalfLen, rotorY,
		rotorCX+rotorHalfLen, rotorY,
		1.0, config.ColourMagenta, false,
	)
}

func (g *Game) drawHelicopters(screen *ebiten.Image) {
	for _, h := range g.helicopters {
		g.drawHelicopter(screen, h)
	}
}

func (g *Game) spawnHelicopter() {
	if rand.Float32() < config.HelicopterSpawnChance {
		startX := -float32(
			config.HelicopterBodyWidth + config.HelicopterTailWidth,
		)
		startY := float32(50 + rand.Intn(50))
		leftToRight := true
		if rand.Intn(2) == 1 {
			startX = config.ScreenWidth - startX
			leftToRight = false
		}
		g.helicopters = append(g.helicopters, &Helicopter{
			x:           startX,
			y:           startY,
			leftToRight: leftToRight,
			lastDrop:    time.Now(),
		})
	}
}

func (g *Game) updateHelicopters() {
	active := make([]*Helicopter, 0, len(g.helicopters))
	for _, h := range g.helicopters {
		vx := float32(config.HelicopterSpeed)
		if !h.leftToRight {
			vx = -vx
		}
		h.x += vx
		timePassed := time.Since(h.lastDrop)
		if timePassed > config.HelicopterDropRate*time.Second &&
			g.canDrop(h.x) && rand.Float32() < config.ParatrooperSpawnChance {
			g.spawnParatrooper(h.x, h.y)
			h.lastDrop = time.Now()
		}
		if h.x > -100 && h.x < config.ScreenWidth+100 {
			active = append(active, h)
		}
	}
	g.helicopters = active
}

func (g *Game) canDrop(x float32) bool {
	if x < config.ParatrooperWidth/2.0 ||
		x > config.ScreenWidth-config.ParatrooperWidth/2.0 {
		return false
	}
	baseX := (config.ScreenWidth - config.BaseWidth) / 2.0
	pX := x - config.ParatrooperWidth/2.0
	if utils.Overlap1D(pX, config.ParatrooperWidth, baseX, config.BaseWidth) {
		return false
	}
	for _, p := range g.paratroopers {
		if utils.Overlap1D(
			pX-1.0,
			config.ParatrooperWidth+2.0,
			p.x-config.ParatrooperWidth/2.0-1.0,
			config.ParatrooperWidth+2.0,
		) {
			return false
		}
	}
	return true
}
