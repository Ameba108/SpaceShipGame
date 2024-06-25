package main

import (
	"fmt"
	calc "game/src/calculation"
	objects "game/src/objects"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	pressedKey []ebiten.Key
	score      int
	Ship       calc.Ship
	Meteor     calc.Meteor
	Bullet     calc.Bullet
}

func (g *Game) Layout(outsideWitdh, outsideHeight int) (int, int) {
	return calc.ScreenWidth, calc.ScreenHeight
}

// Проверка координат пуль и координат метеоритов
func (g *Game) checkCollision() {
	for i := 0; i < len(g.Bullet.Bullets); i++ {
		bullet := g.Bullet.Bullets[i]
		for j := 0; j < len(g.Meteor.Meteorites); j++ {
			meteor := g.Meteor.Meteorites[j]
			// Игрок может попасть в метеорит только в том случае, если метеорит появляется в "поле зрения"
			if bullet.Position.X <= calc.ScreenWidth {
				if bullet.Position.X < meteor.Position.X+calc.MeteorWidth &&
					bullet.Position.X+calc.BulletWidth > meteor.Position.X &&
					bullet.Position.Y < meteor.Position.Y+calc.MeteorHeight &&
					bullet.Position.Y+calc.BulletHeight > meteor.Position.Y {
					// При попадании пули в метеорит, пуля и метеорит удаляется
					g.Meteor.Meteorites = append(g.Meteor.Meteorites[:j], g.Meteor.Meteorites[j+1:]...)
					g.Bullet.Bullets = append(g.Bullet.Bullets[:i], g.Bullet.Bullets[i+1:]...)
					g.score++
					i--
					break
				}
			}
		}
	}
}

// Функция проигрыша
func (g *Game) gameOver() bool {
	for i := 0; i < len(g.Meteor.Meteorites); i++ {
		enemy := g.Meteor.Meteorites[i]
		//Проверяем, совпадают ли координаты корабля, с координатами врага
		if calc.ShipPositionX < enemy.Position.X+calc.MeteorWidth &&
			calc.ShipPositionX+calc.ShipWidth > enemy.Position.X &&
			calc.ShipPositionY+calc.ShipHeight > enemy.Position.Y &&
			calc.ShipPositionY < enemy.Position.Y+calc.MeteorHeight {
			//Если координаты совпадают, поле shipDie в структуре Game изменяется на true, чтобы игра остановилась
			g.Ship.ShipDie = true
			return true
		}
	}
	return false
}

// Функция перезапуска игры
func (g *Game) Restart() {
	//Очищаем списки с метеоритами и пулями
	//чтобы начать игровой процесс заново
	g.Bullet.Bullets = g.Bullet.Bullets[:0]
	g.Meteor.Meteorites = g.Meteor.Meteorites[:0]
	//Обнуляем счетик очков
	g.score = 0
	//"Воскрешаем" корабль, чтобы игра снова заработала
	g.Ship.ShipDie = false
	//Возвращаем корабль на его позицию
	calc.ShipPositionX = float64(calc.ScreenWidth / 2)
	calc.ShipPositionY = float64(calc.ScreenHeight / 2)
	//Заново создаем метеориты
	for i := range g.Meteor.Meteorites {
		g.Meteor.Meteorites[i] = *calc.NewMeteor()
	}
	//Возобновляем скорость пуль
	for i := range g.Bullet.Bullets {
		g.Bullet.Bullets[i].Speed = g.Bullet.Speed
	}

}

func (g *Game) Update() error {
	pressedKeys := inpututil.AppendPressedKeys(g.pressedKey)
	for _, key := range pressedKeys {
		switch key.String() {
		case "R":
			if g.Ship.ShipDie {
				g.Restart()
			}
		}

	}

	// Обновление состояния корабля
	g.Ship.ShipUpdate(calc.ShipWidth, calc.ShipHeight)

	// При проигрыше все объекты останавливаются
	if g.gameOver() {
		g.Ship.ShipDie = true
		g.Meteor.Dead = true
		g.Bullet.Dead = true
	}

	// Если корабль "мертв", то и метеориты и пули замирают
	if g.Ship.ShipDie {
		g.Meteor.MeteorUpdate()
		g.Bullet.BulletUpdate()

	}

	// Обновление состояния метеоритов
	for time.Since(g.Meteor.LastMeteor) > time.Second*1 {
		g.Meteor.Meteorites = append(g.Meteor.Meteorites, *calc.NewMeteor())
		g.Meteor.LastMeteor = time.Now()
	}

	for v := range g.Meteor.Meteorites {
		g.Meteor.Meteorites[v].Position.X += g.Meteor.Meteorites[v].Speed.X
	}

	// Обновление состояния пуль
	if time.Since(g.Bullet.LastShotTime) > time.Millisecond*500 {
		g.Bullet.Bullets = append(g.Bullet.Bullets, *calc.NewBullet())
		g.Bullet.LastShotTime = time.Now()
	}

	for i := range g.Bullet.Bullets {
		g.Bullet.Bullets[i].Position.X += g.Bullet.Bullets[i].Speed
	}

	g.checkCollision()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// отрисовка фона
	spaceInstance := &objects.Objects{}
	spaceInstance.DrawBackground(screen)

	for _, bullet := range g.Bullet.Bullets {
		vector.DrawFilledRect(screen, float32(bullet.Position.X), float32(bullet.Position.Y), float32(calc.BulletWidth), float32(calc.BulletHeight), color.RGBA{255, 0, 0, 255}, false)
	}
	// отрисовка корабля
	spaceInstance.DrawSpaceShip(screen, float64(calc.ShipPositionX), float64(calc.ShipPositionY))

	// отрисовка метеоритов

	for _, meteor := range g.Meteor.Meteorites {
		spaceInstance.DrawMeteorit(screen, meteor.Position.X, meteor.Position.Y)
	}

	// счетчик количества очков
	face := basicfont.Face7x13
	if g.Ship.ShipDie {
		text.Draw(screen, "Game Over", face, calc.ScreenWidth/2-40, calc.ScreenHeight/2, color.White)
		text.Draw(screen, "Press 'R' to restart", face, calc.ScreenWidth/2-60, calc.ScreenHeight/2+16, color.White)
	}
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, face, 5, calc.ScreenHeight-5, color.White)
}

func main() {
	ebiten.SetWindowSize(calc.ScreenWidth*2, calc.ScreenHeight*2)
	ebiten.SetWindowTitle("Spaceship game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
