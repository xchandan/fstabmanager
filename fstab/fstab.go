package fstab

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type FsTabRecord struct {
	Comment        string
	Filesystem     string
	MountPoint     string
	FilesystemType string
	MountOptions   string
	Dump           string
	Pass           string
}

type FsTabDb struct {
	Path    string
	Records []*FsTabRecord
}

func (fr *FsTabRecord) String() string {
	return fmt.Sprintf("Comment: %s, Filesystem: %s, MountPoint: %s FilesystemType: %s, MountOptions: %s, Dump: %s, Pass: %s", fr.Comment, fr.Filesystem, fr.MountPoint, fr.FilesystemType, fr.MountOptions, fr.Dump, fr.Pass)
}

func (f *FsTabDb) Parse() error {
	fd, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			rec := FsTabRecord{
				Comment: line,
			}
			f.Records = append(f.Records, &rec)
		} else {
			tokens := strings.Fields(line)
			log.Println(tokens)
			rec := FsTabRecord{
				Filesystem:     tokens[0],
				MountPoint:     tokens[1],
				FilesystemType: tokens[2],
				MountOptions:   tokens[3],
				Dump:           tokens[4],
				Pass:           tokens[5],
			}
			f.Records = append(f.Records, &rec)
		}
	}

	return nil
}

func (f *FsTabDb) AddMount(filesystem, mountpt, fstype, mountopts, dump, pass string) error {
	var exists bool = false
	log.Println("AddMount", filesystem, mountpt, fstype, mountopts, dump, pass)
	for _, fr := range f.Records {
		if fr.Filesystem == filesystem {
			exists = true
			if fr.MountPoint == mountpt &&
				fr.FilesystemType == fstype &&
				fr.MountOptions == mountopts &&
				fr.Dump == dump &&
				fr.Pass == pass {
				log.Println("  Duplicate entry")
				return nil
			} else if fr.MountPoint == mountpt {
				log.Println("  same fs and mountpt, update the mount opts etc")
				fr.FilesystemType = fstype
				fr.MountOptions = mountopts
				fr.Dump = dump
				fr.Pass = pass
				break

			} else {
				log.Println("  same fs and different mountpt")
				exists = false
				break

			}
		} else if fr.MountPoint == mountpt {
			return fmt.Errorf("duplicate mount point exists %s is mounted on %s", fr.Filesystem, fr.MountPoint)
		}

	}
	if !exists {
		log.Println("  New mount")

		f.Records = append(f.Records, &FsTabRecord{
			Filesystem:     filesystem,
			MountPoint:     mountpt,
			FilesystemType: fstype,
			MountOptions:   mountopts,
			Dump:           dump,
			Pass:           pass,
		})
	}
	return nil
}

func (f *FsTabDb) Save() {
	fd, _ := os.OpenFile(f.Path+".new", os.O_RDWR|os.O_CREATE, 0644)
	defer fd.Close()

	for _, fr := range f.Records {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", fr.Filesystem, fr.MountPoint, fr.FilesystemType, fr.MountOptions, fr.Dump, fr.Pass)
		fd.WriteString(line)
	}
}

func NewFsTabDb(path string) (*FsTabDb, error) {

	newDb := &FsTabDb{Path: path}
	if err := newDb.Parse(); err != nil {
		return nil, err
	}

	return newDb, nil
}
