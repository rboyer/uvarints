package uvarints

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sean-/seed"
	"github.com/stretchr/testify/require"
)

func init() {
	seed.MustInit()
}

func TestReaderAndWriter(t *testing.T) {
	const dataSize = 1000

	data := make([]uint64, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = rand.Uint64()
	}

	tmpDir, err := ioutil.TempDir("", "ints")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	fn := filepath.Join(tmpDir, "ints.tmp")

	w, err := OpenWriter(fn, 0600)
	require.NoError(t, err)
	defer w.Close()

	start := time.Now()
	for _, v := range data {
		require.NoError(t, w.WriteUint(v))
	}

	require.NoError(t, w.Close())
	dur := time.Since(start)
	t.Logf("wrote in %v", dur)

	r, err := OpenReader(fn)
	require.NoError(t, err)
	defer r.Close()

	start = time.Now()
	got := make([]uint64, 0, dataSize)
	for {
		v, err := r.ReadUint()
		if err == io.EOF {
			break
		}
		got = append(got, v)
	}
	dur = time.Since(start)
	t.Logf("read in %v", dur)

	require.ElementsMatch(t, data, got)
}
