package godlocker

type Label interface {
	String() string
	Decompose() (prefix, id string, err error)
	Valid() (ok bool, err error)
}
