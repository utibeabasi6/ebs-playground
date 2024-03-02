package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/adelowo/gulter"
	storage "github.com/adelowo/gulter/storage"
	human "github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/shirou/gopsutil/disk"
)

type EbsDevice struct {
	MountPath  string `json:"mountPath"`
	Size       string `json:"size"`
	Used       string `json:"used"`
	UsePercent string `json:"usePercent"`
	Available  string `json:"available"`
}

func fetchDiskSpace() ([]EbsDevice, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	var lvs []EbsDevice
	for _, p := range parts {
		device := p.Mountpoint
		if strings.HasPrefix(device, "/home/ubuntu/mounts") {
			s, err := disk.Usage(device)
			if err != nil {
				return nil, err
			}
			percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
			dvc := EbsDevice{
				MountPath:  device,
				Size:       human.Bytes(s.Used + s.Free),
				Used:       human.Bytes(s.Used),
				Available:  human.Bytes(s.Free),
				UsePercent: percent,
			}
			lvs = append(lvs, dvc)
		}
	}

	return lvs, nil

}

var mountList = []string{"lv1", "lv2"}

func main() {
	app := fiber.New(fiber.Config{StreamRequestBody: true})

	lv1DiskStore, err := storage.NewDiskStorage(fmt.Sprintf("/home/ubuntu/mounts/%s", mountList[0]))
	if err != nil {
		panic(err.Error())
	}

	lv2DiskStore, err := storage.NewDiskStorage(fmt.Sprintf("/home/ubuntu/mounts/%s", mountList[1]))
	if err != nil {
		panic(err.Error())
	}

	uploadLv1Handler, err := gulter.New(
		gulter.WithMaxFileSize(10<<30),
		gulter.WithStorage(lv1DiskStore),
	)

	if err != nil {
		panic(err.Error())
	}

	uploadLv2Handler, err := gulter.New(
		gulter.WithMaxFileSize(10<<30),
		gulter.WithStorage(lv2DiskStore),
	)

	if err != nil {
		panic(err.Error())
	}

	app.Use("/upload/lv1", adaptor.HTTPMiddleware(uploadLv1Handler.Upload("lv1"))).Post("/upload/lv1", func(c *fiber.Ctx) error {
		devices, err := fetchDiskSpace()
		if err != nil {
			return err
		}
		return c.JSON(devices)
	})
	app.Use("/upload/lv2", adaptor.HTTPMiddleware(uploadLv2Handler.Upload("lv2"))).Post("/upload/lv2", func(c *fiber.Ctx) error {
		devices, err := fetchDiskSpace()
		if err != nil {
			return err
		}
		return c.JSON(devices)
	})

	app.Get("/devices", func(c *fiber.Ctx) error {
		devices, err := fetchDiskSpace()
		if err != nil {
			return err
		}
		return c.JSON(devices)
	})

	log.Println("Starting server on port 3000")
	log.Fatal(app.Listen(":3000"))
}
