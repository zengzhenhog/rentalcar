// use("coolcar");

db.account.createIndex({
  open_id: 1,
}, {
  unique: true
})

// 同一个account最多只能有一个进行中的Trip，通过建立索引实现
db.trip.createIndex({
  "trip.accountid": 1, // 1指从小到大
  "trip.status": 1
}, {
unique: true,
partialFilterExpression: {
    "trip.status": 1 // 1代表status值为 1 只能有一个
  }
})
// db.trip.dropIndex("trip.account_1_trip.status_1")