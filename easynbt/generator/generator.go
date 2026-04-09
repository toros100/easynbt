package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/toros100/easynbt/easynbt/util"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

const nbtPackagePath = "github.com/toros100/easynbt/nbt"
const unmarshallerInterfaceName = "Unmarshaller"
const optionTypeName = "Option"

const debugAllowBrokenOutput = true

type Generator struct {
	buf                     *bytes.Buffer
	nbtPkg                  *packages.Package
	targetPkg               *packages.Package
	unmarshallerInterface   *types.Interface
	unmarshallerMethodNames []string
	optionType              *types.Named
	targetTypes             *util.TypeSet
	opts                    *Options
}

type Options struct {
	PackagePattern string // TODO:
	OutFile        string
	DryRun         bool
	Verbose        bool
	BuildFlags     []string // TODO:
}

type FieldInfo struct {
	FieldName string
	NbtName   string
	Optional  bool
	FieldType types.Type
}

type StructInfo struct {
	Fields []FieldInfo
}

func New() *Generator {
	return &Generator{
		buf: bytes.NewBuffer(nil),
	}
}

// loadPackage loads the package that code will be generated for, as well as
// github.com/toros100/easynbt/nbt, which will be a dependency of the generated code
// and has to be available included in the local go.mod.
// Invoking the main program via "//go:generate easynbt [...]" in some .go-file
// sets the working directory to the directory containing that file, loadPackage
// is usually called with "." as the pattern. (cf. documentation of packages.Load)
func (g *Generator) loadPackage(pattern string) error {

	pkgs, err := packages.Load(&packages.Config{
		BuildFlags: g.opts.BuildFlags,
	}, pattern)

	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return errors.New("failed to load package")
	}

	targetPkgName := pkgs[0].Name

	if g.opts.OutFile == "" {
		g.opts.OutFile = filepath.Join(".", targetPkgName+"_nbt_gen.go")
	}

	overlay := make(map[string][]byte)

	outFileAbs, err := filepath.Abs(g.opts.OutFile)
	if err != nil {
		return err
	}

	// essentially excluding the target file from the loaded package
	// the purpose of this is that we will be checking if a particular named type
	// already satisfies nbt.Unmarshaller, but for that we need to ignore the intended output file
	// because it will be overwritten
	overlay[outFileAbs] = fmt.Appendf(nil, "package %s\n", targetPkgName)

	cfg := &packages.Config{
		// todo: can probably load less
		Mode:       packages.LoadAllSyntax,
		Overlay:    overlay,
		BuildFlags: g.opts.BuildFlags,
	}

	pkgs, err = packages.Load(cfg, pattern, nbtPackagePath)
	if err != nil {
		return err
	}

	if len(pkgs) != 2 {
		return errors.New("failed to load packages")
	}

	// need to figure out which of the two loaded pkgs is the nbt package and which is the target
	// they have to be loaded in the same call of packages.Load s.t. types.Identical and types.Implements
	// work correctly
	//
	// this is a bit awkward but i do not know if there is a better way to do this?
	n, t := pkgs[0], pkgs[1]
	if n.PkgPath != nbtPackagePath {
		if t.PkgPath != nbtPackagePath {
			panic("internal error: failed to load nbt package")
		}
		g.nbtPkg = t
		g.targetPkg = n
	} else {
		g.nbtPkg = n
		g.targetPkg = t
	}

	g.checkErrs(g.targetPkg.PkgPath, g.targetPkg.Errors)
	g.checkErrs(g.nbtPkg.PkgPath, g.nbtPkg.Errors)

	if g.targetPkg.Types == nil {
		g.nbtPkg = nil
		g.targetPkg = nil
		return errors.New("failed to load types from target package")
	}

	if g.nbtPkg.Types == nil {
		panic("internal error: failed to load types from nbt package")
	}

	res := g.nbtPkg.Types.Scope().Lookup(unmarshallerInterfaceName)
	if res == nil {
		panic("internal error: failed to lookup " + unmarshallerInterfaceName)
	}
	if unmarshallerInterface, ok := res.Type().Underlying().(*types.Interface); ok {
		g.unmarshallerInterface = unmarshallerInterface
	} else {
		panic("internal error: " + unmarshallerInterfaceName + " is not an interface")
	}

	methodNames := make([]string, 0, g.unmarshallerInterface.NumMethods())
	for m := range g.unmarshallerInterface.Methods() {
		methodNames = append(methodNames, m.Name())
	}

	g.unmarshallerMethodNames = methodNames

	res = g.nbtPkg.Types.Scope().Lookup(optionTypeName)
	if res == nil {
		panic("internal error: failed to lookup " + optionTypeName)
	}

	if optionType, ok := res.Type().(*types.Named); ok {
		g.optionType = optionType
	} else {
		panic("internal error: " + optionTypeName + " is not a named type")
	}

	return nil
}

