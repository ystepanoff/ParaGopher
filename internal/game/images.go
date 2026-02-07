package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
)

func (g *Game) initTurretImage() {
	w := int(config.BaseWidth)
	h := int(config.BaseWidth/3.0 + config.BaseHeight + 1.0)
	g.turretBaseImage = ebiten.NewImage(w, h)
	g.turretBaseImage.Fill(config.TransparentBlack)
	vector.DrawFilledRect(
		g.turretBaseImage,
		0.0,
		float32(h)-config.BaseHeight,
		config.BaseWidth,
		config.BaseHeight,
		config.ColourWhite,
		true,
	)
	vector.DrawFilledRect(
		g.turretBaseImage,
		config.BaseWidth/3.0,
		0.0,
		config.BaseWidth/3.0,
		config.BaseWidth/3.0,
		config.ColourPink,
		false,
	)
}

func (g *Game) initBarrelImage() {
	w := config.BaseWidth
	g.barrelImage = ebiten.NewImage(int(w), int(w))
	g.barrelImage.Fill(config.TransparentBlack)

	rectX := w/2 - w/12
	rectY := w / 12
	rectW := w / 6
	rectH := w / 3
	vector.DrawFilledRect(
		g.barrelImage,
		rectX,
		rectY,
		rectW,
		rectH,
		config.ColourTeal,
		false,
	)

	circleX := w / 2
	circleY := w / 2
	pinkCircleRadius := w / 6
	tealCircleRaduis := w / 24
	vector.DrawFilledCircle(
		g.barrelImage,
		circleX,
		circleY,
		pinkCircleRadius,
		config.ColourPink,
		true,
	)
	vector.DrawFilledCircle(
		g.barrelImage,
		circleX,
		circleY,
		tealCircleRaduis,
		config.ColourTeal,
		true,
	)
	topCircleX, topCircleY := w/2, w/12
	topCircleRadius := w / 12
	vector.DrawFilledCircle(
		g.barrelImage,
		topCircleX,
		topCircleY,
		topCircleRadius,
		config.ColourTeal,
		true,
	)
}

func (g *Game) initBulletImage() {
	w := int(2 * config.BulletRadius)
	g.bulletImage = ebiten.NewImage(w, w)
	vector.DrawFilledCircle(
		g.bulletImage,
		config.BulletRadius,
		config.BulletRadius,
		config.BulletRadius,
		config.ColourWhite,
		true,
	)
}

func (g *Game) initHelicopterImage() {
	w := int(config.HelicopterBodyWidth + config.HelicopterTailWidth)
	h := int(config.HelicopterBodyHeight) + 6
	g.helicopterImage = ebiten.NewImage(w, h)
	tailX := float32(0.0)
	tailY := float32(h-config.HelicopterTailHeight) / 2
	bodyX := float32(config.HelicopterTailWidth)
	bodyY := float32(h-config.HelicopterBodyHeight) / 2

	vector.DrawFilledRect(
		g.helicopterImage,
		tailX,
		tailY,
		config.HelicopterTailWidth,
		config.HelicopterTailHeight,
		config.ColourTeal,
		false,
	)

	vector.DrawFilledRect(
		g.helicopterImage,
		bodyX,
		bodyY,
		config.HelicopterBodyWidth,
		config.HelicopterBodyHeight,
		config.ColourTeal,
		false,
	)

	bodyCenterX := bodyX + config.HelicopterBodyWidth/2.0
	bodyTopY := bodyY
	rotorStartX := bodyCenterX - config.HelicopterRotorLen/2.0
	rotorStartY := bodyTopY - 2.0
	rotorEndX := bodyCenterX + config.HelicopterRotorLen/2.0
	rotorEndY := rotorStartY
	vector.StrokeLine(
		g.helicopterImage,
		rotorStartX,
		rotorStartY,
		rotorEndX,
		rotorEndY,
		1.0,
		config.ColourMagenta,
		false,
	)
}

