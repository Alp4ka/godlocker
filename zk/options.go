package zk

type Option interface {
	Apply(*Locker)
}

type OptionFunc func(locker *Locker)

func (f OptionFunc) Apply(locker *Locker) {
	f(locker)
}
