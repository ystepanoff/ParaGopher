package config

import (
	"errors"
	"image/color"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	BarrelAngleMin  = -70.0
	BarrelAngleMax  = 70.0
	BarrelAngleStep = 0.5

	HelicopterSpawnChance = 0.004
	HelicopterSpeed       = 0.75
	HelicopterDropRate    = 2
	HelicopterBodyWidth   = 30.0
	HelicopterBodyHeight  = 10.0
	HelicopterTailWidth   = 16.0
	HelicopterTailHeight  = 3.0
	HelicopterRotorLen    = 36.0

	ParatrooperSpawnChance = 0.008
	ParatrooperFallSpeed   = 0.6
	ParatrooperWalkSpeed   = 0.25

	BulletSpeed  = 5.0
	BulletRadius = 2.0
	ShotCooldown = 200

	GroundY = 600
)

type Resolution struct {
	Width  int
	Height int
	Label  string
}

var Resolutions = []Resolution{
	{640, 480, "640 x 480"},
	{800, 600, "800 x 600"},
	{1024, 768, "1024 x 768"},
	{1280, 960, "1280 x 960"},
	{1600, 1200, "1600 x 1200"},
}

const DefaultResolutionIdx = 1

var (
	BaseWidth  = float32(ScreenWidth) / 10.0
	BaseHeight = float32(ScreenHeight) / 10.0

	ParatrooperWidth  = float32(12.0)
	ParatrooperHeight = BaseHeight / 3.0
	ParachuteRadius   = float32(14.0)

	ColourTeal           = color.RGBA{R: 101, G: 247, B: 246, A: 255}
	ColourPink           = color.RGBA{R: 255, G: 82, B: 242, A: 255}
	ColourMagenta        = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	ColourWhite          = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColourBlack          = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	ColourLightGrey      = color.RGBA{R: 50, G: 50, B: 50, A: 255}
	ColourDarkGrey       = color.RGBA{R: 25, G: 25, B: 25, A: 255}
	TransparentBlack     = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	SemiTransparentBlack = color.RGBA{R: 0, G: 0, B: 0, A: 225}

	ErrQuit = errors.New("user quit the game")
)
