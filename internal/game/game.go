package game

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/audio"
	"github.com/ystepanoff/paragopher/internal/config"
	"github.com/ystepanoff/paragopher/internal/utils"
)

type Game struct {
	Score        int
	gameData     *utils.GameData
	soundProfile *audio.SoundProfile

	showIntro     bool
	introStep     int
	lastIntroStep time.Time

	showExitDialog     bool
	showGameOverDialog bool
	showResolutionMenu bool

	resolutionMenuIdx    int
	currentResolutionIdx int
	fullscreen           bool

	barrelAngle            float64
	barrelImage            *ebiten.Image
	turretBaseImage        *ebiten.Image
	helicopterImage        *ebiten.Image
	paratrooperImage       *ebiten.Image
	paratrooperLandedImage *ebiten.Image
	paratrooperFellImage   *ebiten.Image
	bulletImage            *ebiten.Image

	bullets      []*Bullet
	lastShot     time.Time
	helicopters  []*Helicopter
	paratroopers []*Paratrooper
}

func NewGame() *Game {
	gameData, err := utils.LoadData()
	if err != nil {
		log.Println("Error loading game data!")
		gameData = &utils.GameData{}
	}
	game := &Game{
		bullets:              make([]*Bullet, 0),
		lastShot:             time.Now(),
		gameData:             gameData,
		soundProfile:         audio.NewSoundProfile(),
		showIntro:            true,
		resolutionMenuIdx:    config.DefaultResolutionIdx,
		currentResolutionIdx: config.DefaultResolutionIdx,
	}
	game.initTurretImage()
	game.initBarrelImage()
	game.initBulletImage()
	game.initHelicopterImage()
	game.initParatrooperImage()
	game.initParatrooperLandedImage()
	game.initParatrooperFellImage()
	game.initIntro()

	return game
}

// Ebiten Game Interface
func (g *Game) Draw(screen *ebiten.Image) {
	if g.showIntro {
		g.drawIntro(screen)
		if g.showResolutionMenu {
			g.drawResolutionMenu(screen)
		}
		return
	}
	g.drawTurret(screen)
	g.drawBullets(screen)
	g.drawHelicopters(screen)
	g.drawParatroopers(screen)

	// Display Score
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("SCORE: %d    HI-SCORE: %d", g.Score, g.gameData.HiScore),
	)

	if g.showExitDialog {
		showYesNoDialog(screen, "Do you want to exit the game?")
	}

	if g.showGameOverDialog {
		showYesNoDialog(screen, "GAME OVER!\nWould you like to start again?")
	}

	if g.showResolutionMenu {
		g.drawResolutionMenu(screen)
	}
}

func (g *Game) Update() error {
	if g.showResolutionMenu {
		g.updateResolutionMenu()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) &&
		!g.showExitDialog && !g.showGameOverDialog {
		g.resolutionMenuIdx = g.currentResolutionIdx
		g.showResolutionMenu = true
		return nil
	}
	if g.showIntro {
		return nil
	}
	if g.showExitDialog {
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			if err := utils.SaveData(g.gameData); err != nil {
				log.Fatalf("Failed to save game dada: %v", err)
			}
			return config.ErrQuit
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			g.showExitDialog = false
		}
		return nil
	}
	if g.showGameOverDialog {
		if err := utils.SaveData(g.gameData); err != nil {
			log.Fatalf("Failed to save game dada: %v", err)
		}
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			g.Reset()
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			return config.ErrQuit
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.showExitDialog = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.barrelAngle > config.BarrelAngleMin {
			g.barrelAngle -= config.BarrelAngleStep
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.barrelAngle < config.BarrelAngleMax {
			g.barrelAngle += config.BarrelAngleStep
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if time.Since(g.lastShot).Milliseconds() > config.ShotCooldown {
			g.shoot()
		}
	}

	g.updateBullets()
	g.spawnHelicopter()
	g.updateHelicopters()
	g.updateParatroopers()
	g.checkHits()

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) Reset() {
	g.soundProfile.GameOverPlayer.Pause()
	g.Score = 0
	g.showExitDialog = false
	g.showGameOverDialog = false
	g.barrelAngle = 0.0
	g.bullets = nil
	g.helicopters = nil
	g.paratroopers = nil
}

func showYesNoDialog(screen *ebiten.Image, message string) {
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(config.SemiTransparentBlack)
	screen.DrawImage(overlay, nil)

	dialogWidth, dialogHeight := 300, 150
	dialogX := (screen.Bounds().Dx() - dialogWidth) / 2
	dialogY := (screen.Bounds().Dy() - dialogHeight) / 2
	dialog := ebiten.NewImage(dialogWidth, dialogHeight)
	dialog.Fill(config.ColourDarkGrey)

	vector.DrawFilledRect(
		dialog,
		0,
		0,
		float32(dialogWidth),
		5,
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		0,
		float32(dialogHeight-5),
		float32(dialogWidth),
		5,
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		0,
		0,
		5,
		float32(dialogHeight),
		config.ColourBlack,
		false,
	)
	vector.DrawFilledRect(
		dialog,
		float32(dialogWidth-5),
		0,
		5,
		float32(dialogHeight),
		config.ColourBlack,
		false,
	)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dialogX), float64(dialogY))
	screen.DrawImage(dialog, op)

	textX := dialogX + 50
	textY := dialogY + 40
	ebitenutil.DebugPrintAt(screen, message, textX, textY)

	yesText := "Y: Yes"
	noText := "N: No"
	ebitenutil.DebugPrintAt(screen, yesText, dialogX+50, dialogY+90)
	ebitenutil.DebugPrintAt(screen, noText, dialogX+200, dialogY+90)
}
