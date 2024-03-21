// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package try_on

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson569002f9DecodeTryOnInternalPkgDeliveryTryOn(in *jlexer.Lexer, out *tryOnRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "user_image_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserImageID).UnmarshalText(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson569002f9EncodeTryOnInternalPkgDeliveryTryOn(out *jwriter.Writer, in tryOnRequest) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.ClothesID).MarshalText())
	}
	if (in.UserImageID).IsDefined() {
		const prefix string = ",\"user_image_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserImageID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v tryOnRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson569002f9EncodeTryOnInternalPkgDeliveryTryOn(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v tryOnRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson569002f9EncodeTryOnInternalPkgDeliveryTryOn(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *tryOnRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson569002f9DecodeTryOnInternalPkgDeliveryTryOn(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *tryOnRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson569002f9DecodeTryOnInternalPkgDeliveryTryOn(l, v)
}
