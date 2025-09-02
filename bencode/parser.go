// Goal-->Parse .torrent files and Extract metadata like announce URL, piece hashes, file info
package bencode

import (
	"fmt"
	"io"
	"os"
	"unicode"

	"main.go/types"
)

type Parser struct {
	data []byte
	pos  int
}

// Creates a parser instance
func NewParser(data []byte) *Parser {
	return &Parser{data: data, pos: 0}
}

func ParseTorrentFile(filename string) (*types.TorrentFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open torrent file: %v", err)
	}

	data, err := io.ReadAll(file) //data is for SHA-1 Decryption
	if err != nil {
		return nil, fmt.Errorf("Failed to read torrent file: %v", err)
	}

	parser := NewParser(data)
	torrentData, err := parser.Parse() //To extract differnt fields
	if err != nil {
		return nil, fmt.Errorf("Failed to parse torrent file: %v", err)
	}

	return buildTorrentFile(torrentData, data)
}

/*
Bdecode (recursive descent)
If it starts with:
i → parse integer until e
l → parse list until matching e
d → parse dict: alternating string keys and values until e
\d → parse string: <len>:<raw bytes>
Validate as you go (no leading zeros, proper terminators, boundary checks).
This is a Parse method and not a function
interface{} lets you return any type
*/
func (p *Parser) Parse() (interface{}, error) {
	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unexpected end of data")
	}
	switch {
	case unicode.IsDigit(rune(p.data[p.pos])):
		return p.parseString()
	case p.data[p.pos] == 'i':
		return p.parseInteger()
	case p.data[p.pos] == 'l':
		return p.parseList()
	case p.data[p.pos] == 'd':
		return p.parseDictionary()
	default:
		return nil, fmt.Errorf("invalid bencode data at position %d", p.pos)
	}
}

// Parsing String, Interger, List, Dictionary
func (p *Parser) parseString() (string, error) {
	start := p.pos
	for p.pos < len(p.data) && p.data[p.pos] != ':' {
		if !unicode.IsDigit(rune(p.data[p.pos])) {
			return "", fmt.Errorf("invalid string length at position %d", p.pos)
		}
		p.pos++
	}

	if p.pos >= len(p.data) {
		return "", fmt.Errorf("unexpected end of data while parsing string length")
	}

	lengthStr := string(p.data[start:p.pos])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", fmt.Errorf("invalid string length: %v", err)
	}

	p.pos++ // skip ':'

	if p.pos+length > len(p.data) {
		return "", fmt.Errorf("string length exceeds data bounds")
	}

	result := string(p.data[p.pos : p.pos+length])
	p.pos += length

	return result, nil
}

// parseInteger parses bencode integer format: i<integer>e
func (p *Parser) parseInteger() (int, error) {
	if p.data[p.pos] != 'i' {
		return 0, fmt.Errorf("expected 'i' at position %d", p.pos)
	}
	p.pos++

	start := p.pos
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		p.pos++
	}

	if p.pos >= len(p.data) {
		return 0, fmt.Errorf("unexpected end of data while parsing integer")
	}

	intStr := string(p.data[start:p.pos])
	p.pos++ // skip 'e'

	result, err := strconv.Atoi(intStr)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %v", err)
	}

	return result, nil
}

// parseList parses bencode list format: l<elements>e
func (p *Parser) parseList() ([]interface{}, error) {
	if p.data[p.pos] != 'l' {
		return nil, fmt.Errorf("expected 'l' at position %d", p.pos)
	}
	p.pos++

	var result []interface{}
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		element, err := p.Parse()
		if err != nil {
			return nil, err
		}
		result = append(result, element)
	}

	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unexpected end of data while parsing list")
	}

	p.pos++ // skip 'e'
	return result, nil
}

// parseDictionary parses bencode dictionary format: d<key-value pairs>e
func (p *Parser) parseDictionary() (map[string]interface{}, error) {
	if p.data[p.pos] != 'd' {
		return nil, fmt.Errorf("expected 'd' at position %d", p.pos)
	}
	p.pos++

	result := make(map[string]interface{})
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		// Parse key (must be string)
		key, err := p.parseString()
		if err != nil {
			return nil, fmt.Errorf("failed to parse dictionary key: %v", err)
		}

		// Parse value
		value, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse dictionary value: %v", err)
		}

		result[key] = value
	}

	if p.pos >= len(p.data) {
		return nil, fmt.Errorf("unexpected end of data while parsing dictionary")
	}

	p.pos++ // skip 'e'
	return result, nil
}
