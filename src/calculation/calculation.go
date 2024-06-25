package movement

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// размер окна
	ScreenHeight = 300
	ScreenWidth  = 600

	// корабль
	ShipWidth  = 25
	ShipHeight = 25

	// пули
	BulletHeight = ShipHeight / 5
	BulletWidth  = ShipWidth / 2
	bulletSpeedX = float64(3)

	// метеориты
	MeteorWidth  = 20
	MeteorHeight = 20
	MeteorSpeed  = float64(2.4)
)

var (
	// корабль
	ShipPositionX                = float64(ScreenWidth / 2)
	ShipPositionY                = float64(ScreenHeight / 2)
	shipAccewlerationSpeedUpMult = float64(2.5) // ускорение корабля (Пробел)
	shipAccelerationConstant     = float64(0.0000000015)
	shipResistance               = float64(0.975)
	prevUpdateTime               = time.Now()
)

type Vec struct {
	X, Y float64 // оси X и Y
}

// Движение корабля
type Ship struct {
	ShipDie          bool
	Speed            Vec          // скорость
	ShipAcceleration Vec          // движение корабля по оси X и Y
	pressedKey       []ebiten.Key // слайс из клавишь, нажатых пользователем
}

type GameInterface interface {
	Restart()
}

func SomeFunction(game GameInterface) {
	game.Restart()
}

func (s *Ship) ShipUpdate(width, height float64) error {
	timeDelta := float64(time.Since(prevUpdateTime))
	prevUpdateTime = time.Now()

	pressedKeys := inpututil.AppendPressedKeys(s.pressedKey)
	s.ShipAcceleration.X = 0
	s.ShipAcceleration.Y = 0

	// изменение скорости корабля в зависимости от его направления движения
	acc := shipAccelerationConstant

	for _, key := range pressedKeys {
		switch key.String() {
		case "Space":
			acc *= shipAccewlerationSpeedUpMult
		}
	}

	for _, key := range pressedKeys {
		switch key.String() {
		case "ArrowUp":
			s.ShipAcceleration.Y = -acc
		case "ArrowDown":
			s.ShipAcceleration.Y = acc
		case "ArrowLeft":
			s.ShipAcceleration.X = -acc
		case "ArrowRight":
			s.ShipAcceleration.X = acc
		}
	}
	if s.ShipDie {
		s.Speed.X *= 0
		s.Speed.Y *= 0
		s.ShipAcceleration.X *= 0
		s.ShipAcceleration.Y *= 0
	}

	s.Speed.X += s.ShipAcceleration.X
	s.Speed.Y += s.ShipAcceleration.Y

	s.Speed.X *= shipResistance // сопротивление по оси Х
	s.Speed.Y *= shipResistance // сопротивление по оси Y

	ShipPositionX += s.Speed.X * timeDelta
	ShipPositionY += s.Speed.Y * timeDelta

	var minX = float64(0)                     // самая маленькая координата по оси Х
	var minY = float64(0)                     // самая маленькая координата по оси Y
	var maxX = float64(ScreenWidth - width)   // самая большая координата по оси X
	var maxY = float64(ScreenHeight - height) // самая большая координата по оси Y

	// если позиция корабля < или > крайних координат, то корабль не уходит за границы, а "упирается в стену"
	if ShipPositionX >= maxX || ShipPositionX <= minX {
		if ShipPositionX > maxX {
			ShipPositionX = maxX
		} else if ShipPositionX < minX {
			ShipPositionX = minX
		}
		s.Speed.X *= -0.1
	}

	if ShipPositionY >= maxY || ShipPositionY <= minY {
		if ShipPositionY > maxY {
			ShipPositionY = maxY
		} else if ShipPositionY < minY {
			ShipPositionY = minY
		}
		s.Speed.Y *= -0.1
	}
	return nil
}

// Структура метеорита
type Meteor struct {
	Dead       bool
	Meteorites []Meteor
	LastMeteor time.Time
	Position   Vec
	Speed      Vec
}

// Создание метеорита
func NewMeteor() *Meteor {
	return &Meteor{
		Position: Vec{
			X: ScreenWidth + 300,
			Y: float64(rand.Intn(ScreenHeight - MeteorHeight)),
		},
		Speed: Vec{X: -MeteorSpeed}}
}

func (m *Meteor) MeteorUpdate() {
	if m.Dead {
		for i := range m.Meteorites {
			m.Meteorites[i].Speed.X = 0
		}
	}
}

// Структура пули
type Bullet struct {
	Dead         bool
	Position     Vec
	Speed        float64
	LastShotTime time.Time
	Bullets      []Bullet
}

func NewBullet() *Bullet {
	return &Bullet{
		Position: Vec{
			X: ShipPositionX,
			Y: ShipPositionY + 10,
		},
		Speed: bulletSpeedX,
	}
}

func (b *Bullet) BulletUpdate() {
	if b.Dead {
		for i := range b.Bullets {
			b.Bullets[i].Speed = 0
		}
	}
}
