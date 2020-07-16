package ahocorasick

import (
	"compress/gzip"
	"encoding/binary"
	"io"
)

// Encode writes a Trie to w in gzip compressed binary format.
func Encode(w io.Writer, trie *Trie) error {
	enc := newEncoder(w)
	return enc.encode(trie)
}

// Decode reads a Trie in gzip compressed binary format from r.
func Decode(r io.Reader) (*Trie, error) {
	dec := newDecoder(r)
	return dec.decode()
}

type encoder struct {
	w io.Writer
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{
		w,
	}
}

func (enc *encoder) encode(trie *Trie) error {

	w := gzip.NewWriter(enc.w)
	defer w.Close()

	if err := binary.Write(w, binary.LittleEndian, uint64(len(trie.dict))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint64(len(trie.trans))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint64(len(trie.failLink))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint64(len(trie.dictLink))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, trie.dict); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, trie.trans); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, trie.failLink); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, trie.dictLink); err != nil {
		return err
	}

	return nil
}

type decoder struct {
	r io.Reader
}

func newDecoder(r io.Reader) *decoder {
	return &decoder{
		r,
	}
}

func (dec *decoder) decode() (*Trie, error) {

	r, err := gzip.NewReader(dec.r)
	if err != nil {
		return nil, err
	}

	var dictLen, transLen, dictLinkLen, failLinkLen uint64

	if err := binary.Read(r, binary.LittleEndian, &dictLen); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &transLen); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &dictLinkLen); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &failLinkLen); err != nil {
		return nil, err
	}

	dict := make([]int64, dictLen)
	if err := binary.Read(r, binary.LittleEndian, dict); err != nil {
		return nil, err
	}

	trans := make([][256]int64, transLen)
	if err := binary.Read(r, binary.LittleEndian, trans); err != nil {
		return nil, err
	}

	failLink := make([]int64, failLinkLen)
	if err := binary.Read(r, binary.LittleEndian, failLink); err != nil {
		return nil, err
	}

	dictLink := make([]int64, dictLinkLen)
	if err := binary.Read(r, binary.LittleEndian, dictLink); err != nil {
		return nil, err
	}

	trie := Trie{dict, trans, failLink, dictLink}
	return &trie, nil
}
