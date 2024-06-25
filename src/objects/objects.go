package objects

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	spaceImage     *ebiten.Image
	spaceShipImage *ebiten.Image
	meteorImage    *ebiten.Image
)

type Objects struct{}

// отрисовка фона
func init() {
	var err error
	spaceImage, _, err = ebitenutil.NewImageFromFile("textures/space.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (o Objects) DrawBackground(screen *ebiten.Image) {
	screen.DrawImage(spaceImage, nil)
}

// отрисовка корабля
func init() {
	var err error
	spaceShipImage, _, err = ebitenutil.NewImageFromFile("textures/spaceship.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (o Objects) DrawSpaceShip(screen *ebiten.Image, shipPositionX, shipPositionY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shipPositionX-15, shipPositionY-12)
	screen.DrawImage(spaceShipImage, op)
}

// отрисовка метеоритов
func init() {
	var err error
	meteorImage, _, err = ebitenutil.NewImageFromFile("textures/meteor.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (o Objects) DrawMeteorit(screen *ebiten.Image, meteorPositionX, meteorPositionY float64) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(meteorPositionX-5, meteorPositionY-5)
	screen.DrawImage(meteorImage, opt)
}

// отрисовка пуль
func (o Objects) DrawBullet(screen *ebiten.Image, bulletWidth, bulletHeight, bulletPositionX, bulletPositionY float64) {
	vector.DrawFilledRect(screen, float32(bulletPositionX), float32(bulletPositionY), float32(bulletWidth), float32(bulletHeight), color.RGBA{255, 0, 0, 255}, false)
}
