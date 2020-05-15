package main

import (
	"fmt"
)

// type Task interface {
// 	Execute()
// }

// type Pool struct {
// 	mu    sync.Mutex
// 	size  int
// 	tasks chan Task
// 	kill  chan struct{}
// 	wg    sync.WaitGroup
// }

// func NewPool(size int) *Pool {
// 	pool := &Pool{
// 		tasks: make(chan Task, 128),
// 		kill:  make(chan struct{}),
// 	}
// 	pool.Resize(size)
// 	return pool
// }

// func (p *Pool) worker() {
// 	defer p.wg.Done()
// 	for {
// 		select {
// 		case task, ok := <-p.tasks:
// 			if !ok {
// 				return
// 			}
// 			task.Execute()
// 		case <-p.kill:
// 			return
// 		}
// 	}
// }

// func (p *Pool) Resize(n int) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	for p.size < n {
// 		p.size++
// 		p.wg.Add(1)
// 		go p.worker()
// 	}
// 	for p.size > n {
// 		p.size--
// 		p.kill <- struct{}{}
// 	}
// }

// func (p *Pool) Close() {
// 	close(p.tasks)
// }

// func (p *Pool) Wait() {
// 	p.wg.Wait()
// }

// func (p *Pool) Exec(task Task) {
// 	p.tasks <- task
// }

// type ExampleTask string

// func (e ExampleTask) Execute() {
// 	fmt.Println("executing:", string(e))
// }

// func main() {
// 	pool := NewPool(5)

// 	pool.Exec(ExampleTask("foo"))
// 	pool.Exec(ExampleTask("bar"))

// 	// pool.Resize(3)

// 	// pool.Resize(6)

// 	for i := 0; i < 20; i++ {
// 		pool.Exec(ExampleTask(fmt.Sprintf("additional_%d", i+1)))
// 	}

// 	pool.Close()

// 	pool.Wait()
// }

func main() {

	snapshots := []string{"1", "2", "3", "tty"}

	jobs := make(chan string, len(snapshots))
	results := make(chan string, len(snapshots))

	go worker(jobs, results)

	for _, snapshot := range snapshots {
		jobs <- snapshot
	}
	close(jobs)

	for i := 0; i < len(snapshots); i++ {
		fmt.Println(<-results)
		// Do something with the result here
	}
	close(results)
}

func worker(jobs <-chan string, results chan<- string) {
	for snapshot := range jobs {
		newSnapshot := snapshot + "abc"
		results <- newSnapshot
	}
}

// jobs := make(chan *ec2.Snapshot, len(snapshots))
// 	results := make(chan AssetDetails, len(snapshots))

// 	go worker(jobs, results, amis, totalCost)

// 	for _, snapshot := range snapshots {
// 		jobs <- snapshot
// 	}
// 	close(jobs)

// 	for i := 0; i < len(snapshots); i++ {
// 		fmt.Println(<-results)
// 	}
// 	close(results)

// func worker(jobs <-chan *ec2.Snapshot, results chan<- AssetDetails, amis []string, totalCost float64) {
// 	for snapshot := range jobs {
// 		words := strings.Fields(*snapshot.Description)
// 		if len(words) == 7 {
// 			if words[0] == "Created" && words[1] == "by" && words[2][0:11] == "CreateImage" {
// 				amiId := words[4]
// 				amiExists := stringInSlice(amiId, amis)
// 				if amiExists == false {
// 					cost := float64(*snapshot.VolumeSize) * 0.05
// 					totalCost = totalCost + cost
// 					s := AssetDetails{}
// 					s.Id = *snapshot.SnapshotId
// 					s.SizeInGB = strconv.FormatInt(*snapshot.VolumeSize, 10)
// 					s.CostPerMonth = ("$" + fmt.Sprintf("%.2f", cost))
// 					results <- s
// 				}
// 			}
// 		}
// 	}
// }
