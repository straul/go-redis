package redis

import (
	"net"
	"sync"
)

type Pool struct {
	mu       sync.Mutex
	conns    chan net.Conn
	address  string
	password string
}

// NewPool 创建连接池
func NewPool(address, password string, size int) (*Pool, error) {
	pool := &Pool{
		conns:    make(chan net.Conn, size),
		address:  address,
		password: password,
	}

	for i := 0; i < size; i++ {
		conn, err := pool.createConnection()
		if err != nil {
			return nil, err
		}
		pool.conns <- conn
	}

	return pool, nil
}

func (p *Pool) Get() (net.Conn, error) {
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		return p.createConnection()
	}
}

func (p *Pool) Put(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case p.conns <- conn:
	default:
		conn.Close()
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
}

func (p *Pool) createConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", p.address)
	if err != nil {
		return nil, err
	}

	if p.password != "" {
		_, err = auth(conn, p.password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
