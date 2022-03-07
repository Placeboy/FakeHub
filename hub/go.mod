module github.com/Placeboy/FakeHub/hub

go 1.15

replace github.com/Placeboy/FakeHub/rsa => ../rsa

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/Placeboy/FakeHub/rsa v0.0.0-00010101000000-000000000000
)
