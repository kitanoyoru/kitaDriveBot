package app

type App interface {
	Run() error
	Close() error
}
