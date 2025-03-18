# API Spec

## Authentication

All API must use this authentication

Request :

- Header :
  - Authorization : Bearer "your_secret_token_key"

<!-- Panel Admin  -->

<!-- CRUD Department -->

## CRUD Department +

### Create Department 

Request :

- Method : POST
- Endpoint : `/api/departments`
- Header :
  - Content-Type: application/json
  - Accept: application/json
- Body :

```json
{
  "id": "string, unique",
  "parent_department_id": "integer",
  "department_name": "string"
}
```

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "parent_department_id": "integer",
    "department_name": "string"
  }
}
```

## Get Department by id 

Request :

- Method : GET
- Endpoint : `/api/departments/{id_department}`
- Header :
  - Accept: application/json
- Query Param :
  - id : number,

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "parent_department_id": "integer",
    "department_name": "string"
  }
}
```

## List Department 

Request :

- Method : GET
- Endpoint : `/api/departments`
- Header :
  - Accept: application/json
  <!-- - Query Param :
  - size : number,
  - page : number -->

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": [
    {
      "id": "string, unique",
      "parent_department_id": "integer",
      "department_name": "string"
    },
    {
      "id": "string, unique",
      "parent_department_id": "integer",
      "department_name": "string"
    }
  ]
}
```

## Update Department 

Request :

- Method : PUT
- Endpoint : `/api/departments/{id_department}`
- Header :
  - Content-Type: application/json
  - Accept: application/json
- Query Param :
  - id : number,
- Body :

```json
{
  "parent_department_id": "integer",
  "department_name": "string"
}
```

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "parent_department_id": "integer",
    "department_name": "string"
  }
}
```

## Delete Department 

Request :

- Method : DELETE
- Endpoint : `/api/departments/{id_department}`
- Header :
  - Accept: application/json

Response :

```json
{
  "error": "string",
  "message": "string"
}
```

<!-- CRUD Position -->

## CRUD Position +

### Create Position

Request :

- Method : POST
- Endpoint : `/api/positions`
- Header :
  - Content-Type: application/json
  - Accept: application/json
- Body :

```json
{
  "id": "string, unique",
  "department_id": "integer",
  "position_name": "string",
  "is_completed": "boolean"
}
```

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "department_id": "integer",
    "department_name": "string",
    "position_name": "string",
    "is_completed": "boolean"
  }
}
```

### List All Position

Request :

- Method : GET
- Endpoint : `/api/positions`
- Header :
  - Accept: application/json

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": [
    {
      "id": "string, unique",
      "department_id": "integer",
      "department_name": "string",
      "position_name": "string",
      "is_completed": "boolean"
    },
    {
      "id": "string, unique",
      "department_id": "integer",
      "department_name": "string",
      "position_name": "string",
      "is_completed": "boolean"
    }
  ]
}
```

### Get Position By ID

Request :

- Method : GET
- Endpoint : `/api/positions/{id_position}`
- Header :
  - Accept: application/json
