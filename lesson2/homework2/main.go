package main

import (
	"fmt"
	"sync"
	"time"
)

func download(filename string, wg *sync.WaitGroup, result chan<- string, updates chan<- string) {
	defer wg.Done()
	if filename == "lastfile_minit.ios" {
		SlowDownload(filename, result, updates)
	} else {
		FastDownload(filename, result, updates)
	}
}
func FastDownload(filename string, result chan<- string, updates chan<- string) {
	updates <- fmt.Sprintf("⚡正在以疯狂动物城中闪电飙车的速度下载：%s文件...\n", filename)
	time.Sleep(1 * time.Second)
	result <- fmt.Sprintf("✅ %s文件下载完成\n,耗时1秒（超快的）", filename)
}
func SlowDownload(filename string, result chan<- string, updates chan<- string) {
	updates <- fmt.Sprintf("开始下载你非要下的垃圾文件：%s，⚡还是闪电的速度🐌很快的，秒计的...", filename)
	totalSeconds := 24 * 60 * 60
	surfSeconds := 100
	usedSeconds := 0
	for i := 0; i < 6; i++ {
		time.Sleep(1 * time.Second)
		usedSeconds += 1
		progress := (usedSeconds * 100) / surfSeconds
		updates <- fmt.Sprintf("%s:已下载%d%%(只用了%d秒快吧）", filename, progress, usedSeconds)
	}
	progress := 6
	stucktime := 5
	updates <- fmt.Sprintf("系统检查到%s是垃圾文件，闪电要先去确定一下", filename)
	time.Sleep(5 * time.Second)
	for i := 0; i < 5; i++ {
		usedSeconds += stucktime
		updates <- fmt.Sprintf("%s文件卡在了%d%%,已下载%d秒，还剩%d秒，要不等会再看看", filename, progress, usedSeconds, totalSeconds-usedSeconds)
		time.Sleep(time.Duration(stucktime) * time.Second)
	}
	result <- fmt.Sprintf("球了，卡死了，喊你要下勒个垃圾文件，卡死了赛，花了%d秒，结果才下了%d%%,玩完了呗", stucktime, progress)
}
func main() {
	files := []string{
		"file1.doc",
		"file2.pdf",
		"file3.jpg",
		"file4.txt",
		"lastfile_minit.ios",
	}
	result := make(chan string, len(files))
	updates := make(chan string, 50)
	var wg sync.WaitGroup
	fmt.Println("开始下载文件...")
	for _, file := range files {
		wg.Add(1)
		go download(file, &wg, result, updates)
	}
	go func() {
		for update := range updates {
			fmt.Println(update)
		}
	}()
	go func() {
		wg.Wait()
		close(updates)
		close(result)
	}()
	fmt.Println("===最终文件下载结果（都去给我感受1小时下了6%的悲哀）===")
	for result := range result {
		fmt.Println(result)
	}
	fmt.Println("最后一个文件你永远都别想下完，因为我克隆git的时候就是这样")
}
