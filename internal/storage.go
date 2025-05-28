package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	VERSION uint8 = 1
)

var byteOrder = binary.LittleEndian

type Storage struct {
	file *os.File
}

type Header struct {
	Version      uint8
	Timestamp    int64
	MessageCount uint16
}

func NewStorage(fname string) (*Storage, error) {
	isConversationExists := isFileExists(fname)

	file, err := openOrCreateFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation file: %w", err)
	}

	storage := &Storage{
		file: file,
	}

	if !isConversationExists {
		header := Header{
			Version:      1,
			Timestamp:    time.Now().Unix(),
			MessageCount: 0,
		}
		if err := storage.writeHeader(header); err != nil {
			return nil, fmt.Errorf("failed to write header: %w", err)
		}
	}

	return storage, nil
}

func (c *Storage) writeHeader(header Header) error {
	if _, err := c.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to start of file: %w", err)
	}

	if err := binary.Write(c.file, byteOrder, header); err != nil {
		return fmt.Errorf("failed to write header version: %w", err)
	}

	// Sync the file to ensure the header is written to disk
	if err := c.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

func (c *Storage) ReadHeader() (Header, error) {
	var header Header

	if _, err := c.file.Seek(0, io.SeekStart); err != nil {
		return header, fmt.Errorf("failed to seek to start of file: %w", err)
	}

	if err := binary.Read(c.file, byteOrder, &header); err != nil {
		return header, fmt.Errorf("failed to read header version: %w", err)
	}

	return header, nil
}

// AddMessage will add a message to the conversation
// Format: [Length] [Encoded Message]
func (c *Storage) AddMessage(message *Message) error {
	header, err := c.ReadHeader()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(message); err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	if _, err := c.file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}
	if err := binary.Write(c.file, byteOrder, uint64(buf.Len())); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}
	if err := binary.Write(c.file, byteOrder, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	header.MessageCount++
	if err := c.writeHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	return nil
}

func (c *Storage) ReadMessages() ([]Message, error) {
	header, err := c.ReadHeader()
	if err != nil {
		return nil, err
	}

	var messages []Message
	for i := uint16(0); i < uint16(header.MessageCount); i++ {
		var length uint64
		if err := binary.Read(c.file, byteOrder, &length); err != nil {
			return nil, fmt.Errorf("failed to read message length: %w", err)
		}

		buf := make([]byte, length)
		if _, err := io.ReadFull(c.file, buf); err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var message Message
		if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&message); err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (c *Storage) Close() error {
	if c.file == nil {
		return nil
	}
	
	return c.file.Close()
}
