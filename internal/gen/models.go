package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"
)

type Render interface {
	Render(w io.Writer) error
}

const (
	base64Pkg = "encoding/base64"
	errorsPkg = "github.com/pkg/errors"
	fmtPkg    = "fmt"
	jsonPkg   = "encoding/json"
	logrusPkg = "github.com/sirupsen/logrus"
)

type Schema struct {
	Package   string
	Enums     Enums
	Variables Variables
}

func (s Schema) Render(w io.Writer) error {
	file := NewFile(s.Package)
	file.HeaderComment("This file is generated by the ./internal/gen command -- do not edit!")

	s.Enums.Render(file.Group)
	s.Variables.Render(file.Group)

	var buf bytes.Buffer
	err := file.Render(&buf)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}

type Enums []Enum

func (e Enums) Len() int {
	return len(e)
}

func (e Enums) Less(i, j int) bool {
	return strings.Compare(e[i].Name, e[j].Name) < 0
}

func (e Enums) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Enums) Render(g *Group) {
	sort.Stable(e)
	for _, e := range e {
		e.Render(g)
	}
}

type Enum struct {
	Name   string
	Values EnumValues
}

func (v Enum) Render(g *Group) {
	name := v.Name
	totalName := "total" + inflect.Pluralize(name)
	sliceName := inflect.CamelizeDownFirst(inflect.Pluralize(name))
	mapName := sliceName + "Map"
	newFuncName := "New" + name

	/*
		type AccessLevel uint8
	*/
	g.Type().Id(v.Name).Uint8()

	/*
		const (
			AccessLevelViewer AccessLevel = iota
			AccessLevelContributor
			AccessLevelAdmin
			totalAccessLevels
		)
	*/
	g.Const().
		DefsFunc(func(g *Group) {
			for idx, v := range v.Values {
				s := g.Id(name + v.Name)
				if idx == 0 {
					s.Id(name).Op("=").Iota()
				}
			}
			g.Id(totalName)
		})

	/*
		var (
			accessLevels = [totalAccessLevels]string{
				"viewer",
				"contributor",
				"admin",
			}

			accessLevelsMap = map[string]AccessLevel{
				"viewer":      AccessLevelViewer,
				"contributor": AccessLevelContributor,
				"admin":       AccessLevelAdmin,
			}
		)
	*/
	g.Var().
		DefsFunc(func(g *Group) {
			g.Id(sliceName).
				Op("=").
				Index(Id(totalName)).
				String().
				ValuesFunc(func(g *Group) {
					for _, v := range v.Values {
						g.Line().Lit(v.Value)
					}
					g.Line()
				})

			g.Line()

			g.Id(mapName).
				Op("=").
				Map(String()).
				Id(name).
				ValuesFunc(func(g *Group) {
					for _, v := range v.Values {
						g.Line().
							Lit(v.Value).
							Op(":").
							Id(name + v.Name)
					}
					g.Line()
				})
		})

	/*
		func NewAccessLevel(name string) (AccessLevel, error) {
			if v, ok := accessLevelsMap[name]; ok {
				return v, nil
			}

			return 0, fmt.Errorf("unknown AccessLevel name %q", name)
		}
	*/
	g.Func().
		Id(newFuncName).
		Params(Id("name").String()).
		Params(Id(name), Error()).
		Block(
			If(
				List(Id("v"), Id("ok")).Op(":=").Id(mapName).Index(Id("name")),
				Id("ok"),
			).Block(
				Return(Id("v"), Nil()),
			),
			Line(),
			Return(
				Lit(0),
				Qual(fmtPkg, "Errorf").
					Call(
						Lit(fmt.Sprintf("unknown %s name %%q", name)),
						Id("name"),
					),
			),
		).
		Line()

	/*
		func (v AccessLevel) String() string {
			if v < totalAccessLevels {
				return accessLevels[v]
			}

			return fmt.Sprintf("unknown AccessLevel value %02x", uint8(v))
		}
	*/
	g.Func().
		Params(Id("v").Id(name)).
		Id("String").
		Params().
		String().
		Block(
			If(Id("v").Op("<").Id(totalName)).
				Block(
					Return(Id(sliceName).Index(Id("v"))),
				),
			Line(),
			Return(
				Qual(fmtPkg, "Sprintf").
					Call(
						Lit(fmt.Sprintf("unknown %s value %%02x", name)),
						Uint8().Call(Id("v")),
					),
			),
		).
		Line()

	/*
		func (v *AccessLevel) UnmarshalText(text []byte) (err error) {
			*v, err = NewAccessLevel(string(text))
			return err
		}
	*/
	g.Func().
		Params(Id("v").Op("*").Id(name)).
		Id("UnmarshalText").
		Params(Id("text").Index().Byte()).
		Params(Err().Error()).
		Block(
			List(Op("*").Id("v"), Err()).Op("=").
				Id(newFuncName).Call(String().Call(Id("text"))),
			Return(Err()),
		).
		Line()

	/*
		func (v AccessLevel) MarshalText() ([]byte, error) {
			if v < totalAccessLevels {
				return []byte(accessLevels[v]), nil
			}

			return nil, errors.New(v.String())
		}
	*/
	g.Func().
		Params(Id("v").Id(name)).
		Id("MarshalText").
		Params().
		Params(Index().Byte(), Error()).
		Block(
			If(Id("v").Op("<").Id(totalName)).
				Block(
					Return(
						Index().Byte().Call(Id(sliceName).Index(Id("v"))),
						Nil(),
					),
				),
			Line(),
			Return(
				Nil(),
				Qual(errorsPkg, "New").Call(Id("v").Dot("String").Call()),
			),
		).Line()
}