// An ugly hack until vector.DrawFilledPath is available in Ebitengine
func DrawFilledSemicircle(
	screen *ebiten.Image,
	centerX, centerY, radius float32,
	startAngle, endAngle float32,
	clr color.Color,
) {
	segments := 180 // Number of triangles to approximate the semicircle
	angleStep := (endAngle - startAngle) / float32(segments)

	vertices := make([]ebiten.Vertex, (segments+1)*3)
	indices := make([]uint16, segments*3)

	for i := 0; i < segments; i++ {
		theta1 := float64((startAngle + float32(i)*angleStep) * math.Pi / 180)
		theta2 := float64((startAngle + float32(i+1)*angleStep) * math.Pi / 180)

		v0 := ebiten.Vertex{
			DstX:   centerX,
			DstY:   centerY,
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(clr.(color.RGBA).R) / 255,
			ColorG: float32(clr.(color.RGBA).G) / 255,
			ColorB: float32(clr.(color.RGBA).B) / 255,
			ColorA: float32(clr.(color.RGBA).A) / 255,
		}

		v1 := ebiten.Vertex{
			DstX:   centerX + radius*float32(math.Cos(theta1)),
			DstY:   centerY + radius*float32(math.Sin(theta1)),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(clr.(color.RGBA).R) / 255,
			ColorG: float32(clr.(color.RGBA).G) / 255,
			ColorB: float32(clr.(color.RGBA).B) / 255,
			ColorA: float32(clr.(color.RGBA).A) / 255,
		}

		v2 := ebiten.Vertex{
			DstX:   centerX + radius*float32(math.Cos(theta2)),
			DstY:   centerY + radius*float32(math.Sin(theta2)),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(clr.(color.RGBA).R) / 255,
			ColorG: float32(clr.(color.RGBA).G) / 255,
			ColorB: float32(clr.(color.RGBA).B) / 255,
			ColorA: float32(clr.(color.RGBA).A) / 255,
		}

		vertices[i*3] = v0
		vertices[i*3+1] = v1
		vertices[i*3+2] = v2

		indices[i*3] = uint16(i * 3)
		indices[i*3+1] = uint16(i*3 + 1)
		indices[i*3+2] = uint16(i*3 + 2)
	}

	meshImg := ebiten.NewImage(1, 1)
	meshImg.Fill(config.ColourWhite)

	screen.DrawTriangles(vertices, indices, meshImg, nil)
}

func (g *Game) initParatrooperImage() {
	w := int(math.Max(
		float64(config.ParachuteRadius*2.0),
		float64(config.ParatrooperWidth),
	))
	h := int(config.ParachuteRadius*2 + config.ParatrooperHeight)
	g.paratrooperImage = ebiten.NewImage(w, h)

	cw := float32(w)
	canopyR := config.ParachuteRadius
	canopyCX := cw / 2.0
	canopyCY := canopyR

	DrawFilledSemicircle(
		g.paratrooperImage,
		canopyCX, canopyCY, canopyR,
		-180.0, 0.0,
		config.ColourTeal,
	)

	bodyStartY := canopyR * 2.0
	bodyCX := cw / 2.0

	headR := float32(2.0)
	headCY := bodyStartY + headR + 0.5
	vector.DrawFilledCircle(
		g.paratrooperImage,
		bodyCX, headCY, headR,
		config.ColourTeal, true,
	)

	torsoW := float32(4.0)
	torsoH := float32(6.0)
	torsoX := bodyCX - torsoW/2.0
	torsoY := headCY + headR
	vector.DrawFilledRect(
		g.paratrooperImage,
		torsoX, torsoY, torsoW, torsoH,
		config.ColourTeal, false,
	)

	shoulderY := torsoY + 1.5
	leftHandX := bodyCX - 5.0
	rightHandX := bodyCX + 5.0
	handY := headCY - 1.5
	vector.StrokeLine(
		g.paratrooperImage,
		torsoX, shoulderY,
		leftHandX, handY,
		1.5, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperImage,
		torsoX+torsoW, shoulderY,
		rightHandX, handY,
		1.5, config.ColourTeal, false,
	)

	hipY := torsoY + torsoH
	vector.StrokeLine(
		g.paratrooperImage,
		bodyCX-1.0, hipY,
		bodyCX-2.5, float32(h),
		1.5, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperImage,
		bodyCX+1.0, hipY,
		bodyCX+2.5, float32(h),
		1.5, config.ColourTeal, false,
	)

	canopyBottomY := canopyCY
	vector.StrokeLine(
		g.paratrooperImage,
		canopyCX-canopyR+2, canopyBottomY,
		leftHandX, handY,
		0.75, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperImage,
		canopyCX+canopyR-2, canopyBottomY,
		rightHandX, handY,
		0.75, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperImage,
		canopyCX-canopyR/3.0, canopyBottomY,
		torsoX, shoulderY,
		0.75, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperImage,
		canopyCX+canopyR/3.0, canopyBottomY,
		torsoX+torsoW, shoulderY,
		0.75, config.ColourTeal, false,
	)
}

