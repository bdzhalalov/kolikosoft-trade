## Getting started

### Prerequisites
- **Docker** with Docker Compose

### Start up project

- Go to directory with cloned project
- Use command `make start`

#### For subsequent launches of the application, it is enough to use the command `make run`

---

## Migrations

#### All necessary database migrations are applied when the service is first launched.

### Create migration

To create a new migration, you need to use the `make migrate-create name="migration name"` command.

### Run created migration

Use `make migrate-up` command to apply new migration

### Revert created migration

Use `make migrate-down` command to revert new migration

---

## Usage
1. **Items**
   
   - **GET** `/api/v1/items/list`
   
   Get items from external API. Endpoint is cached for 5 minutes.
  
2. **User balance**
   
   - **POST** `/api/v1/users/{id}/balance/withdraw`

   Withdraw money from user balance
   
   ```
   Required body:
   {
     "amount": int
   }
   ```
   !Endpoint is idempotent: The response contain a "Idempotency-Key" header. For idempotency, it should be used in request headers.

   - **GET** `/api/v1/users/{id}/balance/history`

   Get user balance history

---

## Testing

To run tests, simply use the command `make test`
