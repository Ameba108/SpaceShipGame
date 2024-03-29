package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	//Размер окна
	screenHeight = 300
	screenWidth  = 600
	//Размер корабля
	shipHeight = 25
	shipWidth  = 25

	//Размер пуль
	bulletHeight = shipHeight / 5
	bulletWidth  = shipWidth / 2

	//Базовая скорость ускорения
	shipAccelerationConstant = float64(0.0000000015)
	//Множитель ускорения
	shipAccelerationSpeedUpMultiplier = float64(2)
	//Сопротивление, благодаря которому корабль замедляется при движении
	shipResistence = float64(0.975)

	//Размер и скорость метеоритов
	meteorWidth  = 20
	meteorHeight = 20
	meteorSpeed  = float64(2.4)
)

var (
	//Позиция корабля
	shipPositionX = float64(screenWidth / 2)
	shipPositionY = float64(screenHeight / 2)
	//Движение корабля по оси X или Y
	shipMovementX = float64(0)
	shipMovementY = float64(0)

	//Текущее ускорение корабля по оси X и Y
	shipAccelerationX = float64(0)
	shipAccelerationY = float64(0)

	//Последнее обновление игры
	prevUpdateTime = time.Now()

	//Скорость пуль
	bulletSpeedX = float64(3)

	//Переменные, хранящие изображение для объектов
	spaceImage     *ebiten.Image
	spaceshipImage *ebiten.Image
	meteorImage    *ebiten.Image
)

// Функция, загружающая файл с изображением корабля
func init() {
	var err error
	spaceshipImage, _, err = ebitenutil.NewImageFromFile("textures/spaceship.png")
	if err != nil {
		log.Fatal(err)
	}
}

// Функция, загружающая файл с изображением космоса
func init() {
	var err error
	spaceImage, _, err = ebitenutil.NewImageFromFile("textures/space.png")
	if err != nil {
		log.Fatal(err)
	}
}

// Функция, загружающая файл с изображением метеорита
func init() {
	var err error
	meteorImage, _, err = ebitenutil.NewImageFromFile("textures/meteor.png")
	if err != nil {
		log.Fatal(err)
	}
}

// Структура пуль (их позииция, скорость)
type Bullet struct {
	PositionX float64
	PositionY float64
	SpeedX    float64
}

// Структура метеоритов (их позиция, скорость)
type Meteor struct {
	PositionX float64
	PositionY float64
	SpeedX    float64
}

type Game struct {
	//количество очков
	score int
	//смерть корабля (если значение == true - игра заканчивается)
	shipDie     bool
	pressedKeys []ebiten.Key
	//фиксируем время, с которым появляется пуля, чтобы определенное количество пуль появлялось в свое время
	//а не в разнобой
	lastShotTime time.Time
	//Точно такая же фиксация времени, с которым появляется враг
	lastEnemyTime time.Time
	//Список из количества существующих пуль
	bullets []Bullet
	//Список из количества существующиз метеоритов
	meteorites []Meteor
}

