// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

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

func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain(in *jlexer.Lexer, out *UserImage) {
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
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "image":
			out.Image = string(in.String())
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain(out *jwriter.Writer, in UserImage) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.UserID).MarshalText())
	}
	if in.Image != "" {
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserImage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserImage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserImage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserImage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain1(in *jlexer.Lexer, out *Type) {
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
		case "name":
			out.Name = string(in.String())
		case "subtypes":
			if in.IsNull() {
				in.Skip()
				out.Subtypes = nil
			} else {
				in.Delim('[')
				if out.Subtypes == nil {
					if !in.IsDelim(']') {
						out.Subtypes = make([]Subtype, 0, 0)
					} else {
						out.Subtypes = []Subtype{}
					}
				} else {
					out.Subtypes = (out.Subtypes)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Subtype
					(v1).UnmarshalEasyJSON(in)
					out.Subtypes = append(out.Subtypes, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain1(out *jwriter.Writer, in Type) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if len(in.Subtypes) != 0 {
		const prefix string = ",\"subtypes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v2, v3 := range in.Subtypes {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Type) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Type) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Type) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Type) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain1(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain2(in *jlexer.Lexer, out *TryOnResult) {
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
		case "image":
			out.Image = string(in.String())
		case "rating":
			out.Rating = int(in.Int())
		case "user_image_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserImageID).UnmarshalText(data))
			}
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain2(out *jwriter.Writer, in TryOnResult) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Image != "" {
		const prefix string = ",\"image\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Image))
	}
	if in.Rating != 0 {
		const prefix string = ",\"rating\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Rating))
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
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ClothesID).MarshalText())
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TryOnResult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TryOnResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TryOnResult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TryOnResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain2(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain3(in *jlexer.Lexer, out *TryOnResponse) {
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
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "user_image_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserImageID).UnmarshalText(data))
			}
		case "try_on_result_id":
			out.TryOnResultID = string(in.String())
		case "try_on_result_dir":
			out.TryOnResultDir = string(in.String())
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain3(out *jwriter.Writer, in TryOnResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.UserID).MarshalText())
	}
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
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
	if in.TryOnResultID != "" {
		const prefix string = ",\"try_on_result_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TryOnResultID))
	}
	if in.TryOnResultDir != "" {
		const prefix string = ",\"try_on_result_dir\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TryOnResultDir))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TryOnResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TryOnResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TryOnResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TryOnResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain3(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain4(in *jlexer.Lexer, out *TryOnRequest) {
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
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "user_image_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserImageID).UnmarshalText(data))
			}
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "user_image_dir":
			out.UserImageDir = string(in.String())
		case "clothes_dir":
			out.ClothesDir = string(in.String())
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain4(out *jwriter.Writer, in TryOnRequest) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.UserID).MarshalText())
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
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ClothesID).MarshalText())
	}
	if in.UserImageDir != "" {
		const prefix string = ",\"user_image_dir\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.UserImageDir))
	}
	if in.ClothesDir != "" {
		const prefix string = ",\"clothes_dir\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ClothesDir))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TryOnRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TryOnRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TryOnRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TryOnRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain4(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain5(in *jlexer.Lexer, out *Tag) {
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
		case "name":
			out.Name = string(in.String())
		case "use_count":
			out.UseCount = int32(in.Int32())
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain5(out *jwriter.Writer, in Tag) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.UseCount != 0 {
		const prefix string = ",\"use_count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.UseCount))
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Tag) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Tag) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Tag) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Tag) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain5(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain6(in *jlexer.Lexer, out *Subtype) {
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
		case "name":
			out.Name = string(in.String())
		case "type_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.TypeID).UnmarshalText(data))
			}
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain6(out *jwriter.Writer, in Subtype) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if (in.TypeID).IsDefined() {
		const prefix string = ",\"type_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.TypeID).MarshalText())
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Subtype) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Subtype) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Subtype) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Subtype) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain6(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain7(in *jlexer.Lexer, out *Style) {
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
		case "name":
			out.Name = string(in.String())
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain7(out *jwriter.Writer, in Style) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if (in.ID).IsDefined() {
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Style) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Style) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Style) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Style) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain7(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain8(in *jlexer.Lexer, out *Credentials) {
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
		case "name":
			out.Name = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain8(out *jwriter.Writer, in Credentials) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Credentials) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Credentials) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Credentials) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Credentials) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain8(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain9(in *jlexer.Lexer, out *ClothesProcessingResponse) {
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
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "processed_dir":
			out.ProcessedDir = string(in.String())
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain9(out *jwriter.Writer, in ClothesProcessingResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.UserID).MarshalText())
	}
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ClothesID).MarshalText())
	}
	if in.ProcessedDir != "" {
		const prefix string = ",\"processed_dir\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ProcessedDir))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ClothesProcessingResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ClothesProcessingResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ClothesProcessingResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ClothesProcessingResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain9(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain10(in *jlexer.Lexer, out *ClothesProcessingRequest) {
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
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "clothes_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ClothesID).UnmarshalText(data))
			}
		case "clothes_dir":
			out.ClothesDir = string(in.String())
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain10(out *jwriter.Writer, in ClothesProcessingRequest) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.UserID).MarshalText())
	}
	if (in.ClothesID).IsDefined() {
		const prefix string = ",\"clothes_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.ClothesID).MarshalText())
	}
	if in.ClothesDir != "" {
		const prefix string = ",\"clothes_dir\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ClothesDir))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ClothesProcessingRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ClothesProcessingRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ClothesProcessingRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ClothesProcessingRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain10(l, v)
}
func easyjson3e1fa5ecDecodeTryOnInternalPkgDomain11(in *jlexer.Lexer, out *Clothes) {
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
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "name":
			out.Name = string(in.String())
		case "note":
			out.Note = string(in.String())
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]string, 0, 4)
					} else {
						out.Tags = []string{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Tags = append(out.Tags, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "image":
			out.Image = string(in.String())
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "style":
			out.Style = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "subtype":
			out.Subtype = string(in.String())
		case "color":
			out.Color = string(in.String())
		case "seasons":
			if in.IsNull() {
				in.Skip()
				out.Seasons = nil
			} else {
				in.Delim('[')
				if out.Seasons == nil {
					if !in.IsDelim(']') {
						out.Seasons = make([]Season, 0, 4)
					} else {
						out.Seasons = []Season{}
					}
				} else {
					out.Seasons = (out.Seasons)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Season
					v5 = Season(in.String())
					out.Seasons = append(out.Seasons, v5)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson3e1fa5ecEncodeTryOnInternalPkgDomain11(out *jwriter.Writer, in Clothes) {
	out.RawByte('{')
	first := true
	_ = first
	if (in.ID).IsDefined() {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Note != "" {
		const prefix string = ",\"note\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Note))
	}
	if len(in.Tags) != 0 {
		const prefix string = ",\"tags\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v6, v7 := range in.Tags {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	if in.Image != "" {
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	if (in.UserID).IsDefined() {
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserID).MarshalText())
	}
	if in.Style != "" {
		const prefix string = ",\"style\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Style))
	}
	if in.Type != "" {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	if in.Subtype != "" {
		const prefix string = ",\"subtype\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Subtype))
	}
	if in.Color != "" {
		const prefix string = ",\"color\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Color))
	}
	if len(in.Seasons) != 0 {
		const prefix string = ",\"seasons\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.Seasons {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Clothes) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Clothes) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e1fa5ecEncodeTryOnInternalPkgDomain11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Clothes) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Clothes) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e1fa5ecDecodeTryOnInternalPkgDomain11(l, v)
}
