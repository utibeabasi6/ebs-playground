package main

import (
	"fmt"
	"log"

	human "github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/shirou/gopsutil/disk"
)

func fetchDiskSpace() {
	formatter := "%-14s %7s %7s %7s %4s %s\n"
	fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	parts, _ := disk.Partitions(false)
	for _, p := range parts {
		device := p.Mountpoint
		s, _ := disk.Usage(device)

		if s.Total == 0 {
			continue
		}

		percent := fmt.Sprintf("%2.f%%", s.UsedPercent)

		fmt.Printf(formatter,
			s.Fstype,
			human.Bytes(s.Total),
			human.Bytes(s.Used),
			human.Bytes(s.Free),
			percent,
			p.Mountpoint,
		)
	}

}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		fetchDiskSpace()

		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	log.Println("Starting server on port 3000")
	log.Fatal(app.Listen(":3000"))
}
