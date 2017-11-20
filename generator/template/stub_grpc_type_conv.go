package template

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/generator/write_strategy"
	"github.com/devimteam/microgen/util"
	"github.com/vetcher/godecl/types"
)

const (
	GolangProtobufPtypesTimestamp = "github.com/golang/protobuf/ptypes/timestamp"
	JsonbPackage                  = "github.com/sas1024/gorm-jsonb/jsonb"
	GolangProtobufPtypes          = "github.com/golang/protobuf/ptypes"
)

type stubGRPCTypeConverterTemplate struct {
	Info                      *GenerationInfo
	alreadyRenderedConverters []string
	state                     WriteStrategyState
}

func NewStubGRPCTypeConverterTemplate(info *GenerationInfo) Template {
	return &stubGRPCTypeConverterTemplate{
		Info: info,
	}
}

func specialTypeConverter(p *types.Type) *Statement {
	// error -> string
	if p.Name == "error" && p.Import == nil {
		return (&Statement{}).Id("string")
	}
	// time.Time -> timestamp.Timestamp
	if p.Name == "Time" && p.Import != nil && p.Import.Package == "time" {
		return (&Statement{}).Op("*").Qual(GolangProtobufPtypesTimestamp, "Timestamp")
	}
	// jsonb.JSONB -> string
	if p.Name == "JSONB" && p.Import != nil && p.Import.Package == JsonbPackage {
		return (&Statement{}).Id("string")
	}
	return nil
}

func converterToProtoBody(field *types.Variable) Code {
	s := &Statement{}
	switch typeToProto(&field.Type, 0) {
	case "ErrorToProto":
		s.If().Id(util.ToLowerFirst(field.Name)).Op("==").Nil().Block(
			Return().List(Lit(""), Nil()),
		).Line()
		s.Return().List(Id(util.ToLowerFirst(field.Name)).Dot("Error").Call(), Nil())
	case "ByteListToProto":
		s.Return().List(Id(util.ToLowerFirst(field.Name)), Nil())
	case "TimeTimeToProto":
		s.Return().Qual(GolangProtobufPtypes, "TimestampProto").Call(Id(field.Name))
	default:
		s.Panic(Lit("function not provided"))
	}
	return s
}

func converterProtoToBody(field *types.Variable) Code {
	s := &Statement{}
	switch protoToType(&field.Type, 0) {
	case "ProtoToError":
		s.If().Id("proto" + util.ToUpperFirst(field.Name)).Op("==").Lit("").Block(
			Return().List(Nil(), Nil()),
		).Line()
		s.Return().List(Qual("errors", "New").Call(Id("proto"+util.ToUpperFirst(field.Name))), Nil())
	case "ProtoToByteList":
		s.Return().List(Id("proto"+util.ToUpperFirst(field.Name)), Nil())
	case "ProtoToTimeTime":
		s.Return().Qual(GolangProtobufPtypes, "Timestamp").Call(Id("proto" + util.ToUpperFirst(field.Name)))
	default:
		s.Panic(Lit("function not provided"))
	}
	return s
}

// Render whole file with protobuf converters.
//
//		// This file was automatically generated by "microgen" utility.
//		package protobuf
//
//		func IntListToProto(positions []int) (protoPositions []int64, convPositionsErr error) {
//			panic("method not provided")
//		}
//
//		func ProtoToIntList(protoPositions []int64) (positions []int, convPositionsErr error) {
//			panic("method not provided")
//		}
//
func (t *stubGRPCTypeConverterTemplate) Render() write_strategy.Renderer {
	f := &Statement{}

	for _, signature := range t.Info.Iface.Methods {
		args := append(removeContextIfFirst(signature.Args), removeErrorIfLast(signature.Results)...)
		for _, field := range args {
			if _, ok := golangTypeToProto("", &field); !ok && !util.IsInStringSlice(typeToProto(&field.Type, 0), t.alreadyRenderedConverters) {
				f.Line().Add(t.stubConverterToProto(&field)).Line()
				t.alreadyRenderedConverters = append(t.alreadyRenderedConverters, typeToProto(&field.Type, 0))
			}
			if _, ok := protoTypeToGolang("", &field); !ok && !util.IsInStringSlice(protoToType(&field.Type, 0), t.alreadyRenderedConverters) {
				f.Line().Add(t.stubConverterProtoTo(&field)).Line()
				t.alreadyRenderedConverters = append(t.alreadyRenderedConverters, protoToType(&field.Type, 0))
			}
		}
	}

	if t.state == AppendStrat {
		return f
	}

	file := NewFile("protobuf")
	file.PackageComment(FileHeader)
	file.PackageComment(`It is better for you if you do not change functions names!`)
	file.PackageComment(`This file will never be overwritten.`)
	file.Add(f)

	return file
}

