# marsgame
прототип игры для программистов

## Установка
`go build`

## Запуск сервера
`go run main.go`

пример кода, который управляет mech'ом:
```
xelon = getFirstTarget(1)
if empty(xelon) {
   xelon = nearestByType(mech, objects, ObjectTypes:xelon)
   if empty(xelon) {
      return 1
   }
   addTarget(xelon, 1)
}
angleTo = angleToRotate(mech.angle, mech.x, mech.y, xelon.x, xelon.y)
commands.rotate = keepBounds(angleTo, 1.)

commands.move = 1. - absFloat(commands.rotate)

obj = getFirstTarget(2)
if empty(obj) {
   obj = nearestByType(mech, objects, ObjectTypes:spore)
   if empty(obj) {
      return 1
   }
   addTarget(obj, 2)
}

angleSum = mech.angle + mech.cAngle
cAngleTo = angleToRotate(angleSum, mech.x, mech.y, obj.x, obj.y)

if cAngleTo * angleTo < 0. {
   cAngleTo = cAngleTo - angleTo
}
commands.cannon.rotate = keepBounds(cAngleTo, 1.)

dist = distance(mech.x, mech.y, obj.x, obj.y)
toShoot = cAngleTo * cAngleTo * dist
if toShoot < 40. {
   commands.cannon.shoot = 0.1
}
```