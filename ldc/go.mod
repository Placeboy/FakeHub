module github.com/Placeboy/FakeHub/ldc

replace github.com/Placeboy/FakeHub/vs => ../vs

replace github.com/Placeboy/FakeHub/rsa => ../rsa // 这条是必要的,因为vs又依赖于rsa

go 1.15

require (
	github.com/Placeboy/FakeHub/vs v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.0
)