- Query Param :
  - id : number,

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "department_id": "integer",
    "department_name": "string",
    "position_name": "string",
    "is_completed": "boolean"
  }
}
```

### Update Position By ID

Request :

- Method : PUT
- Endpoint : `/api/positions/{id_position}`
- Header :
  - Content-Type: application/json
  - Accept: application/json
- Query Param :
  - id : number,
- Body :

```json
{
  "id": "string, unique",
  "department_id": "integer",
  "position_name": "string",
  "is_completed": "boolean"
}
```

Response :

```json
{
  "error": "false",
  "message": "string",
  "data": {
    "id": "string, unique",
    "department_id": "integer",
    "department_name": "string",
    "position_name": "string",
    "is_completed": "boolean"
  }
}
```

### Delete Position By ID

Request :

- Method : DELETE
- Endpoint : `/api/positions/{id_position}`
- Header :
  - Accept: application/json
- Query Param :
  - id : number,

Response :

```json
{
  "error": "false",
  "message": "string"
}
```

<!-- Panel Admin -->

## Register +

Request :

- Method : POST
- Endpoint : '/api/register'
- Header :
  - Content-Type : application/json
  - Accept: application/json
- Body :

```json
{
  "name": "string, required",
  "email": "string, required, unique",
  "password": "string, required, must least 8 character",
  "phone": "string, required",
  "position": "string, required"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string"
}
```

## Login +

Request :

- Method : POST
- Endpoint : '/api/login'
- Header :
  - Content-Type : application/json
  - Accept: application/json
- Body :

```json
{
  "email": "string, required",
  "password": "string, required"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "loginResult": {
    "id": "string, unique",
    "name": "string",
    "token": "string, unique"
  }
}
```

## Profile Employee +

### Get Profile Employee 

Request :

- Method : GET
- Endpoint : `api/user`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "profile": {
    "id": "string, unique",
    "name": "string",
    "email": "string, unique",
    "phone": "string",
    "position": "string",
    "department": "string",
    "photo": "string_url"
  }
}
```

### Update Profile Employee

Request :

- Method : PUT
- Endpoint : `api/user`
- Header :
  - Authorization : Bearer "token_key"
  - Content-Type: multipart/form-data
  - Accept: application/json
- Body :
  - photo as file, must be a valid image file, max size 2 MB

```json
{
  "name": "string",
  "password": "string",
  "phone": "string",
  "photo": "string"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "profile": {
    "id": "string, unique",
    "name": "string",
    "email": "string, unique",
    "phone": "string",
    "position": "string",
    "department": "string",
    "photo": "string_url"
  }
}
```

<!-- ## Get All Profile User (Panel Admin)
## Delete Profile User (Panel Admin) -->

## Get Attendance Employee by DateNow (Done)

Request :

- Method : GET
- Endpoint : `api/attendance?date={date}`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "attendance": {
    "id": "int, unique",
    "shift_id": "integer",
    "employee_id": "integer",
    "created_by": "integer",
    "schedule_id": "integer",
    "name": "string",
    "position": "string",
    "type": "string",
    "date_schedule": "date",
    "status": "string",
    "date": "date",
    "clock_in": "time",
    "clock_out": "time",
    "duration": "string",
    "clock_in_status": "string",
    "clock_out_status": "string"
  }
}
```

## Get Attendance by Status (Done)

Request :

- Method : GET
- Endpoint : `api/attendance?status={clock_in_status}`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "attendances": [
    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "schedule_id": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date": "date",
      "clock_in": "time",
      "clock_out": "time",
      "duration": "string",
      "clock_in_status": "string",
      "clock_out_status": "string"
    },

    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date": "date",
      "clock_in": "time",
      "clock_out": "time",
      "duration": "string",
      "clock_in_status": "string",
      "clock_out_status": "string"
    }
  ]
}
```

## Get All Attendance (Done)

Request :

- Method : GET
- Endpoint : `api/attendance?status={clock_in_status}`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "attendances": [
    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date": "date",
      "clock_in": "time",
      "clock_out": "time",
      "duration": "string",
      "clock_in_status": "string",
      "clock_out_status": "string"
    },

    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date": "date",
      "clock_in": "time",
      "clock_out": "time",
      "duration": "string",
      "clock_in_status": "string",
      "clock_out_status": "string"
    }
  ]
}
```

<!-- ## Get All attendance users (panel admin) -->

## Get All Schedule Employee (Done)

Request :

- Method : GET
- Endpoint : `/api/schedules`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "schedules": [
    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "name": "string",
      "position": "string",
      "date_schedule": "date",
      "status": "string",
      "type": "string",
      "start_time": "time",
      "end_time": "time"
    },
    {
      "id": "int, unique",
      "shift_id": "integer",
      "employee_id": "integer",
      "created_by": "integer",
      "name": "string",
      "position": "string",
      "date_schedule": "date",
      "status": "string",
      "type": "string",
      "start_time": "time",
      "end_time": "time"
    }
  ]
}
```

