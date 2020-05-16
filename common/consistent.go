package common

import (
	"errors"
	"github.com/unknwon/com"
	"hash/crc32"
	"sort"
	"sync"
)

// 声明新切片类型
// 实现 sort.Interface 接口，可以 sort 提供的排序算法自动进行排序
// 详情介绍请看： https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter03/03.1.html
type units []uint32

func (x units) Len() int {
	return len(x)
}
func (x units) Less(i, j int) bool {
	return x[i] > x[j]
}
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// 无数据错误提示
var errEmpty = errors.New("hash  环上没有数据")

// 定义结构体保存 一致性 hash 信息数据
type ConsistentHash struct {
	// hash 环，key 为哈希值，值存放节点信息
	circle map[uint32]string
	// 已经排序的节点 hash 切片
	sortedHashes units
	// 虚拟节点个数，永安里增加 hash 的平衡性
	VirtualNode int
	//  读写锁
	sync.RWMutex
}

func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

func (ch *ConsistentHash) generateKey(element string, index int) string {
	// 副本，虚拟节点 key 生成逻辑
	return element + com.ToStr(index)
}

// 获取 hash 位置
func (ch *ConsistentHash) hashKey(key string) uint32 {
	if len(key) < 64 {
		// 通常业务逻辑，一般对于 key 值较小的会加点其他字符标识，这里没有只是直接截取返回
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// 往 hash 环中添加节点
func (ch *ConsistentHash) Add(element string) {
	// 加锁
	ch.Lock()
	// 使用完后解锁
	defer ch.Unlock()
	ch.add(element)
}

func (ch *ConsistentHash) add(element string) {
	// 根据 源节点 循环虚拟节点，设置副本
	for i := 0; i < ch.VirtualNode; i++ {
		ch.circle[ch.hashKey(ch.generateKey(element, i))] = element
	}
	// 更新排序
	ch.updateSortHashes()
}

//删除节点
func (ch *ConsistentHash) Remove(element string) {
	ch.Lock()
	defer ch.Unlock()
	ch.remove(element)
}

func (ch *ConsistentHash) remove(element string) {
	// 删除节点以及节点副本
	for i := 0; i < ch.VirtualNode; i++ {
		delete(ch.circle, ch.hashKey(ch.generateKey(element, i)))
	}
	// 更新
	ch.updateSortHashes()
}

// 更新排序
func (ch *ConsistentHash) updateSortHashes() {
	hashes := ch.sortedHashes[:0]
	// 判断切片容量，是否过大，如果过大则重置
	if cap(ch.sortedHashes)/(ch.VirtualNode) > len(ch.circle) {
		hashes = nil
	}

	// 添加 hashes
	for k := range ch.circle {
		hashes = append(hashes, k)
	}

	// hash节点添加完成后，进行排序
	sort.Sort(hashes)
	// 排序后从新赋值
	ch.sortedHashes = hashes
}

// 根据数据标识获取最近服务器节点信息
func (ch *ConsistentHash) Get(name string) (string, error) {
	// 只有读，使用读锁
	ch.RLock()
	defer ch.RUnlock()

	if len(ch.circle) == 0 {
		return "", errEmpty
	}

	// 计算 源数据 对应 hash 环上的位置（值）
	key := ch.hashKey(name)
	// 通过该 key 查找已排序节点切片中对应最近节点 位置
	nodeIndex := ch.search(key)
	// 通过最近节点 获取 节点存储的对应服务器信息
	return ch.circle[ch.sortedHashes[nodeIndex]], nil
}

// 顺时针查找 源数据对应 在已排序节点切片中最近的服务节点的下标位置
func (ch *ConsistentHash) search(key uint32) int {
	// 查找算法
	f := func(x int) bool {
		return ch.sortedHashes[x] > key
	}
	// 使用 "二分查找" 算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(ch.sortedHashes), f)
	//
	if i >= len(ch.sortedHashes) {
		i = 0
	}
	return i
}