func isFatal(err packages.Error) bool {
	// TODO: not considering type errors fatal right now, because type errors can easily be caused by
	// existing code that already uses some of the yet to be generated methods.
	// these methods might even exist in the current code base, in the old target file, which is
	// overlaid with an empty file for the generator.
	// i dont think there is any realistic way of distinguishing which type errors should be fatal here.
	// maybe the best thing to do would be to do another call to packages.Load, where we overlay
	// the target file with a file containing stubs of the to be generated methods. in this case, there
	// should be no errors? but even then, there could be unused imports (which have kind packages.TypeError here,
	// but clearly wont harm the codegen)
	return err.Kind != packages.TypeError
}

func (g *Generator) checkErrs(pkgPath string, errs []packages.Error) {

	for _, e := range errs {
		if isFatal(e) {
			log.Fatalf("encountered fatal error while loading package %s:\n\t%s\n", pkgPath, e.Error())
		}
	}

	if g.opts.Verbose {
		preamble := "encountered the following (not necessarily fatal) errors while loading package %s:\n"
		preambleOnce := false
		for _, e := range errs {
			if !isFatal(e) {
				if !preambleOnce {
					fmt.Fprintf(os.Stderr, preamble, pkgPath)
					preambleOnce = true
				}
				fmt.Fprintf(os.Stderr, "\t%s\n", e.Error())
			}
		}
	}
}

func (g *Generator) checkInterfaceMethodsCollision(typ *types.Named) error {
	if len(g.unmarshallerMethodNames) == 0 {
		panic("internal error")
	}

	for m := range typ.Methods() {
		if slices.Contains(g.unmarshallerMethodNames, m.Name()) {
			return fmt.Errorf("%w: type %s already has a method %q", ErrMethodCollision, util.TypeString(typ), m.Name())
		}
	}
	return nil
}

func (g *Generator) typeStringRelative(typ types.Type) string {
	if g.targetPkg == nil || g.targetPkg.PkgPath == "" {
		panic("internal error")
	}

	return util.TypeStringRelative(typ, g.targetPkg.PkgPath)
}

func (g *Generator) loadTypes(typeNames []string) error {
	if len(typeNames) == 0 {
		return fmt.Errorf("%w (no types specified)", ErrInvalidInput)
	}
	if g.targetPkg == nil {
		panic("internal error: no package loaded")
	}

	scope := g.targetPkg.Types.Scope()
	// targetTypes := make([]*types.Named, len(typeNames))
	g.targetTypes = util.NewTypeSet()

	for _, tn := range typeNames {

		res := scope.Lookup(tn)
		if res == nil {
			return fmt.Errorf("%q: %w", tn, ErrTypeNotFound)
		}

		if named, ok := types.Unalias(res.Type()).(*types.Named); !ok {
			return fmt.Errorf("%w: %s (not a named type)", ErrInvalidInput, res.Name())
		} else {
			if named.Obj().Pkg().Path() != g.targetPkg.PkgPath {
				return fmt.Errorf("%w: %s (not from target package)", ErrUnexpectedType, util.TypeString(res.Type()))
			}
			if named.TypeParams() != nil {
				return fmt.Errorf("%w: %s (generic)", ErrUnexpectedType, util.TypeString(res.Type()))
			}
			g.targetTypes.Add(named)
		}
	}
	return nil
}

