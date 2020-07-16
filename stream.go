package ahocorasick

import (
	"encoding/binary"
	"io"
	"os"
)

func WriteTrie(trie *Trie, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := NewEncoder(f)
	return enc.Encode(trie)
}

func ReadTrie(filename string) (*Trie, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	dec := NewDecoder(f)
	return dec.Decode()
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w,
	}
}

func (enc *Encoder) Encode(trie *Trie) error {

	if err := binary.Write(enc.w, binary.LittleEndian, uint64(len(trie.dict))); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, uint64(len(trie.trans))); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, uint64(len(trie.failLink))); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, uint64(len(trie.dictLink))); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, trie.dict); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, trie.trans); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, trie.failLink); err != nil {
		return err
	}

	if err := binary.Write(enc.w, binary.LittleEndian, trie.dictLink); err != nil {
		return err
	}

	return nil
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r,
	}
}

func (dec *Decoder) Decode() (*Trie, error) {
	var dictLen, transLen, dictLinkLen, failLinkLen uint64

	if err := binary.Read(dec.r, binary.LittleEndian, &dictLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dec.r, binary.LittleEndian, &transLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dec.r, binary.LittleEndian, &dictLinkLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dec.r, binary.LittleEndian, &failLinkLen); err != nil {
		return nil, err
	}

	dict := make([]int64, dictLen)
	if err := binary.Read(dec.r, binary.LittleEndian, dict); err != nil {
		return nil, err
	}

	trans := make([][256]int64, transLen)
	if err := binary.Read(dec.r, binary.LittleEndian, trans); err != nil {
		return nil, err
	}

	failLink := make([]int64, failLinkLen)
	if err := binary.Read(dec.r, binary.LittleEndian, failLink); err != nil {
		return nil, err
	}

	dictLink := make([]int64, dictLinkLen)
	if err := binary.Read(dec.r, binary.LittleEndian, dictLink); err != nil {
		return nil, err
	}

	trie := Trie{dict, trans, failLink, dictLink}
	return &trie, nil
}
