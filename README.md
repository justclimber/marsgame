# marsgame
прототип игры для программистов

пример кода, который управляет mech'ом:
```
ifempty obj = nearestByType(mech, objects, 3) {
   return 1
}
angleObj = angle(mech.x, mech.y, obj.x, obj.y)
angleMech = mech.angle
angleTo = angleObj - angleMech
if angleTo < -PI {
   angleTo = 2. * PI + angleTo
}
if angleTo > PI {
   angleTo = angleTo - 2. * PI
}

switch angleTo {
case > 1.:
   commands.rotate = 1.
case < -1.:
   commands.rotate = -1.
default:
   commands.rotate = angleTo
}

dist = distance(mech.x, mech.y, obj.x, obj.y)
if obj.type == 3 {
   commands.move = 1.
   return 1
}
if dist > 200. {
   commands.move = distance / 1000.
   if commands.move > 1. {
      commands.move = 1.
   }
}
if angleTo * angleTo * dist < 70. {
   commands.cannon.shoot = 0.1
   return 1
}
```