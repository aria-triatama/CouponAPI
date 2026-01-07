# CouponAPI

## Prerequisites

- Docker Desktop

## How to run

```bash
$ docker compose up -d --build
```

## How to test

```bash
$ curl -X POST http://localhost:1323/api/coupons -H "Content-Type: application/json" -d '{"name": "PROMO_SUPER", "amount": 100}'
```

```bash
$ curl -X POST http://localhost:1323/api/coupons/claim -H "Content-Type: application/json" -d '{"user_id": "user_12345", "coupon_name": "PROMO_SUPER"}'
```

```bash
$ curl -X GET http://localhost:1323/api/coupons/PROMO_SUPER
```

## Architecture note

Database design is simple, Coupons collection stores coupon documents and Claims collection stores claim documents. In MongoDB I use compound unique index to prevent duplicate claims and using transaction for claiming coupon to prevent race condition. I use aggregation pipeline to get coupon details with claimed by users.