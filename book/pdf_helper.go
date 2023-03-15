package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var mode = flag.String("mode", "union", "operate file mode")

func main() {
	flag.Parse()

	var (
		sourceFile = "Kubernetes in action.pdf"
	)

	switch *mode {
	case "union":
		if err := unionFile(sourceFile); err != nil {
			log.Fatalf("[F] Union file failure, nest error: %v", err)
		}
	case "split":
		parts, err := splitFile(sourceFile, 20*1024*1024)
		if err != nil {
			log.Fatalf("[F] Split file failure, nest error: %v", err)
		}
		for _, part := range parts {
			fmt.Printf("=> %s\r\n", part)
		}
	default:
		log.Printf("[F] Not support mode, just include [union/split]")
	}
}

func splitFile(name string, size int) ([]string, error) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		buf   [20 * 1024 * 1024]byte
		i     int
		parts = make([]string, 0, 8)

		partFile = "Kubernetes in action.pdf.part.%d"
	)
	for {
		i++
		n, err := file.Read(buf[0:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var part = fmt.Sprintf(partFile, i)
		if err := os.WriteFile(part, buf[:n], 0644); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
}

func unionFile(to string) error {
	var (
		partFile = "Kubernetes in action.pdf.part.%d"
		i        = 1
	)

	file, err := os.OpenFile(to, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		var part = fmt.Sprintf(partFile, i)
		fi, err := os.Stat(part)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("panic: part is one dir, path: %v", part)
		}
		i++

		buf, err := os.ReadFile(part)
		if err != nil {
			return err
		}
		if _, err := file.Write(buf); err != nil {
			return err
		}
	}

	return nil
}
