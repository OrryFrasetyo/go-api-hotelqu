# HotelQu API

## Overview
HotelQu API is a REST API-based hotel management system built with Go and the Gin framework. This API is designed to manage various aspects of hotel operations, including:

- 👥 Employee Management
- 📅 Shift and Schedule
- ⏰ Presence System
- 🏢 Department Management
- 👔 Position Management

## Technology Used
- **Backend:** Go (Golang) with Gin Framework
- **Database:** MySQL
- **Authentication:** JWT (JSON Web Token)
- **Architecture:** Domain-Driven Design (DDD)

## Main Features
- ✅ Login
- ✅ Register
- ✅ Department Management
- ✅ Position Management
- ✅ Shift Management
- ✅ Presence System
- ✅ Schedule Employee
- ✅ Profile Employee Management

## Endpoints
### Department
- **GET /api/departments** : Endpoint to get all department data.
- **POST /api/departments** : Endpoint to add new department data.
- **GET /api/departments/:id** : Endpoint to get department data by ID.
- **PUT /api/departments/:id** : Endpoint to update department data by ID.
- **DELETE /api/departments/:id** : Endpoint to delete department data by ID.

### Position

- **GET /api/positions** : Endpoint to get all positions data.
- **POST /api/positions** : Endpoint to add new position data.
- **GET /api/positions/:id** : Endpoint to get position data by ID.
- **PUT /api/positions/:id** : Endpoint to update position data by ID.
- **DELETE /api/positions/:id** : Endpoint to delete position data by ID.

### Shift

- **GET /api/shifts** : Endpoint to get all shifts data.
- **POST /api/shifts** : Endpoint to add new shift data.
- **GET /api/shifts/:id** : Endpoint to get shift data by ID.
- **PUT /api/shifts/:id** : Endpoint to update shift data by ID.
- **DELETE /api/shifts/:id** : Endpoint to delete shift data by ID.

### Login-Register
- **POST /api/register** : register account.
- **POST /api/login** : login account.

### Employee
- **GET /api/user** : get profile employee.
- **GET /uploads/{name_photo}** : get photo profile.
- **PUT /api/user** : update profile.

### Schedule
- **GET /api/schedules/department?date={set date(ex: 03-04-2025)}** : Displays all employee schedule data in one department in the hotel according to the selected date. Access api for manajer or supervisor position
- **POST /api/schedules** : create schedule employee (access api for manajer or supervisor position)
- **PUT /api/schedules/:id** : update schedule employee (access api for manajer or supervisor position)
- **DELETE /api/schedules/:id** : delete schedule employee (access api for manajer or supervisor position)
- **GET /api/schedules** : Displays all hotel employee work schedules in each department.
- **GET /api/schedules/today** : get schedule for today.

### Attendance
- **POST /api/attendance** : clockin presence
- **PUT /api/attendance** : clockout presence
- **GET /api/attendance** : get attendance 3 days ago
- **GET /api/attendance/today** : get attendance for today
- **GET /api/attendance/month** : get attendance for this month
- **GET /api/attendance/status?{clock_in_status=value} or {clock_out_status=value}** : get attendance by status
- **GET /api/employees** : get presence employee

**Note:** All the above endpoints require authentication, except for `POST api/register` , `POST api/login`, shift, department, and position. To use endpoints that require authentication, you need to send the authentication token in the request header with the format `Authorization: Bearer <token>`.

## Deployment Link
API Hotelqu : https://backend-pkl-orry.up.railway.app/

## Postman Documentation
[**Documentation Link**](https://documenter.getpostman.com/view/33562686/2sB2x6kBhN)


