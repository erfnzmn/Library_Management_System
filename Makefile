#  Makefile for Library Management System 

run:
	@echo "🚀 Starting the server..."
	go run cmd/server/main.go

login-test:
	@echo "🔐 Testing login endpoint..."
	http POST :8080/users/login email="fanzm1316@gmail.com" password="Erfnzmn1316"

reserve-test:
	@echo "📦 Testing reserve endpoint..."
	http POST :8080/api/loans/reserve "Authorization: Bearer $(TOKEN)" book_id:=1
