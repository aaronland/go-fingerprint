package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"sync"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		im, _, err := image.Decode(r)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		bounds := im.Bounds()

		w := bounds.Dx()
		h := bounds.Dy()

		wg := new(sync.WaitGroup)
		mu := new(sync.RWMutex)
		colours := new(sync.Map)

		for x := 0; x < w; x++ {

			for y := 0; y < h; y++ {

				wg.Add(1)

				go func(im image.Image, x int, y int) {

					defer wg.Done()
					c := im.At(x, y)

					r, g, b, a := c.RGBA()
					key := fmt.Sprintf("%d,%d,%d,%d", r, g, b, a)

					mu.Lock()
					defer mu.Unlock()

					v, exists := colours.Load(key)

					var points [][2]int

					if exists {
						points = v.([][2]int)
					} else {
						points = make([][2]int, 0)
					}

					points = append(points, [2]int{x, y})
					colours.Store(key, points)

				}(im, x, y)
			}
		}

		wg.Wait()

		count := 0

		colours.Range(func(k interface{}, v interface{}) bool {

			key := k.(string)
			points := v.([][2]int)

			log.Println(key, len(points))
			return true
		})

		log.Println(count)

	}

}