func (stubGRPCTypeConverterTemplate) DefaultPath() string {
	return "./transport/converter/protobuf/type_converters.go"
}

func (t *stubGRPCTypeConverterTemplate) Prepare() error {
	if t.Info.ProtobufPackage == "" {
		return fmt.Errorf("protobuf package is empty")
	}
	return nil
}

func (t *stubGRPCTypeConverterTemplate) ChooseStrategy() (write_strategy.Strategy, error) {
	if err := util.StatFile(t.Info.AbsOutPath, t.DefaultPath()); os.IsNotExist(err) {
		t.state = FileStrat
		return write_strategy.NewCreateFileStrategy(t.Info.AbsOutPath, t.DefaultPath()), nil
	}
	file, err := util.ParseFile(filepath.Join(t.Info.AbsOutPath, t.DefaultPath()))
	if err != nil {
		return nil, err
	}

	for i := range file.Functions {
		t.alreadyRenderedConverters = append(t.alreadyRenderedConverters, file.Functions[i].Name)
	}

	t.state = AppendStrat
	return write_strategy.NewAppendToFileStrategy(t.Info.AbsOutPath, t.DefaultPath()), nil
}

// Render stub method for golang to protobuf converter.
//
//		func IntListToProto(positions []int) (protoPositions []int64, convPositionsErr error) {
//			return
//		}
//
func (t *stubGRPCTypeConverterTemplate) stubConverterToProto(field *types.Variable) *Statement {
	return Func().Id(typeToProto(&field.Type, 0)).
		Params(Id(util.ToLowerFirst(field.Name)).Add(fieldType(&field.Type))).
		Params(Add(t.protoFieldType(&field.Type)), Error()).
		Block(converterToProtoBody(field))
}

// Render stub method for protobuf to golang converter.
//
//		func ProtoToIntList(protoPositions []int64) (positions []int, convPositionsErr error) {
//			return
//		}
//
func (t *stubGRPCTypeConverterTemplate) stubConverterProtoTo(field *types.Variable) *Statement {
	return Func().Id(protoToType(&field.Type, 0)).
		Params(Id("proto"+util.ToUpperFirst(field.Name)).Add(t.protoFieldType(&field.Type))).
		Params(Add(fieldType(&field.Type)), Error()).
		Block(converterProtoToBody(field))
}

// Render protobuf field type for given func field.
//
//  	*repository.Visit
//
func (t *stubGRPCTypeConverterTemplate) protoFieldType(field *types.Type) *Statement {
	c := &Statement{}

	if field.IsArray {
		c.Index()
	}

	if field.IsPointer {
		c.Op("*")
	}

	if field.IsMap {
		m := field.Map
		return c.Map(t.protoFieldType(&m.Key)).Add(t.protoFieldType(&m.Value))
	}
	protoType := field.Name
	if tmp, ok := goToProtoTypesMap[field.Name]; ok {
		protoType = tmp
	}
	if code := specialTypeConverter(field); code != nil {
		return c.Add(code)
	}
	if field.Import != nil {
		c.Qual(t.Info.ProtobufPackage, protoType)
	} else {
		c.Id(protoType)
	}
	if field.IsInterface {
		c.Interface()
	}

	return c
}
