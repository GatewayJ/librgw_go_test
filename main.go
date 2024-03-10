package main

import (
	"fmt"

	"github.com/google/uuid"
)

import "C"

// g++   -fPIC -shared -o golibrgw/librgw.so  -I /home/jhw/ceph/src/include/rados/*
func main() {
	ret, rgwS := Create()
	if ret != 0 {
		fmt.Printf("RGW Create failed: %v", ret)
	}

	ret, rgwfs := Mount(rgwS,
		"test",
		"T575J0KWSSTWNLO9N3R1",
		"jJtAa3hjg2uUMy6MzTP0xjIfX1yu58idoHzvAsk3",
		0)
	if ret == 0 {
		fmt.Printf("RGW Mounted: %v\n", rgwfs)
	} else {
		fmt.Printf("Failed to mount: %v\n", ret)
	}

	ret, statvfs := StatFs(rgwfs, rgwfs.root_fh, 0)
	if ret == 0 {
		fmt.Printf("Statfs: %+v\n", statvfs)
	} else {
		fmt.Printf("Statfs failed: %v", ret)
	}

	stat := NewStat(0, 0, 0755)
	createMask := SetAttrUID | SetAttrGID | SetAttrMode
	newBucketName := uuid.NewString()
	ret, bucketFs := Mkdir(rgwfs, rgwfs.root_fh, newBucketName, stat,
		createMask, 0) // mb
	if ret == 0 {
		fmt.Printf("Created new bucket: %v  %+v %v %v\n", newBucketName, stat, bucketFs, ret)
	} else {
		fmt.Printf("Failed to create %v: %v\n", newBucketName, ret)
	}
	newDirName := uuid.NewString()
	ret, dirFh := Mkdir(rgwfs, bucketFs, newDirName, stat,
		createMask, 0) //make dir,actually ,a object willappear
	if ret == 0 {
		fmt.Printf("Created new directory: %v  %+v %v %v\n", newDirName, stat, dirFh, ret)
	} else {
		fmt.Printf("Failed to create %v: %v\n", newDirName, ret)
	}

	stat = NewStat(0, 0, 0644)

	newFileName := uuid.NewString()
	ret, fh := CreateFile(rgwfs, dirFh, newFileName, stat,
		createMask, 0, 0) //a object willappear
	if ret == 0 {
		fmt.Printf("Created new file %v in %v: %+v %v %v\n", newFileName, newDirName, stat, fh, ret)
	} else {
		fmt.Printf("Failed to create %v: %v", newFileName, ret)
	}

	fmt.Println("Read first entry:")
	err, eof := ReadDir(rgwfs, bucketFs, "", ReaddirFlagDotDot, func(name string, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Readdir: %v\n", name)
		return nil
	})
	fmt.Println(err, eof)

	fmt.Println("Read seconed entry:")
	err, eof = ReadDir(rgwfs, dirFh, "", ReaddirFlagDotDot, func(name string, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Readdir: %v\n", name)
		return nil
	})
	fmt.Println(err, eof)

	ret = Umount(rgwfs, 0)
	fmt.Println(ret, eof)
}
