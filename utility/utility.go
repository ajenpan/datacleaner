package utility

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type FileInfoWrapper struct {
	Path string
	Info os.FileInfo
}

//7z sopport:
//Packing / unpacking: 7z, XZ, BZIP2, GZIP, TAR, ZIP and WIM
//Unpacking only: AR, ARJ, CAB, CHM, CPIO, CramFS, DMG, EXT, FAT, GPT, HFS, IHEX, ISO, LZH, LZMA, MBR, MSI, NSIS, NTFS, QCOW2, RAR, RPM, SquashFS, UDF, UEFI, VDI, VHD, VHDX, VMDK, WIM, XAR and Z.

var CompressedExts = mapset.NewSet(".7z", ".xz", ".bzip2", ".gzip", ".tar", ".zip", ".wim", ".ar", ".arj", ".cab",
	".chm", ".cpio", ".cramfs", ".dmg", ".ext", ".fat", ".gpt", ".hfs", ".ihex", ".iso", ".lzh",
	".lzma", ".mbr", ".msi", ".nsis", ".ntfs", ".qcow2", ".rar", ".rpm", ".squashfs",
	".udf", ".uefi", ".vdi", ".vhd", ".vhdx", ".vmdk", ".wim", ".xar", ".z")

func IsCompressed(info fs.FileInfo) bool {
	return CompressedExts.Contains(strings.ToLower(filepath.Ext(info.Name())))
}

func AllFiles(targetpath string) ([]*FileInfoWrapper, error) {
	ret := []*FileInfoWrapper{}
	err := filepath.Walk(targetpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ret = append(ret, &FileInfoWrapper{
			Path: path,
			Info: info,
		})
		return nil
	})
	return ret, err
}

func AllCompressed(tagetDir string) ([]*FileInfoWrapper, error) {
	ret := []*FileInfoWrapper{}
	err := WalkCompressed(tagetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		ret = append(ret, &FileInfoWrapper{path, info})
		return nil
	})
	return ret, err
}

func WalkCompressed(dir string, onFile filepath.WalkFunc) error {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if IsCompressed(info) {
			return onFile(path, info, err)
		}
		return nil
	})
	return err
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func DecompressionAll(tagetDir string, failedDir string) {
	files, err := AllCompressed(tagetDir)
	if err != nil {
		panic(err)
	}
	total := len(files)

	for i, file := range files {
		path := file.Path

		ext := filepath.Ext(path)
		suffix := []byte(ext)
		suffix[0] = '_'
		outPath := strings.TrimSuffix(path, ext) + string(suffix)

		cmd := exec.Command("C:\\Program Files\\7-Zip\\7z.exe", "x", path, "-o"+outPath)

		if err := cmd.Run(); err != nil {

			//move to other folder?
			fmt.Printf("[%d/%d] failed path:%v, size:%v, err:%v\n", i, total, path, ByteCountIEC(file.Info.Size()), err.Error())

			newName := filepath.Join(failedDir, file.Info.Name())
			err := os.Rename(path, newName)
			if err != nil {
				fmt.Println("rename err: ", err)
			}

		} else {

			fmt.Printf("[%d/%d] success path:%v, size:%v\n", i, total, path, ByteCountIEC(file.Info.Size()))

			if err := os.Remove(path); err != nil {
				fmt.Println("remove error:", err)
			}
		}

	}
}

func RemoveEmptyDirAndFiles(dir string) {
	paths := []*FileInfoWrapper{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if (info.IsDir()) || (!info.IsDir() && info.Size() == 0) {
			paths = append(paths, &FileInfoWrapper{path, info})
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for i := len(paths) - 1; i >= 0; i-- {
		targetPath := paths[i].Path
		info := paths[i].Info
		needRemove := false
		if info.IsDir() {
			needRemove, _ = IsDirEmpty(targetPath)
		} else {
			needRemove = info.Size() == 0
		}

		if needRemove {
			if err := os.Remove(targetPath); err != nil {
				panic(err)
			}
			fmt.Println("removed:", targetPath, "size:", info.Size())
		}
	}
}

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func Sha1File(filePath string) []byte {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}

func RemoveSameFiles(targetpath string) {
	allfiles, err := AllFiles(targetpath)
	if err != nil {
		panic(err)
	}

	store := make(map[string][]*FileInfoWrapper)

	for _, file := range allfiles {
		info := file.Info
		store[info.Name()] = append(store[info.Name()], &FileInfoWrapper{
			Path: file.Path,
			Info: info,
		})
	}

	temp := make(map[string][]*FileInfoWrapper)
	for k, v := range store {
		if len(v) > 1 {
			temp[k] = v
		}
	}

	for _, v := range temp {
		keep := v[0]
		sizeMap := make(map[int64]int)
		sizeMap[keep.Info.Size()] = 0

		for i := 1; i < len(v); i++ {
			if repeati, has := sizeMap[v[i].Info.Size()]; has {
				fmt.Println("remove: ", v[i].Path, "; 重复: ", v[repeati].Path)
				err := os.Remove(v[i].Path)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				sizeMap[v[i].Info.Size()] = i
			}
		}
	}
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