## Create Attendance Employee (Done)

Request :

- Method : POST
- Endpoint : `/api/attendance`
- Header :
  - Authorization : Bearer "token_key"
  - Content-Type: application/json
  - Accept: application/json
- Body :

```json
{
  "shift_id": "integer",
  "employee_id": "integer",
  "created_by": "integer",
  "schedule_id": "integer",
  "name": "string",
  "position": "string",
  "type": "string",
  "date_schedule": "date",
  "status": "string",
  "date": "date",
  "clock_in": "time",
  "clock_out": "time",
  "duration": "string",
  "clock_in_status": "string",
  "clock_out_status": "string"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "attendance": {
    "id": "int, unique",
    "shift_id": "integer",
    "employee_id": "integer",
    "created_by": "integer",
    "schedule_id": "integer",
    "name": "string",
    "position": "string",
    "type": "string",
    "date_schedule": "date",
    "status": "string",
    "date": "date",
    "clock_in": "time",
    "clock_out": "time",
    "duration": "string",
    "clock_in_status": "string",
    "clock_out_status": "string"
  }
}
```

<!-- CRUD Schedule -->

## Create Schedule Employee

Request :

- Method : POST
- Endpoint : `/api/schedule`
- Header :
  - Authorization : Bearer "token_key"
  - Content-Type: application/json
  - Accept: application/json
- Body :

```json
{
  "name": "string",
  "type": "string",
  "date_schedule": "date"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "schedule": {
    "id": "int, unique",
    "created_by": "integer",
    "shift_id": "integer",
    "employee_id": "integer",
    "name": "string",
    "position": "string",
    "type": "string",
    "date_schedule": "date",
    "status": "string",
    "date_schedule": "date",
    "clock_in": "time",
    "clock_out": "time"
  }
}
```

## List Schedule Employee (In Department)

Request :

- Method : GET
- Endpoint : `/api/schedules`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "schedules": [
    {
      "id": "int, unique",
      "created_by": "integer",
      "shift_id": "integer",
      "employee_id": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date_schedule": "date",
      "clock_in": "time",
      "clock_out": "time"
    },
    {
      "id": "int, unique",
      "created_by": "integer",
      "shift_id": "integer",
      "employee_id": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date_schedule": "date",
      "clock_in": "time",
      "clock_out": "time"
    },
    {
      "id": "int, unique",
      "created_by": "integer",
      "shift_id": "integer",
      "employee_id": "integer",
      "name": "string",
      "position": "string",
      "type": "string",
      "date_schedule": "date",
      "status": "string",
      "date_schedule": "date",
      "clock_in": "time",
      "clock_out": "time"
    }
  ]
}
```

## Update Schedule Employee

Request :

- Method : PUT
- Endpoint : `/api/schedules/{id}`
- Header :
  - Authorization : Bearer "token_key"
  - Content-Type: application/json
  - Accept: application/json
- Body :

```json
{
  "name": "string",
  "type": "string",
  "date_schedule": "date",
  "status": "string"
}
```

Response :

```json
{
  "error": "boolean",
  "message": "string",
  "schedule": {
    "id": "int, unique",
    "created_by": "integer",
    "shift_id": "integer",
    "employee_id": "integer",
    "name": "string",
    "position": "string",
    "type": "string",
    "date_schedule": "date",
    "status": "string",
    "date_schedule": "date",
    "clock_in": "time",
    "clock_out": "time"
  }
}
```

## Delete Schedule Employee

Request :

- Method : DELETE
- Endpoint : `/api/schedules/{id}`
- Header :
  - Authorization : Bearer "token_key"
  - Accept: application/json

Response :

```json
{
  "error": "boolean",
  "message": "string"
}
```
