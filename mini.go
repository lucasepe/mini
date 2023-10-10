package mini

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	rootSection = "_ROOT_SECTION_"
)

func RootSection() string { return rootSection }

// steal it from stdandard lib's flag pkg, :)

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() any { return bool(*b) }

func (b *boolValue) String() string { return fmt.Sprintf("%v", *b) }

func (b *boolValue) IsBoolFlag() bool { return true }

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Get() any { return int(*i) }

func (i *intValue) String() string { return fmt.Sprintf("%v", *i) }

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() any { return int64(*i) }

func (i *int64Value) String() string { return fmt.Sprintf("%v", *i) }

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() any { return uint(*i) }

func (i *uintValue) String() string { return fmt.Sprintf("%v", *i) }

// -- uint64 Value
type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() any { return uint64(*i) }

func (i *uint64Value) String() string { return fmt.Sprintf("%v", *i) }

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() any { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() any { return float64(*f) }

func (f *float64Value) String() string { return fmt.Sprintf("%v", *f) }

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() any { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// hand made...

// -- []string slice Value
type stringSlice []string

func newStringSliceValue(val []string, p *[]string) *stringSlice {
	*p = val
	return (*stringSlice)(p)
}

func (sa *stringSlice) Set(v string) error {
	*sa = append(*sa, v)
	return nil
}

func (sa *stringSlice) Get() any { return []string(*sa) }

func (sa *stringSlice) String() string {
	return strings.Join(*sa, ",")
}

// -- enum Value
type enumValue struct {
	target  *string
	options []string
}

// New returns a flag.Value implementation for parsing flags with a one-of-a-set value
func newEnumValue(target *string, options ...string) *enumValue {
	return &enumValue{target: target, options: options}
}

func (ef *enumValue) String() string {
	return *ef.target
}

func (ef *enumValue) Set(value string) error {
	for _, v := range ef.options {
		if v == value {
			*ef.target = value
			return nil
		}
	}

	return fmt.Errorf("expected one of %q", ef.options)
}

func (ef *enumValue) Get() any { return []string(ef.options) }

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
type Value interface {
	String() string
	Set(string) error
}

// Getter is an interface that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
type Getter interface {
	Value
	Get() any
}

type ConfSet struct {
	parsed   bool
	sections map[string]*Section
}

type Section struct {
	Name string
	Vals map[string]*Item
}

type Item struct {
	SectionName string
	Name        string
	Val         Value
}

func NewConf() *ConfSet {
	return &ConfSet{false, make(map[string]*Section)}
}

func (c *ConfSet) Var(val Value, sectionName, name string) {
	if strings.TrimSpace(sectionName) == "" {
		sectionName = rootSection
	}

	item := &Item{sectionName, name, val}

	if c.sections == nil {
		c.sections = make(map[string]*Section)
	}

	s, ok := c.sections[sectionName]
	if !ok {
		s = &Section{sectionName, make(map[string]*Item)}
	}

	_, exists := s.Vals[name]
	if exists {
		panic(fmt.Sprintf("item %s already exists", name))
	}

	if s.Vals == nil {
		s.Vals = make(map[string]*Item)
	}

	s.Vals[name] = item
	c.sections[sectionName] = s
}

func (c *ConfSet) BoolVar(p *bool, sectionName, name string, value bool) {
	c.Var(newBoolValue(value, p), sectionName, name)
}

func (c *ConfSet) Bool(sectionName, name string, value bool) *bool {
	p := new(bool)
	c.BoolVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) IntVar(p *int, sectionName, name string, value int) {
	c.Var(newIntValue(value, p), sectionName, name)
}

func (c *ConfSet) Int(sectionName, name string, value int) *int {
	p := new(int)
	c.IntVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) Int64Var(p *int64, sectionName, name string, value int64) {
	c.Var(newInt64Value(value, p), sectionName, name)
}

func (c *ConfSet) Int64(sectionName, name string, value int64) *int64 {
	p := new(int64)
	c.Int64Var(p, sectionName, name, value)
	return p
}

func (c *ConfSet) UintVar(p *uint, sectionName, name string, value uint) {
	c.Var(newUintValue(value, p), sectionName, name)
}

func (c *ConfSet) Uint(sectionName, name string, value uint) *uint {
	p := new(uint)
	c.UintVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) Uint64Var(p *uint64, sectionName, name string, value uint64) {
	c.Var(newUint64Value(value, p), sectionName, name)
}

func (c *ConfSet) Uint64(sectionName, name string, value uint64) *uint64 {
	p := new(uint64)
	c.Uint64Var(p, sectionName, name, value)
	return p
}

func (c *ConfSet) StringVar(p *string, sectionName, name string, value string) {
	c.Var(newStringValue(value, p), sectionName, name)
}

func (c *ConfSet) String(sectionName, name string, value string) *string {
	p := new(string)
	c.StringVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) Float64Var(p *float64, sectionName, name string, value float64) {
	c.Var(newFloat64Value(value, p), sectionName, name)
}

func (c *ConfSet) Float64(sectionName, name string, value float64) *float64 {
	p := new(float64)
	c.Float64Var(p, sectionName, name, value)
	return p
}

func (c *ConfSet) DurationVar(p *time.Duration, sectionName, name string, value time.Duration) {
	c.Var(newDurationValue(value, p), sectionName, name)
}

func (c *ConfSet) Duration(sectionName, name string, value time.Duration) *time.Duration {
	p := new(time.Duration)
	c.DurationVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) StringSliceVar(p *[]string, sectionName, name string, value []string) {
	c.Var(newStringSliceValue(value, p), sectionName, name)
}

func (c *ConfSet) StringSlice(sectionName, name string, value []string) *[]string {
	p := new([]string)
	c.StringSliceVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) EnumVar(p *string, sectionName, name string, value []string) {
	c.Var(newEnumValue(p, value...), sectionName, name)
}

func (c *ConfSet) Enum(sectionName, name string, value []string) *string {
	p := new(string)
	c.EnumVar(p, sectionName, name, value)
	return p
}

func (c *ConfSet) parseOne(sectionName string, line string) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid line format, should be 'key = value'")
	}

	name, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	s, sectionExists := c.sections[sectionName]
	if sectionExists {
		if v, valExists := s.Vals[name]; valExists {
			if err := v.Val.Set(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ConfSet) Parse(in io.Reader) error {
	c.parsed = true
	currentSection := rootSection

	ln := 0
	reader := bufio.NewReader(in)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		ln = ln + 1
		if len(line) == 0 {
			continue
		}

		l := strings.TrimSpace(string(line))
		if l[0] == '#' || l[0] == ';' {
			continue
		}

		// parse section
		if l[0] == '[' {
			l := strings.TrimSpace(l)
			if l[len(l)-1] == ']' {
				currentSection = l[1 : len(l)-1]
				continue
			}
		}

		// parse item
		err = c.parseOne(currentSection, l)
		if err != nil {
			return fmt.Errorf("ln %d: %w", ln, err)
		}
	}

	return nil
}
