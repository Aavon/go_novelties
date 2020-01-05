package nonull

import (
	"reflect"
)

func Make(ptr interface{}) {
	p := reflect.TypeOf(ptr)
	if p.Kind() != reflect.Ptr {
		return
	}
	v := reflect.ValueOf(ptr)
	t := p.Elem()
	if v.IsNil() {
		//if v.CanSet() {
		//	newValue := reflect.New(t)
		//	fmt.Printf("new: %v\n", newValue)
		//	v.Set(newValue)
		//	//fmt.Printf("%v\n", v.Interface())
		//} else {
		//	return
		//}
		return
	}
	v = v.Elem()
	if t.Kind() == reflect.Struct {
		fieldCnt := t.NumField()
		for i := 0; i < fieldCnt; i++ {
			f := v.Field(i)
			ft := t.Field(i)
			if !f.CanSet() {
				continue
			}
			switch f.Kind() {
			case reflect.Ptr:
				if f.IsNil() {
					ftt := ft.Type.Elem()
					//if ftt.Kind() == reflect.Slice {
					//	sli := reflect.MakeSlice(ftt, 5, 10)
					//	sli = reflect.Append(sli, reflect.ValueOf(7))
					//	fmt.Println(sli.Interface())
					//	sli.CanAddr()
					//	nv := reflect.ValueOf(sli.Addr())
					//	fmt.Println("new:::", nv.Interface(), nv.IsNil(), nv.Interface())
					//	f.Set(nv)
					//	fmt.Println("new field:::", f.Interface(), f.IsNil())
					//} else if ftt.Kind() == reflect.Map {
					//	nv := reflect.NewAt(ftt, unsafe.Pointer(reflect.MakeMap(ftt).Pointer()))
					//	fmt.Println("new:::", nv.Interface(), nv.IsNil(), nv.Interface())
					//	f.Set(nv)
					//	fmt.Println("new field:::", f.Interface(), f.IsNil())
					//} else {
					f.Set(reflect.New(ftt))
					Make(f.Interface())
					//}
				}
			case reflect.Struct:
				Make(f.Addr().Interface())
			case reflect.Slice:
				if f.IsNil() {
					f.Set(reflect.MakeSlice(f.Type(), 0, 0))
				}
			case reflect.Map:
				if f.IsNil() {
					f.Set(reflect.MakeMap(f.Type()))
				}
			default:
				continue
			}
		}
	} else if t.Kind() == reflect.Ptr {
		Make(v.Interface())
	} else {
		return
	}
}
