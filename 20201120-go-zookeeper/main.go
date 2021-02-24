package main

import (
	"log"
	"path"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var wg *sync.WaitGroup

func main() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	wg = &sync.WaitGroup{}
	// watchDemoNode("/demo", conn)
	watchDir("/demo", conn)
	wg.Wait()
}
func watchDemoNode(path string, conn *zk.Conn) {
	wg.Add(1)
	//创建
	watchNodeCreated(path, conn)
	//改值
	go watchNodeDataChange(path, conn)
	//子节点变化「增删」
	go watchChildrenChanged(path, conn)
	//删除节点
	watchNodeDeleted(path, conn)
	wg.Done()
}
func watchNodeCreated(path string, conn *zk.Conn) {
	log.Println("watchNodeCreated")
	for {
		_, _, ch, _ := conn.ExistsW(path)
		e := <-ch
		log.Println("ExistsW:", e.Type, "Event:", e)
		if e.Type == zk.EventNodeCreated {
			log.Println("NodeCreated ")
			return
		}
	}
}
func watchNodeDeleted(path string, conn *zk.Conn) {
	log.Println("watchNodeDeleted")
	for {
		_, _, ch, _ := conn.ExistsW(path)
		e := <-ch
		log.Println("ExistsW:", e.Type, "Event:", e)
		if e.Type == zk.EventNodeDeleted {
			log.Println("NodeDeleted ")
			return
		}
	}
}
func watchNodeDataChange(path string, conn *zk.Conn) {
	for {
		_, _, ch, _ := conn.GetW(path)
		e := <-ch
		log.Println("GetW('"+path+"'):", e.Type, "Event:", e)
	}
}
func watchChildrenChanged(path string, conn *zk.Conn) {
	for {
		_, _, ch, _ := conn.ChildrenW(path)
		e := <-ch
		log.Println("ChildrenW:", e.Type, "Event:", e)
	}
}

func watchDir(key string, conn *zk.Conn) {
	wg.Add(1)
	for {
		// get current children for a key
		children, _, childEventCh, err := conn.ChildrenW(key)
		if err != nil {
			log.Println("ChildrenW err:", err)
			return
		}

		select {
		case e := <-childEventCh:
			if e.Type != zk.EventNodeChildrenChanged {
				continue
			}

			newChildren, _, err := conn.Children(e.Path)
			if err != nil {
				log.Println("Children err:", err)
				return
			}
			log.Println("watchDir:", key, e.Type, "Event:", e)

			// a node was added -- watch the new node
			for _, i := range newChildren {
				if contains(children, i) {
					continue
				}
				newNode := childPath(e.Path, i)

				go watchKey(newNode, conn)
				// s, _, err := zw.client.Get(newNode)
				e.Type = zk.EventNodeCreated
				log.Println("newChildren:", e.Path, e.Type, "Event:", newNode)
			}
		}
	}
	wg.Done()
}

func watchKey(key string, conn *zk.Conn) {
	wg.Add(1)
	for {
		_, _, keyEventCh, err := conn.GetW(key)
		if err != nil {
			log.Println("GetW err:", err)
			return
		}

		select {
		case e := <-keyEventCh:
			switch e.Type {
			case zk.EventNodeDataChanged, zk.EventNodeCreated, zk.EventNodeDeleted:
				log.Println("watchKey:", key, e.Type, "Event:", e)
			}
			if e.Type == zk.EventNodeDeleted {
				//The Node was deleted - stop watching
				return
			}
		}
	}
	wg.Done()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func nodePath(service, node string) string {
	return path.Join(service, node)
}

func childPath(parent, child string) string {
	return path.Join(parent, child)
}
