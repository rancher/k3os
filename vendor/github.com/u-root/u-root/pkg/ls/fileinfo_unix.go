// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
	"golang.org/x/sys/unix"
)

// Matches characters which would interfere with ls's formatting.
var unprintableRe = regexp.MustCompile("[[:cntrl:]\n]")

// Since `os.FileInfo` is an interface, it is difficult to tweak some of its
// internal values. For example, replacing the starting directory with a dot.
// `extractImportantParts` populates our own struct which we can modify at will
// before printing.
type FileInfo struct {
	Name          string
	Mode          os.FileMode
	Rdev          uint64
	UID, GID      uint32
	Size          int64
	MTime         time.Time
	SymlinkTarget string
}

func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	var link string

	s := fi.Sys().(*syscall.Stat_t)
	if fi.Mode()&os.ModeType == os.ModeSymlink {
		if l, err := os.Readlink(path); err != nil {
			link = err.Error()
		} else {
			link = l
		}
	}

	return FileInfo{
		Name:          fi.Name(),
		Mode:          fi.Mode(),
		Rdev:          uint64(s.Rdev),
		UID:           s.Uid,
		GID:           s.Gid,
		Size:          fi.Size(),
		MTime:         fi.ModTime(),
		SymlinkTarget: link,
	}
}

// Name returns a printable file name.
func (fi FileInfo) PrintableName() string {
	return unprintableRe.ReplaceAllLiteralString(fi.Name, "?")
}

// Without this cache, `ls -l` is orders of magnitude slower.
var (
	uidCache = map[uint32]string{}
	gidCache = map[uint32]string{}
)

// Convert uid to username, or return uid on error.
func lookupUserName(id uint32) string {
	if s, ok := uidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if u, err := user.LookupId(s); err == nil {
		s = u.Username
	}
	uidCache[id] = s
	return s
}

// Convert gid to group name, or return gid on error.
func lookupGroupName(id uint32) string {
	if s, ok := gidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if g, err := user.LookupGroupId(s); err == nil {
		s = g.Name
	}
	gidCache[id] = s
	return s
}

// Stringer provides a consistent way to format FileInfo.
type Stringer interface {
	// FileString formats a FileInfo.
	FileString(fi FileInfo) string
}

// NameStringer is a Stringer implementation that just prints the name.
type NameStringer struct{}

// FileString implements Stringer.FileString and just returns fi's name.
func (ns NameStringer) FileString(fi FileInfo) string {
	return fi.PrintableName()
}

// QuotedStringer is a Stringer that returns the file name surrounded by qutoes
// with escaped control characters.
type QuotedStringer struct{}

// FileString returns the name surrounded by quotes with escaped control characters.
func (qs QuotedStringer) FileString(fi FileInfo) string {
	return fmt.Sprintf("%#v", fi.Name)
}

// LongStringer is a Stringer that returns the file info formatted in `ls -l`
// long format.
type LongStringer struct {
	Human bool
	Name  Stringer
}

// FileString implements Stringer.FileString.
func (ls LongStringer) FileString(fi FileInfo) string {
	// Golang's FileMode.String() is almost sufficient, except we would
	// rather use b and c for devices.
	replacer := strings.NewReplacer("Dc", "c", "D", "b")

	// Ex: crw-rw-rw-  root  root  1, 3  Feb 6 09:31  null
	pattern := "%[1]s\t%[2]s\t%[3]s\t%[4]d, %[5]d\t%[7]v\t%[8]s"
	if fi.Mode&os.ModeDevice == 0 && fi.Mode&os.ModeCharDevice == 0 {
		// Ex: -rw-rw----  myuser  myuser  1256  Feb 6 09:31  recipes.txt
		pattern = "%[1]s\t%[2]s\t%[3]s\t%[6]s\t%[7]v\t%[8]s"
	}

	var size string
	if ls.Human {
		size = humanize.Bytes(uint64(fi.Size))
	} else {
		size = strconv.FormatInt(fi.Size, 10)
	}

	s := fmt.Sprintf(pattern,
		replacer.Replace(fi.Mode.String()),
		lookupUserName(fi.UID),
		lookupGroupName(fi.GID),
		unix.Major(fi.Rdev),
		unix.Minor(fi.Rdev),
		size,
		fi.MTime.Format("Jan _2 15:04"),
		ls.Name.FileString(fi))

	if fi.Mode&os.ModeType == os.ModeSymlink {
		s += fmt.Sprintf(" -> %v", fi.SymlinkTarget)
	}
	return s
}
