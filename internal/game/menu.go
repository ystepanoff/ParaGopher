package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ystepanoff/paragopher/internal/config"
)

func (g *Game) updateResolutionMenu() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.resolutionMenuIdx > 0 {
			g.resolutionMenuIdx--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.resolutionMenuIdx < len(config.Resolutions)-1 {
			g.resolutionMenuIdx++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		res := config.Resolutions[g.resolutionMenuIdx]
		ebiten.SetWindowSize(res.Width, res.Height)
		g.currentResolutionIdx = g.resolutionMenuIdx
		if g.fullscreen {
			g.fullscreen = false
			ebiten.SetFullscreen(false)
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}
		g.showResolutionMenu = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
		if g.fullscreen {
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}
		g.showResolutionMenu = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showResolutionMenu = false
	}
}

func (g *Game) drawResolutionMenu(screen *ebiten.Image) {
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(config.SemiTransparentBlack)
	screen.DrawImage(overlay, nil)

	dialogWidth := 300
	dialogHeight := 50 + len(config.Resolutions)*20 + 60
	dialogX := (screen.Bounds().Dx() - dialogWidth) / 2
	dialogY := (screen.Bounds().Dy() - dialogHeight) / 2

	dialog := ebiten.NewImage(dialogWidth, dialogHeight)
	dialog.Fill(config.ColourDarkGrey)

	borderSize := float32(5)
	vector.DrawFilledRect(dialog, 0, 0, float32(dialogWidth), borderSize, config.ColourBlack, false)
	vector.DrawFilledRect(dialog, 0, float32(dialogHeight)-borderSize, float32(dialogWidth), borderSize, config.ColourBlack, false)
	vector.DrawFilledRect(dialog, 0, 0, borderSize, float32(dialogHeight), config.ColourBlack, false)
	vector.DrawFilledRect(dialog, float32(dialogWidth)-borderSize, 0, borderSize, float32(dialogHeight), config.ColourBlack, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dialogX), float64(dialogY))
	screen.DrawImage(dialog, op)

	ebitenutil.DebugPrintAt(screen, "RESOLUTION", dialogX+105, dialogY+15)

	for i, res := range config.Resolutions {
		cursor := "  "
		if i == g.resolutionMenuIdx {
			cursor = "> "
		}
		active := ""
		if i == g.currentResolutionIdx && !g.fullscreen {
			active = " *"
		}
		label := fmt.Sprintf("%s%s%s", cursor, res.Label, active)
		ebitenutil.DebugPrintAt(screen, label, dialogX+70, dialogY+45+i*20)
	}

	fullscreenLabel := "F: Toggle Fullscreen"
	if g.fullscreen {
		fullscreenLabel = "F: Toggle Fullscreen *"
	}
	footerY := dialogY + 45 + len(config.Resolutions)*20 + 10
	ebitenutil.DebugPrintAt(screen, fullscreenLabel, dialogX+55, footerY)
	ebitenutil.DebugPrintAt(screen, "ENTER: Select  ESC: Back", dialogX+40, footerY+20)
}
