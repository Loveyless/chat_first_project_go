package main

func main() {

	//因为这个和server.go都属于mian包 所以不需要引入
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
