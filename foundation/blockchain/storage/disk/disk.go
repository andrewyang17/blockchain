// Package disk implements the ability to read and write blocks to disk
// writing each block to a separate block numbered file.
package disk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/andrewyang17/blockchain/foundation/blockchain/database"
)

// Disk represents the serialization implementation for reading and storing
// blocks in their own separate files on disk. THis implements the database.Storage
// interface.
type Disk struct {
	dbPath string
}

// New constructs a Disk value for use.
func New(dbPath string) (*Disk, error) {
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}

	return &Disk{dbPath: dbPath}, nil
}

// Close in this implementation has nothing to do since a new file is
// written to disk for each new block and then immediately closed.
func (d *Disk) Close() error {
	return nil
}

// Write takes the specified database blocks and stores it on disk in a
// file labeled with the block number.
func (d *Disk) Write(blockData database.BlockData) error {
	data, err := json.MarshalIndent(blockData, "", "  ")
	if err != nil {
		return nil
	}

	// Create a new file for this block and name it based on the block number.
	f, err := os.OpenFile(d.getPath(blockData.Header.Number), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// GetBlock searches the blockchain on disk to locate and return the
// contents of the specified block by number.
func (d *Disk) GetBlock(num uint64) (database.BlockData, error) {
	f, err := os.OpenFile(d.getPath(num), os.O_RDONLY, 0600)
	if err != nil {
		return database.BlockData{}, err
	}
	defer f.Close()

	var blockData database.BlockData
	if err := json.NewDecoder(f).Decode(&blockData); err != nil {
		return database.BlockData{}, err
	}

	return blockData, nil
}

// ForEach returns an iterator to walk through all the blocks
// starting with block number 1.
func (d *Disk) ForEach() database.Iterator {
	return &diskIterator{storage: d}
}

// Reset will clear out the blockchain on disk.
func (d *Disk) Reset() error {
	if err := os.RemoveAll(d.dbPath); err != nil {
		return err
	}

	return os.MkdirAll(d.dbPath, 0755)
}

// getPath forms the path to the specified block.
func (d *Disk) getPath(blockNum uint64) string {
	name := strconv.FormatUint(blockNum, 10)
	return path.Join(d.dbPath, fmt.Sprintf("%s.json", name))
}

// =============================================================================

// diskIterator represents the iteration implementation for walking
// through and reading blocks on disk. This implements the database
// Iterator interface.
type diskIterator struct {
	storage            *Disk
	currentBlockNumber uint64
	endOfChain         bool
}

// Next retrieves  the next block from disk.
func (di *diskIterator) Next() (database.BlockData, error) {
	if di.endOfChain {
		return database.BlockData{}, errors.New("end of chain")
	}
	di.currentBlockNumber++
	blockData, err := di.storage.GetBlock(di.currentBlockNumber)
	if errors.Is(err, fs.ErrNotExist) {
		di.endOfChain = true
	}

	return blockData, err
}

// Done returns the end of chain value.
func (di *diskIterator) Done() bool {
	return di.endOfChain
}
