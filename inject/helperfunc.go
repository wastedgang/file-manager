package inject

import "reflect"

// getPkgPath 计算包名
func getPkgPath(objectType reflect.Type) string {
	for objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}
	return objectType.PkgPath()
}

// getFullIdentifier 计算类型的标识全名
func getFullIdentifier(objectType reflect.Type) string {
	pkgPath := getPkgPath(objectType)
	if pkgPath == "" {
		return objectType.String()
	}
	return pkgPath + "/" + objectType.String()
}
