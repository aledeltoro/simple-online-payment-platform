# API resources

## Online Payment Platform API

### Ping

<details>
 <summary><code>GET</code> <code><b>/</b></code> <code>(Check if the service is healthy)</code></summary>

#### Parameters

> None

#### Responses

##### HTTP Code 200

```json
Hello World!
```

</details>

### Create payment

**Disclaimer**: When a new payment is created its initial status is intentionally set to `pending`. In order to mock use case where it takes X amount of time to charge a payment. Therefore, the final status will be given by the event receved by webhooks.

<details>
 <summary><code>POST</code> <code><b>/</b></code> <code>(creates a successful or unsuccessful payment, depending on the test card used)</code></summary>

#### Parameters

> | name            |  type     | data type               | description                                              |
> |-----------------|-----------|-------------------------|----------------------------------------------------------|
> | amount          |  required | string (urlencoded)     | Amount to perform a payment                              |
> | currency        |  required | string (urlencoded)     | Currency to perform a payment                            |
> | payment_method  |  required | string (urlencoded)     | Method to perform payment, refers to Stripe's test cards |
> | description     |  optional | string (urlencoded)     | Description on what the payment is about                 |

#### Responses

##### HTTP Code 200

Succesful payment

```json
{
  "transaction_id": "TXN_01HP06ZRSNFDPKN3ZBSWS4Z0KT",
  "status": "pending",
  "description": "Sample transaction",
  "payment_provider": "stripe",
  "amount": 2000,
  "currency": "eur",
  "type": "charge",
  "additional_fields": {
      "charge_id": "ch_3OgwgvGVGHB8I6rc1Etj264n",
      "payment_intent_id": "pi_3OgwgvGVGHB8I6rc1ZC8RNGK"
  }
}
```

Failed payment

```json
{
  "transaction_id": "TXN_01HP07FBXYJJPG7PQVRF5N1MWT",
  "status": "failure",
  "description": "Sample transaction",
  "failure_reason": "card_declined",
  "payment_provider": "stripe",
  "amount": 2000,
  "currency": "eur",
  "type": "charge",
  "additional_fields": {
      "charge_id": "ch_3OgwpAGVGHB8I6rc1HXVKnqH",
      "payment_intent_id": "pi_3OgwpAGVGHB8I6rc1uUXNS1K"
  }
}
```

##### HTTP Code 400

```json
{
  "code": "invalid_request",
  "status_code": 400,
  "message": "Invalid request: invalid amount"
}
```

##### HTTP Code 500

```json
{
  "code": "invalid_server_error",
  "status_code": 500,
  "message": "Internal server error"
}
```

</details>

### Query payment

<details>
 <summary><code>GET</code> <code><b>/{transaction_id}</b></code> <code>(Queries a payment given its transaction_id)</code></summary>

#### Parameters

> | name            |  type     | data type               | description                                              |
> |-----------------|-----------|-------------------------|----------------------------------------------------------|
> | id              |  required | string (path parameter) | Identifier to the gien transaction_id                    |

#### Responses

##### HTTP Code 200

```json
{
  "transaction_id": "TXN_01HP06ZRSNFDPKN3ZBSWS4Z0KT",
  "status": "succeeded",
  "description": "Sample transaction",
  "payment_provider": "stripe",
  "amount": 2000,
  "currency": "eur",
  "type": "charge",
  "additional_fields": {
      "charge_id": "ch_3OgwgvGVGHB8I6rc1Etj264n",
      "payment_intent_id": "pi_3OgwgvGVGHB8I6rc1ZC8RNGK"
  }
}
```

##### HTTP Code 400

```json
{
  "code": "resource_not_found",
  "status_code": 404,
  "message": "Resource 'transaction' not found"
}
```

##### HTTP Code 500

```json
{
  "code": "invalid_server_error",
  "status_code": 500,
  "message": "Internal server error"
}
```

</details>

### Refund payment

**Disclaimer**: When a payment is refunded its initial status is intentionally set to `pending`. In order to mock use case where it takes X amount of time to charge a payment. Therefore, the final status will be given by the event receved by webhooks.

<details>
 <summary><code>POST</code> <code><b>/{transaction_id}</b></code> <code>(Refunds a payment given its transaction_id)</code></summary>

#### Parameters

> | name            |  type     | data type               | description                                              |
> |-----------------|-----------|-------------------------|----------------------------------------------------------|
> | id              |  required | string (path parameter) | Identifier to the gien transaction_id                    |

#### Responses

##### HTTP Code 200

Succesful refund 

```json
{
  "transaction_id": "TXN_01HP06ZRSNFDPKN3ZBSWS4Z0KT",
  "status": "pending",
  "description": "Sample transaction",
  "payment_provider": "stripe",
  "amount": 2000,
  "currency": "eur",
  "type": "refund",
  "additional_fields": {
      "charge_id": "ch_3OgwgvGVGHB8I6rc1Etj264n",
      "payment_intent_id": "pi_3OgwgvGVGHB8I6rc1ZC8RNGK",
      "refund_id": "re_3OgwgvGVGHB8I6rc1rBOb2uO"
  }
}
```

##### HTTP Code 400

```json
{
  "code": "invalid_request",
  "status_code": 400,
  "message": "Invalid request: charge already refunded"
}
```

##### HTTP Code 500

```json
{
  "code": "invalid_server_error",
  "status_code": 500,
  "message": "Internal server error"
}
```

</details>

### Online Payment Webhooks

### Ping

<details>
 <summary><code>GET</code> <code><b>/</b></code> <code>(Check if the service is healthy)</code></summary>

#### Parameters

> None

#### Responses

##### HTTP Code 200

```json
Hello World!
```

</details>