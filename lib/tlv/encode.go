package tlv

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"log"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct && reflect.TypeOf(v).Kind() != reflect.Ptr {
		return nil, errors.New("tlv need struct value to encode")
	}
	return buildTLV(0, v)
}

func makeTLV(tag int, value []byte) []byte {
	buf := new(bytes.Buffer)
	if tag != 0 {
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("makeTLV binary write type error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(value))); err != nil {
			log.Printf("makeTLV binary write value error: %+v", err)
		}
	}

	if err := binary.Write(buf, binary.BigEndian, value); err != nil {
		log.Printf("makeTLV write value error %+v", err)
	}

	return buf.Bytes()
}

func buildTLV(tag int, v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	value := reflect.ValueOf(v)

	if marshaler, ok := value.Interface().(encoding.BinaryMarshaler); ok {
		bin, err := marshaler.MarshalBinary()
		if err != nil {
			return nil, err
		}

		return makeTLV(tag, bin), nil
	}

	value = reflect.Indirect(value)
	switch value.Kind() {
	case reflect.Int8:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(1)); err != nil {
			log.Printf("len(1) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int16:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(2)); err != nil {
			log.Printf("len(2) write error %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int32:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(4)); err != nil {
			log.Printf("len(4) write error %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Int64:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(8)); err != nil {
			log.Printf("len(8) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint8:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(1)); err != nil {
			log.Printf("len(1) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint16:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(2)); err != nil {
			log.Printf("len(2) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint32:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(4)); err != nil {
			log.Printf("len(4) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Uint64:
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(8)); err != nil {
			log.Printf("len(8) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.String:
		str := v.(string)
		if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
			log.Printf("tag write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(str))); err != nil {
			log.Printf("len(str) write error: %+v", err)
		}
		if err := binary.Write(buf, binary.BigEndian, []byte(str)); err != nil {
			log.Printf("value write error: %+v", err)
		}
		return buf.Bytes(), nil
	case reflect.Struct:
		for i := 0; i < value.Type().NumField(); i++ {
			field := value.Field(i)
			if !hasValue(field) {
				continue
			}
			structField := value.Type().Field(i)
			tlvTag, hasTLV := structField.Tag.Lookup("tlv")
			if !hasTLV {
				return nil, errors.New("field " + structField.Name + " need tag `tlv`")
			}
			tagVal, err := strconv.Atoi(tlvTag)
			if err != nil {
				return nil, err
			}
			subValue, err := buildTLV(tagVal, field.Interface())
			if err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.BigEndian, subValue); err != nil {
				return nil, err
			}
		}

		return makeTLV(tag, buf.Bytes()), nil
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			if err := binary.Write(buf, binary.BigEndian, uint16(tag)); err != nil {
				log.Printf("Binary write error: %+v", err)
			}
			if err := binary.Write(buf, binary.BigEndian, uint16(value.Len())); err != nil {
				log.Printf("Binary write error: %+v", err)
			}
			if err := binary.Write(buf, binary.BigEndian, v); err != nil {
				log.Printf("Binary write error: %+v", err)
			}
			return buf.Bytes(), nil
		} else if value.Type().Elem().Kind() == reflect.Ptr || value.Type().Elem().Kind() == reflect.Struct ||
			isNumber(value.Type().Elem()) {
			for i := 0; i < value.Len(); i++ {
				elem := value.Index(i)
				if value.Type().Elem().Kind() == reflect.Struct {
					elem = elem.Addr()
				}
				elemBuf, err := buildTLV(tag, elem.Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(elemBuf)
			}
			return buf.Bytes(), nil
		} else {
			return nil, errors.New("value type `Slice of " + value.Type().Elem().Name() + "` is not support TLV encode")
		}
	default:
		return nil, errors.New("value type " + value.Type().String() + " is not support TLV encode")
	}
}