func getReceiverName(typ *types.Named) string {
	for m := range typ.Methods() {
		if m := m.Signature().Recv(); m != nil && m.Name() != "" && m.Name() != "_" {
			// TODO: technically this could collide with variable names in the generated code?
			// but that would be very annoying to prevent in a way that is actually robust
			return m.Name()
		}
	}
	return strings.ToLower(typ.Obj().Name())[:1]
}

// Generate generates the nbt.Unmarshaller interfaces methods for the specified types, if possible.
// Invoking the main program via "//go:generate easynbt [...]" in some .go-file
// sets the working directory to the directory containing that file, thus Generate
// is usually called with "." as the pattern. (cf. documentation of packages.Load)
func (g *Generator) Generate(opts *Options, packagePattern string, typeNames []string) error {

	if opts == nil {
		// zero values are sane defaults except for output file name
		// but we cant set that yet, need target package name
		g.opts = new(Options)
	} else {
		g.opts = opts
	}

	if err := g.loadPackage(packagePattern); err != nil {
		return err
	}

	if g.opts.Verbose {
		fmt.Println("loaded package", g.targetPkg.PkgPath)
	}

	if err := g.loadTypes(typeNames); err != nil {
		return err
	}

	if g.targetPkg == nil {
		return errors.New("no package loaded")
	}

	fmt.Fprintf(g.buf, "package %s\n", g.targetPkg.Name)

	g.writeImports()

	fmt.Fprintln(g.buf)
	fmt.Fprintln(g.buf, "var (")

	for t := range g.targetTypes.Values() {
		named, ok := t.(*types.Named)
		if !ok {
			// should only have named types from the target package in g.targetTypes
			panic("programmer error")
		}
		if err := g.checkInterfaceMethodsCollision(named); err != nil {
			return err
		}

		// essentially assertions to ensure the code fails to build if the types dont implement the
		// nbt.Unmarshaller interface as promised
		fmt.Fprintf(g.buf, "_ nbt.Unmarshaller = (*%s)(nil)\n", named.Obj().Name())
	}

	fmt.Fprintf(g.buf, ")\n")

	for tn := range g.targetTypes.Values() {
		named := tn.(*types.Named)

		tagType, err := g.getExpectedTagExpression(named.Underlying())
		if err != nil {
			return fmt.Errorf("%w: %s", err, util.TypeString(tn))
		}

		if !strings.HasPrefix(tagType, "nbt.") {
			// must be a known constant from the nbt package at this point
			panic("programmer error")
		}

		receiverVar := getReceiverName(named)

		fmt.Fprintf(g.buf, `func (%s *%s) UnmarshalPayload(data []byte) (int, error) {
			off := 0
			`, receiverVar, named.Obj().Name())

		err = g.writeUnmarshal(receiverVar, named.Underlying(), 0)

		if err != nil {
			return fmt.Errorf("%w: %s", err, util.TypeString(tn))
		}

		fmt.Fprint(g.buf, `
		return off, nil
		}
		`)

		fmt.Fprintf(g.buf, `
		func (%s *%s) TagType() byte {
			return %s
		}
		

		`, receiverVar, named.Obj().Name(), tagType)
	}

	outBytes, err := imports.Process(g.opts.OutFile, g.buf.Bytes(), nil)

	if err != nil {
		fmtErr := fmt.Errorf("failed to format generated code: %w", err)
		if debugAllowBrokenOutput {
			fmt.Fprintln(os.Stderr, fmtErr.Error())
			outBytes = g.buf.Bytes()
		} else {
			return fmtErr
		}
	}

	if !g.opts.DryRun {
		err = os.WriteFile(g.opts.OutFile, outBytes, 0o644)
		if err != nil {
			return fmt.Errorf("failed to write generated code to file: %w", err)
		}
	}

	return nil
}

