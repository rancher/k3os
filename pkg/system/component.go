package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/mount"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

type VersionName string

const (
	VersionCurrent  VersionName = "current"
	VersionPrevious VersionName = "previous"
)

func StatComponentVersion(root, key string, alias VersionName) (os.FileInfo, error) {
	currentPath := filepath.Join(root, key, string(alias))
	currentInfo, err := os.Stat(currentPath)
	if err != nil {
		return nil, err
	}
	if !currentInfo.IsDir() {
		return nil, fmt.Errorf("stat %s: not a directory", currentPath)
	}
	version, err := os.Readlink(currentPath)
	if err != nil {
		return nil, err
	}
	versionPath := filepath.Join(root, key, version)
	versionInfo, err := os.Stat(versionPath)
	if err != nil {
		return nil, err
	}
	if !versionInfo.IsDir() {
		return versionInfo, fmt.Errorf("stat %s: not a directory", versionPath)
	}
	return versionInfo, nil
}

func CopyComponent(src, dst string, remount bool, key string) (bool, error) {
	srcInfo, err := StatComponentVersion(src, key, VersionCurrent)
	if err != nil {
		return false, err
	}
	dstInfo, _ := StatComponentVersion(dst, key, VersionCurrent)
	if dstInfo != nil && dstInfo.Name() == srcInfo.Name() {
		logrus.Infof("skipping %q because destination version matches source: %s", key, dstInfo.Name())
		return false, nil
	}
	if remount {
		if err := mount.Mount("", dst, "none", "remount,rw"); err != nil {
			return false, err
		}
	}

	srcPath := filepath.Join(src, key, srcInfo.Name())
	dstPath := filepath.Join(dst, key, srcInfo.Name())
	dstPrevPath := filepath.Join(dst, key, string(VersionPrevious))
	dstCurrPath := filepath.Join(dst, key, string(VersionCurrent))

	dstCurrTemp := dstCurrPath + `.tmp`
	if err := os.Symlink(filepath.Base(dstPath), dstCurrTemp); err != nil {
		return false, err
	}
	logrus.Debugf("created symlink: %v", dstCurrTemp)
	defer os.Remove(dstCurrTemp) // if this fails, that means it's gone which is correct

	dstTemp, err := ioutil.TempDir(filepath.Split(dstPath))
	if err != nil {
		return false, err
	}
	logrus.Debugf("created temporary dir: %v", dstTemp)
	defer os.RemoveAll(dstTemp) // if this fails, that means it's gone which is correct

	logrus.Debugf("copying: %v -> %v", srcPath, dstTemp)
	if err := copy.Copy(srcPath, dstTemp); err != nil {
		return false, err
	}

	logrus.Debugf("renaming: %v -> %v", dstTemp, dstPath)
	if err := os.Rename(dstTemp, dstPath); err != nil {
		return false, err
	}

	logrus.Debugf("chmod %s %s", srcInfo.Mode(), dstPath)
	if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("removing: %v", dstPrevPath)
	if err := os.Remove(dstPrevPath); err != nil {
		logrus.Warn(err)
	}

	logrus.Debugf("copying: %v -> %v", dstCurrPath, dstPrevPath)
	if err := copy.Copy(dstCurrPath, dstPrevPath); err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("renaming: %v -> %v", dstCurrTemp, dstCurrPath)
	if err := os.Rename(dstCurrTemp, dstCurrPath); err != nil {
		return false, err
	}

	return true, nil
}
