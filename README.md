# Go JWT Clean Architecture 
 
A JWT authentication system built with Go using clean architecture principles. 
 
## Features 
 
- JWT token generation and validation 
- Protected routes with middleware 
- Clean architecture with separated concerns 
- Environment-based configuration 
- User context extraction 
 
## API Endpoints 
 
### Public Endpoints 
- `GET /` - Health check 
- `GET /check` - Server status 
- `GET /generate-token?username=<user>` - Generate JWT token 
- `GET /parse-token?token=<jwt>` - Parse and validate JWT token 
 
### Protected Endpoints (Requires Authorization: Bearer <token>) 
- `GET /dashboard` - User dashboard 
- `GET /profile` - User profile 
- `GET /settings` - User settings 
 
## Setup 
 
1. Clone the repository 
```bash 
git clone https://github.com/YOUR_USERNAME/go-jwt-clean.git 
cd go-jwt-clean 
``` 
 
2. Install dependencies 
```bash 
go mod tidy 
``` 
 
3. Create `.env` file 
```env 
BACKEND_SERVER_PORT=8080 
JWT_SECRET_SIGNING_KEY=your_super_secret_jwt_signing_key_at_least_32_characters_long 
``` 
 
4. Run the server 
```bash 
go run cmd/server/main.go 
``` 
 
## Tech Stack 
 
- Go 1.21+ 
- Gorilla Mux - HTTP router 
- golang-jwt/jwt - JWT handling 
- godotenv - Environment variables 