func (g *Generator) writeImports() {
	// TODO: this is just grabbing ALL imports from the package and writing them into the buffer as import statements.
	// later, the goimports tool/package is used to remove unused imports.
	// since this just grabs import on a per-package basis, all files, this could theoretically break in several ways:
	//  1. it could introduce collisions, e.g. if one file in the package imports compress/zlib and another file imports
	//     github.com/klauspost/compress/zlib, then one would need to be aliased to have them in the same file.
	//  2. this doesn't do anything to handle imports that are just done for side-effects (_ alias)
	// ---
	// should look into existing codegen stuff like dave/jennifer, which might also make my codegen less unhinged in general

	if g.targetPkg == nil {
		panic("internal error (writeImports with no loaded package)")
	}

	fmt.Fprintln(g.buf, "import (")

	for k := range g.targetPkg.Imports {
		fmt.Fprintf(g.buf, "%q\n", k)
	}
	fmt.Fprintf(g.buf, "%q\n", nbtPackagePath)
	fmt.Fprintln(g.buf, ")")
}

func (g *Generator) isOptionOf(t types.Type) (types.Type, bool) {
	if named, ok := types.Unalias(t).(*types.Named); ok && types.Identical(named.Origin(), g.optionType) {

		str, ok := named.Underlying().(*types.Struct)
		if !ok || str.NumFields() != 3 {
			panic("nbt.Option type has unexpected shape (programmer error)")
		}
		return str.Field(1).Type(), true
	}
	return nil, false
}

func (g *Generator) isTargetType(t types.Type) bool {
	return g.targetTypes.Contains(t)
}

func (g *Generator) implementsUnmarshaller(t *types.Pointer) bool {
	return types.Implements(t, g.unmarshallerInterface)
}

func (g *Generator) writeUnmarshal(into string, typ types.Type, level uint) error {

	switch t := typ.(type) {
	case *types.Alias:
		u := types.Unalias(t)
		return g.writeUnmarshal(into, u, level)
	case *types.Basic:
		return g.writeUnmarshalBasic(into, t, level)
	case *types.Struct:
		return g.writeUnmarshalStruct(into, t, level)
	case *types.Named:
		return g.writeUnmarshalNamed(into, t, level)
	case *types.Pointer:
		return g.writeUnmarshalPointer(into, t)
	case *types.Slice:
		return g.writeUnmarshalSlice(into, t, level)
	default:
		// future work?
		// maybe array types for lists of known length, generics, more flexibility on pointers
		return fmt.Errorf("%w: %s", ErrUnexpectedType, util.TypeString(typ))
	}
}

