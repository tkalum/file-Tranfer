// internal/transfer/transfer.go
package transfer


import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
	"bufio"
	
)

const TransferPort = 24243


func ReadMessage(conn net.Conn) (string, error) {
	reader  :=bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	message = strings.TrimSpace(message)

	return message, nil
}

func Listener(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	defer l.Close()

	fmt.Printf("Server listening on port %d\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()

			filename, err := ReadMessage(c)
			if err != nil {
				log.Printf("Error reading filename: %v", err)
				return
			}
			
			filesize , err := ReadMessage(c)
			if err != nil {
				log.Printf("Error reading filesize: %v", err)
				return
			}
			log.Printf("Receiving file: %s (%s bytes)\n", filename, filesize)

			err = Receivefile(c, filename, filesize)
			if err != nil {
				log.Fatal(err)
			}
			
			fmt.Println("File received successfully")
		}(conn)
		
	}
}

func Dialer(filename string, host string) {
	var d net.Dialer
	port := TransferPort // Default port
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Correctly combine the host and port string for the dialer.
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	file , err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	filesize := fileinfo.Size()

	conn.Write([]byte(filename + "\n"))
	conn.Write([]byte(fmt.Sprintf("%d\n", filesize)))
	err = SendFile(conn, filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File sent successfully")
}

func SendFile(conn net.Conn, filename string) error {
	var transferred int64
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	filesize := fileInfo.Size()
	if err != nil {
		return err
	}

	buffer := make([]byte, 64*1024)
	lastPrint := time.Now()

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			return err
		}
		transferred += int64(n)
		if time.Since(lastPrint) > 100*time.Millisecond {
			percent := float64(transferred) / float64(filesize) * 100
			fmt.Printf("\rTransferring: %.1f%% (%d/%d MB)", 
				percent, transferred/(1024*1024), filesize/(1024*1024))
			lastPrint = time.Now()
		}

	}
	fmt.Printf("\rTransfer complete: 100.0%% (%d/%d MB)\n", 
		filesize/(1024*1024), filesize/(1024*1024))
	return nil
}

func Receivefile(conn net.Conn, filename string, filesize string) error {
	var received int64
	filesizeInt, err := strconv.ParseInt(filesize, 10, 64)
	if err != nil {
		return err
	}
	file, err := os.Create("receive_" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 64*1024)
	lastPrint := time.Now()

	for {
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}
		received += int64(n)
		if time.Since(lastPrint) > 100*time.Millisecond {
			percent := float64(received) / float64(filesizeInt) * 100
			fmt.Printf("\rReceiving: %.1f%% (%d/%d MB)", 
				percent, received/(1024*1024), filesizeInt/(1024*1024))
			lastPrint = time.Now()
		}
	}
	fmt.Printf("\rReceive complete: 100.0%% (%d/%d MB)\n", 
		filesizeInt/(1024*1024), filesizeInt/(1024*1024))
	return nil
}