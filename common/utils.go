package common

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// func formatSize(file os.FileInfo) string {
// 	if file.IsDir() {
// 		return "-"
// 	}
// 	size := file.Size()
// 	switch {
// 	case size > 1024*1024:
// 		return fmt.Sprintf("%.1f MB", float64(size)/1024/1024)
// 	case size > 1024:
// 		return fmt.Sprintf("%.1f KB", float64(size)/1024)
// 	default:
// 		return strconv.Itoa(int(size)) + " B"
// 	}
// 	return ""
// }

func GetRealIP(req *http.Request) string {
	xip := req.Header.Get("X-Real-IP")
	if xip == "" {
		xip = strings.Split(req.RemoteAddr, ":")[0]
	}
	return xip
}

func SublimeContains(s, substr string) bool {
	rs, rsubstr := []rune(s), []rune(substr)
	if len(rsubstr) > len(rs) {
		return false
	}

	var ok = true
	var i, j = 0, 0
	for ; i < len(rsubstr); i++ {
		found := -1
		for ; j < len(rs); j++ {
			if rsubstr[i] == rs[j] {
				found = j
				break
			}
		}
		if found == -1 {
			ok = false
			break
		}
		j += 1
	}
	return ok
}

// getLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Convert path to normal paths
func CleanPath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}

func getSeparatedPath(path string) []string {
	pathSeparated := []string{}
	tmpStr := ""
	index := 0
	for index < len(path) {
		if path[index] == '*' {
			if index+1 < len(path) && path[index+1] == '*' {
				if tmpStr != "" {
					pathSeparated = append(pathSeparated, tmpStr)
				}
				pathSeparated = append(pathSeparated, "**")
				index += 2
				tmpStr = ""
			} else {
				if tmpStr != "" {
					pathSeparated = append(pathSeparated, tmpStr)
				}
				pathSeparated = append(pathSeparated, "*")
				index += 1
				tmpStr = ""
			}
		} else {
			tmpStr += string(path[index])
		}
		index++
	}
	if tmpStr != "" {
		pathSeparated = append(pathSeparated, tmpStr)
	}
	return pathSeparated
}

func CheckPath(path1, path2 string) bool {
	if len(path1) == 0 || len(path2) == 0 {
		return true
	}
	path1 = strings.TrimPrefix(path1, "/")
	path2 = strings.TrimPrefix(path2, "/")
	if (path1 != "" && path2 == "") || (path1 == "" && path2 != "") {
		if path1 == "*" || path1 == "**" {
			return true
		}
		return false
	}
	path1Separated := getSeparatedPath(path1)

	for len(path1Separated) > 0 && len(path2) > 0 {
		p1 := path1Separated[0]
		if p1 == "*" {
			path1Separated = path1Separated[1:]
			path2 = strings.TrimPrefix(path2, "/")
			index := strings.Index(path2, "/")
			if index == -1 {
				for p1 == "*" || p1 == "**" {
					if len(path1Separated) == 0 {
						return true
					}
					p1 = path1Separated[0]
					path1Separated = path1Separated[1:]
				}
				return false
			}
			path2 = path2[index:]
			continue
		}
		if p1 == "**" {
			for p1 == "*" || p1 == "**" {
				if len(path1Separated) == 1 {
					return true
				}
				path1Separated = path1Separated[1:]
				p1 = path1Separated[0]
			}

			index := strings.Index(path2, p1)
			if index == -1 {
				return false
			}
			path2 = path2[index+len(p1):]
			path1Separated = path1Separated[1:]
			continue
		}
		path2 = strings.TrimPrefix(path2, "/")
		path2 = strings.TrimSuffix(path2, "/")
		p1 = strings.TrimPrefix(p1, "/")
		p1 = strings.TrimSuffix(p1, "/")
		if strings.HasPrefix(strings.TrimPrefix(path2, "/"), strings.TrimPrefix(p1, "/")) {
			path2 = path2[len(p1):]
			path1Separated = path1Separated[1:]
			continue
		}
		return false
	}

	if len(path2) == 0 && len(path1Separated) > 0 {
		path := strings.Join(path1Separated, "")
		path = strings.ReplaceAll(path, "*", "")
		if path != "" {
			return false
		}
	}

	return true
}
