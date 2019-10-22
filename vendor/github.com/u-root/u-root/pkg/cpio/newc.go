// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/uio"
)

const (
	newcMagic = "070701"
	magicLen  = 6
)

var (
	// Newc is the newc CPIO record format.
	Newc RecordFormat = newc{magic: newcMagic}
)

type header struct {
	Ino        uint32
	Mode       uint32
	UID        uint32
	GID        uint32
	NLink      uint32
	MTime      uint32
	FileSize   uint32
	Major      uint32
	Minor      uint32
	Rmajor     uint32
	Rminor     uint32
	NameLength uint32
	CRC        uint32
}

func headerFromInfo(i Info) header {
	var h header
	h.Ino = uint32(i.Ino)
	h.Mode = uint32(i.Mode)
	h.UID = uint32(i.UID)
	h.GID = uint32(i.GID)
	h.NLink = uint32(i.NLink)
	h.MTime = uint32(i.MTime)
	h.FileSize = uint32(i.FileSize)
	h.Major = uint32(i.Major)
	h.Minor = uint32(i.Minor)
	h.Rmajor = uint32(i.Rmajor)
	h.Rminor = uint32(i.Rminor)
	h.NameLength = uint32(len(i.Name)) + 1
	return h
}

func (h header) Info() Info {
	var i Info
	i.Ino = uint64(h.Ino)
	i.Mode = uint64(h.Mode)
	i.UID = uint64(h.UID)
	i.GID = uint64(h.GID)
	i.NLink = uint64(h.NLink)
	i.MTime = uint64(h.MTime)
	i.FileSize = uint64(h.FileSize)
	i.Major = uint64(h.Major)
	i.Minor = uint64(h.Minor)
	i.Rmajor = uint64(h.Rmajor)
	i.Rminor = uint64(h.Rminor)
	return i
}

// newc implements RecordFormat for the newc format.
type newc struct {
	magic string
}

// round4 returns the next multiple of 4 close to n.
func round4(n int64) int64 {
	return (n + 3) &^ 0x3
}

type writer struct {
	n   newc
	w   io.Writer
	pos int64
}

// Writer implements RecordFormat.Writer.
func (n newc) Writer(w io.Writer) RecordWriter {
	return NewDedupWriter(&writer{n: n, w: w})
}

func (w *writer) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		return 0, err
	}
	w.pos += int64(n)
	return n, nil
}

func (w *writer) pad() error {
	if o := round4(w.pos); o != w.pos {
		var pad [3]byte
		if _, err := w.Write(pad[:o-w.pos]); err != nil {
			return err
		}
	}
	return nil
}

// WriteRecord writes newc cpio records. It pads the header+name write to 4
// byte alignment and pads the data write as well.
func (w *writer) WriteRecord(f Record) error {
	// Write magic.
	if _, err := w.Write([]byte(w.n.magic)); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	hdr := headerFromInfo(f.Info)
	if f.ReaderAt == nil {
		hdr.FileSize = 0
	}
	hdr.CRC = 0
	if err := binary.Write(buf, binary.BigEndian, hdr); err != nil {
		return err
	}

	hexBuf := make([]byte, hex.EncodedLen(buf.Len()))
	n := hex.Encode(hexBuf, buf.Bytes())
	// It's much easier to debug if we match GNU output format.
	hexBuf = bytes.ToUpper(hexBuf)

	// Write header.
	if _, err := w.Write(hexBuf[:n]); err != nil {
		return err
	}

	// Append NULL char.
	cstr := append([]byte(f.Info.Name), 0)
	// Write name.
	if _, err := w.Write(cstr); err != nil {
		return err
	}

	// Pad to a multiple of 4.
	if err := w.pad(); err != nil {
		return err
	}

	// Some files do not have any content.
	if f.ReaderAt == nil {
		return nil
	}

	// Write file contents.
	m, err := io.Copy(w, uio.Reader(f))
	if err != nil {
		return err
	}
	if m != int64(f.Info.FileSize) {
		return fmt.Errorf("WriteRecord: %s: wrote %d bytes of file instead of %d bytes; archive is now corrupt", f.Info.Name, m, f.Info.FileSize)
	}
	if c, ok := f.ReaderAt.(io.Closer); ok {
		if err := c.Close(); err != nil {
			return err
		}
	}
	if m > 0 {
		return w.pad()
	}
	return nil
}

type reader struct {
	n   newc
	r   io.ReaderAt
	pos int64
}

// Reader implements RecordFormat.Reader.
func (n newc) Reader(r io.ReaderAt) RecordReader {
	return EOFReader{&reader{n: n, r: r}}
}

func (r *reader) read(p []byte) error {
	n, err := r.r.ReadAt(p, r.pos)
	if err == io.EOF {
		return io.EOF
	}
	if err != nil || n != len(p) {
		return fmt.Errorf("ReadAt(pos = %d): got %d, want %d bytes; error %v", r.pos, n, len(p), err)
	}
	r.pos += int64(n)
	return nil
}

func (r *reader) readAligned(p []byte) error {
	err := r.read(p)
	r.pos = round4(r.pos)
	return err
}

// ReadRecord implements RecordReader for the newc cpio format.
func (r *reader) ReadRecord() (Record, error) {
	hdr := header{}
	recPos := r.pos

	Debug("Next record: pos is %d\n", r.pos)

	buf := make([]byte, hex.EncodedLen(binary.Size(hdr))+magicLen)
	if err := r.read(buf); err != nil {
		return Record{}, err
	}

	// Check the magic.
	if magic := string(buf[:magicLen]); magic != r.n.magic {
		return Record{}, fmt.Errorf("reader: magic got %q, want %q", magic, r.n.magic)
	}
	Debug("Header is %v\n", buf)

	// Decode hex header fields.
	dst := make([]byte, binary.Size(hdr))
	if _, err := hex.Decode(dst, buf[magicLen:]); err != nil {
		return Record{}, fmt.Errorf("reader: error decoding hex: %v", err)
	}
	if err := binary.Read(bytes.NewReader(dst), binary.BigEndian, &hdr); err != nil {
		return Record{}, err
	}
	Debug("Decoded header is %v\n", hdr)

	// Get the name.
	nameBuf := make([]byte, hdr.NameLength)
	if err := r.readAligned(nameBuf); err != nil {
		return Record{}, err
	}

	info := hdr.Info()
	info.Name = string(nameBuf[:hdr.NameLength-1])

	recLen := uint64(r.pos - recPos)
	filePos := r.pos
	content := io.NewSectionReader(r.r, r.pos, int64(hdr.FileSize))
	r.pos = round4(r.pos + int64(hdr.FileSize))
	return Record{
		Info:     info,
		ReaderAt: content,
		RecLen:   recLen,
		RecPos:   recPos,
		FilePos:  filePos,
	}, nil
}

func init() {
	formatMap["newc"] = Newc
}