func (g *Game) initParatrooperLandedImage() {
	w := int(config.ParatrooperWidth)
	h := int(config.ParatrooperHeight)
	g.paratrooperLandedImage = ebiten.NewImage(w, h)

	cx := float32(w) / 2.0

	headR := float32(2.0)
	headCY := headR + 0.5
	vector.DrawFilledCircle(
		g.paratrooperLandedImage,
		cx, headCY, headR,
		config.ColourTeal, true,
	)

	torsoW := float32(4.0)
	torsoH := float32(6.0)
	torsoX := cx - torsoW/2.0
	torsoY := headCY + headR + 0.5
	vector.DrawFilledRect(
		g.paratrooperLandedImage,
		torsoX, torsoY, torsoW, torsoH,
		config.ColourTeal, false,
	)

	shoulderY := torsoY + 1.5
	vector.StrokeLine(
		g.paratrooperLandedImage,
		torsoX, shoulderY,
		1.0, torsoY+torsoH-1,
		1.5, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperLandedImage,
		torsoX+torsoW, shoulderY,
		float32(w)-1.0, torsoY+torsoH-1,
		1.5, config.ColourTeal, false,
	)

	hipY := torsoY + torsoH
	vector.StrokeLine(
		g.paratrooperLandedImage,
		cx-1.0, hipY,
		1.5, float32(h),
		1.5, config.ColourTeal, false,
	)
	vector.StrokeLine(
		g.paratrooperLandedImage,
		cx+1.0, hipY,
		float32(w)-1.5, float32(h),
		1.5, config.ColourTeal, false,
	)
}

func (g *Game) initParatrooperFellImage() {
	w := int(config.ParatrooperWidth)
	h := int(config.ParatrooperHeight)
	g.paratrooperFellImage = ebiten.NewImage(w, h)

	cx := float32(w) / 2.0
	clr := config.ColourLightGrey

	headR := float32(2.0)
	headCY := headR + 0.5
	vector.DrawFilledCircle(
		g.paratrooperFellImage,
		cx, headCY, headR,
		clr, true,
	)

	torsoW := float32(4.0)
	torsoH := float32(6.0)
	torsoX := cx - torsoW/2.0
	torsoY := headCY + headR + 0.5
	vector.DrawFilledRect(
		g.paratrooperFellImage,
		torsoX, torsoY, torsoW, torsoH,
		clr, false,
	)

	shoulderY := torsoY + 1.5
	vector.StrokeLine(
		g.paratrooperFellImage,
		torsoX, shoulderY,
		0, torsoY+3,
		1.5, clr, false,
	)
	vector.StrokeLine(
		g.paratrooperFellImage,
		torsoX+torsoW, shoulderY,
		float32(w), torsoY+3,
		1.5, clr, false,
	)

	hipY := torsoY + torsoH
	vector.StrokeLine(
		g.paratrooperFellImage,
		cx-1.0, hipY,
		0, float32(h),
		1.5, clr, false,
	)
	vector.StrokeLine(
		g.paratrooperFellImage,
		cx+1.0, hipY,
		float32(w), float32(h),
		1.5, clr, false,
	)
}
