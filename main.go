package main

func main() {
	repository := MinioRepository{}

	App{repository: repository}.run()
}