type EnumValues []EnumValue

type EnumValue struct {
	Name  string
	Value string
}

type Variables []Variable

func (v Variables) Len() int {
	return len(v)
}

func (v Variables) Less(i, j int) bool {
	return strings.Compare(v[i].Name, v[j].Name) < 0
}

func (v Variables) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Variables) Render(g *Group) {
	sort.Sort(v)
	for _, v := range v {
		v.Render(g)
	}
}

type Variable struct {
	Name           string
	NoPrefix       bool
	Aliases        []string
	DecodedType    string
	DecodedPointer bool
}

func (v Variable) Render(g *Group) {
	lookupName := "Lookup" + strcase.ToCamel(v.Name)
	getName := "Get" + strcase.ToCamel(v.Name)

	rType := Null()
	if v.DecodedPointer {
		rType.Op("*")
	}
	if v.DecodedType == "" {
		rType.String()
	} else {
		rType.Id(v.DecodedType)
	}

	receiver := Id("e").Op("*").Id("Environment")

	/*
		func LookupApplication(p PlatformProvider) (*Application, bool) {
			name := p.Prefix() + "APPLICATION"
			value, ok := p.Lookup(name)
			if !ok {
				return nil, false
			}
			data, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				logrus.WithError(err).Warn("unable to decode value")
				return nil, false
			}
			obj := Application{}
			err = json.Unmarshal(data, &obj)
			if err != nil {
				logrus.WithError(err).Warn("unable to unmarshal value")
				return nil, false
			}
			return &app, true
		}
	*/
	g.Func().
		Id(lookupName).
		Params(Id("p").Id("PlatformProvider")).
		Params(rType, Bool()).
		BlockFunc(func(g *Group) {
			s := g.Id("name").Op(":=")
			if !v.NoPrefix {
				s.Id("p").Dot("Prefix").Call().Op("+")
			}
			s.Lit(strcase.ToScreamingSnake(v.Name))

			g.List(Id("value"), Id("ok")).Op(":=").Id("p").Dot("Lookup").Call(Id("name"))

			if v.DecodedType == "" {
				g.Return(Id("value"), Id("ok"))
				return
			}

			g.If(Op("!").Id("ok")).Block(Return(Nil(), False()))

			g.List(Id("data"), Err()).
				Op(":=").
				Qual(base64Pkg, "StdEncoding").
				Dot("DecodeString").
				Call(Id("value"))

			errNotNil(g, "unable to decode value")

			g.Id("obj").Op(":=").Id(v.DecodedType).Values()

			g.Err().Op("=").Qual(jsonPkg, "Unmarshal").Call(Id("data"), Op("&").Id("obj"))

			errNotNil(g, "unable to unmarshal value")

			val := Null()
			if v.DecodedPointer {
				val.Op("&")
			}
			val.Id("obj")

			g.Return(val, True())
		}).
		Line()

	/*
		func GetApplication(p PlatformProvider) *Application {
			v, _ := LookupApplication(p)
			return v
		}
	*/
	g.Func().
		Id(getName).
		Params(Id("p").Id("PlatformProvider")).
		Add(rType).
		Block(
			List(Id("v"), Id("_")).Op(":=").Id(lookupName).Call(Id("p")),
			Return(Id("v")),
		).
		Line()

	/*
		func (e *Environment) LookupApplication() (*Application, bool) {
			return LookupApplication(e)
		}
	*/
	g.Func().
		Params(receiver).
		Id(lookupName).
		Params().
		Params(rType, Bool()).
		Block(
			Return(Id(lookupName).Call(Id("e"))),
		).
		Line()

	/*
		func (e *Environment) GetApplication() *Application {
			return GetApplication(e)
		}
	*/
	g.Func().
		Params(receiver).
		Id(getName).
		Params().
		Add(rType).
		Block(
			Return(Id(getName).Call(Id("e"))),
		).
		Line()

	for _, a := range v.Aliases {
		lookupAlias := "Lookup" + strcase.ToCamel(a)
		getAlias := "Get" + strcase.ToCamel(a)

		/*
			func LookupApp(p PlatformProvider) (*Application, bool) {
				return LookupApplication(p)
			}
		*/
		g.Func().
			Id(lookupAlias).
			Params(Id("p").Id("PlatformProvider")).
			Params(rType, Bool()).
			Block(
				Return(Id(lookupName).Call(Id("p"))),
			).
			Line()

		/*
			func GetApp(p PlatformProvider) *Application {
				return GetApplication(p)
			}
		*/
		g.Func().
			Id(getAlias).
			Params(Id("p").Id("PlatformProvider")).
			Add(rType).
			Block(
				Return(Id(getName).Call(Id("p"))),
			).
			Line()

		/*
			func (e *Environment) LookupApp() (*Application, bool) {
				return LookupApp(e)
			}
		*/
		g.Func().
			Params(receiver).
			Id(lookupAlias).
			Params().
			Params(rType, Bool()).
			Block(
				Return(Id(lookupAlias).Call(Id("e"))),
			).
			Line()

		/*
			func (e *Environment) GetApp() *Application {
				return GetApp(e)
			}
		*/
		g.Func().
			Params(receiver).
			Id(getAlias).
			Params().
			Add(rType).
			Block(
				Return(Id(getAlias).Call(Id("e"))),
			).
			Line()
	}
}

func errNotNil(g *Group, msg string) {
	g.If(Err().Op("!=").Nil()).
		Block(
			Qual(logrusPkg, "WithError").
				Call(Err()).
				Dot("Warn").
				Call(Lit(msg)),
			Return(Nil(), False()),
		)
}