func (g *Generator) writeUnmarshalBasic(into string, typ *types.Basic, level uint) error {

	// a bit hacky:
	// helper methods in nbt package like
	// nbt.BytePayloadFromBytes take a pointer argument due to restrictions on
	// acceptable input types, into is already a pointer if and only if this
	// is writing the unmarshaller for a named type like "type MyByte int8"

	var intoExpr string
	if level == 0 {
		intoExpr = into
	} else {
		intoExpr = "&" + into
	}

	// my reasoning for calling helper functions from nbt instead of just inlining the code:
	// 1. i get handling of alias and type names for free with ~T type parameters on the generic helpers
	// 2. might add features like unsafe or interned strings, which could be enabled with build tags
	//		without even having to touch the generated code at all with this approach
	// 3. the functions helpers are deliberately small (length checks done externally) such that the compiler
	//		will inline them for me, so there should be no runtime cost for the above conveniences

	switch typ.Kind() {
	case types.Int8:
		fmt.Fprintf(g.buf, `
			if off+1 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.BytePayloadFromBytes(%s, data[off:])
			off += 1

			`, intoExpr)

	case types.Int16:
		fmt.Fprintf(g.buf, `
			if off+2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.ShortPayloadFromBytes(%s, data[off:])
			off += 2

			`, intoExpr)

	case types.Int32:
		fmt.Fprintf(g.buf, `
			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.IntPayloadFromBytes(%s, data[off:])
			off += 4

			`, intoExpr)

	case types.Int64:
		fmt.Fprintf(g.buf, `
			if off+8 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.LongPayloadFromBytes(%s, data[off:])
			off += 8

			`, intoExpr)

	case types.Float32:
		fmt.Fprintf(g.buf, `
			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.FloatPayloadFromBytes(%s, data[off:])
			off += 4

			`, intoExpr)

	case types.Float64:
		fmt.Fprintf(g.buf, `
			if off+8 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.DoublePayloadFromBytes(%s, data[off:])
			off += 8

			`, intoExpr)

	case types.String:
		fmt.Fprintf(g.buf, `
			if off+2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			strLen := int(binary.BigEndian.Uint16(data[off:]))
			off += 2

			if off+strLen > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			if !utf8.Valid(data[off:off+strLen]) {
				return 0, nbt.ErrInvalidUTF8
			}

			nbt.StringFromBytes(%s, data[off:off+strLen])
			off += strLen

			`, intoExpr)

	default:
		return fmt.Errorf("%w: %s", ErrUnexpectedType, util.TypeString(typ))
	}
	return nil
}

func (g *Generator) processStruct(t *types.Struct) (*StructInfo, error) {
	tInfo := new(StructInfo)
	numRelevantFields := 0
	fieldNames := make(map[string]struct{})

	for i := range t.NumFields() {
		f := t.Field(i)
		if f.Embedded() {
			return nil, fmt.Errorf("%w (embedded field)", ErrUnexpectedType)
		}

		if f.Name() == "_" {
			// ignoring blank fields, there are potentially good reasons for having them
			// (alignment, noCopy, meta-information, ...)
			continue
		}

		tag := reflect.StructTag(t.Tag(i))

		if _, ignore := tag.Lookup("nbtignore"); ignore {
			continue
		}

		numRelevantFields++

		nbtName, ok := tag.Lookup("nbt")

		if ok {
			if !utf8.ValidString(nbtName) {
				// you would have to try pretty hard to smuggle non-UTF-8 into the name
				// but it is possible, e.g.
				// type Data struct {
				// 		someField int32 `nbt:"FIELD\xff"`
				// }
				return nil, fmt.Errorf("%w (field name not valid UTF-8)", ErrFieldNames)
			}
		} else {
			// using field name if no name is given via struct tag
			// this will always be UTF-8
			nbtName = f.Name()
		}

		_, hasOptionalTag := tag.Lookup("nbtoptional")
		_, isOptionType := g.isOptionOf(f.Type())

		fieldNames[nbtName] = struct{}{}

		tInfo.Fields = append(tInfo.Fields, FieldInfo{
			FieldName: f.Name(),
			NbtName:   nbtName,
			Optional:  hasOptionalTag || isOptionType,
			FieldType: f.Type(),
		})
	}

	if len(fieldNames) != numRelevantFields {
		return nil, fmt.Errorf("%w (field name collision)", ErrFieldNames)
	}

	return tInfo, nil
}