// Функция, создающая метеориты
func NewMeteor() *Meteor {
	return &Meteor{
		//Метеориты появляется дальше видимого игрового поля
		//чтобы у игрока было время на подготовку к игре
		PositionX: screenWidth + 300,
		PositionY: float64(rand.Intn(screenHeight - meteorHeight)),
		SpeedX:    -meteorSpeed,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Функция, при которой метеорит, в которого попадает пуля, удалется
func (g *Game) checkCollision() {
	for i := 0; i < len(g.bullets); i++ {
		bullet := g.bullets[i]
		for j := 0; j < len(g.meteorites); j++ {
			meteor := g.meteorites[j]
			//Проверяем, не совпадают ли координаты пулиь и метеорита
			//только в том случае, если пуля находится в пределах игрового поля
			//чтобы игрок не мог разрушать метеориты еще до того, как они появятся на окне
			if bullet.PositionX <= screenWidth {
				if bullet.PositionX < meteor.PositionX+meteorWidth &&
					bullet.PositionX+bulletWidth > meteor.PositionX &&
					bullet.PositionY < meteor.PositionY+meteorHeight &&
					bullet.PositionY+bulletHeight > meteor.PositionY {
					//Если координаты совпали, то мы удаляем и пулю, и метеорит, в который попала пуля, из списка
					//Чтобы они исчезли с игровго экрана
					g.meteorites = append(g.meteorites[:j], g.meteorites[j+1:]...)
					g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
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
	for i := 0; i < len(g.meteorites); i++ {
		enemy := g.meteorites[i]
		//Проверяем, совпадают ли координаты корабля, с координатами врага
		if shipPositionX < enemy.PositionX+meteorWidth &&
			shipPositionX+shipWidth > enemy.PositionX &&
			shipPositionY+shipHeight > enemy.PositionY &&
			shipPositionY < enemy.PositionY+meteorHeight {
			//Если координаты совпадают, поле shipDie в структуре Game изменяется на true, чтобы игра остановилась
			g.shipDie = true
			return true
		}
	}
	return false
}

// Функция перезапуска игры
func (g *Game) restart() {
	//Очищаем списки с метеоритами и пулями
	//чтобы начать игровой процесс заново
	g.bullets = g.bullets[:0]
	g.meteorites = g.meteorites[:0]
	//Обнуляем счетик очков
	g.score = 0
	//"Воскрешаем" корабль, чтобы игра снова заработала
	g.shipDie = false
	//Возвращаем корабль на его позицию
	shipPositionX = float64(screenWidth / 2)
	shipPositionY = float64(screenHeight / 2)
	//Заново создаем метеориты
	for i := range g.meteorites {
		g.meteorites[i] = *NewMeteor()
	}
	//Возобновляем скорость пуль
	for i := range g.bullets {
		g.bullets[i].SpeedX = bulletSpeedX
	}

}

func (g *Game) Update() error {
	//Дельта времени последнего обновления
	timeDelta := float64(time.Since(prevUpdateTime))
	prevUpdateTime = time.Now()

	//Определяем, какие клавиши нажимает игрок
	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])
	shipAccelerationX = 0
	shipAccelerationY = 0

	//Базовая скорость корабля
	acc := shipAccelerationConstant

	//Уравление кораблём
	for _, key := range g.pressedKeys {
		switch key.String() {
		//Если игрок нажимает пробел, то корабль ускоряется
		case "Space":
			acc *= shipAccelerationSpeedUpMultiplier
		}
	}
	for _, key := range g.pressedKeys {
		switch key.String() {
		//Если игрок нажимает клавишу стрелки вверх, то и корабль соответственно поднимается
		case "ArrowDown":
			shipAccelerationY = acc
		//"Стрелка вниз"-корабль летит вниз
		case "ArrowUp":
			shipAccelerationY = -acc
		//"Стрелка вправо"-корабль летит вправо
		case "ArrowRight":
			shipAccelerationX = acc
		//"Стрелка влево"-корабль летит влево
		case "ArrowLeft":
			shipAccelerationX = -acc
		//Клавиша "R" отвечает за рестарт игры. Но она работает только в том случае, если игрок проиграл
		case "R":
			if g.shipDie {
				g.restart()
			}
		}
	}
	//Если игрок проиграл, изменяем состояние корабля
	if g.gameOver() {
		g.shipDie = true
	}

	//Если корабль мертв, обнуляем скорость у всех объектов и не даем кораблю двигаться, чтобы вся игра замерла
	if g.shipDie {
		shipMovementX = 0
		shipMovementY = 0
		shipAccelerationX = 0
		shipAccelerationY = 0
		for i := range g.meteorites {
			g.meteorites[i].SpeedX = 0
		}
		for i := range g.bullets {
			g.bullets[i].SpeedX = 0
		}
	}

	//Скорость движения по оси X и Y
	shipMovementY += shipAccelerationY
	shipMovementX += shipAccelerationX
	shipMovementX *= shipResistence
	shipMovementY *= shipResistence
	shipPositionX += shipMovementX * timeDelta
	shipPositionY += shipMovementY * timeDelta

	//Константы с минимальным значением координат X и Y (соответственно - 0)
	//и с максимальным значением (конец игрового поля)
	const minX = 0
	const minY = 0
	const maxX = screenWidth - shipWidth
	const maxY = screenHeight - shipHeight

	//Если корабль "соприкосается" с границами окна, то дальше он пройти не сможет
	if shipPositionX >= maxX || shipPositionX <= minX {
		if shipPositionX > maxX {
			shipPositionX = maxX
		} else if shipPositionX < minX {
			shipPositionX = minX
		}
		shipMovementX *= -1
	}

	if shipPositionY >= maxY || shipPositionY <= minY {
		if shipPositionY > maxY {
			shipPositionY = maxY
		} else if shipPositionY < minY {
			shipPositionY = minY
		}
		shipMovementY *= -1
	}

	//Создаем пули
	if time.Since(g.lastShotTime) > time.Millisecond*500 {
		g.bullets = append(g.bullets, Bullet{
			PositionX: shipPositionX,
			PositionY: shipPositionY + 10,
			SpeedX:    bulletSpeedX,
		})

		g.lastShotTime = time.Now()
	}

	for i := range g.bullets {
		g.bullets[i].PositionX += g.bullets[i].SpeedX
	}

	//Создаем врагов
	for time.Since(g.lastEnemyTime) > time.Second*1 {
		g.meteorites = append(g.meteorites, *NewMeteor())
		g.lastEnemyTime = time.Now()
	}

	//Позиция метеоритов по оси меняется в зависимости от их скорости
	for v := range g.meteorites {
		g.meteorites[v].PositionX += g.meteorites[v].SpeedX
	}

	//Проверяем, попала ли пуля в метеорит
	g.checkCollision()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Прикрепляем фоновое изображение с космосом
	screen.DrawImage(spaceImage, nil)
	//Ставим изображение корабля под его координаты
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shipPositionX-15, shipPositionY-12)
	//Рисуем пули с помощью функции из библиотеки
	for _, bullet := range g.bullets {
		vector.DrawFilledRect(screen, float32(bullet.PositionX), float32(bullet.PositionY), float32(bulletWidth), float32(bulletHeight), color.RGBA{255, 0, 0, 255}, false)
	}
	//Рисуем корабль поверх пуль
	screen.DrawImage(spaceshipImage, op)

	//Рисуем метеориты
	for _, meteor := range g.meteorites {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(meteor.PositionX-5, meteor.PositionY-5)
		screen.DrawImage(meteorImage, opt)
	}

	//определяем шрифт
	face := basicfont.Face7x13

	//При проигрыше появляется соответствующая надпись
	if g.shipDie {
		text.Draw(screen, "Game Over", face, screenWidth/2-40, screenHeight/2, color.White)
		text.Draw(screen, "Press 'R' to restart", face, screenWidth/2-60, screenHeight/2+16, color.White)
	}

	//Выводим счетчик очков на кэкран
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, face, 5, screenHeight-5, color.White)
}

func main() {
	//Создаем окно
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Spaceship game")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
