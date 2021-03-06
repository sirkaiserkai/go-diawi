package godiawi

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type formWriter struct {
	buff *bytes.Buffer
	mw   *multipart.Writer
}

func newformWriter() formWriter {
	fw := formWriter{}
	fw.buff = &bytes.Buffer{}
	fw.mw = multipart.NewWriter(fw.buff)

	return fw
}

func (fw *formWriter) AddFormFile(fieldName, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fieldWriter, err := fw.mw.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fieldWriter, f); err != nil {
		return err
	}

	return nil
}

func (fw *formWriter) AddField(fieldName string, fieldValue interface{}) error {
	switch val := fieldValue.(type) {
	case string:
		return fw.AddStringField(fieldName, val)
	case bool:
		return fw.AddBoolField(fieldName, val)
	case []string:
		return fw.AddStringSliceField(fieldName, val)
	default:
		return fmt.Errorf("Unexpected type %T\n", val)
	}
}

func (fw *formWriter) AddStringField(fieldName string, fieldValue string) error {
	fieldWriter, err := fw.mw.CreateFormField(fieldName)
	if err != nil {
		return err
	}

	if _, err = fieldWriter.Write([]byte(fieldValue)); err != nil {
		return err
	}

	return nil
}

func (fw *formWriter) AddBoolField(fieldName string, fieldValue bool) error {
	fieldWriter, err := fw.mw.CreateFormField(fieldName)
	if err != nil {
		return err
	}

	var val string
	if fieldValue {
		val = "1"
	} else {
		val = "0"
	}

	if _, err = fieldWriter.Write([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (fw *formWriter) AddStringSliceField(fieldName string, fieldValue []string) error {
	fieldWriter, err := fw.mw.CreateFormField(fieldName)
	if err != nil {
		return err
	}

	val := ""
	for i, s := range fieldValue {
		val += s
		if i != len(fieldValue) {
			val += ", "
		}
	}

	if _, err = fieldWriter.Write([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (fw *formWriter) GetBuffer() *bytes.Buffer {
	return fw.buff
}

func (fw *formWriter) Close() {
	fw.mw.Close()
}
