package migraterice

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/golang-migrate/migrate/v4/source"
)


func init() {
	source.Register("rice", &RiceSourceDriver{})
}

type RiceSourceDriver struct {
	path       string
	box        *rice.Box
	migrations *source.Migrations
}

func WithInstance(box *rice.Box) (source.Driver, error) {
	riceSourceDriver := &RiceSourceDriver{
		path:       "<rice>",
		box:        box,
		migrations: source.NewMigrations(),
	}

	box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		m, err := source.DefaultParse(path)
		if err != nil {
			return nil
		}

		if !riceSourceDriver.migrations.Append(m) {
			return nil
		}
		return nil
	})

	return riceSourceDriver, nil
}

func (r *RiceSourceDriver) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (r *RiceSourceDriver) Close() error {
	return nil
}

func (r *RiceSourceDriver) First() (uint, error) {
	if v, ok := r.migrations.First(); !ok {
		return 0, &os.PathError{
			Op:   "first",
			Path: r.path,
			Err:  os.ErrNotExist,
		}
	} else {
		return v, nil
	}
}

func (r *RiceSourceDriver) Prev(version uint) (uint, error) {
	if v, ok := r.migrations.Prev(version); !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("prev for version %v", version),
			Path: r.path,
			Err:  os.ErrNotExist,
		}
	} else {
		return v, nil
	}
}

func (r *RiceSourceDriver) Next(version uint) (uint, error) {
	if v, ok := r.migrations.Next(version); !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("next for version %v", version),
			Path: r.path,
			Err:  os.ErrNotExist,
		}
	} else {
		return v, nil
	}
}

func (r *RiceSourceDriver) ReadUp(version uint) (io.ReadCloser, string, error) {
	if m, ok := r.migrations.Up(version); ok {
		file, err := r.box.Bytes(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return ioutil.NopCloser(bytes.NewReader(file)), m.Identifier, nil
	}
	return nil, "", &os.PathError{
		Op:   fmt.Sprintf("read version %v", version),
		Path: r.path,
		Err:  os.ErrNotExist,
	}
}

func (r *RiceSourceDriver) ReadDown(version uint) (io.ReadCloser, string, error) {
	if m, ok := r.migrations.Down(version); ok {
		file, err := r.box.Bytes(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return ioutil.NopCloser(bytes.NewReader(file)), m.Identifier, nil
	}
	return nil, "", &os.PathError{
		Op:   fmt.Sprintf("read version %v", version),
		Path: r.path,
		Err:  os.ErrNotExist,
	}
}
