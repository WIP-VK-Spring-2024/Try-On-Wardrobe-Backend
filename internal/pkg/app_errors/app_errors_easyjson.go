// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package app_errors

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

func easyjson5d450e75DecodeTryOnInternalPkgAppErrors(in *jlexer.Lexer, out *ResponseError) {
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
		case "msg":
			out.Msg = string(in.String())
		case "errors":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Errors = make(map[string][]string)
				} else {
					out.Errors = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 []string
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if v1 == nil {
							if !in.IsDelim(']') {
								v1 = make([]string, 0, 4)
							} else {
								v1 = []string{}
							}
						} else {
							v1 = (v1)[:0]
						}
						for !in.IsDelim(']') {
							var v2 string
							v2 = string(in.String())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					(out.Errors)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
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
func easyjson5d450e75EncodeTryOnInternalPkgAppErrors(out *jwriter.Writer, in ResponseError) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Msg != "" {
		const prefix string = ",\"msg\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Msg))
	}
	if len(in.Errors) != 0 {
		const prefix string = ",\"errors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v3First := true
			for v3Name, v3Value := range in.Errors {
				if v3First {
					v3First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v3Name))
				out.RawByte(':')
				if v3Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v4, v5 := range v3Value {
						if v4 > 0 {
							out.RawByte(',')
						}
						out.String(string(v5))
					}
					out.RawByte(']')
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ResponseError) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5d450e75EncodeTryOnInternalPkgAppErrors(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ResponseError) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5d450e75EncodeTryOnInternalPkgAppErrors(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ResponseError) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5d450e75DecodeTryOnInternalPkgAppErrors(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ResponseError) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5d450e75DecodeTryOnInternalPkgAppErrors(l, v)
}