func (g *Generator) writeUnmarshalStruct(into string, typ *types.Struct, level uint) error {

	info, err := g.processStruct(typ)

	if err != nil {
		return err
	}

	foundFieldsVar := fmt.Sprintf("foundFields%d", level)
	tagVar := fmt.Sprintf("tag%d", level)

	anyOptional := false
	allOptional := true
	requiredFieldsLiteralBuilder := strings.Builder{}
	fieldNamesLiteralBuilder := strings.Builder{}
	fmt.Fprintf(&requiredFieldsLiteralBuilder, "[%d]bool{", len(info.Fields))
	fmt.Fprintf(&fieldNamesLiteralBuilder, "[%d]string{", len(info.Fields))
	fieldNameStrs := make([]string, len(info.Fields))

	for i, f := range info.Fields {
		if f.Optional {
			anyOptional = true
		} else {
			allOptional = false
		}
		requiredFieldsLiteralBuilder.WriteString(strconv.FormatBool(!f.Optional))

		if f.FieldName == f.NbtName {
			fieldNameStrs[i] = f.FieldName
		} else {
			fieldNameStrs[i] = fmt.Sprintf("%s (nbt: %s)", f.FieldName, f.NbtName)
		}

		fieldNamesLiteralBuilder.WriteString("\"" + fieldNameStrs[i] + "\"")

		if i < len(info.Fields) {
			requiredFieldsLiteralBuilder.WriteString(",")
			fieldNamesLiteralBuilder.WriteString(",")
		}
	}
	requiredFieldsLiteralBuilder.WriteString("}")
	fieldNamesLiteralBuilder.WriteString("}")

	fmt.Fprintf(g.buf, `
		%s := [%d]bool{}
		`, foundFieldsVar, len(info.Fields))

	fmt.Fprintf(g.buf, `
		for {
			if off >= len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			%s := data[off]
			off += 1

			if %s == nbt.TagEnd {
				break
			}

			if off + 2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			strLen := int(binary.BigEndian.Uint16(data[off:]))
			off += 2

			if off + strLen > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			strData := data[off:off+strLen]
			off += strLen

			switch string(strData) {
		`, tagVar, tagVar)

	for i, field := range info.Fields {

		fmt.Fprintf(g.buf, `
		case %q:
			if %s[%d] {
				return 0, fmt.Errorf("on field %s: %%w", nbt.ErrDuplicateValue)
			}
			`, field.NbtName, foundFieldsVar, i, fieldNameStrs[i])

		nextInto := fmt.Sprintf("%s.%s", into, field.FieldName)
		tagExpr, err := g.getExpectedTagExpression(field.FieldType)

		if err != nil {
			return fmt.Errorf("field %s: %w", field.FieldName, err)
		}

		fmt.Fprintf(g.buf, `
		if %s != %s {
			return 0, nbt.ErrUnexpectedTag
		}
		`, tagVar, tagExpr)

		if err := g.writeUnmarshal(nextInto, field.FieldType, level+1); err != nil {
			return fmt.Errorf("field %s: %w", field.FieldName, err)
		}

		fmt.Fprintf(g.buf, "%s[%d] = true\n\n", foundFieldsVar, i)
	}

	fmt.Fprintf(g.buf, `
			default:
				k, err := nbt.SkipPayload(%s, data[off:])
				if err != nil {
					return 0, err
				}
				off += k
			}
		}
		`, tagVar)

	switch {
	case anyOptional && !allOptional:
		requiredFieldsVar := fmt.Sprintf("required%d", level)

		fmt.Fprintf(g.buf, `
		%s := %s
		for i := range %s {
			if !%s[i] && %s[i] {
				return 0, fmt.Errorf("on field %%s: %%w", %s[i], nbt.ErrMissingValue)
			}
		}
		`,
			requiredFieldsVar,
			requiredFieldsLiteralBuilder.String(),
			foundFieldsVar,
			foundFieldsVar,
			requiredFieldsVar,
			fieldNamesLiteralBuilder.String(),
		)

	case !anyOptional:
		fmt.Fprintf(g.buf, `
		for i := range %s {
			if !%s[i] {
				return 0, fmt.Errorf("on field %%s: %%w", %s[i], nbt.ErrMissingValue)
			}
		}
		`,
			foundFieldsVar,
			foundFieldsVar,
			fieldNamesLiteralBuilder.String(),
		)
	}

	return nil
}

