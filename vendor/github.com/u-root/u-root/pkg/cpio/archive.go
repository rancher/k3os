// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"io"
	"strings"
)

// Archive is an in-memory list of files.
//
// Archive itself is a RecordWriter, and Archive.Reader() returns a new
// RecordReader for the archive starting from the first file.
type Archive struct {
	// Files is a map of relative archive path -> record.
	Files map[string]Record

	// Order is a list of relative archive paths and represents the order
	// in which Files were added.
	Order []string
}

// InMemArchive returns an in-memory file archive.
func InMemArchive() *Archive {
	return &Archive{
		Files: make(map[string]Record),
	}
}

// ArchiveFromRecords creates a new Archive from the records.
func ArchiveFromRecords(rs []Record) *Archive {
	a := InMemArchive()
	for _, r := range rs {
		a.WriteRecord(r)
	}
	return a
}

// ArchiveFromReader reads records from r into a new Archive in memory.
func ArchiveFromReader(r RecordReader) (*Archive, error) {
	a := InMemArchive()
	if err := Concat(a, r, nil); err != nil {
		return nil, err
	}
	return a, nil
}

// WriteRecord implements RecordWriter and adds a record to the archive.
//
// WriteRecord uses Normalize to deduplicate paths.
func (a *Archive) WriteRecord(r Record) error {
	r.Name = Normalize(r.Name)
	a.Files[r.Name] = r
	a.Order = append(a.Order, r.Name)
	return nil
}

// Empty returns whether the archive has any files in it.
func (a *Archive) Empty() bool {
	return len(a.Files) == 0
}

// Contains returns true if a record matching r is in the archive.
func (a *Archive) Contains(r Record) bool {
	r.Name = Normalize(r.Name)
	if s, ok := a.Files[r.Name]; ok {
		return Equal(r, s)
	}
	return false
}

// Get returns a record for the normalized path or false if there is none.
//
// The path is normalized using Normalize, so Get("/bin/bar") is the same as
// Get("bin/bar") is the same as Get("bin//bar").
func (a *Archive) Get(path string) (Record, bool) {
	r, ok := a.Files[Normalize(path)]
	return r, ok
}

// String implements fmt.Stringer.
//
// String lists files like ls would.
func (a *Archive) String() string {
	var b strings.Builder
	r := a.Reader()
	for {
		record, err := r.ReadRecord()
		if err != nil {
			return b.String()
		}
		b.WriteString(record.String())
		b.WriteString("\n")
	}
}

type archiveReader struct {
	a   *Archive
	pos int
}

// Reader returns a RecordReader for the archive that starts at the first
// record.
func (a *Archive) Reader() RecordReader {
	return &EOFReader{&archiveReader{a: a}}
}

// ReadRecord implements RecordReader.
func (ar *archiveReader) ReadRecord() (Record, error) {
	if ar.pos >= len(ar.a.Order) {
		return Record{}, io.EOF
	}

	path := ar.a.Order[ar.pos]
	ar.pos++
	return ar.a.Files[path], nil
}
