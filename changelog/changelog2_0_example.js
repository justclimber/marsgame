const example = {
  timeIds: [1112, 1215, 1318, 1415, 1520, 1625],
  objs: [
    {
      id: 1,
      type: "player",
      times: [
        {
          timeId: 1112,
          x: 550,
          y: 10,
          rotate: 1,
          // вектор изменений в секунду с учетом поворота вектора
          posFuture: {
            x: 100,
            y: 0,
            rotate: 0,
            untilTimeId: 1415
          },
        },
        {
          timeId: 1215,
          cannonRotate : 2,
          // изменения поворота башни в будущем
          cannonFuture: {
            rotate: 0.2,
            untilTimeId: 1625
          }
        },
        {
          timeId: 1318,
          fire: true
        },
        {
          timeId: 1520,
          x: 610,
          y: 15,
          deleteOtherObjectId: 15
        },
      ]
    },
    {
      id: 12,
      type: "missile",
      times: [
        {
          timeId: 1112,
          x: 150,
          y: 20,
          rotate: 1.3,
          // вектор изменений в секунду с учетом поворота вектора
          posChangeVector: {
            x: 100,
            y: 100,
            rotate: 0,
            untilTimeId: 1415
          },
        },
        {
          timeId: 1415,
          delete: true,
          explodeOtherObjectId: 123
        },
      ]
    }
  ]
};