func (g *Generator) writeUnmarshalNamed(into string, t *types.Named, level uint) error {

	if elem, ok := g.isOptionOf(t); ok {
		// special-cased behavior for the builtin nbt.Option type
		// given the restrictions on allowed input types this can only occur as a struct field

		if err := g.writeUnmarshal(into+".Value", elem, level); err != nil {
			return fmt.Errorf("%w: %s", err, util.TypeString(t))
		}

		fmt.Fprintf(g.buf, `
		%s.Ok = true
		`, into)
		return nil
	}

	if !g.isTargetType(t) && !g.implementsUnmarshaller(types.NewPointer(t)) {
		// todo: maybe add a flag --force or similar to allow this error to be ignored and just generate
		// the code calling interface methods. i could imagine that sometimes you might want to do that in dev,
		// even though it produces code that wont even compile
		return fmt.Errorf("%w (does not implement nbt.Unmarshaller)", ErrUnexpectedType)
	}

	fmt.Fprintf(g.buf, `
			k, err := %s.UnmarshalPayload(data[off:])
			if err != nil {
				return 0, err
			}
			off += k
			`, into)
	return nil
}

func (g *Generator) writeUnmarshalPointer(into string, typ *types.Pointer) error {
	// the only allowed pointer types are pointers to named types that are either
	// targets of the current codegen, or already implement nbt.Unmarshaller

	named, ok := types.Unalias(typ.Elem()).(*types.Named)
	if !ok {
		return ErrUnexpectedType
	}

	if _, ok := g.isOptionOf(named); ok {
		// in particular, pointers of nbt.Option are not allowed
		return fmt.Errorf("%w (pointer to nbt.Option)", ErrUnexpectedType)
	}

	if !g.isTargetType(named) && !g.implementsUnmarshaller(typ) {
		return fmt.Errorf("%w (%s does not implement nbt.Unmarshaller)", ErrUnexpectedType, util.TypeString(typ))
	}

	fmt.Fprintf(g.buf, `
	%s = new(%s)

	k, err := %s.UnmarshalPayload(data[off:])
	if err != nil {
		return 0, err
	}
	off += k
	`, into, g.typeStringRelative(typ.Elem()), into)

	return nil
}

func (g *Generator) writeUnmarshalSlice(into string, t *types.Slice, level uint) error {
	elemTyp := t.Elem()
	listLenVar := fmt.Sprintf("listLen%d", level)
	listVar := fmt.Sprintf("list%d", level)
	idxVar := fmt.Sprintf("i%d", level)
	elemTagVar := fmt.Sprintf("elemTag%d", level)

	fmt.Fprintf(g.buf, `
		if off + 5 > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		%s := data[off]
		off += 1
		`, elemTagVar)

	var intoExpr string
	// a bit hacky but due to the restrictions imposed on input types,
	// the only possibility for the expression we unmarshal into to be a pointer that
	// needs to be dereferenced before assignments is on the top-level call
	if level == 0 {
		intoExpr = "*" + into
	} else {
		intoExpr = into
	}

	fmt.Fprintf(g.buf, `
		%s := binary.BigEndian.Uint32(data[off:])
		off += 4

		if %s > math.MaxInt32 {
			return 0, nbt.ErrInvalidLength
		}
		`, listLenVar, listLenVar)

	var escape string
	if level == 0 {
		escape = "return off, nil"
	} else {
		escape = "continue"
	}

	tagExpr, err := g.getExpectedTagExpression(elemTyp)

	if err != nil {
		return fmt.Errorf("slice element type: %w", err)
	}

	// apparently it's a bit inconsistent which element type a nbt list of length 0 is supposed to have?
	// i am allowing either the tag implied by the go type or nbt.TagEnd
	// TODO: look into this more, maybe do a build tag for "strict lists" if i can figure out
	// what the "official" correct behaviour is supposed to be (e.g. if it ALWAYS should be nbt.TagEnd for empty lists)
	fmt.Fprintf(g.buf, `
		if %s != %s && !(%s == 0 && %s == nbt.TagEnd) {
			return 0, nbt.ErrUnexpectedTag
		}
		if %s == 0 {
			%s
		}
		`, elemTagVar, tagExpr, listLenVar, elemTagVar, listLenVar, escape)

	// pain and suffering

	fmt.Fprintf(g.buf, `
		if %s == nil {
			%s = make(%s, int(%s))
		} else {
			%s = slices.Grow(%s, int(%s))[:int(%s)]
		}
		%s := %s
		`,
		intoExpr,
		intoExpr,
		g.typeStringRelative(t),
		listLenVar,
		intoExpr,
		intoExpr,
		listLenVar,
		listLenVar,
		listVar,
		intoExpr,
	)

	fmt.Fprintf(g.buf, `
		for %s := range %s {
		`, idxVar, listVar)

	nextInto := fmt.Sprintf("%s[%s]", listVar, idxVar)

	err = g.writeUnmarshal(nextInto, elemTyp, level+1)
	if err != nil {
		return fmt.Errorf("slice element type: %w", err)
	}

	fmt.Fprintln(g.buf, "}")
	return nil
}

