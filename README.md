# marsgame
прототип игры для программистов

пример кода, который управляет mech'ом:
```
obj = nearest(mech, objects)
angleObj = angle(mech.x, mech.y, obj.x, obj.y)
angleMech = mech.angle
angleTo = angleObj - angleMech
if angleTo < -PI {
   angleTo = 2. * PI + angleTo
}
if angleTo > PI {
   angleTo = angleTo - 2. * PI
}

switch {
case angleTo > 1.:
   mrThr = 1.
case angleTo < -1.:
   mrThr = -1.
default:
   mrThr = angleTo
}

distance = distance(mech.x, mech.y, obj.x, obj.y)
if distance > 200. {
   mThr = distance / 1000.
   if mThr > 1. {
      mThr = 1.
   }
}

if mrThr * mrThr * distance < 70. {
   shoot = 0.1
   return 1
}
```