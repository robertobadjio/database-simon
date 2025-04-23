package e2e

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EWAL(t *testing.T) {
	buffer := make([]byte, 1024)
	const serverAddress = "localhost:3223"

	cmd := exec.Command("../../storage_server", "-config", "./config_wal.yml")
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second)

	connection, clientErr := net.Dial("tcp", serverAddress)
	require.NoError(t, clientErr)

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr := connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[not found]", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("SET key1 value1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "[ok]", string(buffer[:size]))

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "value1", string(buffer[:size]))

	time.Sleep(time.Second)

	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	time.Sleep(time.Second)

	cmd = exec.Command("../../storage_server", "-config", "./config_wal.yml")
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second)

	connection, clientErr = net.Dial("tcp", serverAddress)
	require.NoError(t, clientErr)

	_, clientErr = connection.Write([]byte("GET key1"))
	require.NoError(t, clientErr)

	size, clientErr = connection.Read(buffer)
	require.NoError(t, clientErr)
	assert.Equal(t, "value1", string(buffer[:size]))

	require.NoError(t, connection.Close())
	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	files, errGlob := filepath.Glob("./data/wal/wal_*")
	require.NoError(t, errGlob)

	for _, f := range files {
		if errRemove := os.Remove(f); errRemove != nil {
			require.NoError(t, errRemove)
		}
	}

	time.Sleep(time.Second)
}
