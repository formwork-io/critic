/*
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

See http://formwork-io.github.io/ for more.
*/

package main

import "os"
import "fmt"
import "strconv"
import "strings"

func main() {
	dirs := []string{}
	for i, dir := range os.Args {
		if i == 0 {
			continue
		}
		dirs = append(dirs, dir)
	}
	fanFD, err := Initialize(0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s - are you root?", err)
		os.Exit(1)
	}
	modflags := FAN_MARK_ADD | FAN_MARK_MOUNT
	maskflags := uint64(FAN_MODIFY)
	for _, dir := range dirs {
		fmt.Printf("%s\n", dir)
		err := fanFD.Mark(modflags, maskflags, 0, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
	}
	for {
		evMeta, _ := fanFD.GetEvent()
		selfStr := "/proc/self/fd"
		fdStr := strconv.Itoa(int(evMeta.File.Fd()))
		symlink := strings.Join([]string{selfStr, fdStr}, "/")
		path, err := os.Readlink(symlink)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", path)
		evMeta.File.Close()
	}
}