func (g *Generator) getExpectedTagExpression(typ types.Type) (string, error) {
	// a lot of error checking is "hidden" in this method
	// i don't really like this, might refactor, but unclear how
	switch t := types.Unalias(typ).(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int8:
			return "nbt.TagByte", nil
		case types.Int16:
			return "nbt.TagShort", nil
		case types.Int32:
			return "nbt.TagInt", nil
		case types.Int64:
			return "nbt.TagLong", nil
		case types.Float32:
			return "nbt.TagFloat", nil
		case types.Float64:
			return "nbt.TagDouble", nil
		case types.String:
			return "nbt.TagString", nil
		default:
			return "", fmt.Errorf("%w (unexpected basic type %s)", ErrUnexpectedType, util.TypeString(t))
		}
	case *types.Struct:
		return "nbt.TagCompound", nil
	case *types.Slice:
		return "nbt.TagList", nil
	case *types.Named:
		if s, ok := g.isOptionOf(t); ok {
			if _, ok := g.isOptionOf(s); ok {
				return "", fmt.Errorf("%w (nested nbt.Option)", ErrUnexpectedType)
			}
			if ex, err := g.getExpectedTagExpression(s); err != nil {
				return "", fmt.Errorf("%w (unexpected nbt.Option type argument %s)", ErrUnexpectedType, util.TypeString(s))
			} else {
				return ex, nil
			}
		}

		p := types.NewPointer(t)
		if !g.isTargetType(t) && !g.implementsUnmarshaller(p) {
			return "", fmt.Errorf("%w (%s does not implement nbt.Unmarshaller)", ErrUnexpectedType, util.TypeString(p))
		}

		return fmt.Sprintf("(*%s)(nil).TagType()", g.typeStringRelative(typ)), nil
	case *types.Pointer:
		elem := t.Elem()
		if named, ok := types.Unalias(elem).(*types.Named); !ok {
			return "", fmt.Errorf("%w (pointer to unexpected type %s)", ErrUnexpectedType, util.TypeString(elem))
		} else {
			if _, ok := g.isOptionOf(named); ok {
				return "", fmt.Errorf("%w (pointer to nbt.Option)", ErrUnexpectedType)
			}
			if !g.isTargetType(named) && !g.implementsUnmarshaller(t) {
				return "", fmt.Errorf("%w (%s does not implement nbt.Unmarshaller)", ErrUnexpectedType, util.TypeString(t))
			}
		}
		return fmt.Sprintf("(%s)(nil).TagType()", g.typeStringRelative(typ)), nil
	default:
		return "", ErrUnexpectedType
	}
}
