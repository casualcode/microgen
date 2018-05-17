package template

import (
	"context"
	"path"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	mstrings "github.com/devimteam/microgen/generator/strings"
	"github.com/devimteam/microgen/generator/write_strategy"
	"github.com/devimteam/microgen/util"
	"github.com/vetcher/go-astra/types"
)

const (
	defaultHTTPMethod = "POST"

	HttpMethodTag  = "http-method"
	HttpMethodPath = "http-path"
)

type httpServerTemplate struct {
	info    *GenerationInfo
	methods map[string]string
	paths   map[string]string
}

func NewHttpServerTemplate(info *GenerationInfo) Template {
	return &httpServerTemplate{
		info: info,
	}
}

func (t *httpServerTemplate) DefaultPath() string {
	return filenameBuilder(PathTransport, "http", "server")
}

func (t *httpServerTemplate) ChooseStrategy(ctx context.Context) (write_strategy.Strategy, error) {
	return write_strategy.NewCreateFileStrategy(t.info.OutputFilePath, t.DefaultPath()), nil
}

func (t *httpServerTemplate) Prepare(ctx context.Context) error {
	t.methods = make(map[string]string)
	t.paths = make(map[string]string)
	for _, fn := range t.info.Iface.Methods {
		t.methods[fn.Name] = FetchHttpMethodTag(fn.Docs)
		t.paths[fn.Name] = buildMethodPath(fn)
	}
	return nil
}

func FetchHttpMethodTag(rawString []string) string {
	tags := util.FetchTags(rawString, TagMark+HttpMethodTag)
	if len(tags) == 1 {
		return strings.ToTitle(tags[0])
	}
	return defaultHTTPMethod
}

func buildMethodPath(fn *types.Function) string {
	url := strings.Replace(mstrings.FetchMetaInfo(TagMark+HttpMethodPath, fn.Docs), " ", "", -1)
	if url == "" {
		return buildDefaultMethodPath(fn)
	}
	return url
}

func buildDefaultMethodPath(fn *types.Function) string {
	edges := []string{util.ToURLSnakeCase(fn.Name)} // parts of full path
	if FetchHttpMethodTag(fn.Docs) == "GET" {
		edges = append(edges, gorillaMuxUrlTemplateVarList(RemoveContextIfFirst(fn.Args))...)
	}
	return path.Join(edges...)
}

func gorillaMuxUrlTemplateVarList(vars []types.Variable) []string {
	var list []string
	for i := range vars {
		list = append(list, "{"+util.ToURLSnakeCase(vars[i].Name)+"}")
	}
	return list
}

// Render http server constructor.
//		// This file was automatically generated by "microgen" utility.
//		// DO NOT EDIT.
//		package transporthttp
//
//		import (
//			svc "github.com/devimteam/microgen/examples/svc"
//			http2 "github.com/devimteam/microgen/examples/svc/transport/converter/http"
//			http "github.com/go-kit/kit/transport/http"
//			http1 "net/http"
//		)
//
//		func NewHTTPHandler(endpoints *svc.Endpoints, opts ...http.ServerOption) http1.Handler {
//			handler := http1.NewServeMux()
//			handler.Handle("/test_case", http.NewServer(
//				endpoints.TestCaseEndpoint,
//				http2.DecodeHTTPTestCaseRequest,
//				http2.EncodeHTTPTestCaseResponse,
//				opts...))
//			handler.Handle("/empty_req", http.NewServer(
//				endpoints.EmptyReqEndpoint,
//				http2.DecodeHTTPEmptyReqRequest,
//				http2.EncodeHTTPEmptyReqResponse,
//				opts...))
//			handler.Handle("/empty_resp", http.NewServer(
//				endpoints.EmptyRespEndpoint,
//				http2.DecodeHTTPEmptyRespRequest,
//				http2.EncodeHTTPEmptyRespResponse,
//				opts...))
//			return handler
//		}
//
func (t *httpServerTemplate) Render(ctx context.Context) write_strategy.Renderer {
	f := NewFile("transporthttp")
	f.ImportAlias(t.info.SourcePackageImport, serviceAlias)
	f.HeaderComment(t.info.FileHeader)

	f.Func().Id("NewHTTPHandler").ParamsFunc(func(p *Group) {
		p.Id("endpoints").Op("*").Qual(t.info.OutputPackageImport+"/transport", EndpointsSetName)
		if Tags(ctx).Has(TracingMiddlewareTag) {
			p.Id("logger").Qual(PackagePathGoKitLog, "Logger")
		}
		if Tags(ctx).Has(TracingMiddlewareTag) {
			p.Id("tracer").Qual(PackagePathOpenTracingGo, "Tracer")
		}
		p.Id("opts").Op("...").Qual(PackagePathGoKitTransportHTTP, "ServerOption")
	}).Params(
		Qual(PackagePathHttp, "Handler"),
	).BlockFunc(func(g *Group) {
		g.Id("mux").Op(":=").Qual(PackagePathGorillaMux, "NewRouter").Call()
		for _, fn := range t.info.Iface.Methods {
			g.Id("mux").Dot("Methods").Call(Lit(t.methods[fn.Name])).Dot("Path").
				Call(Lit("/" + t.paths[fn.Name])).Dot("Handler").Call(
				Line().Qual(PackagePathGoKitTransportHTTP, "NewServer").Call(
					Line().Id("endpoints").Dot(endpointStructName(fn.Name)),
					Line().Id(decodeRequestName(fn)),
					Line().Id(encodeResponseName(fn)),
					Line().Add(t.serverOpts(ctx, fn)).Op("...")),
			)
		}
		g.Return(Id("mux"))
	})

	return f
}

func (t *httpServerTemplate) serverOpts(ctx context.Context, fn *types.Function) *Statement {
	s := &Statement{}
	if Tags(ctx).Has(TracingMiddlewareTag) {
		s.Op("append(")
		defer s.Op(")")
	}
	s.Id("opts")
	if Tags(ctx).Has(TracingMiddlewareTag) {
		s.Op(",").Qual(PackagePathGoKitTransportHTTP, "ServerBefore").Call(
			Line().Qual(PackagePathGoKitTracing, "HTTPToContext").Call(Id("tracer"), Lit(fn.Name), Id("logger")),
		)
	}
	return s
}

func pathToHttpConverter(servicePath string) string {
	return filepath.Join(servicePath, "transport/converter/http")
}
